package service

import (
	"context"

	"sekolah-platform/services/tu-core/internal/repository"
)

type StatistikService struct {
	repo repository.StatistikRepository
}

func NewStatistikService(repo repository.StatistikRepository) *StatistikService {
	return &StatistikService{repo: repo}
}

func (s *StatistikService) GetDashboardStats(ctx context.Context) (*repository.StatistikData, error) {
	return s.repo.GetDashboardStats(ctx)
}
