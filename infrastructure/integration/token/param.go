package token

import (
	"context"
	"mq/infrastructure/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

func WithUrl(url string) Option {
	return func(r *Token) {
		r.url = url
	}
}

func WithAppid(s string) Option {
	return func(r *Token) {
		r.appid = s
	}
}

func WithCtx(ctx context.Context) Option {
	return func(r *Token) {
		r.ctx = ctx
		r.logger = logx.WithContext(ctx)
	}
}

func WithSvc(svc *svc.ServiceContext) Option {
	return func(r *Token) {
		r.svcCtx = svc
	}
}

func WithVersion(s string) Option {
	return func(r *Token) {
		r.version = s
	}
}

func WithStrictCheck(s string) Option {
	return func(r *Token) {
		r.strictCheck = s
	}
}

func WithAPPSecret(s string) Option {
	return func(r *Token) {
		r.appSecret = s
	}
}
