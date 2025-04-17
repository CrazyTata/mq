package provider

import (
	"mq/domain/appointment"
	"mq/domain/health"
	"mq/domain/operation"
	"mq/domain/patient"
	"mq/domain/statistics"
	"mq/domain/subject"
	"mq/domain/user"
	"mq/infrastructure/integration/qiniu"
	appointmentModel "mq/infrastructure/persistence/model/appointments"
	healthModel "mq/infrastructure/persistence/model/health_records"
	operationModel "mq/infrastructure/persistence/model/operation_records"
	patientModel "mq/infrastructure/persistence/model/patients"
	statisticsModel "mq/infrastructure/persistence/model/statistics"
	subjectModel "mq/infrastructure/persistence/model/subject"
	userModel "mq/infrastructure/persistence/model/user"
	"mq/infrastructure/svc"

	"github.com/google/wire"
)

// RepositoryProviderSet 仓储层依赖提供者集合
var RepositoryProviderSet = wire.NewSet(
	ProviderUserRepo,
	ProviderSubjectRepo,
	ProviderAppSharingRepo,
	ProviderOperationRepo,
	ProviderHealthRepo,
	ProviderPatientRepo,
	ProviderAppointmentRepo,
	ProviderStatisticsRepo,
	ProviderUploadRepo,
)

// ProviderUserRepo 提供用户仓储实现
func ProviderUserRepo(svcCtx *svc.ServiceContext) user.UserRepository {
	return userModel.NewUserModel(svcCtx.GetDB(), svcCtx.GetCache())
}

// ProviderSubjectRepo 提供应用仓储实现
func ProviderSubjectRepo(svcCtx *svc.ServiceContext) subject.SubjectRepository {
	return subjectModel.NewSubjectModel(svcCtx.GetDB(), svcCtx.GetCache())
}

// ProviderAppSharingRepo 提供应用共享仓储实现
func ProviderAppSharingRepo(svcCtx *svc.ServiceContext) subject.AppSharingRepository {
	return subjectModel.NewAppSharingModel(svcCtx.GetDB(), svcCtx.GetCache())
}

// ProviderOperationRepo 提供操作记录仓储实现
func ProviderOperationRepo(svcCtx *svc.ServiceContext) operation.OperationRepository {
	return operationModel.NewOperationRecordsModel(svcCtx.GetDB(), svcCtx.GetCache())
}

// ProviderHealthRepo 提供健康记录仓储实现
func ProviderHealthRepo(svcCtx *svc.ServiceContext) health.HealthRepository {
	return healthModel.NewHealthRecordsModel(svcCtx.GetDB(), svcCtx.GetCache())
}

// ProviderPatientRepo 提供患者仓储实现
func ProviderPatientRepo(svcCtx *svc.ServiceContext) patient.PatientRepository {
	return patientModel.NewPatientsModel(svcCtx.GetDB(), svcCtx.GetCache())
}

// ProviderAppointmentRepo 提供预约仓储实现
func ProviderAppointmentRepo(svcCtx *svc.ServiceContext) appointment.AppointmentRepository {
	return appointmentModel.NewAppointmentsModel(svcCtx.GetDB(), svcCtx.GetCache())
}

// ProviderStatisticsRepo 提供统计仓储实现
func ProviderStatisticsRepo(svcCtx *svc.ServiceContext) statistics.StatisticsRepository {
	return statisticsModel.NewStatisticsModel(svcCtx.GetDB(), svcCtx.GetCache())
}

// ProviderUploadRepo 提供上传仓储实现
func ProviderUploadRepo(svcCtx *svc.ServiceContext) qiniu.UploadInterface {
	return qiniu.NewUploader(svcCtx)
}
