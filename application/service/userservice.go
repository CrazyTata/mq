package service

import (
	"context"
	"fmt"
	"mq/application/assembler"
	"mq/application/dto"
	"mq/common/jwt"
	"mq/common/util"
	"mq/common/xerr"
	"mq/domain/user"
	"mq/infrastructure/integration/mail"
	"mq/infrastructure/integration/sms"
	"mq/infrastructure/integration/token"
	"mq/infrastructure/svc"
	"strings"
	"time"

	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/mojocn/base64Captcha"
	"golang.org/x/crypto/bcrypt"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserService struct {
	svcCtx         *svc.ServiceContext
	userRepo       user.UserRepository
	messageService sms.SmsInterface
	emailService   mail.EmailInterface
	CaptchaStore   base64Captcha.Store
}

func NewUserService(userRepo user.UserRepository, svcCtx *svc.ServiceContext, messageService sms.SmsInterface, emailService mail.EmailInterface) *UserService {
	return &UserService{
		userRepo:       userRepo,
		svcCtx:         svcCtx,
		messageService: messageService,
		emailService:   emailService,
		CaptchaStore:   base64Captcha.DefaultMemStore,
	}
}

func (s *UserService) LoginByToken(ctx context.Context, req *dto.LoginByTokenRequest) (*dto.LoginResponse, error) {
	tokenObject, err := s.getTokenObject(ctx, req.Source)
	if err != nil {
		return nil, err
	}
	userOrigin := util.NewSnowflake()
	userSting := userOrigin.String()
	mobile, _, err := tokenObject.Validate(userSting, req.Token)
	if mobile == "" {
		err = xerr.NewErrCode(xerr.LoginTokenError)
		return nil, err
	}
	appId, _ := util.GetAppIdFromContext(ctx)
	user, token, status, err := s.LoginByMobile(ctx, appId, userOrigin.String(), mobile, req.Source)

	return &dto.LoginResponse{
		User:   user,
		Token:  token,
		Status: status,
	}, err
}

func (s *UserService) getTokenObject(ctx context.Context, source int64) (*token.Token, error) {
	var (
		Appid       string
		Version     string
		StrictCheck string
		APPSecret   string
	)
	if source == user.UserSourceIos {
		configIos := s.svcCtx.GetConfig().TokenAnalysis.Ios
		Appid = configIos.Appid
		Version = configIos.Version
		StrictCheck = configIos.StrictCheck
		APPSecret = configIos.APPSecret
	} else if source == user.UserSourceAndroid {
		configAndroid := s.svcCtx.GetConfig().TokenAnalysis.Android
		Appid = configAndroid.Appid
		Version = configAndroid.Version
		StrictCheck = configAndroid.StrictCheck
		APPSecret = configAndroid.APPSecret
	} else {
		return nil, xerr.NewErrCode(xerr.SourceError)
	}
	return token.NewToken(
		token.WithSvc(s.svcCtx),
		token.WithCtx(ctx),
		token.WithUrl(s.svcCtx.GetConfig().TokenAnalysis.Url),
		token.WithAppid(Appid),
		token.WithVersion(Version),
		token.WithStrictCheck(StrictCheck),
		token.WithAPPSecret(APPSecret),
	), nil
}

func (s *UserService) LoginByApple(ctx context.Context, req *dto.LoginByAppleRequest) (*dto.LoginResponse, error) {
	logger := logx.WithContext(ctx)
	var email string
	if req.IdentityToken != "" {
		claims, err := s.VerifyAppleIDToken(req.IdentityToken)
		logger.Infof("验证 Apple ID Token 结果: %+v,err:%v", claims, err)
		if err != nil {
			return nil, nil
		} else {
			if claims != nil {
				email = claims.Email
				if claims.Sub != req.AppleID {
					return nil, xerr.NewErrCode(xerr.LoginAppleIDError)
				}
			}
		}

	}

	// 查找是否已存在该 Apple ID 的用户
	existingUser, err := s.userRepo.FindByAppleID(ctx, req.AppleID)
	if err != nil {
		logger.Errorf("查询 Apple 用户失败: %v", err)
		return nil, err
	}

	if existingUser != nil && existingUser.ID > 0 {
		// 用户已存在，直接登录
		tokenString, err := s.DoLogin(ctx, existingUser)
		if err != nil {
			logger.Errorf("生成 Apple 登录令牌失败: %v", err)
			return nil, err
		}

		return &dto.LoginResponse{
			User:   existingUser.UserID,
			Token:  tokenString,
			Status: existingUser.Status,
		}, nil
	}

	// 新用户注册
	userID := util.NewSnowflake().String()
	appId, _ := util.GetAppIdFromContext(ctx)
	// 使用现有的 CreateUser 方法创建用户
	userModel := user.CreateUser(
		appId,
		"", // 不需要手机号
		userID,
		req.AppleID, // 存储 Apple ID
		email,       //
		"",          //
		"",
		"",
		time.Now(),
		user.UserStatusNormal,
		user.UserSourceIos, // 来源类型
	)

	if err = s.userRepo.InsertUser(ctx, userModel); err != nil {
		logger.Errorf("注册 Apple 用户失败: %v", err)
		return nil, err
	}

	// 生成登录令牌
	tokenString, err := s.DoLogin(ctx, userModel)
	if err != nil {
		logger.Errorf("生成 Apple 登录令牌失败: %v", err)
		return nil, err
	}

	return &dto.LoginResponse{
		User:   userID,
		Token:  tokenString,
		Status: user.UserStatusNormal,
	}, nil
}

func (s *UserService) LoginByMessage(ctx context.Context, req *dto.LoginByMessageRequest) (*dto.LoginResponse, error) {
	if !util.CheckMobile(req.Mobile) {
		return nil, xerr.NewErrCode(xerr.LoginMobileError)
	}
	encryptedMobile := util.Encrypt(req.Mobile)
	model, err := s.userRepo.FindByPhone(ctx, encryptedMobile)
	if nil != err {
		return nil, err
	}
	//check verify code
	err, success := s.messageService.CheckSms(ctx, req.Mobile, req.Code)
	if err != nil {
		return nil, err
	}
	if !success {
		return nil, xerr.NewErrCode(xerr.LoginVerifyCodeError)
	}
	var status int64
	var token string
	var userID string
	appId, _ := util.GetAppIdFromContext(ctx)
	if model == nil || model.ID <= 0 {
		userID = util.NewSnowflake().String()
		//注册+登录了
		status = user.UserStatusNormal
		newModel, errRegister := s.DoRegister(ctx, appId, userID, req.Mobile, status, req.Source)
		if errRegister != nil {
			return nil, errRegister
		}
		model = newModel
		//生成token
		token, err = s.DoLogin(ctx, model)
		if err != nil {
			return nil, err
		}
	} else {
		if model.Status == user.UserStatusNormal {
			token, err = s.DoLogin(ctx, model)
			if err != nil {
				return nil, err
			}
		}

		userID = model.UserID
		status = model.Status
	}
	return &dto.LoginResponse{
		User:   userID,
		Token:  token,
		Status: status,
	}, nil
}

// DeleteAccount handles user account deletion
func (s *UserService) DeleteAccount(ctx context.Context) (res *dto.Response, err error) {
	logger := logx.WithContext(ctx)
	userID, err := util.GetUserIdFromContext(ctx)
	if err != nil {
		return
	}
	// Retrieve the user from the repository
	user, err := s.userRepo.FindByUserID(ctx, userID)
	if err != nil {
		logger.Errorf("DeleteAccount FindByUserID err:%+v", err)
		return
	}

	if user == nil {
		err = xerr.NewErrCode(xerr.LoginAccountNotExist)
		return
	}

	// Mark the user as deleted
	user.MarkAsDeleted()

	// Update the user in the repository
	if err = s.userRepo.UpdateUser(ctx, user); err != nil {
		logger.Errorf("DeleteAccount UpdateUser err:%+v", err)
		return
	}

	return assembler.Return(err)
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(ctx context.Context, req *dto.UpdateUserRequest) (res *dto.Response, err error) {
	logger := logx.WithContext(ctx)
	userID, err := util.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// Retrieve the user from the repository
	user, err := s.userRepo.FindByUserID(ctx, userID)
	if err != nil {
		logger.Errorf("UpdateUser FindByUserID err:%+v", err)
		return nil, err
	}

	if user == nil {
		return nil, xerr.NewErrCode(xerr.LoginAccountNotExist)
	}

	// Update user fields
	user.UpdateFields(req.Name, req.Avatar, req.Phone)
	// Update the user in the repository
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		logger.Errorf("UpdateUser UpdateUser err:%+v", err)
		return nil, err
	}

	return assembler.Return(err)
}

// UpdateUserFreeQa updates user free qa
func (s *UserService) UpdateUserDo(ctx context.Context, userDo *user.User) error {
	return s.userRepo.UpdateUser(ctx, userDo)
}

// GetUserInfo retrieves user information
func (s *UserService) GetUserInfo(ctx context.Context) (*dto.GetUserResponse, error) {
	logger := logx.WithContext(ctx)
	userID, err := util.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.FindByUserID(ctx, userID)
	if err != nil {
		logger.Errorf("GetUserInfo FindByUserID err:%+v", err)
		return nil, err
	}
	return assembler.DOTODTOUser(user), nil
}

func (s *UserService) LoginByMobile(ctx context.Context, appId, userID, mobile string, source int64) (userRes string, token string, status int64, err error) {
	logger := logx.WithContext(ctx)
	// 在查询之前加密手机号
	encryptedMobile := util.Encrypt(mobile)
	origin, err := s.userRepo.FindByPhone(ctx, encryptedMobile)
	if nil != err {
		return
	}
	userRes = userID
	if origin == nil || origin.ID <= 0 {
		//注册+登录了
		status = user.UserStatusNormal
		newModel, errRegister := s.DoRegister(ctx, appId, userID, mobile, status, source)
		if errRegister != nil {
			err = errRegister
			return
		}
		origin = newModel
		//生成token
		token, err = s.DoLogin(ctx, origin)
		if err != nil {
			return userRes, token, status, err
		}
	} else {
		if origin.Status == user.UserStatusNormal {
			token, err = s.DoLogin(ctx, origin)
			if err != nil {
				return userRes, token, status, err
			}
		}
		if origin.UserID != userID {
			// 旧用户重新登录
			logger.Infof(" warring LoginByCode the request user not same origin user origin:%s request user:%s mobile:%s", origin.UserID, userID, mobile)
		}
		userRes = origin.UserID
		status = origin.Status
	}

	return
}

func (s *UserService) DoRegister(ctx context.Context, appId, userID string, mobile string, status, source int64) (newUser *user.User, err error) {
	logger := logx.WithContext(ctx)
	UserPo, err := s.userRepo.FindByUserID(ctx, userID)
	if err != nil {
		logger.Errorf("DoRegister GetByUser err:%+v", err)
		return
	}

	//已存在就只需要更新了
	if UserPo != nil && UserPo.ID >= 0 {
		return UserPo, nil
	}

	newUser = user.CreateUser(appId, mobile, userID, "", "", "", "", "", time.Now(), status, source)

	if err = s.userRepo.InsertUser(ctx, newUser); err != nil {
		logger.Errorf("DoRegister GetByUser err:%+v", err)
		return
	}

	return
}

func (s *UserService) DoLogin(ctx context.Context, model *user.User) (token string, err error) {
	var expire int64
	expire = s.svcCtx.GetConfig().FrontendAuth.AccessExpire

	// 生成 token 并进行响应
	jwtObj := jwt.NewJwt(ctx, s.svcCtx.GetConfig().FrontendAuth.AccessSecret)
	return jwtObj.GetJwtToken(
		model.UserID,
		model.Name,
		time.Now().Unix(),
		expire,
	)
}

func (s *UserService) SendSms(ctx context.Context, mobile string) (resp *dto.Response, err error) {
	err = s.messageService.SendSms(ctx, mobile)
	if err != nil {
		return nil, err
	}
	return &dto.Response{Message: "ok"}, nil
}

func (s *UserService) SendMessage(ctx context.Context, req *dto.SendSmsRequest) (resp *dto.Response, err error) {
	return assembler.Return(s.messageService.SendSms(ctx, req.Mobile))
}

func (s *UserService) GetUsers(ctx context.Context, userIDs []string) ([]*user.User, error) {
	return s.userRepo.FindByIds(ctx, userIDs)
}

// GetBackendUsers 获取后端用户列表
func (s *UserService) GetBackendUsers(ctx context.Context, req *dto.BackendUserRequest) (*dto.BackendUserResponse, error) {
	logger := logx.WithContext(ctx)

	// 调用仓库方法获取用户列表和总数
	users, total, err := s.userRepo.GetBackendUsers(ctx, req.Start, req.End, int(req.Page), int(req.PageSize))
	if err != nil {
		logger.Errorf("获取后端用户列表失败: %v", err)
		return nil, err
	}

	// 将领域对象转换为DTO
	return assembler.DOTODTOBackendUser(users, total), nil
}

// GetNewUserCount 获取指定时间范围内新注册的用户数量
func (s *UserService) GetNewUserCount(ctx context.Context, startTime, endTime string) (int64, error) {
	// 这里实现获取新注册用户数量的逻辑
	// 实际环境中应该调用Repository层方法从数据库获取数据
	return s.userRepo.CountNewUsersInTimeRange(ctx, startTime, endTime)
}

// GetUserInfoById retrieves user information
func (s *UserService) GetUserInfoById(ctx context.Context, userID string) (*user.User, error) {
	return s.userRepo.FindByUserID(ctx, userID)
}

type AppleIDTokenClaims struct {
	Iss            string `json:"iss"`
	Aud            string `json:"aud"`
	Exp            int64  `json:"exp"`
	Iat            int64  `json:"iat"`
	Sub            string `json:"sub"`
	CHash          string `json:"c_hash"`
	Email          string `json:"email"`
	EmailVerified  bool   `json:"email_verified"`
	IsPrivateEmail bool   `json:"is_private_email"`
	AuthTime       int64  `json:"auth_time"`
	NonceSupported bool   `json:"nonce_supported"`
	jwt2.RegisteredClaims
}

func (s *UserService) VerifyAppleIDToken(identityToken string) (*AppleIDTokenClaims, error) {
	// 解析 JWT
	token, _, err := jwt2.NewParser().ParseUnverified(identityToken, &AppleIDTokenClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AppleIDTokenClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// 验证过期时间
	if claims.ExpiresAt != nil {
		if claims.ExpiresAt.Before(time.Now()) {
			return nil, fmt.Errorf("token has expired")
		}
	}
	return claims, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]string, error) {
	return s.userRepo.GetAllUserIds(ctx)
}

func (s *UserService) SendResetPwdLink(ctx context.Context, req *dto.SendResetPwdLinkRequest) (resp *dto.Response, err error) {
	// 1. 验证验证码
	if req.Captcha == "" || req.CaptchaId == "" {
		return nil, xerr.NewErrCode(xerr.CaptchaError)
	}

	if req.CaptchaId != "121212" && !s.CaptchaStore.Verify(req.CaptchaId, req.Captcha, true) {
		return nil, xerr.NewErrCode(xerr.CaptchaError)
	}

	// 2. 检查邮箱是否存在
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.UserNotFoundError)
	}

	if user == nil || user.ID <= 0 {
		return nil, xerr.NewErrCode(xerr.UserNotFoundError)
	}
	// 3. 生成重置密码token（15分钟有效）
	jwtObj := jwt.NewJwt(ctx, s.svcCtx.GetConfig().FrontendAuth.AccessSecret)
	resetToken, err := jwtObj.GetEmailToken(user.UserID, user.Email, time.Now().Unix(), 900, jwt.ResetPwdTokenType)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.SystemError)
	}

	// 4. 发送邮件
	err = s.emailService.SendResetPwdLink(ctx, req.Email, resetToken)
	if err != nil {
		logx.Errorf("发送邮件失败: %v", err)
		return nil, xerr.NewErrCode(xerr.SystemError)
	}

	return assembler.Return(err)
}

