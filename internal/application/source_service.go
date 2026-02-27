package application

import (
	"context"
	"strings"

	"github.com/wentf9/MyGoFileHub/config"
	"github.com/wentf9/MyGoFileHub/internal/domain/model"
	"github.com/wentf9/MyGoFileHub/internal/domain/repository"
	"github.com/wentf9/MyGoFileHub/internal/infrastructure/crypto"
	"github.com/wentf9/MyGoFileHub/internal/infrastructure/drivers"
)

type SourceService struct {
	sourceRepo  repository.SourceRepository
	fileService *FileService
}

func NewSourceService(repo repository.SourceRepository, fileService *FileService) *SourceService {
	return &SourceService{sourceRepo: repo, fileService: fileService}
}

func (s *SourceService) ListSources(ctx context.Context) ([]*model.StorageSource, error) {
	sources, err := s.sourceRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, source := range sources {
		s.maskSensitiveConfig(source)
	}
	return sources, nil
}

func (s *SourceService) GetSource(ctx context.Context, id uint) (*model.StorageSource, error) {
	source, err := s.sourceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if source != nil {
		s.maskSensitiveConfig(source)
	}
	return source, nil
}

func (s *SourceService) CreateSource(ctx context.Context, source *model.StorageSource) error {
	// 验证存储源 key 的合法性
	if err := model.ValidateStorageKey(source.Key); err != nil {
		return err
	}
	if err := s.encryptSensitiveConfig(source); err != nil {
		return err
	}
	return s.sourceRepo.Save(ctx, source)
}

func (s *SourceService) UpdateSource(ctx context.Context, source *model.StorageSource) error {
	// 获取旧信息，用于清理缓存（特别是 Key 变更的情况）以及恢复未修改的密码
	oldSource, err := s.sourceRepo.FindByID(ctx, source.ID)
	if err == nil && oldSource != nil {
		// 如果 key 发生变更，验证新 key 的合法性
		if oldSource.Key != source.Key {
			if err := model.ValidateStorageKey(source.Key); err != nil {
				return err
			}
			s.fileService.ClearDriverCache(oldSource.Key)
		}
		s.restoreUnchangedSensitiveConfig(source, oldSource)
	}

	if err := s.encryptSensitiveConfig(source); err != nil {
		return err
	}

	if err := s.sourceRepo.Save(ctx, source); err != nil {
		return err
	}

	// 清理新 Key 的缓存
	s.fileService.ClearDriverCache(source.Key)
	return nil
}

func (s *SourceService) DeleteSource(ctx context.Context, id uint) error {
	source, err := s.sourceRepo.FindByID(ctx, id)
	if err == nil && source != nil {
		s.fileService.ClearDriverCache(source.Key)
	}
	return s.sourceRepo.Delete(ctx, id)
}

// getSensitiveFields 获取某种驱动类型中定义为 password 的字段名列表
func (s *SourceService) getSensitiveFields(driverType string) []string {
	var sensitiveFields []string
	schemas := drivers.GetRegisteredSchemas()
	for _, schema := range schemas {
		if schema.Type == driverType {
			for _, item := range schema.Config {
				if item.Type == "password" {
					sensitiveFields = append(sensitiveFields, item.Name)
				}
			}
			break
		}
	}
	return sensitiveFields
}

// encryptSensitiveConfig 加密敏感字段，存入数据库前调用
func (s *SourceService) encryptSensitiveConfig(source *model.StorageSource) error {
	if source.Config == nil {
		return nil
	}
	sensitiveFields := s.getSensitiveFields(source.Type)
	for _, field := range sensitiveFields {
		if val, ok := source.Config[field]; ok {
			strVal, isStr := val.(string)
			// 如果是字符串且没有被加密过（不以 ENC: 开头），则进行加密
			if isStr && strVal != "" && !strings.HasPrefix(strVal, "ENC:") {
				encrypted, err := crypto.Encrypt(strVal, config.AppConfig.SecretKey)
				if err != nil {
					return err
				}
				source.Config[field] = "ENC:" + encrypted
			}
		}
	}
	return nil
}

// maskSensitiveConfig 脱敏敏感字段，返回给前端前调用
func (s *SourceService) maskSensitiveConfig(source *model.StorageSource) {
	if source.Config == nil {
		return
	}
	sensitiveFields := s.getSensitiveFields(source.Type)
	for _, field := range sensitiveFields {
		if val, ok := source.Config[field]; ok {
			strVal, isStr := val.(string)
			if isStr && strVal != "" {
				source.Config[field] = "********" // 返回给前端统一打码
			}
		}
	}
}

// restoreUnchangedSensitiveConfig 在更新时，如果前端传来的密码是打码的，则恢复为数据库旧的密文
func (s *SourceService) restoreUnchangedSensitiveConfig(newSource *model.StorageSource, oldSource *model.StorageSource) {
	if newSource.Config == nil || oldSource.Config == nil {
		return
	}
	sensitiveFields := s.getSensitiveFields(newSource.Type)
	for _, field := range sensitiveFields {
		if newVal, ok := newSource.Config[field]; ok {
			strVal, isStr := newVal.(string)
			// 如果前端传上来的是占位符，说明用户没有修改密码，从旧数据恢复密文
			if isStr && strVal == "********" {
				if oldVal, oldOk := oldSource.Config[field]; oldOk {
					newSource.Config[field] = oldVal
				}
			}
		}
	}
}
