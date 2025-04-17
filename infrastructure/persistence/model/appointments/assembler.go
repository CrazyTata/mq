package appointments

import (
	"mq/domain/appointment"
)

func POTODOGetAppointment(res *Appointments) *appointment.Appointments {
	return &appointment.Appointments{
		Id:          res.Id,
		PatientId:   res.PatientId,
		PatientName: res.PatientName,
		Date:        res.Date,
		Time:        res.Time,
		Duration:    res.Duration,
		Type:        res.Type,
		Status:      res.Status,
		Notes:       res.Notes,
		UserId:      res.UserId,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}
}

func DOTOPOAppointment(res *appointment.Appointments) *Appointments {
	return &Appointments{
		Id:          res.Id,
		PatientId:   res.PatientId,
		PatientName: res.PatientName,
		Date:        res.Date,
		Time:        res.Time,
		Duration:    res.Duration,
		Type:        res.Type,
		Status:      res.Status,
		Notes:       res.Notes,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
		IsDeleted:   res.IsDeleted,
		UserId:      res.UserId,
	}
}

func POTODOGetAppointments(res []Appointments) []*appointment.Appointments {
	resp := make([]*appointment.Appointments, 0, len(res))
	for _, v := range res {
		resp = append(resp, POTODOGetAppointment(&v))
	}
	return resp
}