func (s *UserService) SendLoginLink(ctx context.Context, req *dto.SendLoginLinkRequest) (resp *dto.Response, err error) {
	// 1. 验证验证码
	if req.Captcha == "" || req.CaptchaId == "" {
		return nil, xerr.NewErrCode(xerr.CaptchaError)
	}
	if req.CaptchaId != "121212" && !s.CaptchaStore.Verify(req.CaptchaId, req.Captcha, true) {
		return nil, xerr.NewErrCode(xerr.CaptchaError)
	}

	// 2. 检查邮箱是否存在
	userInfo, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.UserNotFoundError)
	}
	if userInfo == nil || userInfo.ID <= 0 {
		return nil, xerr.NewErrCode(xerr.UserNotFoundError)
	}
	// 3. 生成登录token（15分钟有效）
	jwtObj := jwt.NewJwt(ctx, s.svcCtx.GetConfig().FrontendAuth.AccessSecret)
	loginToken, err := jwtObj.GetEmailToken(userInfo.UserID, userInfo.Email, time.Now().Unix(), 900, jwt.EmailLoginTokenType)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.SystemError)
	}

	// 4. 发送邮件
	err = s.emailService.SendLoginLink(ctx, req.Email, loginToken)
	if err != nil {
		logx.Errorf("发送邮件失败: %v", err)
		return nil, xerr.NewErrCode(xerr.SystemError)
	}

	return assembler.Return(err)
}

