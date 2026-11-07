package service

import (
	"context"

	"sekolah-platform/services/tu-core/internal/repository"
)

type PtkPerformanceService struct {
	repo repository.PtkPerformanceRepository
}

func NewPtkPerformanceService(repo repository.PtkPerformanceRepository) *PtkPerformanceService {
	return &PtkPerformanceService{repo: repo}
}

func (s *PtkPerformanceService) ListPerformances(ctx context.Context) ([]repository.Performance, error) {
	return s.repo.ListPerformances(ctx)
}

func (s *PtkPerformanceService) CreatePerformance(ctx context.Context, req repository.Performance) error {
	return s.repo.CreatePerformance(ctx, &req)
}

func (s *PtkPerformanceService) DeletePerformance(ctx context.Context, id string) error {
	return s.repo.DeletePerformance(ctx, id)
}
