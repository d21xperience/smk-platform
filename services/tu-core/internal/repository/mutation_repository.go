package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrMutationNotFound = errors.New("mutation not found")

type StudentMutation struct {
	ID               string
	SchoolID         string
	AcademicPeriodID string
	StudentID        *string
	MutationType     string
	NIS              string
	NISN             string
	Nama             string
	JenisKelamin     string
	TempatLahir      string
	TanggalLahir     *time.Time
	AsalSekolah      string
	TujuanSekolah    string
	Alasan           string
	DokumenURL       string
	Status           string
	Catatan          string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type MutationFilter struct {
	SchoolID         string
	AcademicPeriodID string
	Status           string
	Type             string
	Limit            int
	Offset           int
}

type MutationRepository interface {
	Create(ctx context.Context, m *StudentMutation) error
	GetByID(ctx context.Context, schoolID, academicPeriodID, mutationID string) (*StudentMutation, error)
	List(ctx context.Context, filter MutationFilter) ([]StudentMutation, int, error)
	UpdateStatus(ctx context.Context, mutationID, status, catatan string) error
	GetHistoryByStudent(ctx context.Context, schoolID, studentID string) ([]StudentMutation, error)
}

type pgxMutationRepository struct {
	db *pgxpool.Pool
}

func NewMutationRepository(db *pgxpool.Pool) MutationRepository {
	return &pgxMutationRepository{db: db}
}

func (r *pgxMutationRepository) Create(ctx context.Context, m *StudentMutation) error {
	query := `
		INSERT INTO student_mutations (
			id, school_id, academic_period_id, student_id, mutation_type,
			nis, nisn, nama, jenis_kelamin, tempat_lahir, tanggal_lahir,
			asal_sekolah, tujuan_sekolah, alasan, dokumen_url, status, catatan, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`
	_, err := r.db.Exec(ctx, query,
		m.ID, m.SchoolID, m.AcademicPeriodID, m.StudentID, m.MutationType,
		m.NIS, m.NISN, m.Nama, m.JenisKelamin, m.TempatLahir, m.TanggalLahir,
		m.AsalSekolah, m.TujuanSekolah, m.Alasan, m.DokumenURL, m.Status, m.Catatan, m.CreatedAt, m.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxMutationRepository) GetByID(ctx context.Context, schoolID, academicPeriodID, mutationID string) (*StudentMutation, error) {
	query := `
		SELECT id, school_id, academic_period_id, student_id, mutation_type,
		       COALESCE(nis, ''), COALESCE(nisn, ''), nama, COALESCE(jenis_kelamin, ''), 
		       COALESCE(tempat_lahir, ''), tanggal_lahir, COALESCE(asal_sekolah, ''), 
		       COALESCE(tujuan_sekolah, ''), COALESCE(alasan, ''), COALESCE(dokumen_url, ''), 
		       status, COALESCE(catatan, ''), created_at, updated_at
		FROM student_mutations
		WHERE id = $1 AND school_id = $2 AND academic_period_id = $3
	`
	var m StudentMutation
	err := r.db.QueryRow(ctx, query, mutationID, schoolID, academicPeriodID).Scan(
		&m.ID, &m.SchoolID, &m.AcademicPeriodID, &m.StudentID, &m.MutationType,
		&m.NIS, &m.NISN, &m.Nama, &m.JenisKelamin, &m.TempatLahir, &m.TanggalLahir,
		&m.AsalSekolah, &m.TujuanSekolah, &m.Alasan, &m.DokumenURL, &m.Status, &m.Catatan, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMutationNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return &m, nil
}

func (r *pgxMutationRepository) List(ctx context.Context, filter MutationFilter) ([]StudentMutation, int, error) {
	whereClause := "WHERE school_id = $1 AND academic_period_id = $2"
	args := []any{filter.SchoolID, filter.AcademicPeriodID}
	argIndex := 3

	if filter.Status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}
	if filter.Type != "" {
		whereClause += fmt.Sprintf(" AND mutation_type = $%d", argIndex)
		args = append(args, filter.Type)
		argIndex++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM student_mutations %s", whereClause)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if filter.Limit <= 0 { filter.Limit = 20 }
	if filter.Offset < 0 { filter.Offset = 0 }

	dataQuery := fmt.Sprintf(`
		SELECT id, school_id, academic_period_id, student_id, mutation_type,
		       COALESCE(nis, ''), COALESCE(nisn, ''), nama, COALESCE(jenis_kelamin, ''), 
		       COALESCE(tempat_lahir, ''), tanggal_lahir, COALESCE(asal_sekolah, ''), 
		       COALESCE(tujuan_sekolah, ''), COALESCE(alasan, ''), COALESCE(dokumen_url, ''), 
		       status, COALESCE(catatan, ''), created_at, updated_at
		FROM student_mutations %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var mutations []StudentMutation
	for rows.Next() {
		var m StudentMutation
		if err := rows.Scan(
			&m.ID, &m.SchoolID, &m.AcademicPeriodID, &m.StudentID, &m.MutationType,
			&m.NIS, &m.NISN, &m.Nama, &m.JenisKelamin, &m.TempatLahir, &m.TanggalLahir,
			&m.AsalSekolah, &m.TujuanSekolah, &m.Alasan, &m.DokumenURL, &m.Status, &m.Catatan, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		mutations = append(mutations, m)
	}
	return mutations, total, rows.Err()
}

func (r *pgxMutationRepository) UpdateStatus(ctx context.Context, mutationID, status, catatan string) error {
	query := `UPDATE student_mutations SET status = $1, catatan = $2, updated_at = NOW() WHERE id = $3`
	result, err := r.db.Exec(ctx, query, status, catatan, mutationID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrMutationNotFound
	}
	return nil
}

func (r *pgxMutationRepository) GetHistoryByStudent(ctx context.Context, schoolID, studentID string) ([]StudentMutation, error) {
	query := `
		SELECT id, school_id, academic_period_id, student_id, mutation_type,
		       COALESCE(nis, ''), COALESCE(nisn, ''), nama, COALESCE(jenis_kelamin, ''), 
		       COALESCE(tempat_lahir, ''), tanggal_lahir, COALESCE(asal_sekolah, ''), 
		       COALESCE(tujuan_sekolah, ''), COALESCE(alasan, ''), COALESCE(dokumen_url, ''), 
		       status, COALESCE(catatan, ''), created_at, updated_at
		FROM student_mutations
		WHERE school_id = $1 AND student_id = $2
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, schoolID, studentID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var mutations []StudentMutation
	for rows.Next() {
		var m StudentMutation
		if err := rows.Scan(
			&m.ID, &m.SchoolID, &m.AcademicPeriodID, &m.StudentID, &m.MutationType,
			&m.NIS, &m.NISN, &m.Nama, &m.JenisKelamin, &m.TempatLahir, &m.TanggalLahir,
			&m.AsalSekolah, &m.TujuanSekolah, &m.Alasan, &m.DokumenURL, &m.Status, &m.Catatan, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		mutations = append(mutations, m)
	}
	return mutations, rows.Err()
}
