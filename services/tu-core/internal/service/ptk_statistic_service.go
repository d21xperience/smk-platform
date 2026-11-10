package service

import (
	"context"
	"time"

	"sekolah-platform/services/tu-core/internal/repository"
)

type PtkStatisticService struct {
	repo repository.PtkStatisticRepository
}

func NewPtkStatisticService(repo repository.PtkStatisticRepository) *PtkStatisticService {
	return &PtkStatisticService{repo: repo}
}

func (s *PtkStatisticService) GetDashboardData(ctx context.Context, schoolID string, startDate, endDate string) (*repository.StatisticDashboardData, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		start = time.Now().AddDate(0, -1, 0) // Default 1 bulan ke belakang
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		end = time.Now()
	}

	return s.repo.GetDashboardData(ctx, schoolID, start, end)
}
