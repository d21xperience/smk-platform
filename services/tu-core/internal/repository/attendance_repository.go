package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"sekolah-platform/services/tu-core/internal/domain/attendance/models"
	"sekolah-platform/services/tu-core/internal/domain/attendance/valueobjects"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrSessionNotFound = errors.New("session not found")
var ErrRecordNotFound = errors.New("record not found")

type SessionFilter struct {
	SchoolID         string
	AcademicPeriodID string
	ClassID          string
	Status           string
	StartDate        time.Time
	EndDate          time.Time
	Limit            int
	Offset           int
}

type AttendanceRepository interface {
	CreateSession(ctx context.Context, s *models.Session) error
	GetSessionByID(ctx context.Context, schoolID, academicPeriodID, sessionID string) (*models.Session, error)
	ListSessions(ctx context.Context, filter SessionFilter) ([]models.Session, int, error)
	UpdateSession(ctx context.Context, s *models.Session) error
	CreateRecord(ctx context.Context, r *models.Record) error
	GetRecordsBySession(ctx context.Context, sessionID string) ([]models.Record, error)
	GetRecordsByStudent(ctx context.Context, schoolID, academicPeriodID, studentID string, month, year int) ([]models.Record, error)
}

type pgxAttendanceRepository struct {
	db *pgxpool.Pool
}

func NewAttendanceRepository(db *pgxpool.Pool) AttendanceRepository {
	return &pgxAttendanceRepository{db: db}
}

func (r *pgxAttendanceRepository) CreateSession(ctx context.Context, s *models.Session) error {
	query := `
		INSERT INTO attendance_sessions (
			id, school_id, academic_period_id, class_id, subject_id, teacher_id,
			session_date, start_time, end_time, status, notes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.Exec(ctx, query,
		s.ID, s.SchoolID, s.AcademicPeriodID, s.ClassID, s.SubjectID, s.TeacherID,
		s.SessionDate, s.StartTime, s.EndTime, s.Status.String(), s.Notes, s.CreatedAt, s.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxAttendanceRepository) GetSessionByID(ctx context.Context, schoolID, academicPeriodID, sessionID string) (*models.Session, error) {
	query := `
		SELECT id, school_id, academic_period_id, class_id, subject_id, teacher_id,
		       session_date, start_time, end_time, status, notes, created_at, updated_at
		FROM attendance_sessions
		WHERE id = $1 AND school_id = $2 AND academic_period_id = $3
	`

	var s models.Session
	var statusStr string

	err := r.db.QueryRow(ctx, query, sessionID, schoolID, academicPeriodID).Scan(
		&s.ID, &s.SchoolID, &s.AcademicPeriodID, &s.ClassID, &s.SubjectID, &s.TeacherID,
		&s.SessionDate, &s.StartTime, &s.EndTime, &statusStr, &s.Notes, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	s.Status, _ = valueobjects.NewSessionStatus(statusStr)

	// Load records
	records, err := r.GetRecordsBySession(ctx, s.ID)
	if err != nil {
		return nil, err
	}
	s.Records = make([]*models.Record, len(records))
	for i := range records {
		s.Records[i] = &records[i]
	}

	return &s, nil
}

func (r *pgxAttendanceRepository) ListSessions(ctx context.Context, filter SessionFilter) ([]models.Session, int, error) {
	whereClause := "WHERE school_id = $1 AND academic_period_id = $2"
	args := []interface{}{filter.SchoolID, filter.AcademicPeriodID}
	argIndex := 3

	if filter.ClassID != "" {
		whereClause += fmt.Sprintf(" AND class_id = $%d", argIndex)
		args = append(args, filter.ClassID)
		argIndex++
	}

	if filter.Status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	if !filter.StartDate.IsZero() {
		whereClause += fmt.Sprintf(" AND session_date >= $%d", argIndex)
		args = append(args, filter.StartDate)
		argIndex++
	}

	if !filter.EndDate.IsZero() {
		whereClause += fmt.Sprintf(" AND session_date <= $%d", argIndex)
		args = append(args, filter.EndDate)
		argIndex++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM attendance_sessions %s", whereClause)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	dataQuery := fmt.Sprintf(`
		SELECT id, school_id, academic_period_id, class_id, subject_id, teacher_id,
		       session_date, start_time, end_time, status, notes, created_at, updated_at
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

	var sessions []models.Session
	for rows.Next() {
		var s models.Session
		var statusStr string

		if err := rows.Scan(
			&s.ID, &s.SchoolID, &s.AcademicPeriodID, &s.ClassID, &s.SubjectID, &s.TeacherID,
			&s.SessionDate, &s.StartTime, &s.EndTime, &statusStr, &s.Notes, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}

		s.Status, _ = valueobjects.NewSessionStatus(statusStr)
		sessions = append(sessions, s)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return sessions, total, nil
}

func (r *pgxAttendanceRepository) UpdateSession(ctx context.Context, s *models.Session) error {
	query := `
		UPDATE attendance_sessions
		SET status = $1, notes = $2, updated_at = $3
		WHERE id = $4 AND school_id = $5 AND academic_period_id = $6
	`
	result, err := r.db.Exec(ctx, query,
		s.Status.String(), s.Notes, s.UpdatedAt,
		s.ID, s.SchoolID, s.AcademicPeriodID,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrSessionNotFound
	}
	return nil
}

func (r *pgxAttendanceRepository) CreateRecord(ctx context.Context, rec *models.Record) error {
	query := `
		INSERT INTO attendance_records (
			id, session_id, student_id, status, note, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query,
		rec.ID, rec.SessionID, rec.StudentID, rec.Status.String(), rec.Note, rec.CreatedAt, rec.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxAttendanceRepository) GetRecordsBySession(ctx context.Context, sessionID string) ([]models.Record, error) {
	query := `
		SELECT id, session_id, student_id, status, note, created_at, updated_at
		FROM attendance_records
		WHERE session_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var records []models.Record
	for rows.Next() {
		var rec models.Record
		var statusStr string

		if err := rows.Scan(
			&rec.ID, &rec.SessionID, &rec.StudentID, &statusStr, &rec.Note, &rec.CreatedAt, &rec.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}

		rec.Status, _ = valueobjects.NewAttendanceStatus(statusStr)
		records = append(records, rec)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return records, nil
}

func (r *pgxAttendanceRepository) GetRecordsByStudent(ctx context.Context, schoolID, academicPeriodID, studentID string, month, year int) ([]models.Record, error) {
	query := `
		SELECT r.id, r.session_id, r.student_id, r.status, r.note, r.created_at, r.updated_at
		FROM attendance_records r
		JOIN attendance_sessions s ON r.session_id = s.id
		WHERE s.school_id = $1 AND s.academic_period_id = $2 AND r.student_id = $3
		  AND EXTRACT(MONTH FROM s.session_date) = $4
		  AND EXTRACT(YEAR FROM s.session_date) = $5
		ORDER BY s.session_date DESC
	`

	rows, err := r.db.Query(ctx, query, schoolID, academicPeriodID, studentID, month, year)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var records []models.Record
	for rows.Next() {
		var rec models.Record
		var statusStr string

		if err := rows.Scan(
			&rec.ID, &rec.SessionID, &rec.StudentID, &statusStr, &rec.Note, &rec.CreatedAt, &rec.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}

		rec.Status, _ = valueobjects.NewAttendanceStatus(statusStr)
		records = append(records, rec)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return records, nil
}
