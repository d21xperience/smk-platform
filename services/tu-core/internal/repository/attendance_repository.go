package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ===== STRUCTS =====

type AttendanceSession struct {
	ID               string
	SchoolID         string
	AcademicPeriodID string
	ClassID          string
	SubjectID        string
	TeacherID        string
	SessionDate      time.Time
	StartTime        string
	EndTime          string
	Status           string
	Notes            string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type AttendanceRecord struct {
	ID        string
	SessionID string
	StudentID string
	Status    string
	Note      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SessionFilter struct {
	SchoolID         string
	AcademicPeriodID string
	ClassID          string
	Status           string
	StartDate        string
	EndDate          string
	Limit            int
	Offset           int
}

type AttendanceSummary struct {
	TotalSessions        int
	PresentCount         int
	AbsentCount          int
	LateCount            int
	SickCount            int
	PermissionCount      int
	AttendancePercentage float64
}

// ===== INTERFACE =====

type AttendanceRepository interface {
	// Session CRUD
	CreateSession(ctx context.Context, s *AttendanceSession) error
	GetSessionByID(ctx context.Context, schoolID, periodID, sessionID string) (*AttendanceSession, error)
	ListSessions(ctx context.Context, filter SessionFilter) ([]AttendanceSession, int, error)
	UpdateSessionStatus(ctx context.Context, sessionID, status string) error

	// Records CRUD
	UpsertRecord(ctx context.Context, r *AttendanceRecord) error
	GetRecordsBySession(ctx context.Context, sessionID string) ([]AttendanceRecord, error)

	// Summary
	GetStudentAttendanceSummary(ctx context.Context, schoolID, academicPeriodID, studentID string) (*AttendanceSummary, error)
}

type pgxAttendanceRepository struct {
	db *pgxpool.Pool
}

func NewAttendanceRepository(db *pgxpool.Pool) AttendanceRepository {
	return &pgxAttendanceRepository{db: db}
}

// ===== SESSION IMPLEMENTATIONS =====

func (r *pgxAttendanceRepository) CreateSession(ctx context.Context, s *AttendanceSession) error {
	query := `
		INSERT INTO attendance_sessions (
			id, school_id, academic_period_id, class_id, subject_id, teacher_id,
			session_date, start_time, end_time, status, notes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.Exec(ctx, query,
		s.ID, s.SchoolID, s.AcademicPeriodID, s.ClassID, s.SubjectID, s.TeacherID,
		s.SessionDate, s.StartTime, s.EndTime, s.Status, s.Notes, s.CreatedAt, s.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxAttendanceRepository) GetSessionByID(ctx context.Context, schoolID, periodID, sessionID string) (*AttendanceSession, error) {
	query := `
		SELECT id, school_id, academic_period_id, class_id, COALESCE(subject_id::text, ''), COALESCE(teacher_id::text, ''),
		       session_date, COALESCE(start_time::text, ''), COALESCE(end_time::text, ''), status, COALESCE(notes, ''), created_at, updated_at
		FROM attendance_sessions
		WHERE id = $1 AND school_id = $2 AND academic_period_id = $3
	`
	var s AttendanceSession
	var startTime, endTime *string
	err := r.db.QueryRow(ctx, query, sessionID, schoolID, periodID).Scan(
		&s.ID, &s.SchoolID, &s.AcademicPeriodID, &s.ClassID, &s.SubjectID, &s.TeacherID,
		&s.SessionDate, &startTime, &endTime, &s.Status, &s.Notes, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if startTime != nil { s.StartTime = *startTime }
	if endTime != nil { s.EndTime = *endTime }
	return &s, nil
}

func (r *pgxAttendanceRepository) ListSessions(ctx context.Context, filter SessionFilter) ([]AttendanceSession, int, error) {
	whereClause := "WHERE school_id = $1 AND academic_period_id = $2"
	args := []any{filter.SchoolID, filter.AcademicPeriodID}
	argIndex := 3

	if filter.ClassID != "" {
		whereClause += fmt.Sprintf(" AND class_id::text = $%d", argIndex)
		args = append(args, filter.ClassID)
		argIndex++
	}
	if filter.Status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM attendance_sessions %s", whereClause)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if filter.Limit <= 0 { filter.Limit = 20 }
	if filter.Offset < 0 { filter.Offset = 0 }

	dataQuery := fmt.Sprintf(`
		SELECT id, school_id, academic_period_id, class_id, COALESCE(subject_id::text, ''), COALESCE(teacher_id::text, ''),
		       session_date, COALESCE(start_time::text, ''), COALESCE(end_time::text, ''), status, COALESCE(notes, ''), created_at, updated_at
		FROM attendance_sessions %s
		ORDER BY session_date DESC, created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var sessions []AttendanceSession
	for rows.Next() {
		var s AttendanceSession
		var startTime, endTime *string
		if err := rows.Scan(
			&s.ID, &s.SchoolID, &s.AcademicPeriodID, &s.ClassID, &s.SubjectID, &s.TeacherID,
			&s.SessionDate, &startTime, &endTime, &s.Status, &s.Notes, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		if startTime != nil { s.StartTime = *startTime }
		if endTime != nil { s.EndTime = *endTime }
		sessions = append(sessions, s)
	}
	return sessions, total, rows.Err()
}

func (r *pgxAttendanceRepository) UpdateSessionStatus(ctx context.Context, sessionID, status string) error {
	query := `UPDATE attendance_sessions SET status = $1, updated_at = NOW() WHERE id = $2`
	result, err := r.db.Exec(ctx, query, status, sessionID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// ===== RECORD IMPLEMENTATIONS =====

func (r *pgxAttendanceRepository) UpsertRecord(ctx context.Context, rec *AttendanceRecord) error {
	query := `
		INSERT INTO attendance_records (id, session_id, student_id, status, note, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (session_id, student_id) DO UPDATE SET
			status = EXCLUDED.status,
			note = EXCLUDED.note,
			updated_at = NOW()
	`
	_, err := r.db.Exec(ctx, query,
		rec.ID, rec.SessionID, rec.StudentID, rec.Status, rec.Note, rec.CreatedAt, rec.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxAttendanceRepository) GetRecordsBySession(ctx context.Context, sessionID string) ([]AttendanceRecord, error) {
	query := `
		SELECT id, session_id, student_id, status, COALESCE(note, ''), created_at, updated_at
		FROM attendance_records
		WHERE session_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.Query(ctx, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var records []AttendanceRecord
	for rows.Next() {
		var rec AttendanceRecord
		if err := rows.Scan(&rec.ID, &rec.SessionID, &rec.StudentID, &rec.Status, &rec.Note, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		records = append(records, rec)
	}
	return records, rows.Err()
}

// ===== SUMMARY (Existing) =====

func (r *pgxAttendanceRepository) GetStudentAttendanceSummary(ctx context.Context, schoolID, academicPeriodID, studentID string) (*AttendanceSummary, error) {
	query := `
		SELECT
			COALESCE(SUM(tk.sakit), 0) as sakit,
			COALESCE(SUM(tk.izin), 0) as izin,
			COALESCE(SUM(tk.tanpa_keterangan), 0) as alpa
		FROM tabel_kehadiran tk
		JOIN tabel_anggotakelas tak ON tk.anggota_rombel_id = tak.anggota_rombel_id
		WHERE tak.peserta_didik_id = $1
		  AND tak.semester_id = $2
		  AND tak.rombongan_belajar_id IN (
			  SELECT rombongan_belajar_id FROM tabel_kelas WHERE sekolah_id = $3
		  )
	`
	var sakit, izin, alpa int
	err := r.db.QueryRow(ctx, query, studentID, academicPeriodID, schoolID).Scan(&sakit, &izin, &alpa)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &AttendanceSummary{}, nil
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	total := sakit + izin + alpa
	summary := &AttendanceSummary{
		TotalSessions: total, SickCount: sakit, PermissionCount: izin, AbsentCount: alpa,
		PresentCount: 0, LateCount: 0,
	}
	if summary.TotalSessions > 0 {
		summary.AttendancePercentage = float64(summary.PresentCount+summary.LateCount) / float64(summary.TotalSessions) * 100
	}
	return summary, nil
}
