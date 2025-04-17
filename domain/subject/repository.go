package subject

import (
	"context"
)

// Repository defines the interface for user data access
type SubjectRepository interface {
	FindByAppId(ctx context.Context, appId string) (*Subject, error)
	UpdateSubject(ctx context.Context, subject *Subject) error
	InsertSubject(ctx context.Context, subject *Subject) error
	FindAllActive(ctx context.Context) ([]*Subject, error)
}

type AppSharingRepository interface {
	GetByAppId(ctx context.Context, appId string) ([]*AppSharing, error)
	InsertAppSharing(ctx context.Context, appSharing *AppSharing) (int64, error)
	UpdateAppSharing(ctx context.Context, appSharing *AppSharing) error
}
