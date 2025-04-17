package subject

import (
	"database/sql"
	"time"
)

// Subject represents the subject entity
type Subject struct {
	Id           int64
	AppId        string         // API密钥标识
	AppSecret    string         // API密钥(加密存储)
	MerchantId   string         // 商户ID
	MerchantName string         // 商户名称
	IsActive     int64          // 是否激活 1: 激活 0: 禁用
	ExpireAt     sql.NullTime   // 过期时间
	AllowedIps   sql.NullString // 允许的IP地址列表,用逗号分隔
	AllowedPaths sql.NullString // 允许访问的路径,用逗号分隔
	CreatedAt    time.Time      // 创建时间
	UpdatedAt    time.Time      // 更新时间
}

type AppSharing struct {
	Id          int64
	AppId       string    // 发起共享的应用
	SharedAppId string    // 共享目标应用
	IsDeleted   int64     // 是否删除
	CreatedAt   time.Time // 创建时间
	UpdatedAt   time.Time // 更新时间
}

const (
	SubjectStatusActive   = 1
	SubjectStatusInactive = 0
)

func CreateSubject(appId, appSecret, merchantId, merchantName string) *Subject {
	return &Subject{
		AppId:        appId,
		AppSecret:    appSecret,
		MerchantId:   merchantId,
		MerchantName: merchantName,
		IsActive:     SubjectStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func CreateAppSharing(appId, sharedAppId string) *AppSharing {
	return &AppSharing{
		AppId:       appId,
		SharedAppId: sharedAppId,
	}
}

func (s *AppSharing) MarkAsDeleted() {
	s.IsDeleted = 1
}