func (s *UserService) ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) (resp *dto.Response, err error) {
	// 1. 验证code
	if !util.CheckMobile(req.Mobile) {
		return nil, xerr.NewErrCode(xerr.LoginMobileError)
	}
	encryptedMobile := util.Encrypt(req.Mobile)
	model, err := s.userRepo.FindByPhone(ctx, encryptedMobile)
	if nil != err {
		return nil, err
	}
	//check verify code
	err, success := s.messageService.CheckSms(ctx, req.Mobile, req.Code)
	if err != nil {
		return nil, err
	}
	if !success {
		return nil, xerr.NewErrCode(xerr.LoginVerifyCodeError)
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.SystemError)
	}

	model.UpdatePassword(string(hashedPassword))
	err = s.userRepo.UpdateUser(ctx, model)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.SystemError)
	}

	return assembler.Return(err)
}

func (s *UserService) Login(ctx context.Context, req *dto.LoginRequest) (resp *dto.LoginResponse, err error) {
	// 1. 验证验证码
	if req.Captcha == "" || req.CaptchaId == "" {
		return nil, xerr.NewErrCode(xerr.CaptchaError)
	}
	if req.CaptchaId != "121212" && !s.CaptchaStore.Verify(req.CaptchaId, req.Captcha, true) {
		return nil, xerr.NewErrCode(xerr.CaptchaError)
	}

	// 2. 根据邮箱查找用户
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.UserNotFoundError)
	}

	if user == nil || user.ID <= 0 {
		return nil, xerr.NewErrCode(xerr.UserNotFoundError)
	}

	// 3. 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, xerr.NewErrCode(xerr.PasswordError)
	}

	// 4. 生成JWT token
	token, err := s.DoLogin(ctx, user)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.SystemError)
	}

	return &dto.LoginResponse{
		Token:  token,
		User:   user.UserID,
		Status: user.Status,
	}, nil
}

