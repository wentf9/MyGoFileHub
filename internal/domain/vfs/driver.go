package vfs

import (
	"context"
	"io"
	"io/fs"
)

// StorageDriver 定义了所有存储源必须具备的行为
// Context 用于处理超时和取消（如下载中断）
type StorageDriver interface {
	// DriverName 返回驱动名称，如 "smb", "local"
	DriverName() string

	// Init 初始化连接
	// config 是 JSON 解析后的 map，包含 host, user, password 等
	Init(ctx context.Context, config map[string]any) error

	// List 列出指定路径下的文件
	List(ctx context.Context, path string) ([]FileInfo, error)

	// Open 打开文件流进行读取 (用于在线播放、下载)
	Open(ctx context.Context, path string) (io.ReadCloser, error)

	OpenFile(ctx context.Context, path string, flag int, perm fs.FileMode) (File, error)

	// Create 创建或覆盖文件 (用于上传、离线下载写入)
	// reader 是输入流，size 是预估大小（某些协议如 WebDAV 需要预知大小）
	Create(ctx context.Context, path string, reader io.Reader, size int64) error

	// Stat 获取单个文件详情
	Stat(ctx context.Context, path string) (FileInfo, error)

	// Delete 删除文件或目录
	Delete(ctx context.Context, path string) error

	// Rename 重命名或移动（在同一源内）
	Rename(ctx context.Context, srcPath, dstPath string) error

	// Close 释放资源（如断开 SMB 连接）
	Close() error
}
