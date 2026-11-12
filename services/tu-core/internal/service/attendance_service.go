package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

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
	return &AttendanceService{repo: repo, studentRepo: studentRepo, publisher: publisher}
}

// ===== DTOs =====

type CreateSessionRequest struct {
	ClassID     string    `json:"class_id"`
	SubjectID   string    `json:"subject_id"`
	TeacherID   string    `json:"teacher_id"`
	SessionDate time.Time `json:"session_date"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Notes       string    `json:"notes"`
}

type SubmitSessionAttendanceRequest struct {
	SessionID string                  `json:"session_id"`
	Records   []AttendanceRecordInput `json:"records"`
}

type AttendanceRecordInput struct {
	StudentID string `json:"student_id"`
	Status    string `json:"status"`
	Note      string `json:"note"`
}

type SessionResponse struct {
	ID               string           `json:"id"`
	SchoolID         string           `json:"school_id"`
	AcademicPeriodID string           `json:"academic_period_id"`
	ClassID          string           `json:"class_id"`
	SubjectID        string           `json:"subject_id,omitempty"`
	TeacherID        string           `json:"teacher_id,omitempty"`
	SessionDate      time.Time        `json:"session_date"`
	StartTime        time.Time        `json:"start_time,omitempty"`
	EndTime          time.Time        `json:"end_time,omitempty"`
	Status           string           `json:"status"`
	Notes            string           `json:"notes,omitempty"`
	Records          []RecordResponse `json:"records,omitempty"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}

