package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Enrollment struct {
	ID        string
	StudentID string
	SchoolID  string
	PeriodID  string
	ClassID   string
	Status    string
	Reason    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type EnrollmentFilter struct {
	SchoolID  string
	PeriodID  string
	StudentID string
	ClassID   string
	Status    string
	Limit     int
	Offset    int
}

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

type pgxEnrollmentRepository struct {
	db *pgxpool.Pool
}

func NewEnrollmentRepository(db *pgxpool.Pool) EnrollmentRepository {
	return &pgxEnrollmentRepository{db: db}
}

func (r *pgxEnrollmentRepository) Create(ctx context.Context, enrollment *Enrollment) error {
	if enrollment.SchoolID == "" || enrollment.PeriodID == "" {
		return ErrMissingOperationalContext
	}

	query := `
	INSERT INTO tabel_anggotakelas (
		anggota_rombel_id, peserta_didik_id, rombongan_belajar_id, semester_id
	) VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(ctx, query,
		enrollment.ID, enrollment.StudentID, enrollment.ClassID, enrollment.PeriodID,
	)

	if err != nil {
		if isDuplicateKeyError(err) {
			return fmt.Errorf("%w: %v", ErrDuplicateEntry, err)
		}
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxEnrollmentRepository) GetByID(ctx context.Context, schoolID, periodID, enrollmentID string) (*Enrollment, error) {
	if schoolID == "" || periodID == "" {
		return nil, ErrMissingOperationalContext
	}

	// Karena tabel_anggotakelas tidak punya school_id, kita JOIN ke tabel_kelas
	query := `
	SELECT
		tak.anggota_rombel_id, tak.peserta_didik_id, tk.sekolah_id, tak.semester_id, 
		tak.rombongan_belajar_id
	FROM tabel_anggotakelas tak
	JOIN tabel_kelas tk ON tak.rombongan_belajar_id = tk.rombongan_belajar_id
	WHERE tak.anggota_rombel_id = $1 AND tak.semester_id = $2 AND tk.sekolah_id = $3
	`

	var e Enrollment
	err := r.db.QueryRow(ctx, query, enrollmentID, periodID, schoolID).Scan(
		&e.ID, &e.StudentID, &e.SchoolID, &e.PeriodID, &e.ClassID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	e.Status = "active"
	return &e, nil
}

func (r *pgxEnrollmentRepository) List(ctx context.Context, filter EnrollmentFilter) ([]Enrollment, int, error) {
	if filter.SchoolID == "" || filter.PeriodID == "" {
		return nil, 0, ErrMissingOperationalContext
	}

	whereClause := "WHERE tak.semester_id = $1 AND tk.sekolah_id = $2"
	args := []interface{}{filter.PeriodID, filter.SchoolID}
	argIndex := 3

	if filter.StudentID != "" {
		whereClause += fmt.Sprintf(" AND tak.peserta_didik_id = $%d", argIndex)
		args = append(args, filter.StudentID)
		argIndex++
	}
	if filter.ClassID != "" {
		whereClause += fmt.Sprintf(" AND tak.rombongan_belajar_id = $%d", argIndex)
		args = append(args, filter.ClassID)
		argIndex++
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM tabel_anggotakelas tak
		JOIN tabel_kelas tk ON tak.rombongan_belajar_id = tk.rombongan_belajar_id
		%s
	`, whereClause)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if filter.Limit <= 0 { filter.Limit = 50 }
	if filter.Offset < 0 { filter.Offset = 0 }

	dataQuery := fmt.Sprintf(`
		SELECT
			tak.anggota_rombel_id, tak.peserta_didik_id, tk.sekolah_id, tak.semester_id, 
			tak.rombongan_belajar_id
		FROM tabel_anggotakelas tak
		JOIN tabel_kelas tk ON tak.rombongan_belajar_id = tk.rombongan_belajar_id
		%s
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
			&e.ID, &e.StudentID, &e.SchoolID, &e.PeriodID, &e.ClassID,
		); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		e.Status = "active"
		enrollments = append(enrollments, e)
	}
	return enrollments, total, rows.Err()
}

func (r *pgxEnrollmentRepository) Update(ctx context.Context, enrollment *Enrollment) error {
	if enrollment.SchoolID == "" || enrollment.PeriodID == "" {
		return ErrMissingOperationalContext
	}

	query := `
	UPDATE tabel_anggotakelas SET
		rombongan_belajar_id = $1
	WHERE anggota_rombel_id = $2 AND semester_id = $3
	`

	result, err := r.db.Exec(ctx, query,
		enrollment.ClassID, enrollment.ID, enrollment.PeriodID,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *pgxEnrollmentRepository) Delete(ctx context.Context, schoolID, periodID, enrollmentID string) error {
	if schoolID == "" || periodID == "" {
		return ErrMissingOperationalContext
	}

	query := `DELETE FROM tabel_anggotakelas WHERE anggota_rombel_id = $1 AND semester_id = $2`
	result, err := r.db.Exec(ctx, query, enrollmentID, periodID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *pgxEnrollmentRepository) GetActiveByStudent(ctx context.Context, schoolID, periodID, studentID string) (*Enrollment, error) {
	if schoolID == "" || periodID == "" {
		return nil, ErrMissingOperationalContext
	}

	query := `
	SELECT
		tak.anggota_rombel_id, tak.peserta_didik_id, tk.sekolah_id, tak.semester_id, 
		tak.rombongan_belajar_id
	FROM tabel_anggotakelas tak
	JOIN tabel_kelas tk ON tak.rombongan_belajar_id = tk.rombongan_belajar_id
	WHERE tak.peserta_didik_id = $1 AND tak.semester_id = $2 AND tk.sekolah_id = $3
	LIMIT 1
	`

	var e Enrollment
	err := r.db.QueryRow(ctx, query, studentID, periodID, schoolID).Scan(
		&e.ID, &e.StudentID, &e.SchoolID, &e.PeriodID, &e.ClassID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	e.Status = "active"
	return &e, nil
}

func (r *pgxEnrollmentRepository) GetHistoryByStudent(ctx context.Context, schoolID, studentID string) ([]Enrollment, error) {
	if schoolID == "" {
		return nil, ErrMissingOperationalContext
	}

	query := `
	SELECT
		tak.anggota_rombel_id, tak.peserta_didik_id, tk.sekolah_id, tak.semester_id, 
		tak.rombongan_belajar_id
	FROM tabel_anggotakelas tak
	JOIN tabel_kelas tk ON tak.rombongan_belajar_id = tk.rombongan_belajar_id
	WHERE tak.peserta_didik_id = $1 AND tk.sekolah_id = $2
	ORDER BY tak.semester_id DESC
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
			&e.ID, &e.StudentID, &e.SchoolID, &e.PeriodID, &e.ClassID,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		e.Status = "active"
		enrollments = append(enrollments, e)
	}
	return enrollments, rows.Err()
}

