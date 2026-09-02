package reporting

import (
	"context"
	"fmt"
)

// ReportingService mengorchestrasi permintaan laporan.
type ReportingService struct {
	repo ReportingRepository
}

// NewReportingService membuat instance baru dari ReportingService.
func NewReportingService(repo ReportingRepository) *ReportingService {
	return &ReportingService{repo: repo}
}

// GetAttendanceDashboard mengambil ringkasan absensi untuk dashboard.
func (s *ReportingService) GetAttendanceDashboard(ctx context.Context, filter ReportFilter) (interface{}, error) {
	if filter.StartDate == "" || filter.EndDate == "" {
		return nil, fmt.Errorf("startDate dan endDate wajib diisi")
	}

	data, err := s.repo.GetAttendanceSummary(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan attendance dashboard: %w", err)
	}

	return map[string]interface{}{
		"status": "success",
		"data":   data,
	}, nil
}

// GetFinanceDashboard mengambil ringkasan keuangan untuk dashboard.
func (s *ReportingService) GetFinanceDashboard(ctx context.Context, filter ReportFilter) (interface{}, error) {
	data, err := s.repo.GetFinanceSummary(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan finance dashboard: %w", err)
	}

	return map[string]interface{}{
		"status": "success",
		"data":   data,
	}, nil
}
