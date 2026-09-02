package service

import (
	"context"
)

// DashboardService mengorchestrasi use case untuk Dashboard Kepala Sekolah.
type DashboardService struct {
	studentRepo    StudentCountRepository
	attendanceRepo AttendanceSummaryRepository
	financeRepo    FinanceSummaryRepository
}

// StudentCountRepository kontrak untuk hitung jumlah siswa.
type StudentCountRepository interface {
	CountActiveStudents(ctx context.Context, schoolID, academicPeriodID string) (int, error)
}

// AttendanceSummaryRepository kontrak untuk ringkasan absensi.
type AttendanceSummaryRepository interface {
	GetTodayAttendanceSummary(ctx context.Context, schoolID, academicPeriodID string) (map[string]interface{}, error)
}

// FinanceSummaryRepository kontrak untuk ringkasan keuangan.
type FinanceSummaryRepository interface {
	GetFinanceSummary(ctx context.Context, schoolID, academicPeriodID string) (map[string]interface{}, error)
}

// NewDashboardService membuat instance baru.
func NewDashboardService(
	studentRepo StudentCountRepository,
	attendanceRepo AttendanceSummaryRepository,
	financeRepo FinanceSummaryRepository,
) *DashboardService {
	return &DashboardService{
		studentRepo:    studentRepo,
		attendanceRepo: attendanceRepo,
		financeRepo:    financeRepo,
	}
}

// DashboardResponse adalah response untuk dashboard.
type DashboardResponse struct {
	TotalStudents   int                    `json:"totalStudents"`
	TotalTeachers   int                    `json:"totalTeachers"`
	AttendanceToday map[string]interface{} `json:"attendanceToday"`
	FinanceSummary  map[string]interface{} `json:"financeSummary"`
}

// GetDashboardSummary mengambil ringkasan dashboard.
func (s *DashboardService) GetDashboardSummary(ctx context.Context, schoolID, academicPeriodID string) (*DashboardResponse, error) {
	if schoolID == "" || academicPeriodID == "" {
		return nil, ErrMissingOperationalContext
	}

	// Ambil total siswa
	totalStudents, err := s.studentRepo.CountActiveStudents(ctx, schoolID, academicPeriodID)
	if err != nil {
		return nil, NewServiceError("GetDashboardSummary", "Dashboard", "", "gagal ambil total siswa", err)
	}

	// Ambil ringkasan absensi hari ini
	attendanceToday, err := s.attendanceRepo.GetTodayAttendanceSummary(ctx, schoolID, academicPeriodID)
	if err != nil {
		attendanceToday = map[string]interface{}{}
	}

	// Ambil ringkasan keuangan
	financeSummary, err := s.financeRepo.GetFinanceSummary(ctx, schoolID, academicPeriodID)
	if err != nil {
		financeSummary = map[string]interface{}{}
	}

	return &DashboardResponse{
		TotalStudents:   totalStudents,
		TotalTeachers:   0, // TODO: implementasi
		AttendanceToday: attendanceToday,
		FinanceSummary:  financeSummary,
	}, nil
}
