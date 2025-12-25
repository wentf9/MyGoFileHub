package smb

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/wentf9/MyGoFileHub/internal/domain/vfs"
	"github.com/wentf9/MyGoFileHub/internal/infrastructure/drivers"

	"github.com/hirochachacha/go-smb2"
)

func init() {
	drivers.Register("smb", NewSMBDriver)
}

type SMBDriver struct {
	conn    net.Conn
	session *smb2.Session
	share   *smb2.Share // 类似于 os 包，拥有 Open, Create, Stat 等方法
}

func NewSMBDriver() vfs.StorageDriver {
	return &SMBDriver{}
}

func (d *SMBDriver) DriverName() string {
	return "smb"
}

// Init 初始化 SMB 连接
// Config 需求: host, port(可选,默认445), user, password, share_name
func (d *SMBDriver) Init(ctx context.Context, config map[string]interface{}) error {
	host, _ := config["host"].(string)
	port, _ := config["port"].(string)
	user, _ := config["user"].(string)
	password, _ := config["password"].(string)
	shareName, _ := config["share_name"].(string)

	if host == "" || user == "" || shareName == "" {
		return fmt.Errorf("smb config missing: host, user, or share_name")
	}
	if port == "" {
		port = "445"
	}
	// 1. 建立 TCP 连接
	conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return fmt.Errorf("tcp dial failed: %v", err)
	}
	d.conn = conn

	// 2. 建立 SMB 会话 (认证)
	dialer := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     user,
			Password: password,
			// Domain: domain, // 如果需要域认证，可扩展配置
		},
	}

	session, err := dialer.Dial(conn)
	if err != nil {
		conn.Close()
		return fmt.Errorf("smb session setup failed: %v", err)
	}
	d.session = session

	// 3. 挂载共享目录
	fs, err := session.Mount(shareName)
	if err != nil {
		session.Logoff()
		conn.Close()
		return fmt.Errorf("mount share '%s' failed: %v", shareName, err)
	}
	d.share = fs

	return nil
}

// normalizePath 处理路径分隔符
// SMB 协议内部通常处理 backslash，但库做了封装。为了保险，去除开头的 /
func (d *SMBDriver) normalizePath(path string) string {
	// 移除开头的 "/" 或 "\"
	path = strings.TrimLeft(path, "/\\")
	// 将 "/" 替换为 "\" (Windows 风格)，虽然 go-smb2 很多时候兼容 /，但转换更稳妥
	path = strings.ReplaceAll(path, "/", "\\")
	if path == "" {
		return "."
	}
	return path
}

func (d *SMBDriver) List(ctx context.Context, path string) ([]vfs.FileInfo, error) {
	normPath := d.normalizePath(path)

	// 使用 share.ReadDir，用法几乎和 os.ReadDir 一样
	entries, err := d.share.ReadDir(normPath)
	if err != nil {
		return nil, err
	}

	files := make([]vfs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		// 过滤掉 "." 和 ".."
		if entry.Name() == "." || entry.Name() == ".." {
			continue
		}

		files = append(files, vfs.FileInfo{
			Name:    entry.Name(),
			Size:    entry.Size(),
			IsDir:   entry.IsDir(),
			ModTime: entry.ModTime(),
		})
	}
	return files, nil
}

func (d *SMBDriver) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	return d.share.Open(d.normalizePath(path))
}

func (d *SMBDriver) OpenFile(ctx context.Context, path string, flag int, perm os.FileMode) (vfs.File, error) {
	normPath := d.normalizePath(path)
	// go-smb2 的 OpenFile 也返回实现了 vfs.File 的对象
	return d.share.OpenFile(normPath, flag, perm)
}

func (d *SMBDriver) Create(ctx context.Context, path string, reader io.Reader, size int64) error {
	normPath := d.normalizePath(path)

	// 确保父目录存在 (SMB 协议不支持 mkdir -p，需要逐级检查，这里简化处理，假设目录存在)
	// 如果要健壮实现，需要在这里写递归创建目录的逻辑

	f, err := d.share.Create(normPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, reader)
	return err
}

func (d *SMBDriver) Mkdir(ctx context.Context, path string, perm os.FileMode) error {
	normPath := d.normalizePath(path)
	return d.share.MkdirAll(normPath, perm)
}

func (d *SMBDriver) Stat(ctx context.Context, path string) (vfs.FileInfo, error) {
	info, err := d.share.Stat(d.normalizePath(path))
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

func (d *SMBDriver) Delete(ctx context.Context, path string) error {
	// go-smb2 的 RemoveAll 类似于 os.RemoveAll
	return d.share.RemoveAll(d.normalizePath(path))
}

func (d *SMBDriver) Rename(ctx context.Context, srcPath, dstPath string) error {
	return d.share.Rename(d.normalizePath(srcPath), d.normalizePath(dstPath))
}

func (d *SMBDriver) Close() error {
	if d.share != nil {
		d.share.Umount()
	}
	if d.session != nil {
		d.session.Logoff()
	}
	if d.conn != nil {
		d.conn.Close()
	}
	return nil
}
