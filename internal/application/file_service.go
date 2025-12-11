package application

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/wentf9/MyGoFileHub/internal/domain/repository"
	"github.com/wentf9/MyGoFileHub/internal/domain/vfs"
	"github.com/wentf9/MyGoFileHub/internal/infrastructure/drivers"
)

type FileService struct {
	sourceRepo repository.SourceRepository
}

// NewFileService 注入 Repository
func NewFileService(repo repository.SourceRepository) *FileService {
	return &FileService{sourceRepo: repo}
}

// ListFiles 列出文件
// sourceIDStr: 接收 string 类型的 ID，在 Service 层内部转为 uint
func (s *FileService) ListFiles(ctx context.Context, sourceIDStr string, path string) ([]vfs.FileInfo, error) {
	// 1. ID 类型转换
	sourceID, err := strconv.ParseUint(sourceIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid source id")
	}

	// 2. 从数据库查询配置
	source, err := s.sourceRepo.FindByID(ctx, uint(sourceID))
	if err != nil {
		return nil, fmt.Errorf("storage source not found: %v", err)
	}

	// 3. 通过工厂创建驱动
	driver, err := drivers.CreateInstance(source.Type)
	if err != nil {
		return nil, err
	}

	// 4. 初始化驱动 (传入从数据库取出的 Config JSONMap)
	// source.Config 本身就是 map[string]interface{}，可以直接传
	if err := driver.Init(ctx, source.Config); err != nil {
		return nil, fmt.Errorf("failed to init driver: %v", err)
	}

	// 5. 执行操作
	return driver.List(ctx, path)
}

// GetFileStream 获取文件流 (用于下载或播放)
// sourceID: 数据库中存储源的ID
// path: 文件在存储源中的相对路径
func (s *FileService) GetFileStream(ctx context.Context, sourceID uint, path string) (io.ReadCloser, error) {
	// 1. 从数据库获取源配置
	sourceConfig, err := s.sourceRepo.FindByID(ctx, sourceID)
	if err != nil {
		return nil, err
	}

	// 2. 利用工厂创建具体的驱动实例 (比如 SMB Driver)
	driver, err := drivers.CreateInstance(sourceConfig.Type)
	if err != nil {
		return nil, err
	}

	// 3. 初始化连接
	// 注意：实际生产中，Driver 实例应该有连接池缓存，不要每次请求都 Connect/Close
	err = driver.Init(ctx, sourceConfig.Config)
	if err != nil {
		return nil, err
	}

	// 4. 调用接口方法
	return driver.Open(ctx, path)
}
