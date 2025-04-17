package operation

import (
	"database/sql"
	"time"
)

// Operation 操作记录领域模型
type OperationRecords struct {
	Id        int64
	Action    string         // 操作
	Target    string         // 目标
	Details   sql.NullString // 详情
	Username  string         // 用户名
	UserId    string         // 用户ID
	IsDeleted int64          // 是否删除
	CreatedAt time.Time      // 创建时间
	UpdatedAt time.Time      // 更新时间
}

// Create 创建操作记录
func Create(action string, target string, details string, username string, userId string) *OperationRecords {
	return &OperationRecords{
		Action:    action,
		Target:    target,
		Details:   sql.NullString{String: details, Valid: details != ""},
		Username:  username,
		UserId:    userId,
		IsDeleted: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Update 更新操作记录
func (o *OperationRecords) Update(action string, target string, details string, username string, userId string) {
	o.Action = action
	o.Target = target
	o.Details = sql.NullString{String: details, Valid: details != ""}
	o.Username = username
	o.UserId = userId
}
