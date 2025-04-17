package appointments

import (
	"context"
	"database/sql"
	"mq/domain/appointment"
	"mq/infrastructure/persistence"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ appointment.AppointmentRepository = (*customAppointmentsModel)(nil)

type (
	// AppointmentsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAppointmentsModel.
	AppointmentsModel interface {
		appointmentsModel
	}

	customAppointmentsModel struct {
		*defaultAppointmentsModel
	}
)

// NewAppointmentsModel returns a model for the database table.
func NewAppointmentsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) appointment.AppointmentRepository {
	return &customAppointmentsModel{
		defaultAppointmentsModel: newAppointmentsModel(conn, c, opts...),
	}
}

//// Create 创建预约
// Create(ctx context.Context, appointment *Appointments) (int64, error)

func (m *customAppointmentsModel) Create(ctx context.Context, appointment *appointment.Appointments) (int64, error) {
	po := DOTOPOAppointment(appointment)
	res, err := m.Insert(ctx, po)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// // GetByID 根据ID获取预约
// GetByID(ctx context.Context, id int64) (*Appointments, error)

func (m *customAppointmentsModel) GetByID(ctx context.Context, id int64) (*appointment.Appointments, error) {
	rowBuilder := persistence.RowBuilder(m.table, appointmentsRows)
	query, values, err := rowBuilder.Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}
	var res Appointments
	err = m.QueryRowNoCacheCtx(ctx, &res, query, values...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return POTODOGetAppointment(&res), nil
}

// // GetByUserID 根据用户ID获取预约列表
// GetByUserID(ctx context.Context, userID int64, page, pageSize int64) ([]*Appointments, int64, error)

func (m *customAppointmentsModel) GetByUserID(ctx context.Context, userID string, page, pageSize int64) ([]*appointment.Appointments, int64, error) {
	rowBuilder := persistence.RowBuilder(m.table, appointmentsRows)
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
	var res []Appointments
	err = m.QueryRowsNoCacheCtx(ctx, &res, query1, values1...)
	if err != nil {
		return nil, 0, err
	}
	return POTODOGetAppointments(res), count, nil
}

// // Update 更新预约信息
func (m *customAppointmentsModel) UpdateAppointment(ctx context.Context, appointment *appointment.Appointments) error {
	po := DOTOPOAppointment(appointment)
	return m.Update(ctx, po)
}

// GetList 获取预约列表
func (m *customAppointmentsModel) GetList(ctx context.Context, userID string, search string, searchType, order, patientId, page, pageSize int64) ([]*appointment.Appointments, int64, error) {
	rowBuilder := persistence.RowBuilder(m.table, appointmentsRows)
	countBuilder := persistence.CountBuilder(m.table)
	if userID != "" {
		rowBuilder = rowBuilder.Where(squirrel.Eq{"user_id": userID})
		countBuilder = countBuilder.Where(squirrel.Eq{"user_id": userID})
	}
	if search != "" {
		rowBuilder = rowBuilder.Where(squirrel.Like{"patient_name": "%" + search + "%"})
		countBuilder = countBuilder.Where(squirrel.Like{"patient_name": "%" + search + "%"})
	}
	//searchType = 1 查询今天的，2查询未来的，3查询过去的
	if searchType == 1 {
		rowBuilder = rowBuilder.Where(squirrel.Eq{"date": time.Now().Format("2006-01-02")})
		countBuilder = countBuilder.Where(squirrel.Eq{"date": time.Now().Format("2006-01-02")})
	} else if searchType == 2 {
		rowBuilder = rowBuilder.Where(squirrel.Gt{"date": time.Now().Format("2006-01-02")})
		countBuilder = countBuilder.Where(squirrel.Gt{"date": time.Now().Format("2006-01-02")})
	} else if searchType == 3 {
		rowBuilder = rowBuilder.Where(squirrel.Lt{"date": time.Now().Format("2006-01-02")})
		countBuilder = countBuilder.Where(squirrel.Lt{"date": time.Now().Format("2006-01-02")})
	}

	if patientId != 0 {
		rowBuilder = rowBuilder.Where(squirrel.Eq{"patient_id": patientId})
		countBuilder = countBuilder.Where(squirrel.Eq{"patient_id": patientId})
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
	var res []Appointments
	err = m.QueryRowsNoCacheCtx(ctx, &res, query1, values1...)
	if err != nil {
		return nil, 0, err
	}
	return POTODOGetAppointments(res), count, nil
}

// CountTodayByUserID 根据用户ID获取今日预约数量
func (m *customAppointmentsModel) CountTodayByUserID(ctx context.Context, userID string) (int64, error) {
	countBuilder := persistence.CountBuilder(m.table).Where(squirrel.Eq{"user_id": userID})
	countBuilder = countBuilder.Where(squirrel.Eq{"date": time.Now().Format("2006-01-02")})
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

// CountUpcomingByUserID 根据用户ID获取即将进行的预约数量
func (m *customAppointmentsModel) CountUpcomingByUserID(ctx context.Context, userID string) (int64, error) {
	countBuilder := persistence.CountBuilder(m.table).Where(squirrel.Eq{"user_id": userID}).Where(squirrel.Gt{"date": time.Now().Format("2006-01-02")})
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
