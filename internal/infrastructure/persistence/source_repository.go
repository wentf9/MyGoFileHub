package persistence

import (
	"context"

	"github.com/wentf9/MyGoFileHub/internal/domain/model"
	"github.com/wentf9/MyGoFileHub/internal/domain/repository"

	"gorm.io/gorm"
)

type SourceRepository struct {
	db *gorm.DB
}

// NewSourceRepository 创建实例
func NewSourceRepository(db *gorm.DB) repository.SourceRepository {
	return &SourceRepository{db: db}
}

func (r *SourceRepository) FindByID(ctx context.Context, id uint) (*model.StorageSource, error) {
	var source model.StorageSource
	// GORM 的 First 方法会自动添加 LIMIT 1
	if err := r.db.WithContext(ctx).First(&source, id).Error; err != nil {
		return nil, err
	}
	return &source, nil
}

func (r *SourceRepository) FindByKey(ctx context.Context, key string) (*model.StorageSource, error) {
	var source model.StorageSource
	if err := r.db.WithContext(ctx).Where("key = ?", key).First(&source).Error; err != nil {
		return nil, err
	}
	return &source, nil
}

func (r *SourceRepository) FindAll(ctx context.Context) ([]*model.StorageSource, error) {
	var sources []*model.StorageSource
	if err := r.db.WithContext(ctx).Find(&sources).Error; err != nil {
		return nil, err
	}
	return sources, nil
}

func (r *SourceRepository) Save(ctx context.Context, source *model.StorageSource) error {
	return r.db.WithContext(ctx).Save(source).Error
}

func (r *SourceRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.StorageSource{}, id).Error
}
