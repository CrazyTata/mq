package assembler

import (
	"mq/application/dto"
	domain "mq/domain/appointment"
	appointmentModel "mq/infrastructure/persistence/model/appointments"
	"time"

	"github.com/samber/lo"
)

func DOTOPOAppointment(do *domain.Appointments) *appointmentModel.Appointments {
	return &appointmentModel.Appointments{
		Id:          do.Id,
		PatientId:   do.PatientId,
		PatientName: do.PatientName,
		Date:        do.Date,
		Time:        do.Time,
		Duration:    do.Duration,
		Type:        do.Type,
		Status:      do.Status,
		Notes:       do.Notes,
		UserId:      do.UserId,
		IsDeleted:   do.IsDeleted,
		CreatedAt:   do.CreatedAt,
		UpdatedAt:   do.UpdatedAt,
	}
}

func POTODOAppointment(po *appointmentModel.Appointments) *domain.Appointments {
	return &domain.Appointments{
		Id:          po.Id,
		PatientId:   po.PatientId,
		PatientName: po.PatientName,
		Date:        po.Date,
		Time:        po.Time,
		Duration:    po.Duration,
		Type:        po.Type,
		Status:      po.Status,
		Notes:       po.Notes,
		UserId:      po.UserId,
		IsDeleted:   po.IsDeleted,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
	}
}

func DOTODTOAppointment(do *domain.Appointments) *dto.AppointmentResponse {
	if do == nil {
		return nil
	}
	return &dto.AppointmentResponse{
		ID:          do.Id,
		PatientID:   do.PatientId,
		PatientName: do.PatientName,
		Date:        do.Date,
		Time:        do.Time,
		Duration:    do.Duration,
		Type:        do.Type,
		Status:      do.Status,
		Notes:       do.Notes.String,
		CreatedAt:   do.CreatedAt.Format(time.DateTime),
		UpdatedAt:   do.UpdatedAt.Format(time.DateTime),
	}
}

func DOTODTOAppointmentList(do []*domain.Appointments, total int64) *dto.AppointmentListResponse {
	return &dto.AppointmentListResponse{
		Total: total,
		List: lo.Map(do, func(do *domain.Appointments, _ int) *dto.AppointmentResponse {
			return DOTODTOAppointment(do)
		}),
	}
}
