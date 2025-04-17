package dto

type (
	LoginByTokenRequest struct {
		Token  string `json:"token" validate:"required"`
		Source int64  `json:"source"`
	}

	LoginByMessageRequest struct {
		User   string `json:"user" validate:"required"`
		Mobile string `json:"mobile" validate:"required"`
		Code   string `json:"code" validate:"required,len=6"`
		Source int64  `json:"source"`
	}

	LoginByAppleRequest struct {
		AppleID       string `json:"apple_id" validate:"required"`
		IdentityToken string `json:"identity_token,optional"`
	}

	LoginResponse struct {
		Token  string `json:"token"`
		User   string `json:"user_id"`
		Status int64  `json:"status"`
	}

	SendSmsRequest struct {
		Mobile string `json:"mobile" validate:"required"`
	}

	GetCaptchaResponse struct {
		CaptchaId string `json:"captcha_id"`
		Captcha   string `json:"captcha"`
	}

	UpdateUserRequest struct {
		Name   string `json:"name,optional"`
		Phone  string `json:"phone,optional"`
		Avatar string `json:"avatar,optional"`
	}

	GetUserResponse struct {
		ID             int64  `json:"id"`
		AppId          string `json:"app_id"`
		UserID         string `json:"user_id"`
		Phone          string `json:"phone"`
		AppleID        string `json:"apple_id"`
		Email          string `json:"email"`
		IsPrivateEmail int64  `json:"is_private_email"`
		Name           string `json:"name"`
		Avatar         string `json:"avatar"`
		Password       string `json:"password"`
		RegisterTime   string `json:"register_time"`
		Status         int64  `json:"status"`
		Source         int64  `json:"source"`
		CreatedAt      string `json:"created_at"`
	}

	AddCourseRequest struct {
		CourseID int64 `json:"course_id"`
	}

	AddCourseResponse struct {
		ID int64 `json:"id"`
	}

	GetMyCoursesResponse struct {
		Data []*MyCourse `json:"data"`
	}

	MyCourse struct {
		Id                int64  `json:"id"`
		CourseID          int64  `json:"course_id"`
		Title             string `json:"title"`
		Summary           string `json:"summary"` // 摘要
		Description       string `json:"description"`
		CompletedDuration int64  `json:"completed_duration"` // 完成时长
		TotalDuration     int64  `json:"total_duration"`     // 总时长
		Total             int64  `json:"total"`              // 总课程数
		Price             int64  `json:"price"`              // 价格
		IsFree            int64  `json:"is_free"`            // 是否免费
		Category          int64  `json:"category"`           // 1 Blog文章 2 书籍片段 3 学术文献 4 其他
		CourseType        int64  `json:"course_type"`        // 授课形式:1 视频 2 音频 3 图文
		CoverUrl          string `json:"cover_url"`          // 封面URL
		IsActive          int64  `json:"is_active"`          // 是否上架
		ContentUrl        string `json:"content_url"`        // 内容URL
		CreatedAt         string `json:"created_at"`         // 创建时间
	}

	MyFavoriteCountResponse struct {
		Count int64 `json:"count"`
	}

	GetNotificationResponse struct {
		Data []*Notification `json:"data"`
	}

	Notification struct {
		Id        int64  `json:"id"`
		Type      int64  `json:"type"`       // 1 course update 2 subscription expired
		RelatedId int64  `json:"related_id"` // 关联ID
		UserId    string `json:"user_id"`
		Title     string `json:"title"`      // 标题
		Content   string `json:"content"`    // 内容
		ReadTime  string `json:"read_time"`  // 阅读时间
		Status    int64  `json:"status"`     // 1 unread 2 read
		CreatedAt string `json:"created_at"` // 创建时间
		UpdatedAt string `json:"updated_at"` // 更新时间
	}

	ReadNotificationRequest struct {
		NotificationID int64 `json:"notification_id"`
	}

	BackendUser struct {
		UserID       string `json:"user_id"`
		Phone        string `json:"phone"`
		AppleID      string `json:"apple_id"`
		Email        string `json:"email"`
		Name         string `json:"name"`
		Avatar       string `json:"avatar"`
		RegisterTime string `json:"register_time"`
		Source       int64  `json:"source"`
		CreatedAt    string `json:"created_at"`
	}

	BackendUserResponse struct {
		Data  []*BackendUser `json:"data"`
		Total int64          `json:"total"`
	}

	BackendUserRequest struct {
		Page     int64  `json:"page"`
		PageSize int64  `json:"page_size"`
		Start    string `json:"start"`
		End      string `json:"end"`
	}

	RegisterRequest struct {
		Email     string `json:"email"`     // 用户邮箱地址，用于登录和找回密码
		Password  string `json:"password"`  // 用户密码，长度至少6位
		Name      string `json:"name"`      // 用户名
		Captcha   string `json:"captcha"`   // 图形验证码内容
		CaptchaId string `json:"captchaId"` // 验证码唯一标识
	}

	RegisterResponse struct {
		UserID string `json:"user"`  // 用户ID
		Email  string `json:"email"` // 用户邮箱
	}

	LoginRequest struct {
		Email     string `json:"email"`     // 用户邮箱地址
		Password  string `json:"password"`  // 用户密码
		Captcha   string `json:"captcha"`   // 图形验证码内容
		CaptchaId string `json:"captchaId"` // 验证码唯一标识
	}

	CaptchaRequest struct {
		Width  int64 `json:"width,optional"`  // 验证码图片宽度，默认240
		Height int64 `json:"height,optional"` // 验证码图片高度，默认80
	}

	CaptchaResponse struct {
		CaptchaId   string `json:"captchaId"`   // 验证码唯一标识，用于验证
		ImageBase64 string `json:"imageBase64"` // Base64编码的验证码图片
	}

	SendLoginLinkRequest struct {
		Email     string `json:"email"`     // 用户邮箱地址
		Captcha   string `json:"captcha"`   // 图形验证码内容
		CaptchaId string `json:"captchaId"` // 验证码唯一标识
	}

	SendLoginLinkResponse struct {
		Message string `json:"message"` // 发送结果消息
	}

	EmailLoginRequest struct {
		Token string `json:"token"` // 邮箱登录令牌
	}

	EmailLoginResponse struct {
		Id    int64  `json:"id"`    // 用户ID
		Email string `json:"email"` // 用户邮箱
		Token string `json:"token"` // JWT令牌
	}

	SendResetPwdLinkRequest struct {
		Email     string `json:"email"`     // 用户邮箱地址
		Captcha   string `json:"captcha"`   // 图形验证码内容
		CaptchaId string `json:"captchaId"` // 验证码唯一标识
	}

	SendResetPwdLinkResponse struct {
		Message string `json:"message"` // 发送结果消息
	}

	ResetPasswordRequest struct {
		Code     string `json:"code"` // 重置密码验证码
		Mobile   string `json:"mobile"`
		Password string `json:"password"` // 新密码
	}

	UpdatePasswordRequest struct {
		Password    string `json:"password"`     // 新密码
		OldPassword string `json:"old_password"` // 旧密码
	}

	ResetPasswordResponse struct {
		Message string `json:"message"` // 重置结果消息
	}

	UserProfileRequest struct {
	}

	UserProfileResponse struct {
		Id      int64  `json:"id"` // 用户ID
		Mobile  string `json:"mobile"`
		Email   string `json:"email"`   // 用户邮箱
		Diamond int64  `json:"diamond"` // 钻石数量
	}
)
