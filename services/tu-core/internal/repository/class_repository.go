package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Class merepresentasikan entitas Rombongan Belajar di database.
type Class struct {
	ID               string
	SchoolID         string
	AcademicPeriodID string
	Name             string
	Grade            int // 10, 11, 12
	WaliTeacherID    string
	Level            int // 1, 2, 3
	Major            string
	HomeroomID       string // PTK ID
	Capacity         int
	Status           string // active, inactive
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}

// ClassFilter untuk query dengan filter.
type ClassFilter struct {
	SchoolID         string
	AcademicPeriodID string
	Grade            int
	Status           string
	Keyword          string
	Limit            int
	Offset           int
}

// ClassRepository adalah kontrak untuk akses data class.
type ClassRepository interface {
	Create(ctx context.Context, class *Class) error
	GetByID(ctx context.Context, schoolID, periodID, classID string) (*Class, error)
	List(ctx context.Context, filter ClassFilter) ([]Class, int, error)
	Update(ctx context.Context, class *Class) error
	Delete(ctx context.Context, schoolID, periodID, classID string) error
	CountActiveMembers(ctx context.Context, schoolID, periodID, classID string) (int, error)
	GetByGradeAndPeriod(ctx context.Context, schoolID, periodID string, grade int) ([]Class, error)
}

// pgxClassRepository mengimplementasikan ClassRepository.
type pgxClassRepository struct {
	db *pgxpool.Pool
}

// NewClassRepository membuat instance baru dari ClassRepository.
func NewClassRepository(db *pgxpool.Pool) ClassRepository {
	return &pgxClassRepository{db: db}
}

