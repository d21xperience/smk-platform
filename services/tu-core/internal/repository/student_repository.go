package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Student merepresentasikan entitas Student di database.
type Student struct {
	ID               string
	SchoolID         string
	AcademicPeriodID string
	NIS              string
	NISN             string
	Name             string
	Gender           string
	BirthDate        time.Time
	BirthPlace       string
	Address          string
	Phone            string
	Email            string
	ParentName       string
	ParentPhone      string
	ParentEmail      string
	ClassName        string
	Status           string // ACTIVE, GRADUATED, TRANSFERRED, DROPOUT
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}

// StudentFilter untuk query dengan filter.
type StudentFilter struct {
	SchoolID         string
	AcademicPeriodID string
	ClassID          string
	Status           string
	Keyword          string
	Limit            int
	Offset           int
}

// StudentRepository adalah kontrak untuk akses data student.
type StudentRepository interface {
	Create(ctx context.Context, student *Student) error
	GetByID(ctx context.Context, schoolID, academicPeriodID, studentID string) (*Student, error)
	GetByNIS(ctx context.Context, schoolID, academicPeriodID, nis string) (*Student, error)
	List(ctx context.Context, filter StudentFilter) ([]Student, int, error)
	Update(ctx context.Context, student *Student) error
	Delete(ctx context.Context, schoolID, academicPeriodID, studentID string) error
	CountByClass(ctx context.Context, schoolID, academicPeriodID, classID string) (int, error)
}

// pgxStudentRepository mengimplementasikan StudentRepository.
type pgxStudentRepository struct {
	db *pgxpool.Pool
}

// NewStudentRepository membuat instance baru dari StudentRepository.
func NewStudentRepository(db *pgxpool.Pool) StudentRepository {
	return &pgxStudentRepository{db: db}
}

