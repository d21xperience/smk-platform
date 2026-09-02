package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"sekolah-platform/services/tu-core/internal/domain/attendance/engine"
	"sekolah-platform/services/tu-core/internal/domain/attendance/models"
	"sekolah-platform/services/tu-core/internal/domain/attendance/valueobjects"
	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
)

type AttendanceService struct {
	repo        repository.AttendanceRepository
	studentRepo repository.StudentRepository
	attEngine   *engine.AttendanceEngine
	publisher   EventPublisher
}

func NewAttendanceService(
	repo repository.AttendanceRepository,
	studentRepo repository.StudentRepository,
	attEngine *engine.AttendanceEngine,
	publisher EventPublisher,
) *AttendanceService {
	if publisher == nil {
		publisher = NewNoopEventPublisher()
	}
	return &AttendanceService{
		repo:        repo,
		studentRepo: studentRepo,
		attEngine:   attEngine,
		publisher:   publisher,
	}
}

type CreateSessionRequest struct {
	ClassID     string
	SubjectID   string
	TeacherID   string
	SessionDate time.Time
	StartTime   time.Time
	EndTime     time.Time
	Notes       string
}

func (s *AttendanceService) CreateSession(ctx context.Context, req CreateSessionRequest) (*SessionResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	data := models.SessionData{
		ClassID:     req.ClassID,
		SubjectID:   req.SubjectID,
		TeacherID:   req.TeacherID,
		SessionDate: req.SessionDate,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Notes:       req.Notes,
	}

	result := s.attEngine.CreateSession(data, opCtx)
	if result.Error != nil {
		return nil, NewServiceError("CreateSession", "AttendanceSession", "", "invalid session data", result.Error)
	}

	if err := s.repo.CreateSession(ctx, result.Session); err != nil {
		return nil, NewServiceError("CreateSession", "AttendanceSession", result.Session.ID, "failed to save to database", err)
	}

	event := DomainEvent{
		EventType:        "SessionCreated",
		AggregateType:    "AttendanceSession",
		AggregateID:      result.Session.ID,
		Payload: map[string]interface{}{
			"sessionId": result.Session.ID,
			"classId":   result.Session.ClassID,
			"date":      result.Session.SessionDate.Format("2006-01-02"),
		},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return toSessionResponse(result.Session), nil
}

// ✅ PERBAIKAN: Rename agar tidak bentrok dengan dto.go
type SubmitSessionAttendanceRequest struct {
	SessionID string
	Records   []AttendanceRecordInput
}

type AttendanceRecordInput struct {
	StudentID string
	Status    string
	Note      string
}

func (s *AttendanceService) SubmitAttendance(ctx context.Context, req SubmitSessionAttendanceRequest) (*SessionResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	session, err := s.repo.GetSessionByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, req.SessionID)
	if err != nil {
		return nil, NewServiceError("SubmitAttendance", "AttendanceSession", req.SessionID, "failed to get session", err)
	}

	// Create records
	for _, recInput := range req.Records {
		status, err := valueobjects.NewAttendanceStatus(recInput.Status)
		if err != nil {
			return nil, NewServiceError("SubmitAttendance", "AttendanceRecord", "", "invalid attendance status", err)
		}

		recData := models.RecordData{
			SessionID: session.ID,
			StudentID: recInput.StudentID,
			Status:    status,
			Note:      recInput.Note,
		}

		record, err := models.NewRecord(recData)
		if err != nil {
			return nil, NewServiceError("SubmitAttendance", "AttendanceRecord", "", "failed to create record", err)
		}

		if err := session.AddRecord(record); err != nil {
			return nil, NewServiceError("SubmitAttendance", "AttendanceSession", session.ID, "failed to add record", err)
		}

		if err := s.repo.CreateRecord(ctx, record); err != nil {
			return nil, NewServiceError("SubmitAttendance", "AttendanceRecord", record.ID, "failed to save record", err)
		}

		// Publish event per record
		event := DomainEvent{
			EventType:        "AttendanceRecorded",
			AggregateType:    "AttendanceRecord",
			AggregateID:      record.ID,
			Payload: map[string]interface{}{
				"sessionId": session.ID,
				"studentId": record.StudentID,
				"status":    record.Status.String(),
			},
			SchoolID:         opCtx.SchoolID,
			AcademicPeriodID: opCtx.AcademicPeriodID,
			CorrelationID:    uuid.New().String(),
			Version:          "v1",
		}
		_ = s.publisher.Publish(ctx, event)
	}

	// Submit session
	result := s.attEngine.SubmitSession(session, opCtx)
	if result.Error != nil {
		return nil, NewServiceError("SubmitAttendance", "AttendanceSession", session.ID, "failed to submit session", result.Error)
	}

	if err := s.repo.UpdateSession(ctx, result.Session); err != nil {
		return nil, NewServiceError("SubmitAttendance", "AttendanceSession", session.ID, "failed to update session", err)
	}

	// Reload session with records
	session, err = s.repo.GetSessionByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, session.ID)
	if err != nil {
		return nil, NewServiceError("SubmitAttendance", "AttendanceSession", session.ID, "failed to reload session", err)
	}

	return toSessionResponse(session), nil
}

