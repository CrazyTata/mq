package user

import (
	"context"
	"database/sql"
	"time"

	domainUser "mq/domain/user"
	"mq/infrastructure/persistence"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ domainUser.UserRepository = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) domainUser.UserRepository {
	return &customUserModel{
		defaultUserModel: newUserModel(conn, c, opts...),
	}
}

func (m *customUserModel) InsertUser(ctx context.Context, user *domainUser.User) error {
	po := DOTOPOUser(user)
	_, err := m.Insert(ctx, po)
	return err
}

func (m *customUserModel) FindByUserID(ctx context.Context, userID string) (*domainUser.User, error) {
	query, values, err := persistence.RowBuilder(m.table, userRows).Where(squirrel.Eq{"user_id": userID}).ToSql()
	if err != nil {
		return nil, err
	}

	var res User
	err = m.QueryRowNoCacheCtx(ctx, &res, query, values...)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	resp := POTODOGetUser(&res)
	return resp, nil
}

func (m *customUserModel) UpdateUser(ctx context.Context, user *domainUser.User) error {
	po := DOTOPOUser(user)
	return m.Update(ctx, po)
}

func (m *customUserModel) FindByPhone(ctx context.Context, phone string) (*domainUser.User, error) {
	query, values, err := persistence.RowBuilder(m.table, userRows).Where(squirrel.Eq{"phone": phone}).ToSql()
	if err != nil {
		return nil, err
	}

	var res User
	err = m.QueryRowNoCacheCtx(ctx, &res, query, values...)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	resp := POTODOGetUser(&res)
	return resp, nil
}

func (m *customUserModel) FindByAppleID(ctx context.Context, appleID string) (*domainUser.User, error) {
	query, values, err := persistence.RowBuilder(m.table, userRows).Where(squirrel.Eq{"apple_id": appleID}).ToSql()
	if err != nil {
		return nil, err
	}

	var res User
	err = m.QueryRowNoCacheCtx(ctx, &res, query, values...)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	resp := POTODOGetUser(&res)
	return resp, nil
}

func (m *customUserModel) FindByIds(ctx context.Context, userIDs []string) ([]*domainUser.User, error) {
	query, values, err := persistence.RowBuilder(m.table, userRows).Where(squirrel.Eq{"user_id": userIDs}).ToSql()
	if err != nil {
		return nil, err
	}

	var res []User
	err = m.QueryRowsNoCacheCtx(ctx, &res, query, values...)
	if err != nil {
		return nil, err
	}

	return POTODOGetUsers(res), nil
}

func (m *customUserModel) FindBySubscriptionExpired(ctx context.Context, subscriptionType []uint8, subscriptionEnd time.Time) ([]*domainUser.User, error) {
	query, values, err := persistence.RowBuilder(m.table, userRows).Where(squirrel.Eq{"subscription_type": subscriptionType}).Where(squirrel.Lt{"subscription_end": subscriptionEnd}).ToSql()
	if err != nil {
		return nil, err
	}

	var res []User
	err = m.QueryRowsNoCacheCtx(ctx, &res, query, values...)
	if err != nil {
		return nil, err
	}

	return POTODOGetUsers(res), nil
}

func (m *customUserModel) CountNewUsersInTimeRange(ctx context.Context, startTime, endTime string) (int64, error) {
	query, values, err := persistence.CountBuilder(m.table).Where(squirrel.Gt{"created_at": startTime}).Where(squirrel.Lt{"created_at": endTime}).ToSql()
	if err != nil {
		return 0, err
	}

	var count int64
	err = m.QueryRowNoCacheCtx(ctx, &count, query, values...)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return count, nil
}

func (m *customUserModel) GetBackendUsers(ctx context.Context, startTime, endTime string, page int, pageSize int) ([]*domainUser.User, int64, error) {
	rowBuilder := persistence.RowBuilder(m.table, userRows)
	countBuilder := persistence.CountBuilder(m.table)
	if startTime != "" {
		rowBuilder = rowBuilder.Where(squirrel.Gt{"created_at": startTime})
		countBuilder = countBuilder.Where(squirrel.Gt{"created_at": startTime})
	}
	if endTime != "" {
		rowBuilder = rowBuilder.Where(squirrel.Lt{"created_at": endTime})
		countBuilder = countBuilder.Where(squirrel.Lt{"created_at": endTime})
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
	var res []User
	err = m.QueryRowsNoCacheCtx(ctx, &res, query1, values1...)
	if err != nil {
		return nil, 0, err
	}
	return POTODOGetUsers(res), count, nil
}

func (m *customUserModel) GetAllUserIds(ctx context.Context) ([]string, error) {
	query, values, err := persistence.RowBuilder(m.table, userRows).ToSql()
	if err != nil {
		return nil, err
	}

	var res []User
	err = m.QueryRowsNoCacheCtx(ctx, &res, query, values...)
	if err != nil {
		return nil, err
	}

	userIDs := make([]string, 0, len(res))
	for _, user := range res {
		userIDs = append(userIDs, user.UserId)
	}

	return userIDs, nil
}

func (m *customUserModel) FindByEmail(ctx context.Context, email string) (*domainUser.User, error) {
	query, values, err := persistence.RowBuilder(m.table, userRows).Where(squirrel.Eq{"email": email}).ToSql()
	if err != nil {
		return nil, err
	}

	var res User
	err = m.QueryRowNoCacheCtx(ctx, &res, query, values...)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	resp := POTODOGetUser(&res)
	return resp, nil
}
