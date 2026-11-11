package service

import (
	"context"
	"time"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
)

type AttendanceService struct {
	repo        repository.AttendanceRepository
	studentRepo repository.StudentRepository
	publisher   EventPublisher
}

func NewAttendanceService(
	repo repository.AttendanceRepository,
	studentRepo repository.StudentRepository,
	publisher EventPublisher,
) *AttendanceService {
	if publisher == nil {
		publisher = NewNoopEventPublisher()
	}
	return &AttendanceService{
		repo:        repo,
		studentRepo: studentRepo,
		publisher:   publisher,
	}
}

// --- Type definitions needed by handler (preserve API contracts) ---

type CreateSessionRequest struct {
	ClassID     string
	SubjectID   string
	TeacherID   string
	SessionDate time.Time
	StartTime   time.Time
	EndTime     time.Time
	Notes       string
}

type SubmitSessionAttendanceRequest struct {
	SessionID string
	Records   []AttendanceRecordInput
}

type AttendanceRecordInput struct {
	StudentID string
	Status    string
	Note      string
}

type SessionResponse struct {
	ID               string
	SchoolID         string
	AcademicPeriodID string
	ClassID          string
	SubjectID        string
	TeacherID        string
	SessionDate      time.Time
	StartTime        time.Time
	EndTime          time.Time
	Status           string
	Notes            string
	Records          []RecordResponse
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type RecordResponse struct {
	ID        string
	SessionID string
	StudentID string
	Status    string
	Note      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type StudentAttendanceResponse struct {
	Records []RecordResponse
	Summary AttendanceSummaryResponse
}

// --- Legacy session methods (stubs - to be implemented later) ---

// CreateSession membuat sesi absensi baru (belum diimplementasikan untuk DAPODIK)
func (s *AttendanceService) CreateSession(ctx context.Context, req CreateSessionRequest) (*SessionResponse, error) {
	return nil, NewServiceError("CreateSession", "AttendanceSession", "", "fitur sesi absensi belum diimplementasikan untuk skema DAPODIK", nil)
}

// SubmitAttendance menyimpan absensi (belum diimplementasikan untuk DAPODIK)
func (s *AttendanceService) SubmitAttendance(ctx context.Context, req SubmitSessionAttendanceRequest) (*SessionResponse, error) {
	return nil, NewServiceError("SubmitAttendance", "AttendanceSession", "", "fitur submit absensi belum diimplementasikan untuk skema DAPODIK", nil)
}

// ListSessions mengambil daftar sesi (belum diimplementasikan untuk DAPODIK)
func (s *AttendanceService) ListSessions(ctx context.Context, classID, status string, startDate, endDate time.Time) ([]*SessionResponse, error) {
	return nil, NewServiceError("ListSessions", "AttendanceSession", "", "fitur list sesi belum diimplementasikan untuk skema DAPODIK", nil)
}

// GetSession mengambil detail sesi (belum diimplementasikan untuk DAPODIK)
func (s *AttendanceService) GetSession(ctx context.Context, sessionID string) (*SessionResponse, error) {
	return nil, NewServiceError("GetSession", "AttendanceSession", sessionID, "fitur detail sesi belum diimplementasikan untuk skema DAPODIK", nil)
}

// GetStudentAttendance mengambil absensi per siswa (belum diimplementasikan untuk DAPODIK)
func (s *AttendanceService) GetStudentAttendance(ctx context.Context, studentID string, month, year int) (*StudentAttendanceResponse, error) {
	// Stub: kembalikan data kosong dengan summary nol
	return &StudentAttendanceResponse{
		Records: []RecordResponse{},
		Summary: AttendanceSummaryResponse{},
	}, nil
}

// --- Rekap Kehadiran (IMPLEMENTASI NYATA) ---

// GetAttendanceSummary mengembalikan rekap kehadiran per siswa berdasarkan kelas, bulan, dan tahun.
// Menggunakan tabel_kehadiran DAPODIK (sakit, izin, tanpa_keterangan).
func (s *AttendanceService) GetAttendanceSummary(ctx context.Context, classID string, month, year int) ([]*StudentAttendanceSummaryResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	// 1. Ambil daftar siswa aktif di kelas tersebut
	studentFilter := repository.StudentFilter{
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		ClassID:          classID,
		Limit:            1000,
		Offset:           0,
	}

	students, _, err := s.studentRepo.List(ctx, studentFilter)
	if err != nil {
		return nil, NewServiceError("GetAttendanceSummary", "Student", classID, "gagal mengambil daftar siswa", err)
	}

	// 2. Loop setiap siswa, ambil rekap kehadiran dari tabel_kehadiran
	var summaries []*StudentAttendanceSummaryResponse
	for _, student := range students {
		attSummary, err := s.repo.GetStudentAttendanceSummary(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, student.ID)
		if err != nil {
			// Log dan lanjutkan ke siswa berikutnya (tidak abort seluruh rekap karena 1 siswa error)
			continue
		}

		summaries = append(summaries, &StudentAttendanceSummaryResponse{
			StudentID:   student.ID,
			StudentName: student.Name,
			Summary: AttendanceSummaryResponse{
				TotalSessions:        attSummary.TotalSessions,
				PresentCount:         attSummary.PresentCount,
				AbsentCount:          attSummary.AbsentCount,
				LateCount:            attSummary.LateCount,
				SickCount:            attSummary.SickCount,
				PermissionCount:      attSummary.PermissionCount,
				AttendancePercentage: attSummary.AttendancePercentage,
			},
		})
	}

	return summaries, nil
}

// --- Response DTOs ---

type StudentAttendanceSummaryResponse struct {
	StudentID   string
	StudentName string
	Summary     AttendanceSummaryResponse
}

type AttendanceSummaryResponse struct {
	TotalSessions        int
	PresentCount         int
	AbsentCount          int
	LateCount            int
	SickCount            int
	PermissionCount      int
	AttendancePercentage float64
}
