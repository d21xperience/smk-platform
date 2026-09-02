package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AcademicPeriod merepresentasikan periode akademik (tahun ajaran).
type AcademicPeriod struct {
	ID        string
	SchoolID  string
	Name      string // e.g., "2026/2027"
	StartDate time.Time
	EndDate   time.Time
	IsActive  bool
	CreatedAt time.Time
}

// Class merepresentasikan kelas.
type Class struct {
	ID               string
	SchoolID         string
	AcademicPeriodID string
	Name             string // e.g., "7A", "10-IPA-1"
	Grade            string // e.g., "7", "10"
	WaliTeacherID    string
	Capacity         int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Subject merepresentasikan mata pelajaran.
type Subject struct {
	ID       string
	SchoolID string
	Name     string // e.g., "Matematika"
	Code     string // e.g., "MTK"
	Credits  int
}

// AcademicRepository adalah kontrak untuk akses data akademik.
type AcademicRepository interface {
	// Academic Period
	GetActivePeriod(ctx context.Context, schoolID string) (*AcademicPeriod, error)
	GetPeriodByID(ctx context.Context, schoolID, periodID string) (*AcademicPeriod, error)

	// Class
	GetClassByID(ctx context.Context, schoolID, academicPeriodID, classID string) (*Class, error)
	ListClasses(ctx context.Context, schoolID, academicPeriodID string) ([]Class, error)

	// Subject
	GetSubjectByID(ctx context.Context, schoolID, subjectID string) (*Subject, error)
	ListSubjects(ctx context.Context, schoolID string) ([]Subject, error)
}

// pgxAcademicRepository mengimplementasikan AcademicRepository.
type pgxAcademicRepository struct {
	db *pgxpool.Pool
}

// NewAcademicRepository membuat instance baru.
func NewAcademicRepository(db *pgxpool.Pool) AcademicRepository {
	return &pgxAcademicRepository{db: db}
}

// GetActivePeriod mengambil periode akademik yang sedang aktif.
func (r *pgxAcademicRepository) GetActivePeriod(ctx context.Context, schoolID string) (*AcademicPeriod, error) {
	if schoolID == "" {
		return nil, ErrMissingOperationalContext
	}

	query := `
		SELECT id, school_id, name, start_date, end_date, is_active, created_at
		FROM academic_periods
		WHERE school_id = $1 AND is_active = true
		LIMIT 1
	`

	var p AcademicPeriod
	err := r.db.QueryRow(ctx, query, schoolID).Scan(
		&p.ID, &p.SchoolID, &p.Name, &p.StartDate, &p.EndDate, &p.IsActive, &p.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return &p, nil
}

// GetPeriodByID mengambil periode akademik berdasarkan ID.
func (r *pgxAcademicRepository) GetPeriodByID(ctx context.Context, schoolID, periodID string) (*AcademicPeriod, error) {
	if schoolID == "" {
		return nil, ErrMissingOperationalContext
	}

	query := `
		SELECT id, school_id, name, start_date, end_date, is_active, created_at
		FROM academic_periods
		WHERE id = $1 AND school_id = $2
	`

	var p AcademicPeriod
	err := r.db.QueryRow(ctx, query, periodID, schoolID).Scan(
		&p.ID, &p.SchoolID, &p.Name, &p.StartDate, &p.EndDate, &p.IsActive, &p.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return &p, nil
}

// GetClassByID mengambil kelas berdasarkan ID.
func (r *pgxAcademicRepository) GetClassByID(ctx context.Context, schoolID, academicPeriodID, classID string) (*Class, error) {
	if schoolID == "" || academicPeriodID == "" {
		return nil, ErrMissingOperationalContext
	}

	query := `
		SELECT id, school_id, academic_period_id, name, grade, wali_teacher_id, capacity, created_at, updated_at
		FROM classes
		WHERE id = $1 AND school_id = $2 AND academic_period_id = $3
	`

	var c Class
	err := r.db.QueryRow(ctx, query, classID, schoolID, academicPeriodID).Scan(
		&c.ID, &c.SchoolID, &c.AcademicPeriodID, &c.Name, &c.Grade, &c.WaliTeacherID,
		&c.Capacity, &c.CreatedAt, &c.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return &c, nil
}

// ListClasses mengambil daftar kelas dalam periode akademik.
func (r *pgxAcademicRepository) ListClasses(ctx context.Context, schoolID, academicPeriodID string) ([]Class, error) {
	if schoolID == "" || academicPeriodID == "" {
		return nil, ErrMissingOperationalContext
	}

	query := `
		SELECT id, school_id, academic_period_id, name, grade, wali_teacher_id, capacity, created_at, updated_at
		FROM classes
		WHERE school_id = $1 AND academic_period_id = $2
		ORDER BY grade ASC, name ASC
	`

	rows, err := r.db.Query(ctx, query, schoolID, academicPeriodID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var classes []Class
	for rows.Next() {
		var c Class
		if err := rows.Scan(
			&c.ID, &c.SchoolID, &c.AcademicPeriodID, &c.Name, &c.Grade, &c.WaliTeacherID,
			&c.Capacity, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		classes = append(classes, c)
	}

	return classes, rows.Err()
}

// GetSubjectByID mengambil mata pelajaran berdasarkan ID.
func (r *pgxAcademicRepository) GetSubjectByID(ctx context.Context, schoolID, subjectID string) (*Subject, error) {
	if schoolID == "" {
		return nil, ErrMissingOperationalContext
	}

	query := `
		SELECT id, school_id, name, code, credits
		FROM subjects
		WHERE id = $1 AND school_id = $2
	`

	var s Subject
	err := r.db.QueryRow(ctx, query, subjectID, schoolID).Scan(
		&s.ID, &s.SchoolID, &s.Name, &s.Code, &s.Credits,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return &s, nil
}

// ListSubjects mengambil daftar mata pelajaran di sekolah.
func (r *pgxAcademicRepository) ListSubjects(ctx context.Context, schoolID string) ([]Subject, error) {
	if schoolID == "" {
		return nil, ErrMissingOperationalContext
	}

	query := `
		SELECT id, school_id, name, code, credits
		FROM subjects
		WHERE school_id = $1
		ORDER BY name ASC
	`

	rows, err := r.db.Query(ctx, query, schoolID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var subjects []Subject
	for rows.Next() {
		var s Subject
		if err := rows.Scan(&s.ID, &s.SchoolID, &s.Name, &s.Code, &s.Credits); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		subjects = append(subjects, s)
	}

	return subjects, rows.Err()
}