type RecordResponse struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	StudentID string    `json:"student_id"`
	Status    string    `json:"status"`
	Note      string    `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StudentAttendanceResponse struct {
	Records []RecordResponse          `json:"records"`
	Summary AttendanceSummaryResponse `json:"summary"`
}

type AttendanceSummaryResponse struct {
	TotalSessions        int     `json:"total_sessions"`
	PresentCount         int     `json:"present_count"`
	AbsentCount          int     `json:"absent_count"`
	LateCount            int     `json:"late_count"`
	SickCount            int     `json:"sick_count"`
	PermissionCount      int     `json:"permission_count"`
	AttendancePercentage float64 `json:"attendance_percentage"`
}

type StudentAttendanceSummaryResponse struct {
	StudentID   string                    `json:"student_id"`
	StudentName string                    `json:"student_name"`
	Summary     AttendanceSummaryResponse `json:"summary"`
}

// ===== SESSION IMPLEMENTATIONS =====

func (s *AttendanceService) CreateSession(ctx context.Context, req CreateSessionRequest) (*SessionResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	session := &repository.AttendanceSession{
		ID:               uuid.New().String(),
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		ClassID:          req.ClassID,
		SubjectID:        req.SubjectID,
		TeacherID:        req.TeacherID,
		SessionDate:      req.SessionDate,
		StartTime:        req.StartTime.Format("15:04"),
		EndTime:          req.EndTime.Format("15:04"),
		Status:           "OPEN",
		Notes:            req.Notes,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, NewServiceError("CreateSession", "AttendanceSession", "", "gagal simpan sesi", err)
	}

	return toSessionResponse(session, nil), nil
}

func (s *AttendanceService) SubmitAttendance(ctx context.Context, req SubmitSessionAttendanceRequest) (*SessionResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	// 1. Validasi session
	_, err := s.repo.GetSessionByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, req.SessionID)
	if err != nil {
		return nil, NewServiceError("SubmitAttendance", "AttendanceSession", req.SessionID, "sesi tidak ditemukan", err)
	}

	// 2. Upsert setiap record
	now := time.Now().UTC()
	for _, rec := range req.Records {
		record := &repository.AttendanceRecord{
			ID:        uuid.New().String(),
			SessionID: req.SessionID,
			StudentID: rec.StudentID,
			Status:    rec.Status,
			Note:      rec.Note,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := s.repo.UpsertRecord(ctx, record); err != nil {
			return nil, NewServiceError("SubmitAttendance", "AttendanceRecord", rec.StudentID, "gagal simpan record", err)
		}
	}

	// 3. Tandai sesi selesai (COMPLETED)
	if err := s.repo.UpdateSessionStatus(ctx, req.SessionID, "COMPLETED"); err != nil {
		return nil, NewServiceError("SubmitAttendance", "AttendanceSession", req.SessionID, "gagal update status sesi", err)
	}

	// 4. Ambil data terbaru untuk response
	session, _ := s.repo.GetSessionByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, req.SessionID)
	records, _ := s.repo.GetRecordsBySession(ctx, req.SessionID)

	return toSessionResponse(session, records), nil
}

func (s *AttendanceService) ListSessions(ctx context.Context, classID, status string, startDate, endDate time.Time) ([]*SessionResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	filter := repository.SessionFilter{
		SchoolID: opCtx.SchoolID, AcademicPeriodID: opCtx.AcademicPeriodID,
		ClassID: classID, Status: status, Limit: 100, Offset: 0,
	}

	sessions, _, err := s.repo.ListSessions(ctx, filter)
	if err != nil {
		return nil, NewServiceError("ListSessions", "AttendanceSession", "", "gagal ambil daftar sesi", err)
	}

	var responses []*SessionResponse
	for i := range sessions {
		responses = append(responses, toSessionResponse(&sessions[i], nil))
	}
	return responses, nil
}

func (s *AttendanceService) GetSession(ctx context.Context, sessionID string) (*SessionResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	session, err := s.repo.GetSessionByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, sessionID)
	if err != nil {
		return nil, NewServiceError("GetSession", "AttendanceSession", sessionID, "sesi tidak ditemukan", err)
	}

	records, _ := s.repo.GetRecordsBySession(ctx, sessionID)
	return toSessionResponse(session, records), nil
}

func (s *AttendanceService) GetStudentAttendance(ctx context.Context, studentID string, month, year int) (*StudentAttendanceResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	summary, err := s.repo.GetStudentAttendanceSummary(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, studentID)
	if err != nil {
		return nil, NewServiceError("GetStudentAttendance", "Student", studentID, "gagal ambil summary", err)
	}

	return &StudentAttendanceResponse{
		Records: []RecordResponse{},
		Summary: AttendanceSummaryResponse{
			TotalSessions:        summary.TotalSessions,
			PresentCount:         summary.PresentCount,
			AbsentCount:          summary.AbsentCount,
			LateCount:            summary.LateCount,
			SickCount:            summary.SickCount,
			PermissionCount:      summary.PermissionCount,
			AttendancePercentage: summary.AttendancePercentage,
		},
	}, nil
}

// ===== REKAP (EXISTING) =====

func (s *AttendanceService) GetAttendanceSummary(ctx context.Context, classID string, month, year int) ([]*StudentAttendanceSummaryResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	studentFilter := repository.StudentFilter{
		SchoolID: opCtx.SchoolID, AcademicPeriodID: opCtx.AcademicPeriodID,
		ClassID: classID, Limit: 1000, Offset: 0,
	}
	students, _, err := s.studentRepo.List(ctx, studentFilter)
	if err != nil {
		return nil, NewServiceError("GetAttendanceSummary", "Student", classID, "gagal ambil siswa", err)
	}

	var summaries []*StudentAttendanceSummaryResponse
	for _, student := range students {
		attSummary, err := s.repo.GetStudentAttendanceSummary(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, student.ID)
		if err != nil {
			continue
		}
		summaries = append(summaries, &StudentAttendanceSummaryResponse{
			StudentID: student.ID, StudentName: student.Name,
			Summary: AttendanceSummaryResponse{
				TotalSessions: attSummary.TotalSessions, PresentCount: attSummary.PresentCount,
				AbsentCount: attSummary.AbsentCount, LateCount: attSummary.LateCount,
				SickCount: attSummary.SickCount, PermissionCount: attSummary.PermissionCount,
				AttendancePercentage: attSummary.AttendancePercentage,
			},
		})
	}
	return summaries, nil
}

// ===== HELPERS =====

func toSessionResponse(s *repository.AttendanceSession, records []repository.AttendanceRecord) *SessionResponse {
	if s == nil {
		return nil
	}
	var recResponses []RecordResponse
	for _, r := range records {
		recResponses = append(recResponses, RecordResponse{
			ID: r.ID, SessionID: r.SessionID, StudentID: r.StudentID,
			Status: r.Status, Note: r.Note, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
		})
	}

	// Parsing kembali string "15:04" ke time.Time untuk kompatibilitas gRPC
	startTime, _ := time.Parse("15:04", s.StartTime)
	endTime, _ := time.Parse("15:04", s.EndTime)

	return &SessionResponse{
		ID: s.ID, SchoolID: s.SchoolID, AcademicPeriodID: s.AcademicPeriodID,
		ClassID: s.ClassID, SubjectID: s.SubjectID, TeacherID: s.TeacherID,
		SessionDate: s.SessionDate, StartTime: startTime, EndTime: endTime,
		Status: s.Status, Notes: s.Notes, Records: recResponses,
		CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
	}
}

var _ = fmt.Sprintf
