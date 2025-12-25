package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// JSONMap 自定义类型，用于将 map[string]interface{} 存储为数据库的 JSON/Text 字段
type JSONMap map[string]any

// Value 实现 driver.Valuer 接口，将 Map 转为 JSON 字符串存入数据库
func (m JSONMap) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

// Scan 实现 sql.Scanner 接口，将数据库取出的 JSON 字符串转为 Map
func (m *JSONMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &m)
}

// StorageSource 存储源实体
type StorageSource struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"size:32;uniqueIndex;not null" json:"key"` // 唯一标识符,同时也作为访问数据源的路径
	Name      string    `gorm:"size:64;not null" json:"name"`            // 例如: "我的NAS", "本地磁盘"
	Type      string    `gorm:"size:16;not null" json:"type"`            // 例如: "smb", "local"
	Config    JSONMap   `gorm:"type:text" json:"config"`                 // 存储具体的连接配置
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (StorageSource) TableName() string {
	return "storage_sources"
}
