package operation

import (
	"context"
)

// OperationRepository 操作记录仓储接口
type OperationRepository interface {
	// Create 创建操作记录
	Create(ctx context.Context, patient *OperationRecords) (int64, error)
	// GetByUserID 根据用户ID获取操作记录列表
	GetByUserID(ctx context.Context, userID string, page, pageSize int64) ([]*OperationRecords, int64, error)
	// Update 更新操作记录信息
	UpdateOperation(ctx context.Context, patient *OperationRecords) error
	// Search 搜索操作记录
	GetList(ctx context.Context, userID string, keyword string, page, pageSize int64) ([]*OperationRecords, int64, error)
	// CountRecentByUserID 根据用户ID获取操作记录数量
	CountRecentByUserID(ctx context.Context, userID string, days int) (int64, error)
}
