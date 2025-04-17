package appointment

import (
	"context"
)

// AppointmentRepository 预约仓储接口
type AppointmentRepository interface {
	// Create 创建预约
	Create(ctx context.Context, appointment *Appointments) (int64, error)
	// GetByID 根据ID获取预约
	GetByID(ctx context.Context, id int64) (*Appointments, error)
	// GetByUserID 根据用户ID获取预约列表
	GetByUserID(ctx context.Context, userID string, page, pageSize int64) ([]*Appointments, int64, error)
	// Update 更新预约信息
	UpdateAppointment(ctx context.Context, appointment *Appointments) error
	// GetList 获取预约列表
	GetList(ctx context.Context, userID string, search string, searchType, order, patientId, page, pageSize int64) ([]*Appointments, int64, error)
	// CountTodayByUserID 根据用户ID获取今日预约数量
	CountTodayByUserID(ctx context.Context, userID string) (int64, error)
	// CountUpcomingByUserID 根据用户ID获取即将进行的预约数量
	CountUpcomingByUserID(ctx context.Context, userID string) (int64, error)
}
