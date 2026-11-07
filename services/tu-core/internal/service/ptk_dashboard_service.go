package service

import (
	"context"

	"sekolah-platform/services/tu-core/internal/repository"
)

type PtkDashboardService struct {
	repo repository.PtkDashboardRepository
}

func NewPtkDashboardService(repo repository.PtkDashboardRepository) *PtkDashboardService {
	return &PtkDashboardService{repo: repo}
}

func (s *PtkDashboardService) GetDashboardStats(ctx context.Context) (*repository.DashboardStats, error) {
	return s.repo.GetDashboardStats(ctx)
}
