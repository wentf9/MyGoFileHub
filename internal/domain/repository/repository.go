package repository

import (
	"context"

	"github.com/wentf9/MyGoFileHub/internal/domain/model"
)

// UserRepository 用户数据存取
type UserRepository interface {
	FindByID(ctx context.Context, id uint) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	Save(ctx context.Context, user *model.User) error
}

// SourceRepository 存储源配置存取
type SourceRepository interface {
	FindAll(ctx context.Context) ([]*model.StorageSource, error)
	FindByID(ctx context.Context, id uint) (*model.StorageSource, error)
	FindByKey(ctx context.Context, key string) (*model.StorageSource, error)
	Save(ctx context.Context, source *model.StorageSource) error
	Delete(ctx context.Context, id uint) error
}

// // TaskRepository 下载任务存取
// type TaskRepository interface {
// 	Save(ctx context.Context, task *model.DownloadTask) error
// 	UpdateStatus(ctx context.Context, id uint, status int, progress float64) error
// 	FindPending(ctx context.Context) ([]*model.DownloadTask, error)
// }

// PermissionRepository 用户权限存取
type PermissionRepository interface {
	FindByUserAndSource(ctx context.Context, userID, sourceID uint) ([]*model.UserPermission, error)
	Save(ctx context.Context, perm *model.UserPermission) error
}
