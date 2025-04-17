package sms

import "context"

type SmsInterface interface {
	SendSms(ctx context.Context, mobile string) error
	CheckSms(ctx context.Context, mobile, code string) (error, bool)
}

var NotCheckMobile = map[string]struct{}{
	"18888888888": {},
	"18913533664": {},
}
var NotCheckCode string = "888888"
