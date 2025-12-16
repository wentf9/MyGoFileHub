package local

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/wentf9/MyGoFileHub/internal/domain/vfs"
	"github.com/wentf9/MyGoFileHub/internal/infrastructure/drivers"
)

// init 注册驱动到工厂
func init() {
	drivers.Register("local", NewLocalDriver)
}

// LocalDriver 本地文件系统实现
type LocalDriver struct {
	rootPath string // 挂载的根目录，例如 /data/files
}

func NewLocalDriver() vfs.StorageDriver {
	return &LocalDriver{}
}

func (d *LocalDriver) DriverName() string {
	return "local"
}

// Init 初始化，解析配置
// config 示例: {"root_path": "/var/www/uploads"}
func (d *LocalDriver) Init(ctx context.Context, config map[string]any) error {
	root, ok := config["root_path"].(string)
	if !ok || root == "" {
		return fmt.Errorf("local driver requires 'root_path' in config")
	}

	// 转换绝对路径并校验是否存在
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return err
	}

	info, err := os.Stat(absRoot)
	if os.IsNotExist(err) {
		// 如果目录不存在，尝试创建（可选，视需求而定）
		if err := os.MkdirAll(absRoot, 0755); err != nil {
			return fmt.Errorf("failed to create root path: %v", err)
		}
	} else if !info.IsDir() {
		return fmt.Errorf("root_path '%s' is not a directory", absRoot)
	}

	d.rootPath = absRoot
	return nil
}

// safePath 安全路径转换：将虚拟路径转换为物理路径，并防止 ../ 越权
func (d *LocalDriver) safePath(virtualPath string) (string, error) {
	// 拼接完整路径
	fullPath := filepath.Join(d.rootPath, virtualPath)

	// Clean 处理掉多余的 .. 和 /
	cleanPath := filepath.Clean(fullPath)

	// 核心安全检查：确保最终路径以前缀 rootPath 开头
	// 这防止了 virtualPath = "../../etc/passwd" 的情况
	if !strings.HasPrefix(cleanPath, d.rootPath) {
		return "", fmt.Errorf("access denied: path traversal attempt")
	}

	return cleanPath, nil
}

func (d *LocalDriver) List(ctx context.Context, path string) ([]vfs.FileInfo, error) {
	realPath, err := d.safePath(path)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(realPath)
	if err != nil {
		return nil, err
	}

	files := make([]vfs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue // 跳过无法获取信息的项
		}
		files = append(files, vfs.FileInfo{
			Name:    entry.Name(),
			Size:    info.Size(),
			IsDir:   entry.IsDir(),
			ModTime: info.ModTime(),
		})
	}
	return files, nil
}

func (d *LocalDriver) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	realPath, err := d.safePath(path)
	if err != nil {
		return nil, err
	}
	return os.Open(realPath)
}

func (d *LocalDriver) OpenFile(ctx context.Context, path string, flag int, perm fs.FileMode) (vfs.File, error) {
	realPath, err := d.safePath(path)
	if err != nil {
		return nil, err
	}
	// os.OpenFile 直接返回 *os.File，它完美实现了 vfs.File 接口
	return os.OpenFile(realPath, flag, perm)
}

func (d *LocalDriver) Create(ctx context.Context, path string, reader io.Reader, size int64) error {
	realPath, err := d.safePath(path)
	if err != nil {
		return err
	}

	// 确保父目录存在
	parentDir := filepath.Dir(realPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return err
	}

	// 创建文件
	out, err := os.Create(realPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 写入数据
	// 这里的 Copy 是流式的，内存占用小
	_, err = io.Copy(out, reader)
	return err
}

func (d *LocalDriver) Stat(ctx context.Context, path string) (vfs.FileInfo, error) {
	realPath, err := d.safePath(path)
	if err != nil {
		return vfs.FileInfo{}, err
	}

	info, err := os.Stat(realPath)
	if err != nil {
		return vfs.FileInfo{}, err
	}

	return vfs.FileInfo{
		Name:    info.Name(),
		Size:    info.Size(),
		IsDir:   info.IsDir(),
		ModTime: info.ModTime(),
	}, nil
}

func (d *LocalDriver) Delete(ctx context.Context, path string) error {
	realPath, err := d.safePath(path)
	if err != nil {
		return err
	}
	return os.RemoveAll(realPath) // RemoveAll 删除文件或目录
}

func (d *LocalDriver) Rename(ctx context.Context, srcPath, dstPath string) error {
	realSrc, err := d.safePath(srcPath)
	if err != nil {
		return err
	}
	realDst, err := d.safePath(dstPath)
	if err != nil {
		return err
	}
	return os.Rename(realSrc, realDst)
}

func (d *LocalDriver) Close() error {
	// 本地文件系统不需要关闭连接，但在 SMB/DB 中这里需要断开 Socket
	return nil
}
