package model

import "time"

// UserPermission 用户权限表
type UserPermission struct {
	ID       uint `gorm:"primaryKey" json:"id"`
	UserID   uint `gorm:"index:idx_user_source" json:"user_id"`
	SourceID uint `gorm:"index:idx_user_source" json:"source_id"`

	// PathPrefix 路径前缀
	// "/" 代表整个源
	// "/work" 代表 work 目录及其子目录
	PathPrefix string `gorm:"size:255;not null" json:"path_prefix"`

	AllowRead  bool `gorm:"default:false" json:"allow_read"`  // 读/列出
	AllowWrite bool `gorm:"default:false" json:"allow_write"` // 写/删/改

	CreatedAt time.Time `json:"created_at"`
}

func (UserPermission) TableName() string {
	return "user_permissions"
}
