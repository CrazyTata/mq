package token

import (
	"context"
	"mq/infrastructure/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type (
	Token struct {
		appid       string
		version     string
		strictCheck string
		appSecret   string
		url         string
		logger      logx.Logger
		ctx         context.Context
		svcCtx      *svc.ServiceContext
	}

	Option func(*Token)
)

func NewToken(opts ...Option) *Token {

	client := &Token{}
	for _, opt := range opts {
		opt(client)
	}
	return client
}
