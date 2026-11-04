package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultJenisRombel = 1
	defaultKurikulumID = 1
)

type Class struct {
	ID               string
	SchoolID         string
	AcademicPeriodID string
	Name             string
	Grade            int
	WaliTeacherID    string // ptk_id (hasil scan)
	Level            int    // VIRTUAL (diisi = Grade saat scan)
	Major            string
	HomeroomID       string // ptk_id (diisi Service saat write)
	Capacity         int
	Status           string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}

type ClassFilter struct {
	SchoolID         string
	AcademicPeriodID string
	Grade            int
	Status           string
	Keyword          string
	Limit            int
	Offset           int
}

type ClassRepository interface {
	Create(ctx context.Context, class *Class) error
	GetByID(ctx context.Context, schoolID, periodID, classID string) (*Class, error)
	List(ctx context.Context, filter ClassFilter) ([]Class, int, error)
	Update(ctx context.Context, class *Class) error
	Delete(ctx context.Context, schoolID, periodID, classID string) error
	CountActiveMembers(ctx context.Context, schoolID, periodID, classID string) (int, error)
	GetByGradeAndPeriod(ctx context.Context, schoolID, periodID string, grade int) ([]Class, error)
}

type pgxClassRepository struct {
	db *pgxpool.Pool
}

func NewClassRepository(db *pgxpool.Pool) ClassRepository {
	return &pgxClassRepository{db: db}
}

func (c *Class) homeroomID() string {
	if c.WaliTeacherID != "" {
		return c.WaliTeacherID
	}
	return c.HomeroomID
}

