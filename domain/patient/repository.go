package patient

import (
	"context"
)

// PatientRepository 患者仓储接口
type PatientRepository interface {
	// Create 创建患者
	Create(ctx context.Context, patient *Patients) (int64, error)
	// GetByID 根据ID获取患者
	GetByID(ctx context.Context, id int64) (*Patients, error)
	// GetByFriendlyID 根据友好ID获取患者
	GetByFriendlyID(ctx context.Context, friendlyID string) (*Patients, error)
	// GetByPhone 根据手机号获取患者
	GetByPhone(ctx context.Context, phone string) (*Patients, error)
	// GetByUserID 根据用户ID获取患者列表
	GetByUserID(ctx context.Context, userID string, page, pageSize int64) ([]*Patients, int64, error)
	// UpdatePatient 更新患者信息
	UpdatePatient(ctx context.Context, patient *Patients) error
	// GetList 获取患者列表
	GetList(ctx context.Context, userID, search, status string, order, page, pageSize int64) ([]*Patients, int64, error)
	// CountByUserID 根据用户ID获取患者数量
	CountByUserID(ctx context.Context, userID string) (int64, error)
	// CountActiveByUserID 根据用户ID获取活跃患者数量
	CountActiveByUserID(ctx context.Context, userID string) (int64, error)
}
