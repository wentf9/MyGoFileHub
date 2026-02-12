package application

import (
	"context"

	"github.com/wentf9/MyGoFileHub/internal/domain/model"
	"github.com/wentf9/MyGoFileHub/internal/domain/repository"
)

type SourceService struct {
	sourceRepo repository.SourceRepository
}

func NewSourceService(repo repository.SourceRepository) *SourceService {
	return &SourceService{sourceRepo: repo}
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
	return s.sourceRepo.Save(ctx, source)
}

func (s *SourceService) DeleteSource(ctx context.Context, id uint) error {
	return s.sourceRepo.Delete(ctx, id)
}
