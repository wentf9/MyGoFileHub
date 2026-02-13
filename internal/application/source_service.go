package application

import (
	"context"

	"github.com/wentf9/MyGoFileHub/internal/domain/model"
	"github.com/wentf9/MyGoFileHub/internal/domain/repository"
)

type SourceService struct {
	sourceRepo  repository.SourceRepository
	fileService *FileService
}

func NewSourceService(repo repository.SourceRepository, fileService *FileService) *SourceService {
	return &SourceService{sourceRepo: repo, fileService: fileService}
}

func (s *SourceService) ListSources(ctx context.Context) ([]*model.StorageSource, error) {
	return s.sourceRepo.FindAll(ctx)
}

func (s *SourceService) GetSource(ctx context.Context, id uint) (*model.StorageSource, error) {
	return s.sourceRepo.FindByID(ctx, id)
}

func (s *SourceService) CreateSource(ctx context.Context, source *model.StorageSource) error {
	return s.sourceRepo.Save(ctx, source)
}

func (s *SourceService) UpdateSource(ctx context.Context, source *model.StorageSource) error {
	// 获取旧信息，用于清理缓存（特别是 Key 变更的情况）
	oldSource, err := s.sourceRepo.FindByID(ctx, source.ID)
	if err == nil && oldSource != nil {
		s.fileService.ClearDriverCache(oldSource.Key)
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
