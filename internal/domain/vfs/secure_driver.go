package vfs

import (
	"context"
	"io"
	"io/fs"
	"os"
)

// PermissionChecker 定义检查接口 (解耦 Application 层)
type PermissionChecker func(ctx context.Context, path string, action string) bool

// SecureDriver 装饰器：为普通 Driver 增加权限检查功能
type SecureDriver struct {
	base    StorageDriver     // 底层驱动 (Local, SMB 等)
	checker PermissionChecker // 检查函数
}

// NewSecureDriver 包装一个驱动
func NewSecureDriver(base StorageDriver, checker PermissionChecker) StorageDriver {
	return &SecureDriver{
		base:    base,
		checker: checker,
	}
}

func (d *SecureDriver) DriverName() string {
	return "secure-" + d.base.DriverName()
}

func (d *SecureDriver) Init(ctx context.Context, config map[string]interface{}) error {
	return d.base.Init(ctx, config)
}

// --- 读操作 (检查 "read") ---

func (d *SecureDriver) List(ctx context.Context, path string) ([]FileInfo, error) {
	if !d.checker(ctx, path, "read") {
		return nil, os.ErrPermission
	}
	return d.base.List(ctx, path)
}

func (d *SecureDriver) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	if !d.checker(ctx, path, "read") {
		return nil, os.ErrPermission
	}
	return d.base.Open(ctx, path)
}

func (d *SecureDriver) Stat(ctx context.Context, path string) (FileInfo, error) {
	// Stat 通常允许 read 权限即可
	if !d.checker(ctx, path, "read") {
		return FileInfo{}, os.ErrPermission
	}
	return d.base.Stat(ctx, path)
}

func (d *SecureDriver) OpenFile(ctx context.Context, path string, flag int, perm fs.FileMode) (File, error) {
	// OpenFile 比较特殊，可能是读，可能是写
	action := "read"
	if flag&os.O_WRONLY != 0 || flag&os.O_RDWR != 0 || flag&os.O_CREATE != 0 || flag&os.O_TRUNC != 0 {
		action = "write"
	}

	if !d.checker(ctx, path, action) {
		return nil, os.ErrPermission
	}
	return d.base.OpenFile(ctx, path, flag, perm)
}

// --- 写操作 (检查 "write") ---

func (d *SecureDriver) Create(ctx context.Context, path string, reader io.Reader, size int64) error {
	if !d.checker(ctx, path, "write") {
		return os.ErrPermission
	}
	return d.base.Create(ctx, path, reader, size)
}

func (d *SecureDriver) Delete(ctx context.Context, path string) error {
	if !d.checker(ctx, path, "write") {
		return os.ErrPermission
	}
	return d.base.Delete(ctx, path)
}

func (d *SecureDriver) Rename(ctx context.Context, srcPath, dstPath string) error {
	// 重命名需要：源路径(写/删权限) + 目标路径(写权限)
	if !d.checker(ctx, srcPath, "write") || !d.checker(ctx, dstPath, "write") {
		return os.ErrPermission
	}
	return d.base.Rename(ctx, srcPath, dstPath)
}

func (d *SecureDriver) Close() error {
	return d.base.Close()
}