func nullableUUID(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func nullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func nullableInt(n int) interface{} {
	if n <= 0 {
		return nil
	}
	return n
}

const classSelectColumns = `
	rombongan_belajar_id,
	sekolah_id,
	semester_id,
	nm_kelas,
	COALESCE(tingkat_pendidikan_id, 0),
	COALESCE(ptk_id::text, ''),
	COALESCE(nama_jurusan_sp, ''),
	COALESCE(capacity, 0),
	COALESCE(created_at, '1970-01-01 00:00:00+00'::timestamptz),
	COALESCE(updated_at, '1970-01-01 00:00:00+00'::timestamptz)
`

func scanClass(row pgx.Row, c *Class) error {
	if err := row.Scan(
		&c.ID,
		&c.SchoolID,
		&c.AcademicPeriodID,
		&c.Name,
		&c.Grade,
		&c.WaliTeacherID,
		&c.Major,
		&c.Capacity,
		&c.CreatedAt,
		&c.UpdatedAt,
	); err != nil {
		return err
	}
	c.HomeroomID = c.WaliTeacherID
	c.Status = "active"
	c.Level = c.Grade // FIX: Level diisi dari Grade agar JSON "level" tidak kosong ("X", "XI", dll via service)
	return nil
}

func (r *pgxClassRepository) Create(ctx context.Context, class *Class) error {
	if class.SchoolID == "" || class.AcademicPeriodID == "" {
		return ErrMissingOperationalContext
	}
	if class.ID == "" || class.Name == "" {
		return fmt.Errorf("%w: rombongan_belajar_id dan nm_kelas wajib diisi", ErrInvalidInput)
	}

	query := `
	INSERT INTO tabel_kelas (
		rombongan_belajar_id, sekolah_id, semester_id, nm_kelas,
		tingkat_pendidikan_id, ptk_id, nama_jurusan_sp, jenis_rombel,
		kurikulum_id, capacity, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
	`

	_, err := r.db.Exec(ctx, query,
		class.ID, class.SchoolID, class.AcademicPeriodID, class.Name,
		class.Grade, nullableUUID(class.homeroomID()), nullableString(class.Major),
		defaultJenisRombel, defaultKurikulumID, nullableInt(class.Capacity),
	)
	if err != nil {
		if isDuplicateKeyError(err) {
			return fmt.Errorf("%w: %v", ErrDuplicateEntry, err)
		}
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxClassRepository) GetByID(ctx context.Context, schoolID, periodID, classID string) (*Class, error) {
	if schoolID == "" || periodID == "" {
		return nil, ErrMissingOperationalContext
	}

	query := `
	SELECT ` + classSelectColumns + `
	FROM tabel_kelas
	WHERE rombongan_belajar_id = $1 AND sekolah_id = $2 AND semester_id = $3
	`

	var c Class
	if err := scanClass(r.db.QueryRow(ctx, query, classID, schoolID, periodID), &c); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return &c, nil
}

func (r *pgxClassRepository) List(ctx context.Context, filter ClassFilter) ([]Class, int, error) {
	if filter.SchoolID == "" || filter.AcademicPeriodID == "" {
		return nil, 0, ErrMissingOperationalContext
	}

	whereClause := "WHERE sekolah_id = $1 AND semester_id = $2"
	args := []interface{}{filter.SchoolID, filter.AcademicPeriodID}
	argIndex := 3

	if filter.Grade > 0 {
		whereClause += fmt.Sprintf(" AND tingkat_pendidikan_id = $%d", argIndex)
		args = append(args, filter.Grade)
		argIndex++
	}
	if filter.Keyword != "" {
		whereClause += fmt.Sprintf(" AND (nm_kelas ILIKE $%d OR nama_jurusan_sp ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filter.Keyword+"%")
		argIndex++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tabel_kelas %s", whereClause)
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
	SELECT `+classSelectColumns+`
	FROM tabel_kelas
	%s
	ORDER BY tingkat_pendidikan_id ASC NULLS LAST, nm_kelas ASC
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
		if err := scanClass(rows, &c); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		classes = append(classes, c)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return classes, total, nil
}

func (r *pgxClassRepository) Update(ctx context.Context, class *Class) error {
	if class.SchoolID == "" || class.AcademicPeriodID == "" {
		return ErrMissingOperationalContext
	}

	query := `
	UPDATE tabel_kelas SET
		nm_kelas = $1, tingkat_pendidikan_id = $2, ptk_id = $3,
		nama_jurusan_sp = $4, capacity = $5, updated_at = NOW()
	WHERE rombongan_belajar_id = $6 AND sekolah_id = $7 AND semester_id = $8
	`

	result, err := r.db.Exec(ctx, query,
		class.Name, class.Grade, nullableUUID(class.homeroomID()),
		nullableString(class.Major), nullableInt(class.Capacity),
		class.ID, class.SchoolID, class.AcademicPeriodID,
	)
	if err != nil {
		if isDuplicateKeyError(err) {
			return fmt.Errorf("%w: %v", ErrDuplicateEntry, err)
		}
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *pgxClassRepository) Delete(ctx context.Context, schoolID, periodID, classID string) error {
	if schoolID == "" || periodID == "" {
		return ErrMissingOperationalContext
	}

	query := `DELETE FROM tabel_kelas WHERE rombongan_belajar_id = $1 AND sekolah_id = $2 AND semester_id = $3`
	result, err := r.db.Exec(ctx, query, classID, schoolID, periodID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *pgxClassRepository) CountActiveMembers(ctx context.Context, schoolID, periodID, classID string) (int, error) {
	if schoolID == "" || periodID == "" {
		return 0, ErrMissingOperationalContext
	}

	query := `SELECT COUNT(*) FROM tabel_anggotakelas WHERE rombongan_belajar_id = $1 AND semester_id = $2`
	var count int
	if err := r.db.QueryRow(ctx, query, classID, periodID).Scan(&count); err != nil {
		return 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return count, nil
}

func (r *pgxClassRepository) GetByGradeAndPeriod(ctx context.Context, schoolID, periodID string, grade int) ([]Class, error) {
	if schoolID == "" || periodID == "" {
		return nil, ErrMissingOperationalContext
	}

	query := `
	SELECT ` + classSelectColumns + `
	FROM tabel_kelas
	WHERE sekolah_id = $1 AND semester_id = $2 AND tingkat_pendidikan_id = $3
	ORDER BY nm_kelas ASC
	`
	rows, err := r.db.Query(ctx, query, schoolID, periodID, grade)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var classes []Class
	for rows.Next() {
		var c Class
		if err := scanClass(rows, &c); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		classes = append(classes, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return classes, nil
}
