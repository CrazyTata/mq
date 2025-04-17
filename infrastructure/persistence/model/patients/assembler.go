package patients

import (
	"mq/domain/patient"
)

func POTODOGetPatients(res *Patients) *patient.Patients {
	return &patient.Patients{
		Id:          res.Id,
		FriendlyId:  res.FriendlyId,
		Name:        res.Name,
		Age:         res.Age,
		Gender:      res.Gender,
		Phone:       res.Phone,
		Status:      res.Status,
		LastVisit:   res.LastVisit,
		Avatar:      res.Avatar,
		History:     res.History,
		Allergies:   res.Allergies,
		Note:        res.Note,
		Attachments: res.Attachments,
		Details:     res.Details,
		UserId:      res.UserId,
		IsDeleted:   res.IsDeleted,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}
}

func DOTOPOPatients(res *patient.Patients) *Patients {
	return &Patients{
		Id:          res.Id,
		FriendlyId:  res.FriendlyId,
		Name:        res.Name,
		Age:         res.Age,
		Gender:      res.Gender,
		Phone:       res.Phone,
		Status:      res.Status,
		LastVisit:   res.LastVisit,
		Avatar:      res.Avatar,
		History:     res.History,
		Allergies:   res.Allergies,
		Note:        res.Note,
		Attachments: res.Attachments,
		Details:     res.Details,
		UserId:      res.UserId,
		IsDeleted:   res.IsDeleted,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}
}

func POTODOGetPatientsList(res []Patients) []*patient.Patients {
	resp := make([]*patient.Patients, 0, len(res))
	for _, v := range res {
		resp = append(resp, POTODOGetPatients(&v))
	}
	return resp
}