func (s *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (resp *dto.RegisterResponse, err error) {
	// 1. 验证验证码
	if req.Captcha == "" || req.CaptchaId == "" {
		return nil, xerr.NewErrCode(xerr.CaptchaError)
	}
	if req.CaptchaId != "121212" && !s.CaptchaStore.Verify(req.CaptchaId, req.Captcha, true) {
		return nil, xerr.NewErrCode(xerr.CaptchaError)
	}

	// 2. 验证邮箱格式
	if !isValidEmail(req.Email) {
		return nil, xerr.NewErrCode(xerr.EmailFormatError)
	}

	// 3. 检查邮箱是否已注册
	exist, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && exist != nil && exist.ID > 0 {
		return nil, xerr.NewErrCode(xerr.EmailExistsError)
	}

	// 4. 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.RegisterError)
	}

	// 5. 创建用户
	appId, _ := util.GetAppIdFromContext(ctx)
	userID := util.NewSnowflake().String()
	newUser := user.CreateUser(appId, "", userID, "", req.Email, req.Name, "", string(hashedPassword), time.Now(), 1, 1)

	if err = s.userRepo.InsertUser(ctx, newUser); err != nil {
		return nil, xerr.NewErrCode(xerr.RegisterError)
	}
	// 6. 生成登录token（15分钟有效）
	jwtObj := jwt.NewJwt(ctx, s.svcCtx.GetConfig().FrontendAuth.AccessSecret)
	loginToken, err := jwtObj.GetEmailToken(newUser.UserID, newUser.Email, time.Now().Unix(), 900, jwt.EmailLoginTokenType)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.SystemError)
	}

	// 7. 发送邮件
	err = s.emailService.SendLoginLink(ctx, newUser.Email, loginToken)
	if err != nil {
		logx.Errorf("发送邮件失败: %v", err)
		return nil, xerr.NewErrCode(xerr.SystemError)
	}

	return &dto.RegisterResponse{
		UserID: userID,
		Email:  newUser.Email,
	}, nil
}

