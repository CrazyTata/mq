package sms

import (
	"context"
	"errors"
	"fmt"
	"mq/common/redis"
	"mq/common/util"
	"mq/common/xerr"
	"mq/infrastructure/svc"

	url1 "net/url"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

const SMSTemplate = "【叮当】您的验证码为{code}，该验证码30分钟内有效，请勿泄漏于他人！"

type Sns struct {
	SvcCtx *svc.ServiceContext
}

func (s *Sns) SendSms(ctx context.Context, mobile string) (err error) {
	if mobile == "" {
		return xerr.NewErrCode(xerr.MobileError)
	}
	if _, ok := NotCheckMobile[mobile]; ok {
		return
	}
	if !util.CheckMobile(mobile) {
		return xerr.NewErrCode(xerr.MobileError)
	}

	code := util.GetRandomNum(6)
	redis.Rdb.Set(ctx, fmt.Sprintf(redis.LoginCodeKey, mobile), code, redis.LoginCodeExpireTime)

	message := strings.Replace(s.SvcCtx.GetConfig().MSG.SMSTemplate, "{code}", code, 1)
	url := util.NewUrlBuilder(s.SvcCtx.GetConfig().MSG.Url).
		AddParam("u", s.SvcCtx.GetConfig().MSG.Account).
		AddParam("p", s.SvcCtx.GetConfig().MSG.ApiKey).
		AddParam("m", mobile).
		AddParam("c", message).
		Build()

	body, err := util.GetV2(url)
	if err != nil {
		return err
	}
	decodedStr, _ := url1.QueryUnescape(url)
	logx.WithContext(ctx).Infof("SendSms url:%s decode-url:%s response:%+v err:%+v", url, decodedStr, body, err)
	if body != "0" {
		logx.WithContext(ctx).Errorf("SendSms error response:%s", body)
		return errors.New("message send error")
	}
	return nil
}

func (s *Sns) CheckSms(ctx context.Context, mobile, code string) (err error, success bool) {
	if _, ok := NotCheckMobile[mobile]; ok && code == NotCheckCode {
		return nil, true
	}
	cacheCode, err1 := redis.Rdb.Get(ctx, fmt.Sprintf(redis.LoginCodeKey, mobile)).Result()
	if err1 == nil && cacheCode == code {
		redis.Rdb.Del(ctx, fmt.Sprintf(redis.LoginCodeKey, mobile))
		return nil, true
	}
	return
}
