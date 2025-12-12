package application

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"sync"

	"github.com/wentf9/MyGoFileHub/internal/domain/repository"
	"github.com/wentf9/MyGoFileHub/internal/domain/vfs"
	"github.com/wentf9/MyGoFileHub/internal/infrastructure/drivers"
)

type FileService struct {
	sourceRepo  repository.SourceRepository
	permService *PermissionService // 注入权限服务
}

// NewFileService 注入 Repository
func NewFileService(repo repository.SourceRepository, perm *PermissionService) *FileService {
	return &FileService{sourceRepo: repo, permService: perm}
}

var dirverCache = sync.Map{} // map[uint64]vfs.StorageDriver
var dirverMu sync.Mutex

// ListFiles 列出文件
// sourceIDStr: 接收 string 类型的 ID，在 Service 层内部转为 uint
func (s *FileService) ListFiles(ctx context.Context, sourceIDStr string, path string) ([]vfs.FileInfo, error) {
	driver, err := s.GetDriver(ctx, sourceIDStr)
	if err != nil {
		return nil, err
	}
	return driver.List(ctx, path)
}

// GetFileStream 获取文件流 (用于下载或播放)
// sourceID: 数据库中存储源的ID
// path: 文件在存储源中的相对路径
func (s *FileService) GetFileStream(ctx context.Context, sourceIDStr string, path string) (io.ReadCloser, error) {
	driver, err := s.GetDriver(ctx, sourceIDStr)
	if err != nil {
		return nil, err
	}
	// 4. 调用接口方法
	return driver.Open(ctx, path)
}

func (s *FileService) GetDriver(ctx context.Context, sourceID string) (vfs.StorageDriver, error) {
	// 1. 先从缓存获取
	dirver, ok := dirverCache.Load(sourceID)
	if ok {
		if value, valid := dirver.(vfs.StorageDriver); valid {
			return value, nil
		} else {
			return nil, fmt.Errorf("invalid driver type in cache")
		}
	}
	dirverMu.Lock()
	defer dirverMu.Unlock()
	// 2. 再次检查缓存，防止并发重复创建
	dirver, ok = dirverCache.Load(sourceID)
	if ok {
		if value, valid := dirver.(vfs.StorageDriver); valid {
			return value, nil
		} else {
			return nil, fmt.Errorf("invalid driver type in cache")
		}
	}
	// 3. 从数据库查询配置
	sID, err := strconv.ParseUint(sourceID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid source id")
	}
	source, err := s.sourceRepo.FindByID(ctx, uint(sID))
	if err != nil {
		return nil, fmt.Errorf("storage source not found: %v", err)
	}
	// 4. 通过工厂创建驱动
	driver, err := drivers.CreateInstance(source.Type)
	if err != nil {
		return nil, err
	}
	// ---------------------------------------------------------
	// 5. 获取当前用户并包裹 SecureDriver
	// ---------------------------------------------------------
	// 定义检查闭包
	checker := func(c context.Context, path string, action string) (bool, error) {
		// 从 Context 中提取 UserID (由 JWT 中间件设置)
		userIDVal := c.Value("userID")
		if userIDVal == nil {
			// 如果没有登录(或者内部调用)，视情况处理。
			// 这里假设必须登录，否则返回 error 或者是仅限 Admin 的 Context
			return false, errors.New("unauthorized: user context missing")
		}
		userID := userIDVal.(uint)
		// 调用 PermissionService 进行真正的数据库校验
		return s.permService.CheckPermission(c, userID, uint(sID), path, action), nil
	}
	secureDriver := vfs.NewSecureDriver(driver, checker)
	// 6. 初始化驱动 (传入从数据库取出的 Config JSONMap)
	// source.Config 本身就是 map[string]interface{}，可以直接传
	if err := secureDriver.Init(ctx, source.Config); err != nil {
		return nil, fmt.Errorf("failed to init driver: %v", err)
	}
	dirverCache.Store(sourceID, secureDriver)
	return secureDriver, nil
}
