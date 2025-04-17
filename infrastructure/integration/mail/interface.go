package mail

import "context"

type EmailInterface interface {
	SendLoginLink(ctx context.Context, email, token string) error
	SendResetPwdLink(ctx context.Context, email, token string) error
}
