package subject

import (
	"context"
	"mq/domain/subject"
	"mq/infrastructure/persistence"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ subject.SubjectRepository = (*customSubjectModel)(nil)

type (
	// SubjectModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSubjectModel.
	SubjectModel interface {
		subjectModel
	}

	customSubjectModel struct {
		*defaultSubjectModel
	}
)

// NewSubjectModel returns a model for the database table.
func NewSubjectModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) subject.SubjectRepository {
	return &customSubjectModel{
		defaultSubjectModel: newSubjectModel(conn, c, opts...),
	}
}

func (m *customSubjectModel) FindByAppId(ctx context.Context, appId string) (*subject.Subject, error) {
	query, values, err := persistence.RowBuilder(m.table, subjectRows).Where(squirrel.Eq{"app_id": appId}).ToSql()
	if err != nil {
		return nil, err
	}

	var res Subject
	err = m.QueryRowNoCacheCtx(ctx, &res, query, values...)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return POTODOGetSubject(&res), nil
}

func (m *customSubjectModel) FindAllActive(ctx context.Context) ([]*subject.Subject, error) {
	query, values, err := persistence.RowBuilder(m.table, subjectRows).Where(squirrel.Eq{"is_active": 1}).ToSql()
	if err != nil {
		return nil, err
	}

	var res []Subject
	err = m.QueryRowsNoCacheCtx(ctx, &res, query, values...)
	if err != nil {
		return nil, err
	}

	return POTODOGetSubjects(res), nil
}

func (m *customSubjectModel) UpdateSubject(ctx context.Context, subject *subject.Subject) error {
	po := DOTOPOSubject(subject)
	return m.Update(ctx, po)
}

func (m *customSubjectModel) InsertSubject(ctx context.Context, subject *subject.Subject) error {
	po := DOTOPOSubject(subject)
	_, err := m.Insert(ctx, po)
	return err
}
