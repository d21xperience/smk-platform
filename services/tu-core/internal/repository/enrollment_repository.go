package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Enrollment merepresentasikan entitas Anggota Rombel di database.
type Enrollment struct {
	ID        string
	StudentID string
	SchoolID  string
	PeriodID  string
	ClassID   string
	Status    string // active, completed, transferred
	Reason    string // alasan perubahan status
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// EnrollmentFilter untuk query dengan filter.
type EnrollmentFilter struct {
	SchoolID  string
	PeriodID  string
	StudentID string
	ClassID   string
	Status    string
	Limit     int
	Offset    int
}

// EnrollmentRepository adalah kontrak untuk akses data enrollment.
type EnrollmentRepository interface {
	Create(ctx context.Context, enrollment *Enrollment) error
	GetByID(ctx context.Context, schoolID, periodID, enrollmentID string) (*Enrollment, error)
	List(ctx context.Context, filter EnrollmentFilter) ([]Enrollment, int, error)
	Update(ctx context.Context, enrollment *Enrollment) error
	Delete(ctx context.Context, schoolID, periodID, enrollmentID string) error
	GetActiveByStudent(ctx context.Context, schoolID, periodID, studentID string) (*Enrollment, error)
	GetHistoryByStudent(ctx context.Context, schoolID, studentID string) ([]Enrollment, error)
	GetMembersByClass(ctx context.Context, schoolID, periodID, classID string) ([]Enrollment, error)
	CountActiveByClass(ctx context.Context, schoolID, periodID, classID string) (int, error)
}

// pgxEnrollmentRepository mengimplementasikan EnrollmentRepository.
type pgxEnrollmentRepository struct {
	db *pgxpool.Pool
}

// NewEnrollmentRepository membuat instance baru dari EnrollmentRepository.
func NewEnrollmentRepository(db *pgxpool.Pool) EnrollmentRepository {
	return &pgxEnrollmentRepository{db: db}
}

// Create menyimpan enrollment baru ke database.
func (r *pgxEnrollmentRepository) Create(ctx context.Context, enrollment *Enrollment) error {
	if enrollment.SchoolID == "" || enrollment.PeriodID == "" {
		return ErrMissingOperationalContext
	}

	query := `
        INSERT INTO enrollments (
            id, student_id, school_id, period_id, class_id, status, reason
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	_, err := r.db.Exec(ctx, query,
		enrollment.ID, enrollment.StudentID, enrollment.SchoolID,
		enrollment.PeriodID, enrollment.ClassID, enrollment.Status, enrollment.Reason,
	)

	if err != nil {
		if isDuplicateKeyError(err) {
			return fmt.Errorf("%w: %v", ErrDuplicateEntry, err)
		}
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return nil
}

// GetByID mengambil enrollment berdasarkan ID.
func (r *pgxEnrollmentRepository) GetByID(ctx context.Context, schoolID, periodID, enrollmentID string) (*Enrollment, error) {
	if schoolID == "" || periodID == "" {
		return nil, ErrMissingOperationalContext
	}

	query := `
        SELECT
            id, student_id, school_id, period_id, class_id, status, reason,
            created_at, updated_at, deleted_at
        FROM enrollments
        WHERE id = $1 AND school_id = $2 AND period_id = $3 AND deleted_at IS NULL
    `

	var e Enrollment
	err := r.db.QueryRow(ctx, query, enrollmentID, schoolID, periodID).Scan(
		&e.ID, &e.StudentID, &e.SchoolID, &e.PeriodID,
		&e.ClassID, &e.Status, &e.Reason,
		&e.CreatedAt, &e.UpdatedAt, &e.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return &e, nil
}

// List mengambil daftar enrollment dengan filter dan pagination.
func (r *pgxEnrollmentRepository) List(ctx context.Context, filter EnrollmentFilter) ([]Enrollment, int, error) {
	if filter.SchoolID == "" || filter.PeriodID == "" {
		return nil, 0, ErrMissingOperationalContext
	}

	// Build WHERE clause
	whereClause := "WHERE school_id = $1 AND period_id = $2 AND deleted_at IS NULL"
	args := []interface{}{filter.SchoolID, filter.PeriodID}
	argIndex := 3

	if filter.StudentID != "" {
		whereClause += fmt.Sprintf(" AND student_id = $%d", argIndex)
		args = append(args, filter.StudentID)
		argIndex++
	}

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

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM enrollments %s", whereClause)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	// Set default limit/offset
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	// Query data
	dataQuery := fmt.Sprintf(`
        SELECT
            id, student_id, school_id, period_id, class_id, status, reason,
            created_at, updated_at, deleted_at
        FROM enrollments
        %s
        ORDER BY created_at DESC
        LIMIT $%d OFFSET $%d
    `, whereClause, argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var enrollments []Enrollment
	for rows.Next() {
		var e Enrollment
		if err := rows.Scan(
			&e.ID, &e.StudentID, &e.SchoolID, &e.PeriodID,
			&e.ClassID, &e.Status, &e.Reason,
			&e.CreatedAt, &e.UpdatedAt, &e.DeletedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		enrollments = append(enrollments, e)
	}

	return enrollments, total, rows.Err()
}

// Update memperbarui data enrollment.
func (r *pgxEnrollmentRepository) Update(ctx context.Context, enrollment *Enrollment) error {
	if enrollment.SchoolID == "" || enrollment.PeriodID == "" {
		return ErrMissingOperationalContext
	}

	query := `
        UPDATE enrollments SET
            class_id = $1, status = $2, reason = $3, updated_at = NOW()
        WHERE id = $4 AND school_id = $5 AND period_id = $6 AND deleted_at IS NULL
    `

	result, err := r.db.Exec(ctx, query,
		enrollment.ClassID, enrollment.Status, enrollment.Reason,
		enrollment.ID, enrollment.SchoolID, enrollment.PeriodID,
	)

	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// Delete melakukan soft delete pada enrollment.
func (r *pgxEnrollmentRepository) Delete(ctx context.Context, schoolID, periodID, enrollmentID string) error {
	if schoolID == "" || periodID == "" {
		return ErrMissingOperationalContext
	}

	query := `
        UPDATE enrollments
        SET deleted_at = NOW(), updated_at = NOW()
        WHERE id = $1 AND school_id = $2 AND period_id = $3 AND deleted_at IS NULL
    `

	result, err := r.db.Exec(ctx, query, enrollmentID, schoolID, periodID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// GetActiveByStudent mengambil enrollment aktif siswa di suatu periode.
func (r *pgxEnrollmentRepository) GetActiveByStudent(ctx context.Context, schoolID, periodID, studentID string) (*Enrollment, error) {
	if schoolID == "" || periodID == "" {
		return nil, ErrMissingOperationalContext
	}

	query := `
        SELECT
            id, student_id, school_id, period_id, class_id, status, reason,
            created_at, updated_at, deleted_at
        FROM enrollments
        WHERE student_id = $1 AND school_id = $2 AND period_id = $3
          AND status = 'active' AND deleted_at IS NULL
        LIMIT 1
    `

	var e Enrollment
	err := r.db.QueryRow(ctx, query, studentID, schoolID, periodID).Scan(
		&e.ID, &e.StudentID, &e.SchoolID, &e.PeriodID,
		&e.ClassID, &e.Status, &e.Reason,
		&e.CreatedAt, &e.UpdatedAt, &e.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return &e, nil
}

// GetHistoryByStudent mengambil seluruh histori enrollment seorang siswa.
func (r *pgxEnrollmentRepository) GetHistoryByStudent(ctx context.Context, schoolID, studentID string) ([]Enrollment, error) {
	if schoolID == "" {
		return nil, ErrMissingOperationalContext
	}

	query := `
        SELECT
            id, student_id, school_id, period_id, class_id, status, reason,
            created_at, updated_at, deleted_at
        FROM enrollments
        WHERE student_id = $1 AND school_id = $2 AND deleted_at IS NULL
        ORDER BY period_id DESC, created_at DESC
    `

	rows, err := r.db.Query(ctx, query, studentID, schoolID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var enrollments []Enrollment
	for rows.Next() {
		var e Enrollment
		if err := rows.Scan(
			&e.ID, &e.StudentID, &e.SchoolID, &e.PeriodID,
			&e.ClassID, &e.Status, &e.Reason,
			&e.CreatedAt, &e.UpdatedAt, &e.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		enrollments = append(enrollments, e)
	}

	return enrollments, rows.Err()
}

// GetMembersByClass mengambil semua enrollment aktif di suatu class.
func (r *pgxEnrollmentRepository) GetMembersByClass(ctx context.Context, schoolID, periodID, classID string) ([]Enrollment, error) {
	if schoolID == "" || periodID == "" {
		return nil, ErrMissingOperationalContext
	}

	query := `
        SELECT
            id, student_id, school_id, period_id, class_id, status, reason,
            created_at, updated_at, deleted_at
        FROM enrollments
        WHERE class_id = $1 AND school_id = $2 AND period_id = $3
          AND status = 'active' AND deleted_at IS NULL
        ORDER BY created_at ASC
    `

	rows, err := r.db.Query(ctx, query, classID, schoolID, periodID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var enrollments []Enrollment
	for rows.Next() {
		var e Enrollment
		if err := rows.Scan(
			&e.ID, &e.StudentID, &e.SchoolID, &e.PeriodID,
			&e.ClassID, &e.Status, &e.Reason,
			&e.CreatedAt, &e.UpdatedAt, &e.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		enrollments = append(enrollments, e)
	}

	return enrollments, rows.Err()
}

// CountActiveByClass menghitung jumlah siswa aktif di suatu class.
func (r *pgxEnrollmentRepository) CountActiveByClass(ctx context.Context, schoolID, periodID, classID string) (int, error) {
	if schoolID == "" || periodID == "" {
		return 0, ErrMissingOperationalContext
	}

	query := `
        SELECT COUNT(*)
        FROM enrollments
        WHERE class_id = $1 AND school_id = $2 AND period_id = $3
          AND status = 'active' AND deleted_at IS NULL
    `

	var count int
	err := r.db.QueryRow(ctx, query, classID, schoolID, periodID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return count, nil
}