func (s *AttendanceService) GetSession(ctx context.Context, sessionID string) (*SessionResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	session, err := s.repo.GetSessionByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, sessionID)
	if err != nil {
		return nil, NewServiceError("GetSession", "AttendanceSession", sessionID, "failed to get session", err)
	}

	return toSessionResponse(session), nil
}

func (s *AttendanceService) GetStudentAttendance(ctx context.Context, studentID string, month, year int) (*StudentAttendanceResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	records, err := s.repo.GetRecordsByStudent(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, studentID, month, year)
	if err != nil {
		return nil, NewServiceError("GetStudentAttendance", "AttendanceRecord", studentID, "failed to get records", err)
	}

	// Calculate summary
	summary := calculateSummary(records)

	var recordResponses []RecordResponse
	for _, r := range records {
		recordResponses = append(recordResponses, toRecordResponse(&r))
	}

	return &StudentAttendanceResponse{
		Records: recordResponses,
		Summary: summary,
	}, nil
}

func (s *AttendanceService) ListSessions(ctx context.Context, classID, status string, startDate, endDate time.Time) ([]*SessionResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	filter := repository.SessionFilter{
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		ClassID:          classID,
		Status:           status,
		StartDate:        startDate,
		EndDate:          endDate,
		Limit:            100,
		Offset:           0,
	}

	sessions, _, err := s.repo.ListSessions(ctx, filter)
	if err != nil {
		return nil, NewServiceError("ListSessions", "AttendanceSession", "", "failed to list sessions", err)
	}

	var responses []*SessionResponse
	for _, session := range sessions {
		responses = append(responses, toSessionResponse(&session))
	}

	return responses, nil
}

func (s *AttendanceService) GetAttendanceSummary(ctx context.Context, classID string, month, year int) ([]*StudentAttendanceSummaryResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	// ✅ PERBAIKAN: Gunakan studentRepo.List dengan filter ClassID (karena ListByClass tidak ada)
	studentFilter := repository.StudentFilter{
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		ClassID:          classID,
		Status:           "ACTIVE",
		Limit:            1000,
		Offset:           0,
	}

	students, _, err := s.studentRepo.List(ctx, studentFilter)
	if err != nil {
		return nil, NewServiceError("GetAttendanceSummary", "Student", classID, "failed to get students", err)
	}

	var summaries []*StudentAttendanceSummaryResponse
	for _, student := range students {
		records, err := s.repo.GetRecordsByStudent(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, student.ID, month, year)
		if err != nil {
			continue // Skip if error
		}

		summary := calculateSummary(records)
		summaries = append(summaries, &StudentAttendanceSummaryResponse{
			StudentID:   student.ID,
			StudentName: student.Name,
			Summary:     summary,
		})
	}

	return summaries, nil
}

// Helper functions
func calculateSummary(records []models.Record) AttendanceSummaryResponse {
	summary := AttendanceSummaryResponse{
		TotalSessions: len(records),
	}

	for _, r := range records {
		switch r.Status {
		case valueobjects.StatusPresent:
			summary.PresentCount++
		case valueobjects.StatusAbsent:
			summary.AbsentCount++
		case valueobjects.StatusLate:
			summary.LateCount++
		case valueobjects.StatusSick:
			summary.SickCount++
		case valueobjects.StatusPermission:
			summary.PermissionCount++
		}
	}

	if summary.TotalSessions > 0 {
		summary.AttendancePercentage = float64(summary.PresentCount+summary.LateCount) / float64(summary.TotalSessions) * 100
	}

	return summary
}

// Response DTOs
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

type AttendanceSummaryResponse struct {
	TotalSessions        int
	PresentCount         int
	AbsentCount          int
	LateCount            int
	SickCount            int
	PermissionCount      int
	AttendancePercentage float64
}

type StudentAttendanceSummaryResponse struct {
	StudentID   string
	StudentName string
	Summary     AttendanceSummaryResponse
}

func toSessionResponse(s *models.Session) *SessionResponse {
	if s == nil {
		return nil
	}

	var records []RecordResponse
	for _, r := range s.Records {
		records = append(records, toRecordResponse(r))
	}

	return &SessionResponse{
		ID:               s.ID,
		SchoolID:         s.SchoolID,
		AcademicPeriodID: s.AcademicPeriodID,
		ClassID:          s.ClassID,
		SubjectID:        s.SubjectID,
		TeacherID:        s.TeacherID,
		SessionDate:      s.SessionDate,
		StartTime:        s.StartTime,
		EndTime:          s.EndTime,
		Status:           s.Status.String(),
		Notes:            s.Notes,
		Records:          records,
		CreatedAt:        s.CreatedAt,
		UpdatedAt:        s.UpdatedAt,
	}
}

func toRecordResponse(r *models.Record) RecordResponse {
	if r == nil {
		return RecordResponse{}
	}
	return RecordResponse{
		ID:        r.ID,
		SessionID: r.SessionID,
		StudentID: r.StudentID,
		Status:    r.Status.String(),
		Note:      r.Note,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
