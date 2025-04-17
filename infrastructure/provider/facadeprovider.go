package provider

import (
	"mq/application/service"
	"mq/application/service/facade"

	"github.com/google/wire"
)

// FacedeServiceProviderSet 门面服务层依赖提供者集合
var FacedeServiceProviderSet = wire.NewSet(
	ProviderAPIKey,
	ProviderUserFacade,
)

// ProviderAPIKey 提供APIKey服务实现
func ProviderAPIKey(subjectService *service.SubjectService) *facade.ApiKeyFacade {
	return facade.NewApiKeyFacade(subjectService)
}

// ProviderUserFacade 提供用户门面服务实现
func ProviderUserFacade(patientService *service.PatientService, appointmentService *service.AppointmentService, healthService *service.HealthService, operationService *service.OperationService, statisticsService *service.StatisticsService, userService *service.UserService) *facade.UserFacade {
	return facade.NewUserFacade(patientService, appointmentService, healthService, operationService, statisticsService, userService)
}
