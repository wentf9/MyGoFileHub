package vfs

import (
	"io"
	"io/fs"
	"time"
)

// File 是我们系统中文件的最高抽象，兼容 os.File 和 smb2.File
// 这也是 WebDAV 库所需要的接口
type File interface {
	io.Closer
	io.Reader
	io.Writer
	io.Seeker
	Readdir(count int) ([]fs.FileInfo, error)
	Stat() (fs.FileInfo, error)
}

// FileInfo 是系统内部通用的文件元数据模型，屏蔽底层差异
type FileInfo struct {
	Name    string
	Size    int64
	IsDir   bool
	ModTime time.Time
	// 可以扩展 MIME type, ETag 等
}

type fileInfoAdapter struct {
	info FileInfo
}

func (f *fileInfoAdapter) Name() string       { return f.info.Name }
func (f *fileInfoAdapter) Size() int64        { return f.info.Size }
func (f *fileInfoAdapter) IsDir() bool        { return f.info.IsDir }
func (f *fileInfoAdapter) ModTime() time.Time { return f.info.ModTime }
func (f *fileInfoAdapter) Mode() fs.FileMode  { return 0 }
func (f *fileInfoAdapter) Sys() any           { return nil }

func ToOSFileInfo(info FileInfo) fs.FileInfo {
	return &fileInfoAdapter{info: info}
}
