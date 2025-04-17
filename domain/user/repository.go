package user

import (
	"context"
	"time"
)

// Repository defines the interface for user data access
type UserRepository interface {
	FindByPhone(ctx context.Context, phone string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByUserID(ctx context.Context, userID string) (*User, error)
	FindByAppleID(ctx context.Context, appleID string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	InsertUser(ctx context.Context, user *User) error
	FindByIds(ctx context.Context, userIDs []string) ([]*User, error)
	FindBySubscriptionExpired(ctx context.Context, subscriptionType []uint8, subscriptionEnd time.Time) ([]*User, error)
	CountNewUsersInTimeRange(ctx context.Context, startTime, endTime string) (int64, error)
	GetBackendUsers(ctx context.Context, startTime, endTime string, page, pageSize int) ([]*User, int64, error)
	GetAllUserIds(ctx context.Context) ([]string, error)
}
