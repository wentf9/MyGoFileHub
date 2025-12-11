package webdav

import (
	"context"
	"os"

	"github.com/wentf9/MyGoFileHub/internal/domain/vfs"

	"golang.org/x/net/webdav"
)

// DriverFileSystem 将我们的 StorageDriver 适配为 webdav.FileSystem
type DriverFileSystem struct {
	Driver vfs.StorageDriver
}

func (fsys *DriverFileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	// 我们的 Driver 需要补充 Mkdir 方法，或者暂时用 Create 模拟
	// 这里为了演示，假设 Driver 已经有了 Mkdir
	return fsys.Driver.Create(ctx, name+"/.dir_placeholder", nil, 0)
}

func (fsys *DriverFileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	// 直接透传调用
	return fsys.Driver.OpenFile(ctx, name, flag, perm)
}

func (fsys *DriverFileSystem) RemoveAll(ctx context.Context, name string) error {
	return fsys.Driver.Delete(ctx, name)
}

func (fsys *DriverFileSystem) Rename(ctx context.Context, oldName, newName string) error {
	return fsys.Driver.Rename(ctx, oldName, newName)
}

func (fsys *DriverFileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	info, err := fsys.Driver.Stat(ctx, name)
	if err != nil {
		return nil, err
	}
	return vfs.ToOSFileInfo(info), nil
}
