package subject

import (
	"context"
	"mq/domain/subject"
	"mq/infrastructure/persistence"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ subject.AppSharingRepository = (*customAppSharingModel)(nil)

type (
	// AppSharingModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAppSharingModel.
	AppSharingModel interface {
		appSharingModel
	}

	customAppSharingModel struct {
		*defaultAppSharingModel
	}
)

// NewAppSharingModel returns a model for the database table.
func NewAppSharingModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) subject.AppSharingRepository {
	return &customAppSharingModel{
		defaultAppSharingModel: newAppSharingModel(conn, c, opts...),
	}
}

func (m *customAppSharingModel) GetByAppId(ctx context.Context, appId string) ([]*subject.AppSharing, error) {
	query, values, err := persistence.RowBuilder(m.table, appSharingRows).Where(squirrel.Eq{"app_id": appId}).ToSql()
	if err != nil {
		return nil, err
	}

	var res []AppSharing
	err = m.QueryRowsNoCacheCtx(ctx, &res, query, values...)
	if err != nil {
		return nil, err
	}

	return POTODOGetAppSharings(res), nil
}

func (m *customAppSharingModel) InsertAppSharing(ctx context.Context, appSharing *subject.AppSharing) (int64, error) {
	po := DOTOPOAppSharing(appSharing)
	result, err := m.Insert(ctx, po)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (m *customAppSharingModel) UpdateAppSharing(ctx context.Context, appSharing *subject.AppSharing) error {
	po := DOTOPOAppSharing(appSharing)
	return m.Update(ctx, po)
}
