package cron

import (
	"context"
	"mq/infrastructure/svc"

	"mq/infrastructure/provider"
)

// StatisticsUserHandler 统计用户
func StatisticsUserHandler(ctx context.Context, svcCtx *svc.ServiceContext) {
	provider.InitializeUserFacade(svcCtx).UpdateStatistics(ctx)
}
