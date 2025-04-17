package provider

import (
	"mq/application/consumer"
	"mq/application/service"
	"mq/application/service/facade"
	"mq/domain/appointment"
	"mq/domain/health"
	"mq/domain/operation"
	"mq/domain/patient"
	"mq/domain/statistics"
	"mq/domain/subject"
	"mq/domain/user"
	"mq/infrastructure/integration/mail"
	"mq/infrastructure/integration/qiniu"
	"mq/infrastructure/integration/sms"
	"mq/infrastructure/queue"
	"mq/infrastructure/svc"

	"github.com/google/wire"
)

// ServiceProviderSet 服务层依赖提供者集合
var ServiceProviderSet = wire.NewSet(

	ProviderUserService,
	ProviderSmsService,
	ProviderSubjectService,
	ProviderAppointmentService,
	ProviderHealthService,
	ProviderPatientService,
	ProviderOperationService,
	ProviderStatisticsService,
	ProviderEmailService,
	ProviderToolService,
	ProviderQueueManager,
	ProviderHealthConsumer,
	ProviderConsumerManager,
)

// ProviderUserService 提供用户服务实现
func ProviderUserService(svcCtx *svc.ServiceContext, userRepo user.UserRepository, smsService sms.SmsInterface, emailService mail.EmailInterface) *service.UserService {
	return service.NewUserService(userRepo, svcCtx, smsService, emailService)
}

// ProviderSubjectService 提供课程服务实现
func ProviderSubjectService(subjectRepo subject.SubjectRepository) *service.SubjectService {
	return service.NewSubjectService(subjectRepo)
}

// ProviderSmsService 提供短信服务实现
func ProviderSmsService(svcCtx *svc.ServiceContext) sms.SmsInterface {
	return &sms.Sns{SvcCtx: svcCtx}
}

// ProviderAppointmentService 提供预约服务实现
func ProviderAppointmentService(appointmentRepo appointment.AppointmentRepository) *service.AppointmentService {
	return service.NewAppointmentService(appointmentRepo)
}

// ProviderHealthService 提供健康服务实现
func ProviderHealthService(healthRepo health.HealthRepository, queueManager queue.QueueManager) *service.HealthService {
	return service.NewHealthService(healthRepo, queueManager)
}

// ProviderPatientService 提供患者服务实现
func ProviderPatientService(patientRepo patient.PatientRepository) *service.PatientService {
	return service.NewPatientService(patientRepo)
}

// ProviderOperationService 提供操作记录服务实现
func ProviderOperationService(operationRepo operation.OperationRepository) *service.OperationService {
	return service.NewOperationService(operationRepo)
}

// ProviderStatisticsService 提供统计服务实现
func ProviderStatisticsService(statisticsRepo statistics.StatisticsRepository) *service.StatisticsService {
	return service.NewStatisticsService(statisticsRepo)
}

// ProviderEmailService 提供邮件服务实现
func ProviderEmailService(svcCtx *svc.ServiceContext) mail.EmailInterface {
	return &mail.EmailSender{SvcCtx: svcCtx}
}

// ProviderToolService 提供工具服务实现
func ProviderToolService(uploadRepo qiniu.UploadInterface) *service.ToolService {
	return service.NewToolService(uploadRepo)
}

// ProviderQueueManager 提供队列管理器实现
func ProviderQueueManager(svcCtx *svc.ServiceContext) queue.QueueManager {
	return queue.NewDQQueue(svcCtx)
}

func ProviderHealthConsumer(queueManager queue.QueueManager, userFacade *facade.UserFacade) *consumer.HealthConsumer {
	return consumer.NewHealthConsumer(queueManager, userFacade)
}

// ProviderConsumerManager 提供消费者管理器实现
func ProviderConsumerManager(queueManager queue.QueueManager) *queue.ConsumerManager {
	return queue.NewConsumerManager(queueManager)
}
