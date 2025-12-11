package persistence

import (
	"context"

	"github.com/wentf9/MyGoFileHub/internal/domain/model"
	"github.com/wentf9/MyGoFileHub/internal/domain/repository"
	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) repository.PermissionRepository {
	return &PermissionRepository{db: db}
}

// FindByUserAndSource 获取用户在指定源下的所有权限规则
func (r *PermissionRepository) FindByUserAndSource(ctx context.Context, userID, sourceID uint) ([]*model.UserPermission, error) {
	var perms []*model.UserPermission
	// 只需要查出该用户、该源的记录
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND source_id = ?", userID, sourceID).
		Find(&perms).Error
	return perms, err
}

// Save 保存权限（后台管理用）
func (r *PermissionRepository) Save(ctx context.Context, perm *model.UserPermission) error {
	return r.db.WithContext(ctx).Save(perm).Error
}