func (r *pgxEnrollmentRepository) GetMembersByClass(ctx context.Context, schoolID, periodID, classID string) ([]Enrollment, error) {
	if schoolID == "" || periodID == "" {
		return nil, ErrMissingOperationalContext
	}

	query := `
	SELECT
		tak.anggota_rombel_id, tak.peserta_didik_id, tk.sekolah_id, tak.semester_id, 
		tak.rombongan_belajar_id
	FROM tabel_anggotakelas tak
	JOIN tabel_kelas tk ON tak.rombongan_belajar_id = tk.rombongan_belajar_id
	WHERE tak.rombongan_belajar_id = $1 AND tak.semester_id = $2 AND tk.sekolah_id = $3
	`

	rows, err := r.db.Query(ctx, query, classID, periodID, schoolID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var enrollments []Enrollment
	for rows.Next() {
		var e Enrollment
		if err := rows.Scan(
			&e.ID, &e.StudentID, &e.SchoolID, &e.PeriodID, &e.ClassID,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		e.Status = "active"
		enrollments = append(enrollments, e)
	}
	return enrollments, rows.Err()
}

func (r *pgxEnrollmentRepository) CountActiveByClass(ctx context.Context, schoolID, periodID, classID string) (int, error) {
	if schoolID == "" || periodID == "" {
		return 0, ErrMissingOperationalContext
	}

	query := `
	SELECT COUNT(*)
	FROM tabel_anggotakelas tak
	JOIN tabel_kelas tk ON tak.rombongan_belajar_id = tk.rombongan_belajar_id
	WHERE tak.rombongan_belajar_id = $1 AND tak.semester_id = $2 AND tk.sekolah_id = $3
	`

	var count int
	err := r.db.QueryRow(ctx, query, classID, periodID, schoolID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return count, nil
}
