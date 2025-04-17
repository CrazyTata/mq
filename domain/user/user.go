package user

import (
	"database/sql"
	"time"

	"mq/common/util"
	"mq/domain"
)

// User represents the user entity
type User struct {
	ID             int64        // 用户ID
	AppId          string       // 归属应用
	UserID         string       // 用户ID
	Phone          string       // 手机号（支持手机号一键登录）
	AppleID        string       // Apple 用户唯一标识符(sub 字段)
	Email          string       // 用户邮箱(Apple登录可能提供)
	IsPrivateEmail int64        // 是否是Apple隐私邮箱
	Name           string       // 用户昵称
	Avatar         string       // 用户头像URL
	Password       string       // 密码
	RegisterTime   sql.NullTime // 注册时间
	Status         int64        // 用户状态:1 正常 2冻结 3 注销
	Source         int64        // 用户来源:4 安卓 5 IOS
	IsDeleted      int64        // 是否删除
	CreatedAt      time.Time    // 创建时间
	UpdatedAt      time.Time
}

const (
	UserSourceAndroid = 4
	UserSourceIos     = 5
)

const (
	UserStatusNormal = 1
	UserStatusMute   = 2
	UserStatusDelete = 3
)

func CreateUser(appId string, phone string, userID string, appleID string, email string, name string, avatar string, password string, registerTime time.Time, status, source int64) *User {
	return &User{
		AppId:        appId,
		Phone:        util.Encrypt(phone),
		UserID:       userID,
		AppleID:      appleID,
		Email:        email,
		Name:         name,
		Avatar:       avatar,
		Password:     password,
		Source:       source,
		Status:       status,
		RegisterTime: sql.NullTime{Time: registerTime, Valid: true},
	}
}

func (u *User) UpdateFields(name string, avatar string, phone string) {
	if name != "" {
		u.Name = name
	}
	if avatar != "" {
		u.Avatar = avatar
	}
	if phone != "" {
		u.Phone = util.Encrypt(phone)
	}
	return
}

// GetDecryptedPhone 获取解密后的手机号
func (u *User) GetDecryptedPhone() string {
	return util.Decrypt(u.Phone)
}

func (u *User) MarkAsDeleted() {
	u.IsDeleted = domain.IsDeleted
	return
}

func (u *User) UpdatePassword(password string) {
	if password != "" {
		u.Password = password
	}
	return
}