// Create menyimpan student baru ke database.
func (r *pgxStudentRepository) Create(ctx context.Context, student *Student) error {
	if student.SchoolID == "" || student.AcademicPeriodID == "" {
		return ErrMissingOperationalContext
	}

	query := `
		INSERT INTO students (
			id, school_id, academic_period_id, nis, nisn, name, gender, birth_date, birth_place,
			address, phone, email, parent_name, parent_phone, parent_email, class_name, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	_, err := r.db.Exec(ctx, query,
		student.ID, student.SchoolID, student.AcademicPeriodID, student.NIS, student.NISN,
		student.Name, student.Gender, student.BirthDate, student.BirthPlace, student.Address,
		student.Phone, student.Email, student.ParentName, student.ParentPhone, student.ParentEmail,
		student.ClassName, student.Status,
	)

	if err != nil {
		if isDuplicateKeyError(err) {
			return fmt.Errorf("%w: %v", ErrDuplicateEntry, err)
		}
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return nil
}

// GetByID mengambil student berdasarkan ID.
func (r *pgxStudentRepository) GetByID(ctx context.Context, schoolID, academicPeriodID, studentID string) (*Student, error) {
	if schoolID == "" || academicPeriodID == "" {
		return nil, ErrMissingOperationalContext
	}

	query := `
		SELECT
			id, school_id, academic_period_id, nis, nisn, name, gender, birth_date, birth_place,
			address, phone, email, parent_name, parent_phone, parent_email, class_name, status,
			created_at, updated_at, deleted_at
		FROM students
		WHERE id = $1 AND school_id = $2 AND academic_period_id = $3 AND deleted_at IS NULL
	`

	var s Student
	err := r.db.QueryRow(ctx, query, studentID, schoolID, academicPeriodID).Scan(
		&s.ID, &s.SchoolID, &s.AcademicPeriodID, &s.NIS, &s.NISN, &s.Name, &s.Gender,
		&s.BirthDate, &s.BirthPlace, &s.Address, &s.Phone, &s.Email, &s.ParentName,
		&s.ParentPhone, &s.ParentEmail, &s.ClassName, &s.Status, &s.CreatedAt, &s.UpdatedAt,
		&s.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return &s, nil
}

// GetByNIS mengambil student berdasarkan NIS.
func (r *pgxStudentRepository) GetByNIS(ctx context.Context, schoolID, academicPeriodID, nis string) (*Student, error) {
	if schoolID == "" || academicPeriodID == "" {
		return nil, ErrMissingOperationalContext
	}

	query := `
		SELECT
			id, school_id, academic_period_id, nis, nisn, name, gender, birth_date, birth_place,
			address, phone, email, parent_name, parent_phone, parent_email, class_name, status,
			created_at, updated_at, deleted_at
		FROM students
		WHERE nis = $1 AND school_id = $2 AND academic_period_id = $3 AND deleted_at IS NULL
	`

	var s Student
	err := r.db.QueryRow(ctx, query, nis, schoolID, academicPeriodID).Scan(
		&s.ID, &s.SchoolID, &s.AcademicPeriodID, &s.NIS, &s.NISN, &s.Name, &s.Gender,
		&s.BirthDate, &s.BirthPlace, &s.Address, &s.Phone, &s.Email, &s.ParentName,
		&s.ParentPhone, &s.ParentEmail, &s.ClassName, &s.Status, &s.CreatedAt, &s.UpdatedAt,
		&s.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return &s, nil
}

// List mengambil daftar student dengan filter dan pagination.
func (r *pgxStudentRepository) List(ctx context.Context, filter StudentFilter) ([]Student, int, error) {
	if filter.SchoolID == "" || filter.AcademicPeriodID == "" {
		return nil, 0, ErrMissingOperationalContext
	}

	// Build WHERE clause
	whereClause := "WHERE school_id = $1 AND academic_period_id = $2 AND deleted_at IS NULL"
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

	if filter.Keyword != "" {
		whereClause += fmt.Sprintf(" AND (name ILIKE $%d OR nis ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filter.Keyword+"%")
		argIndex++
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM students %s", whereClause)
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
			id, school_id, academic_period_id, nis, nisn, name, gender, birth_date, birth_place,
			address, phone, email, parent_name, parent_phone, parent_email, class_name, status,
			created_at, updated_at, deleted_at
		FROM students
		%s
		ORDER BY name ASC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var s Student
		if err := rows.Scan(
			&s.ID, &s.SchoolID, &s.AcademicPeriodID, &s.NIS, &s.NISN, &s.Name, &s.Gender,
			&s.BirthDate, &s.BirthPlace, &s.Address, &s.Phone, &s.Email, &s.ParentName,
			&s.ParentPhone, &s.ParentEmail, &s.ClassName, &s.Status, &s.CreatedAt, &s.UpdatedAt,
			&s.DeletedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		students = append(students, s)
	}

	return students, total, rows.Err()
}

// Update memperbarui data student.
func (r *pgxStudentRepository) Update(ctx context.Context, student *Student) error {
	if student.SchoolID == "" || student.AcademicPeriodID == "" {
		return ErrMissingOperationalContext
	}

	query := `
		UPDATE students SET
			nis = $1, nisn = $2, name = $3, gender = $4, birth_date = $5, birth_place = $6,
			address = $7, phone = $8, email = $9, parent_name = $10, parent_phone = $11,
			parent_email = $12, class_name = $13, status = $14, updated_at = NOW()
		WHERE id = $15 AND school_id = $16 AND academic_period_id = $17 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(ctx, query,
		student.NIS, student.NISN, student.Name, student.Gender, student.BirthDate,
		student.BirthPlace, student.Address, student.Phone, student.Email, student.ParentName,
		student.ParentPhone, student.ParentEmail, student.ClassName, student.Status,
		student.ID, student.SchoolID, student.AcademicPeriodID,
	)

	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// Delete melakukan soft delete pada student.
func (r *pgxStudentRepository) Delete(ctx context.Context, schoolID, academicPeriodID, studentID string) error {
	if schoolID == "" || academicPeriodID == "" {
		return ErrMissingOperationalContext
	}

	query := `
		UPDATE students
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND school_id = $2 AND academic_period_id = $3 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(ctx, query, studentID, schoolID, academicPeriodID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// CountByClass menghitung jumlah siswa dalam suatu kelas.
func (r *pgxStudentRepository) CountByClass(ctx context.Context, schoolID, academicPeriodID, classID string) (int, error) {
	if schoolID == "" || academicPeriodID == "" {
		return 0, ErrMissingOperationalContext
	}

	query := `
		SELECT COUNT(*)
		FROM students
		WHERE school_id = $1 AND academic_period_id = $2 AND class_id = $3
		  AND status = 'ACTIVE' AND deleted_at IS NULL
	`

	var count int
	err := r.db.QueryRow(ctx, query, schoolID, academicPeriodID, classID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return count, nil
}

// isDuplicateKeyError memeriksa apakah error adalah duplicate key violation (PostgreSQL error code 23505).
func isDuplicateKeyError(err error) bool {
	return err != nil && err.Error() != "" &&
		(contains(err.Error(), "duplicate key") || contains(err.Error(), "23505"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
