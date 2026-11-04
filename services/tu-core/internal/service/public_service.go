package service

import (
	"context"
	"sekolah-platform/services/tu-core/internal/repository"
)

type PublicService struct {
	publicRepo repository.PublicRepository
}

func NewPublicService(publicRepo repository.PublicRepository) *PublicService {
	return &PublicService{publicRepo: publicRepo}
}

func (s *PublicService) GetDashboardStats(ctx context.Context) (*repository.PublicStats, error) {
	return s.publicRepo.GetDashboardStats(ctx)
}
