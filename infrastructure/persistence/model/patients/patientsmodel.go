package patients

import (
	"context"
	"database/sql"
	"mq/domain/patient"
	"mq/infrastructure/persistence"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ patient.PatientRepository = (*customPatientsModel)(nil)

type (
	// PatientsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPatientsModel.
	PatientsModel interface {
		patientsModel
	}

	customPatientsModel struct {
		*defaultPatientsModel
	}
)

// NewPatientsModel returns a model for the database table.
func NewPatientsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) patient.PatientRepository {
	return &customPatientsModel{
		defaultPatientsModel: newPatientsModel(conn, c, opts...),
	}
}

// Create 创建患者
func (m *customPatientsModel) Create(ctx context.Context, patient *patient.Patients) (int64, error) {
	po := DOTOPOPatients(patient)
	res, err := m.Insert(ctx, po)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetByID 根据ID获取患者
func (m *customPatientsModel) GetByID(ctx context.Context, id int64) (*patient.Patients, error) {
	rowBuilder := persistence.RowBuilder(m.table, patientsRows)
	query, values, err := rowBuilder.Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}
	var res Patients
	err = m.QueryRowNoCacheCtx(ctx, &res, query, values...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return POTODOGetPatients(&res), nil
}

// GetByFriendlyID 根据友好ID获取患者
func (m *customPatientsModel) GetByFriendlyID(ctx context.Context, friendlyID string) (*patient.Patients, error) {
	rowBuilder := persistence.RowBuilder(m.table, patientsRows)
	query, values, err := rowBuilder.Where(squirrel.Eq{"friendly_id": friendlyID}).ToSql()
	if err != nil {
		return nil, err
	}
	var res Patients
	err = m.QueryRowNoCacheCtx(ctx, &res, query, values...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return POTODOGetPatients(&res), nil
}

// GetByPhone 根据手机号获取患者
func (m *customPatientsModel) GetByPhone(ctx context.Context, phone string) (*patient.Patients, error) {
	rowBuilder := persistence.RowBuilder(m.table, patientsRows)
	query, values, err := rowBuilder.Where(squirrel.Eq{"phone": phone}).ToSql()
	if err != nil {
		return nil, err
	}
	var res Patients
	err = m.QueryRowNoCacheCtx(ctx, &res, query, values...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return POTODOGetPatients(&res), nil
}

// GetByUserID 根据用户ID获取患者列表
func (m *customPatientsModel) GetByUserID(ctx context.Context, userID string, page, pageSize int64) ([]*patient.Patients, int64, error) {
	rowBuilder := persistence.RowBuilder(m.table, patientsRows)
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
	var res []Patients
	err = m.QueryRowsNoCacheCtx(ctx, &res, query1, values1...)
	if err != nil {
		return nil, 0, err
	}
	return POTODOGetPatientsList(res), count, nil
}

// UpdatePatient 更新患者信息
func (m *customPatientsModel) UpdatePatient(ctx context.Context, patient *patient.Patients) error {
	po := DOTOPOPatients(patient)
	return m.Update(ctx, po)
}

// GetList 获取患者列表
func (m *customPatientsModel) GetList(ctx context.Context, userID, search, status string, order, page, pageSize int64) ([]*patient.Patients, int64, error) {
	rowBuilder := persistence.RowBuilder(m.table, patientsRows)
	countBuilder := persistence.CountBuilder(m.table)
	if userID != "" {
		rowBuilder = rowBuilder.Where(squirrel.Eq{"user_id": userID})
		countBuilder = countBuilder.Where(squirrel.Eq{"user_id": userID})
	}
	if search != "" {
		//通过FriendlyId和name来匹配
		rowBuilder = rowBuilder.Where(squirrel.Or{squirrel.Eq{"friendly_id": search}, squirrel.Eq{"name": search}})
		countBuilder = countBuilder.Where(squirrel.Or{squirrel.Eq{"friendly_id": search}, squirrel.Eq{"name": search}})
	}
	if status != "" {
		rowBuilder = rowBuilder.Where(squirrel.Eq{"status": status})
		countBuilder = countBuilder.Where(squirrel.Eq{"status": status})
	}
	if order == 1 {
		rowBuilder = rowBuilder.OrderBy("id DESC")
	} else {
		rowBuilder = rowBuilder.OrderBy("id ASC")
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
	var res []Patients
	err = m.QueryRowsNoCacheCtx(ctx, &res, query1, values1...)
	if err != nil {
		return nil, 0, err
	}
	return POTODOGetPatientsList(res), count, nil
}

// CountByUserID 根据用户ID获取患者数量
func (m *customPatientsModel) CountByUserID(ctx context.Context, userID string) (int64, error) {
	countBuilder := persistence.CountBuilder(m.table).Where(squirrel.Eq{"user_id": userID})
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

// CountActiveByUserID 根据用户ID获取活跃患者数量
func (m *customPatientsModel) CountActiveByUserID(ctx context.Context, userID string) (int64, error) {
	countBuilder := persistence.CountBuilder(m.table).Where(squirrel.Eq{"user_id": userID}).Where(squirrel.Eq{"status": "active"})
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
