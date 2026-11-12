package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AcademicPeriod direpresentasikan sesuai tabel `semesters`.
type AcademicPeriod struct {
	ID           string // id (contoh: "20261")
	Name         string // name (contoh: "Ganjil 2026/2027")
	AcademicYear string // academic_year (contoh: "2026/2027")
	Semester     int    // semester (contoh: 1)
	IsActive     bool   // is_active
	IsPsbActive  bool   // is_psb_active
}

type AcademicRepository interface {
	GetActivePeriod(ctx context.Context, schoolID string) (*AcademicPeriod, error)
	GetPeriodByID(ctx context.Context, schoolID, periodID string) (*AcademicPeriod, error)
	GetClassByID(ctx context.Context, schoolID, academicPeriodID, classID string) (*Class, error)
	ListClasses(ctx context.Context, schoolID, academicPeriodID string) ([]Class, error)
}

type pgxAcademicRepository struct {
	db *pgxpool.Pool
}

func NewAcademicRepository(db *pgxpool.Pool) AcademicRepository {
	return &pgxAcademicRepository{db: db}
}

func (r *pgxAcademicRepository) GetActivePeriod(ctx context.Context, schoolID string) (*AcademicPeriod, error) {
	if schoolID == "" {
		return nil, ErrMissingOperationalContext
	}
	// Tabel semesters global (tidak punya school_id)
	query := `SELECT id, name, academic_year, semester, is_active, COALESCE(is_psb_active, FALSE) FROM semesters WHERE is_active = TRUE LIMIT 1`
	var p AcademicPeriod
	err := r.db.QueryRow(ctx, query).Scan(&p.ID, &p.Name, &p.AcademicYear, &p.Semester, &p.IsActive, &p.IsPsbActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return &p, nil
}

func (r *pgxAcademicRepository) GetPeriodByID(ctx context.Context, schoolID, periodID string) (*AcademicPeriod, error) {
	if schoolID == "" || periodID == "" {
		return nil, ErrMissingOperationalContext
	}
	query := `SELECT id, name, academic_year, semester, is_active, COALESCE(is_psb_active, FALSE) FROM semesters WHERE id = $1`
	var p AcademicPeriod
	err := r.db.QueryRow(ctx, query, periodID).Scan(&p.ID, &p.Name, &p.AcademicYear, &p.Semester, &p.IsActive, &p.IsPsbActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return &p, nil
}

func (r *pgxAcademicRepository) GetClassByID(ctx context.Context, schoolID, academicPeriodID, classID string) (*Class, error) {
	if schoolID == "" || academicPeriodID == "" {
		return nil, ErrMissingOperationalContext
	}

	query := `
	SELECT ` + classSelectColumns + `
	FROM tabel_kelas
	WHERE rombongan_belajar_id = $1 AND sekolah_id = $2 AND semester_id = $3
	`
	var c Class
	if err := scanClass(r.db.QueryRow(ctx, query, classID, schoolID, academicPeriodID), &c); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return &c, nil
}

func (r *pgxAcademicRepository) ListClasses(ctx context.Context, schoolID, academicPeriodID string) ([]Class, error) {
	if schoolID == "" || academicPeriodID == "" {
		return nil, ErrMissingOperationalContext
	}

	query := `
	SELECT ` + classSelectColumns + `
	FROM tabel_kelas
	WHERE sekolah_id = $1 AND semester_id = $2
	ORDER BY tingkat_pendidikan_id ASC, nm_kelas ASC
	`
	rows, err := r.db.Query(ctx, query, schoolID, academicPeriodID)
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
