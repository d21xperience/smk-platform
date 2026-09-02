package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Teacher merepresentasikan entitas Teacher di database.
type Teacher struct {
	ID               string
	SchoolID         string
	AcademicPeriodID string
	NIP              string
	Name             string
	Gender           string
	BirthDate        time.Time
	Email            string
	Phone            string
	Address          string
	SubjectID        string
	Status           string // ACTIVE, INACTIVE
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}

// TeacherFilter untuk query dengan filter.
type TeacherFilter struct {
	SchoolID         string
	AcademicPeriodID string
	SubjectID        string
	Status           string
	Keyword          string
	Limit            int
	Offset           int
}

// TeacherRepository adalah kontrak untuk akses data teacher.
type TeacherRepository interface {
	Create(ctx context.Context, teacher *Teacher) error
	GetByID(ctx context.Context, schoolID, academicPeriodID, teacherID string) (*Teacher, error)
	GetByNIP(ctx context.Context, schoolID, academicPeriodID, nip string) (*Teacher, error)
	List(ctx context.Context, filter TeacherFilter) ([]Teacher, int, error)
	Update(ctx context.Context, teacher *Teacher) error
	Delete(ctx context.Context, schoolID, academicPeriodID, teacherID string) error
}

// pgxTeacherRepository mengimplementasikan TeacherRepository.
type pgxTeacherRepository struct {
	db *pgxpool.Pool
}

// NewTeacherRepository membuat instance baru dari TeacherRepository.
func NewTeacherRepository(db *pgxpool.Pool) TeacherRepository {
	return &pgxTeacherRepository{db: db}
}

// Create menyimpan teacher baru.
func (r *pgxTeacherRepository) Create(ctx context.Context, teacher *Teacher) error {
	if teacher.SchoolID == "" || teacher.AcademicPeriodID == "" {
		return ErrMissingOperationalContext
	}

	query := `
		INSERT INTO teachers (
			id, school_id, academic_period_id, nip, name, gender, birth_date, email, phone,
			address, subject_id, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.Exec(ctx, query,
		teacher.ID, teacher.SchoolID, teacher.AcademicPeriodID, teacher.NIP, teacher.Name,
		teacher.Gender, teacher.BirthDate, teacher.Email, teacher.Phone, teacher.Address,
		teacher.SubjectID, teacher.Status,
	)

	if err != nil {
		if isDuplicateKeyError(err) {
			return fmt.Errorf("%w: %v", ErrDuplicateEntry, err)
		}
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return nil
}

// GetByID mengambil teacher berdasarkan ID.
func (r *pgxTeacherRepository) GetByID(ctx context.Context, schoolID, academicPeriodID, teacherID string) (*Teacher, error) {
	if schoolID == "" || academicPeriodID == "" {
		return nil, ErrMissingOperationalContext
	}

	query := `
		SELECT
			id, school_id, academic_period_id, nip, name, gender, birth_date, email, phone,
			address, subject_id, status, created_at, updated_at, deleted_at
		FROM teachers
		WHERE id = $1 AND school_id = $2 AND academic_period_id = $3 AND deleted_at IS NULL
	`

	var t Teacher
	err := r.db.QueryRow(ctx, query, teacherID, schoolID, academicPeriodID).Scan(
		&t.ID, &t.SchoolID, &t.AcademicPeriodID, &t.NIP, &t.Name, &t.Gender, &t.BirthDate,
		&t.Email, &t.Phone, &t.Address, &t.SubjectID, &t.Status, &t.CreatedAt, &t.UpdatedAt,
		&t.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return &t, nil
}

// GetByNIP mengambil teacher berdasarkan NIP.
func (r *pgxTeacherRepository) GetByNIP(ctx context.Context, schoolID, academicPeriodID, nip string) (*Teacher, error) {
	if schoolID == "" || academicPeriodID == "" {
		return nil, ErrMissingOperationalContext
	}

	query := `
		SELECT
			id, school_id, academic_period_id, nip, name, gender, birth_date, email, phone,
			address, subject_id, status, created_at, updated_at, deleted_at
		FROM teachers
		WHERE nip = $1 AND school_id = $2 AND academic_period_id = $3 AND deleted_at IS NULL
	`

	var t Teacher
	err := r.db.QueryRow(ctx, query, nip, schoolID, academicPeriodID).Scan(
		&t.ID, &t.SchoolID, &t.AcademicPeriodID, &t.NIP, &t.Name, &t.Gender, &t.BirthDate,
		&t.Email, &t.Phone, &t.Address, &t.SubjectID, &t.Status, &t.CreatedAt, &t.UpdatedAt,
		&t.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return &t, nil
}

// List mengambil daftar teacher dengan filter.
func (r *pgxTeacherRepository) List(ctx context.Context, filter TeacherFilter) ([]Teacher, int, error) {
	if filter.SchoolID == "" || filter.AcademicPeriodID == "" {
		return nil, 0, ErrMissingOperationalContext
	}

	whereClause := "WHERE school_id = $1 AND academic_period_id = $2 AND deleted_at IS NULL"
	args := []interface{}{filter.SchoolID, filter.AcademicPeriodID}
	argIndex := 3

	if filter.SubjectID != "" {
		whereClause += fmt.Sprintf(" AND subject_id = $%d", argIndex)
		args = append(args, filter.SubjectID)
		argIndex++
	}

	if filter.Status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	if filter.Keyword != "" {
		whereClause += fmt.Sprintf(" AND (name ILIKE $%d OR nip ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
		argIndex++
	}

	// Count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM teachers %s", whereClause)
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
		SELECT
			id, school_id, academic_period_id, nip, name, gender, birth_date, email, phone,
			address, subject_id, status, created_at, updated_at, deleted_at
		FROM teachers %s ORDER BY name ASC LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var teachers []Teacher
	for rows.Next() {
		var t Teacher
		if err := rows.Scan(
			&t.ID, &t.SchoolID, &t.AcademicPeriodID, &t.NIP, &t.Name, &t.Gender, &t.BirthDate,
			&t.Email, &t.Phone, &t.Address, &t.SubjectID, &t.Status, &t.CreatedAt, &t.UpdatedAt,
			&t.DeletedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		teachers = append(teachers, t)
	}

	return teachers, total, rows.Err()
}

// Update memperbarui data teacher.
func (r *pgxTeacherRepository) Update(ctx context.Context, teacher *Teacher) error {
	if teacher.SchoolID == "" || teacher.AcademicPeriodID == "" {
		return ErrMissingOperationalContext
	}

	query := `
		UPDATE teachers SET
			nip = $1, name = $2, gender = $3, birth_date = $4, email = $5, phone = $6,
			address = $7, subject_id = $8, status = $9, updated_at = NOW()
		WHERE id = $10 AND school_id = $11 AND academic_period_id = $12 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(ctx, query,
		teacher.NIP, teacher.Name, teacher.Gender, teacher.BirthDate, teacher.Email,
		teacher.Phone, teacher.Address, teacher.SubjectID, teacher.Status,
		teacher.ID, teacher.SchoolID, teacher.AcademicPeriodID,
	)

	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// Delete melakukan soft delete pada teacher.
func (r *pgxTeacherRepository) Delete(ctx context.Context, schoolID, academicPeriodID, teacherID string) error {
	if schoolID == "" || academicPeriodID == "" {
		return ErrMissingOperationalContext
	}

	query := `
		UPDATE teachers
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND school_id = $2 AND academic_period_id = $3 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(ctx, query, teacherID, schoolID, academicPeriodID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
