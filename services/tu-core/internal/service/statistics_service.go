package service

import (
	"context"

	"sekolah-platform/services/tu-core/internal/repository"
)

type StatisticsService struct {
	repo repository.StatisticsRepository
}

func NewStatisticsService(repo repository.StatisticsRepository) *StatisticsService {
	return &StatisticsService{repo: repo}
}

func (s *StatisticsService) GetDashboardStats(ctx context.Context, schoolID, periodID string) (*repository.StatisticsData, error) {
	if schoolID == "" || periodID == "" {
		return nil, ErrMissingOperationalContext
	}
	return s.repo.GetStatistics(ctx, schoolID, periodID)
}