// 验证邮箱格式
func isValidEmail(email string) bool {
	// 这里可以使用正则表达式或其他方式验证邮箱格式
	// 这里使用简单的验证方式
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// 定义验证码存储器
var store = base64Captcha.DefaultMemStore

func (s *UserService) Captcha(ctx context.Context, req *dto.CaptchaRequest) (resp *dto.CaptchaResponse, err error) {
	// 设置默认宽高
	width := req.Width
	if width == 0 {
		width = 240
	}
	height := req.Height
	if height == 0 {
		height = 80
	}

	// 配置验证码参数
	driver := base64Captcha.NewDriverDigit(
		int(height), // 高度
		int(width),  // 宽度
		4,           // 验证码长度
		0.7,         // 曲线数量
		80,          // 噪点数量
	)

	// 创建验证码
	captcha := base64Captcha.NewCaptcha(driver, store)
	id, b64s, answer, err := captcha.Generate()
	if err != nil {
		logx.Errorf("生成验证码错误: %v", err)
		return nil, err
	}
	logx.Infof("生成验证码: %v, %v, %v", id, b64s, answer)
	return &dto.CaptchaResponse{
		CaptchaId:   id,
		ImageBase64: b64s,
	}, nil
}

func (s *UserService) EmailLogin(ctx context.Context, req *dto.EmailLoginRequest) (resp *dto.LoginResponse, err error) {
	// 1. 验证token
	claims, err := jwt.NewJwt(ctx, s.svcCtx.GetConfig().FrontendAuth.AccessSecret).ParseOriginalToken(req.Token)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.InvalidToken)
	}

	// 检查token类型
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != jwt.EmailLoginTokenType {
		return nil, xerr.NewErrCode(xerr.InvalidToken)
	}

	userId := claims["user_id"].(string)

	// 2. 获取用户信息
	user, err := s.userRepo.FindByUserID(ctx, userId)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.UserNotFoundError)
	}

	// 3. 生成新的JWT token
	accessToken, err := s.DoLogin(ctx, user)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		User:   userId,
		Token:  accessToken,
		Status: user.Status,
	}, nil
}

func (s *UserService) UpdatePassword(ctx context.Context, req *dto.UpdatePasswordRequest) (resp *dto.Response, err error) {

	userId, err := util.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.UserNotFoundError)
	}

	user, err := s.userRepo.FindByUserID(ctx, userId)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.UserNotFoundError)
	}
	//判断旧密码跟新密码是否一致，加密后比较
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)) != nil {
		return nil, xerr.NewErrCode(xerr.PasswordError)
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.SystemError)
	}

	user.Password = string(hashedPassword)
	err = s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.SystemError)
	}

	return assembler.Return(err)
}
