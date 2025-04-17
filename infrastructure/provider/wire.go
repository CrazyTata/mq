//go:build wireinject
// +build wireinject

package provider

import (
	"mq/application/consumer"
	"mq/application/service"
	"mq/application/service/facade"
	"mq/infrastructure/queue"
	"mq/infrastructure/svc"

	"github.com/google/wire"
)

// userProviderSet 定义用户服务相关的依赖提供者
var providerSet = wire.NewSet(

	// 基础设施层
	RepositoryProviderSet,

	// 应用服务层
	ServiceProviderSet,

	// 门面服务层
	FacedeServiceProviderSet,
)

// InitializeUserFacade 初始化统计门面服务
func InitializeUserFacade(svcCtx *svc.ServiceContext) *facade.UserFacade {
	wire.Build(providerSet)
	return nil // 这个返回值会被wire自动生成的代码替换
}

// InitializeAppointmentService 初始化预约服务
func InitializeAppointmentService(svcCtx *svc.ServiceContext) *service.AppointmentService {
	wire.Build(providerSet)
	return nil // 这个返回值会被wire自动生成的代码替换
}

// InitializeHealthService 初始化健康服务
func InitializeHealthService(svcCtx *svc.ServiceContext) *service.HealthService {
	wire.Build(providerSet)
	return nil // 这个返回值会被wire自动生成的代码替换
}

// InitializeOperationService 初始化操作记录服务
func InitializeOperationService(svcCtx *svc.ServiceContext) *service.OperationService {
	wire.Build(providerSet)
	return nil // 这个返回值会被wire自动生成的代码替换
}

// InitializePatientService 初始化患者服务
func InitializePatientService(svcCtx *svc.ServiceContext) *service.PatientService {
	wire.Build(providerSet)
	return nil // 这个返回值会被wire自动生成的代码替换
}

// InitializeStatisticsService 初始化统计服务
func InitializeStatisticsService(svcCtx *svc.ServiceContext) *service.StatisticsService {
	wire.Build(providerSet)
	return nil // 这个返回值会被wire自动生成的代码替换
}

// InitializeUserService 初始化用户服务
func InitializeUserService(svcCtx *svc.ServiceContext) *service.UserService {
	wire.Build(providerSet)
	return nil // 这个返回值会被wire自动生成的代码替换
}

// InitializeToolService 初始化工具服务
func InitializeToolService(svcCtx *svc.ServiceContext) *service.ToolService {
	wire.Build(providerSet)
	return nil // 这个返回值会被wire自动生成的代码替换
}

// InitializeHealthConsumer 初始化健康消费者
func InitializeHealthConsumer(svcCtx *svc.ServiceContext) *consumer.HealthConsumer {
	wire.Build(providerSet)
	return nil // 这个返回值会被wire自动生成的代码替换
}

// InitializeConsumerManager 初始化消费者管理器
func InitializeConsumerManager(svcCtx *svc.ServiceContext) *queue.ConsumerManager {
	wire.Build(providerSet)
	return nil // 这个返回值会被wire自动生成的代码替换
}
