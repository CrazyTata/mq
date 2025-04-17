package assembler

import (
	"mq/application/dto"
	domain "mq/domain/patient"
	patientModel "mq/infrastructure/persistence/model/patients"
	"time"

	"github.com/samber/lo"
)

func DOTOPOPatient(do *domain.Patients) *patientModel.Patients {
	return &patientModel.Patients{
		Id:          do.Id,
		FriendlyId:  do.FriendlyId,
		UserId:      do.UserId,
		Phone:       do.Phone,
		Name:        do.Name,
		Age:         int64(do.Age),
		Gender:      do.Gender,
		Avatar:      do.Avatar,
		Status:      do.Status,
		LastVisit:   do.LastVisit,
		Attachments: do.Attachments,
		History:     do.History,
		Allergies:   do.Allergies,
		Note:        do.Note,
		Details:     do.Details,
		CreatedAt:   do.CreatedAt,
		UpdatedAt:   do.UpdatedAt,
	}
}

func POTODOPatient(po *patientModel.Patients) *domain.Patients {
	return &domain.Patients{
		Id:          po.Id,
		FriendlyId:  po.FriendlyId,
		UserId:      po.UserId,
		Phone:       po.Phone,
		Name:        po.Name,
		Age:         int64(po.Age),
		Gender:      po.Gender,
		Avatar:      po.Avatar,
		Status:      po.Status,
		LastVisit:   po.LastVisit,
		Attachments: po.Attachments,
		History:     po.History,
		Allergies:   po.Allergies,
		Note:        po.Note,
		Details:     po.Details,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
	}
}

func DOTODTOPatient(do *domain.Patients) *dto.Patient {
	if do == nil {
		return nil
	}
	return &dto.Patient{
		Id:          do.Id,
		FriendlyId:  do.FriendlyId,
		Phone:       do.GetDecryptedPhone(),
		Name:        do.Name,
		Age:         int(do.Age),
		Gender:      do.Gender,
		Avatar:      do.Avatar.String,
		Status:      do.Status,
		History:     do.History.String,
		Allergies:   do.Allergies.String,
		Note:        do.Note.String,
		LastVisit:   do.LastVisit.Time.Format(time.DateTime),
		Attachments: do.Attachments.String,
		Details:     do.Details.String,
		CreatedAt:   do.CreatedAt.Format(time.DateTime),
		UpdatedAt:   do.UpdatedAt.Format(time.DateTime),
	}
}

func DOTODTOPatients(dos []*domain.Patients, total int64) *dto.PatientListResponse {
	return &dto.PatientListResponse{
		Total: total,
		List: lo.Map(dos, func(do *domain.Patients, _ int) *dto.Patient {
			return DOTODTOPatient(do)
		}),
	}
}
