package operation_records

import (
	"context"
	"mq/domain/operation"
	"mq/infrastructure/persistence"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ operation.OperationRepository = (*customOperationRecordsModel)(nil)

type (
	// OperationRecordsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOperationRecordsModel.
	OperationRecordsModel interface {
		operationRecordsModel
	}

	customOperationRecordsModel struct {
		*defaultOperationRecordsModel
	}
)

// NewOperationRecordsModel returns a model for the database table.
func NewOperationRecordsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) operation.OperationRepository {
	return &customOperationRecordsModel{
		defaultOperationRecordsModel: newOperationRecordsModel(conn, c, opts...),
	}
}

// Create 创建操作记录
func (m *customOperationRecordsModel) Create(ctx context.Context, patient *operation.OperationRecords) (int64, error) {
	po := DOTOPOOperationRecords(patient)
	res, err := m.Insert(ctx, po)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetByUserID 根据用户ID获取操作记录列表
func (m *customOperationRecordsModel) GetByUserID(ctx context.Context, userID string, page, pageSize int64) ([]*operation.OperationRecords, int64, error) {
	rowBuilder := persistence.RowBuilder(m.table, operationRecordsRows)
	countBuilder := persistence.CountBuilder(m.table)
	if userID != "" {
		rowBuilder = rowBuilder.Where(squirrel.Eq{"user_id": userID})
		countBuilder = countBuilder.Where(squirrel.Eq{"user_id": userID})
	}
	var count int64
	query, values, err := countBuilder.ToSql()
	if err != nil {
		return nil, 0, err
	}
	err = m.QueryRowNoCacheCtx(ctx, &count, query, values...)
	if err != nil {
		return nil, 0, err
	}
	if count == 0 {
		return nil, 0, nil
	}
	query1, values1, err := rowBuilder.Limit(uint64(pageSize)).Offset(uint64((page - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, 0, err
	}
	var res []OperationRecords
	err = m.QueryRowsNoCacheCtx(ctx, &res, query1, values1...)
	if err != nil {
		return nil, 0, err
	}
	return POTODOGetOperationRecordsList(res), count, nil
}

// Update 更新操作记录信息
func (m *customOperationRecordsModel) UpdateOperation(ctx context.Context, patient *operation.OperationRecords) error {
	po := DOTOPOOperationRecords(patient)
	return m.Update(ctx, po)
}

// Search 搜索操作记录
func (m *customOperationRecordsModel) GetList(ctx context.Context, userID string, keyword string, page, pageSize int64) ([]*operation.OperationRecords, int64, error) {
	rowBuilder := persistence.RowBuilder(m.table, operationRecordsRows)
	countBuilder := persistence.CountBuilder(m.table)
	if userID != "" {
		rowBuilder = rowBuilder.Where(squirrel.Eq{"user_id": userID})
		countBuilder = countBuilder.Where(squirrel.Eq{"user_id": userID})
	}
	var count int64
	query, values, err := countBuilder.ToSql()
	if err != nil {
		return nil, 0, err
	}
	err = m.QueryRowNoCacheCtx(ctx, &count, query, values...)
	if err != nil {
		return nil, 0, err
	}
	if count == 0 {
		return nil, 0, nil
	}
	query1, values1, err := rowBuilder.Limit(uint64(pageSize)).Offset(uint64((page - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, 0, err
	}
	var res []OperationRecords
	err = m.QueryRowsNoCacheCtx(ctx, &res, query1, values1...)
	if err != nil {
		return nil, 0, err
	}
	return POTODOGetOperationRecordsList(res), count, nil
}

// CountRecentByUserID 根据用户ID获取操作记录数量
func (m *customOperationRecordsModel) CountRecentByUserID(ctx context.Context, userID string, days int) (int64, error) {
	countBuilder := persistence.CountBuilder(m.table).Where(squirrel.Eq{"user_id": userID})
	countBuilder = countBuilder.Where(squirrel.Gt{"created_at": time.Now().AddDate(0, 0, -days)})
	var count int64
	query, values, err := countBuilder.ToSql()
	if err != nil {
		return 0, err
	}
	err = m.QueryRowNoCacheCtx(ctx, &count, query, values...)
	if err != nil {
		return 0, err
	}
	return count, nil
}
