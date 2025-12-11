package persistence

import (
	"github.com/wentf9/MyGoFileHub/internal/domain/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitDB 初始化 SQLite 连接并自动迁移表结构
func InitDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移模式：自动创建表、缺少的字段
	err = db.AutoMigrate(&model.StorageSource{}, &model.User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
