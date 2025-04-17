package appointment

import (
	"database/sql"
	"time"
)

// Appointment 预约领域模型
type Appointments struct {
	Id          int64
	PatientId   int64          // 患者ID
	PatientName string         // 患者姓名
	Date        string         // 日期
	Time        string         // 时间
	Duration    string         // 时长
	Type        string         // 类型
	Status      string         // 状态
	Notes       sql.NullString // 备注
	UserId      string         // 用户ID
	IsDeleted   int64          // 是否删除
	CreatedAt   time.Time      // 创建时间
	UpdatedAt   time.Time      // 更新时间
}

func Create(patientId int64, patientName, date, time, duration, patientType, status, notes, userId string) *Appointments {
	return &Appointments{
		PatientId:   patientId,
		PatientName: patientName,
		Date:        date,
		Time:        time,
		Duration:    duration,
		Type:        patientType,
		Status:      status,
		Notes:       sql.NullString{String: notes, Valid: notes != ""},
		UserId:      userId,
	}
}

func (a *Appointments) Update(patientName string, date, time, duration string, patientType string, status string, notes string) {
	a.PatientName = patientName
	a.Date = date
	a.Time = time
	a.Duration = duration
	a.Type = patientType
	a.Status = status
	a.Notes = sql.NullString{String: notes, Valid: notes != ""}
}

func (a *Appointments) MarkAsDeleted() {
	a.IsDeleted = 1
}
