package cron

import (
	"context"
	"mq/infrastructure/svc"

	"github.com/robfig/cron/v3"
)

func RegisterCrontab(svcCtx *svc.ServiceContext) {
	c := cron.New()
	ctx := context.Background()
	//每10分钟执行一次
	_, _ = c.AddFunc("*/3 * * * *", func(ctx context.Context, svcCtx *svc.ServiceContext) func() {
		capturedCtx := ctx
		capturedSvcCtx := svcCtx
		return func() {
			StatisticsUserHandler(capturedCtx, capturedSvcCtx)
		}
	}(ctx, svcCtx))

	c.Start()
}
