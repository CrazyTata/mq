package patient

import (
	"database/sql"
	"errors"
	"mq/common/util"
	"time"
)

var (
	ErrInvalidPatientName   = errors.New("invalid patient name")
	ErrInvalidPatientAge    = errors.New("invalid patient age")
	ErrInvalidPatientGender = errors.New("invalid patient gender")
	ErrInvalidPatientPhone  = errors.New("invalid patient phone")
	ErrInvalidPatientStatus = errors.New("invalid patient status")
)

// Patient 患者领域模型
type Patients struct {
	Id          int64
	FriendlyId  string         // 友好ID
	Name        string         // 姓名
	Age         int64          // 年龄
	Gender      string         // 性别
	Phone       string         // 手机号
	Status      string         // 状态
	LastVisit   sql.NullTime   // 最后就诊时间
	Avatar      sql.NullString // 头像
	History     sql.NullString // 病史
	Allergies   sql.NullString // 过敏史
	Note        sql.NullString // 备注
	Attachments sql.NullString // 附件
	Details     sql.NullString // 详情
	UserId      string         // 用户ID
	IsDeleted   int64          // 是否删除
	CreatedAt   time.Time      // 创建时间
	UpdatedAt   time.Time      // 更新时间
}

// 1 正常 2 需要关注 3 情况紧急
const (
	PatientStatusNormal    = "normal"
	PatientStatusAttention = "attention"
	PatientStatusUrgent    = "urgent"
)

// CreatePatient 创建新患者
func Create(name string, age int, gender, phone string, userID, status, avatar, history, allergies, note, attachments, details string) *Patients {
	now := time.Now()
	friendlyID := util.NewSnowflake().String()
	return &Patients{
		FriendlyId:  friendlyID,
		Name:        name,
		Age:         int64(age),
		Gender:      gender,
		Phone:       util.Encrypt(phone),
		Status:      status,
		UserId:      userID,
		LastVisit:   sql.NullTime{Time: now, Valid: true},
		Avatar:      sql.NullString{String: avatar, Valid: avatar != ""},
		History:     sql.NullString{String: history, Valid: history != ""},
		Allergies:   sql.NullString{String: allergies, Valid: allergies != ""},
		Note:        sql.NullString{String: note, Valid: note != ""},
		Attachments: sql.NullString{String: attachments, Valid: attachments != ""},
		Details:     sql.NullString{String: details, Valid: details != ""},
		CreatedAt:   now,
	}
}

// Validate 验证患者信息
func (p *Patients) Validate() error {
	if p.Name == "" {
		return ErrInvalidPatientName
	}
	if p.Age <= 0 {
		return ErrInvalidPatientAge
	}
	if p.Gender == "" {
		return ErrInvalidPatientGender
	}
	if p.Phone == "" {
		return ErrInvalidPatientPhone
	}
	return nil
}

// UpdateProfile 更新患者信息
func (p *Patients) UpdateProfile(name string, age int, gender, phone, avatar, history, allergies, note, attachments, details string) {
	p.Name = name

	p.Age = int64(age)

	p.Gender = gender

	p.Phone = util.Encrypt(phone)

	p.Avatar = sql.NullString{String: avatar, Valid: avatar != ""}

	p.History = sql.NullString{String: history, Valid: history != ""}

	p.Allergies = sql.NullString{String: allergies, Valid: allergies != ""}

	p.Note = sql.NullString{String: note, Valid: note != ""}

	p.Attachments = sql.NullString{String: attachments, Valid: attachments != ""}

	p.Details = sql.NullString{String: details, Valid: details != ""}
	p.LastVisit = sql.NullTime{Time: time.Now(), Valid: true}
	p.UpdatedAt = time.Now()
}

// UpdateStatus 更新患者状态
func (p *Patients) UpdateStatus(status string) error {
	if status != "active" && status != "inactive" && status != "archived" {
		return ErrInvalidPatientStatus
	}
	p.Status = status
	p.UpdatedAt = time.Now()
	return nil
}

// UpdateLastVisit 更新最后就诊时间
func (p *Patients) UpdateLastVisit() {
	p.LastVisit = sql.NullTime{Time: time.Now()}
	p.UpdatedAt = time.Now()
}

// MarkAsDeleted 标记为删除
func (p *Patients) MarkAsDeleted() {
	p.IsDeleted = 1
	p.UpdatedAt = time.Now()
}

// GetDecryptedPhone 获取解密后的手机号
func (p *Patients) GetDecryptedPhone() string {
	return util.Decrypt(p.Phone)
}