// Create menyimpan class baru ke database.
func (r *pgxClassRepository) Create(ctx context.Context, class *Class) error {
	if class.SchoolID == "" || class.AcademicPeriodID == "" {
		return ErrMissingOperationalContext
	}

	query := `
        INSERT INTO classes (
            id, school_id, period_id, name, grade, level, major, homeroom_id,
            capacity, status
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `

	_, err := r.db.Exec(ctx, query,
		class.ID, class.SchoolID, class.AcademicPeriodID, class.Name, class.Grade,
		class.Level, class.Major, class.HomeroomID, class.Capacity, class.Status,
	)

	if err != nil {
		if isDuplicateKeyError(err) {
			return fmt.Errorf("%w: %v", ErrDuplicateEntry, err)
		}
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return nil
}

// GetByID mengambil class berdasarkan ID.
func (r *pgxClassRepository) GetByID(ctx context.Context, schoolID, periodID, classID string) (*Class, error) {
	if schoolID == "" || periodID == "" {
		return nil, ErrMissingOperationalContext
	}

	query := `
        SELECT
            id, school_id, period_id, name, grade, level, major, homeroom_id,
            capacity, status, created_at, updated_at, deleted_at
        FROM classes
        WHERE id = $1 AND school_id = $2 AND period_id = $3 AND deleted_at IS NULL
    `

	var c Class
	err := r.db.QueryRow(ctx, query, classID, schoolID, periodID).Scan(
		&c.ID, &c.SchoolID, &c.AcademicPeriodID, &c.Name, &c.Grade, &c.Level,
		&c.Major, &c.HomeroomID, &c.Capacity, &c.Status,
		&c.CreatedAt, &c.UpdatedAt, &c.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return &c, nil
}

// List mengambil daftar class dengan filter dan pagination.
func (r *pgxClassRepository) List(ctx context.Context, filter ClassFilter) ([]Class, int, error) {
	if filter.SchoolID == "" || filter.AcademicPeriodID == "" {
		return nil, 0, ErrMissingOperationalContext
	}

	// Build WHERE clause
	whereClause := "WHERE school_id = $1 AND period_id = $2 AND deleted_at IS NULL"
	args := []interface{}{filter.SchoolID, filter.AcademicPeriodID}
	argIndex := 3

	if filter.Grade > 0 {
		whereClause += fmt.Sprintf(" AND grade = $%d", argIndex)
		args = append(args, filter.Grade)
		argIndex++
	}

	if filter.Status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	if filter.Keyword != "" {
		whereClause += fmt.Sprintf(" AND (name ILIKE $%d OR major ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filter.Keyword+"%")
		argIndex++
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM classes %s", whereClause)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	// Set default limit/offset
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	// Query data
	dataQuery := fmt.Sprintf(`
        SELECT
            id, school_id, period_id, name, grade, level, major, homeroom_id,
            capacity, status, created_at, updated_at, deleted_at
        FROM classes
        %s
        ORDER BY grade ASC, name ASC
        LIMIT $%d OFFSET $%d
    `, whereClause, argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var classes []Class
	for rows.Next() {
		var c Class
		if err := rows.Scan(
			&c.ID, &c.SchoolID, &c.AcademicPeriodID, &c.Name, &c.Grade, &c.Level,
			&c.Major, &c.HomeroomID, &c.Capacity, &c.Status,
			&c.CreatedAt, &c.UpdatedAt, &c.DeletedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		classes = append(classes, c)
	}

	return classes, total, rows.Err()
}

// Update memperbarui data class.
func (r *pgxClassRepository) Update(ctx context.Context, class *Class) error {
	if class.SchoolID == "" || class.AcademicPeriodID == "" {
		return ErrMissingOperationalContext
	}

	query := `
        UPDATE classes SET
            name = $1, grade = $2, level = $3, major = $4,
            homeroom_id = $5, capacity = $6, status = $7, updated_at = NOW()
        WHERE id = $8 AND school_id = $9 AND period_id = $10 AND deleted_at IS NULL
    `

	result, err := r.db.Exec(ctx, query,
		class.Name, class.Grade, class.Level, class.Major,
		class.HomeroomID, class.Capacity, class.Status,
		class.ID, class.SchoolID, class.AcademicPeriodID,
	)

	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// Delete melakukan soft delete pada class.
func (r *pgxClassRepository) Delete(ctx context.Context, schoolID, periodID, classID string) error {
	if schoolID == "" || periodID == "" {
		return ErrMissingOperationalContext
	}

	query := `
        UPDATE classes
        SET deleted_at = NOW(), updated_at = NOW()
        WHERE id = $1 AND school_id = $2 AND period_id = $3 AND deleted_at IS NULL
    `

	result, err := r.db.Exec(ctx, query, classID, schoolID, periodID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// CountActiveMembers menghitung jumlah siswa aktif di suatu class (via enrollments).
func (r *pgxClassRepository) CountActiveMembers(ctx context.Context, schoolID, periodID, classID string) (int, error) {
	if schoolID == "" || periodID == "" {
		return 0, ErrMissingOperationalContext
	}

	query := `
        SELECT COUNT(*)
        FROM enrollments
        WHERE class_id = $1 AND period_id = $2 AND status = 'active'
    `

	var count int
	err := r.db.QueryRow(ctx, query, classID, periodID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return count, nil
}

// GetByGradeAndPeriod mengambil class berdasarkan grade dan period (untuk validasi pindah rombel).
func (r *pgxClassRepository) GetByGradeAndPeriod(ctx context.Context, schoolID, periodID string, grade int) ([]Class, error) {
	if schoolID == "" || periodID == "" {
		return nil, ErrMissingOperationalContext
	}

	query := `
        SELECT
            id, school_id, period_id, name, grade, level, major, homeroom_id,
            capacity, status, created_at, updated_at, deleted_at
        FROM classes
        WHERE school_id = $1 AND period_id = $2 AND grade = $3 AND status = 'active' AND deleted_at IS NULL
        ORDER BY name ASC
    `

	rows, err := r.db.Query(ctx, query, schoolID, periodID, grade)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var classes []Class
	for rows.Next() {
		var c Class
		if err := rows.Scan(
			&c.ID, &c.SchoolID, &c.AcademicPeriodID, &c.Name, &c.Grade, &c.Level,
			&c.Major, &c.HomeroomID, &c.Capacity, &c.Status,
			&c.CreatedAt, &c.UpdatedAt, &c.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		classes = append(classes, c)
	}

	return classes, rows.Err()
}
