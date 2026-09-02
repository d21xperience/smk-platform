package repository

import (
	"context"
	"errors"
	"fmt"

	"sekolah-platform/services/tu-core/internal/domain/mutation/models"
	"sekolah-platform/services/tu-core/internal/domain/mutation/valueobjects"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrMutationNotFound = errors.New("mutation not found")

type MutationFilter struct {
	SchoolID         string
	AcademicPeriodID string
	Status           string
	Type             string
	Limit            int
	Offset           int
}

type MutationRepository interface {
	Create(ctx context.Context, m *models.Mutation) error
	GetByID(ctx context.Context, schoolID, academicPeriodID, mutationID string) (*models.Mutation, error)
	List(ctx context.Context, filter MutationFilter) ([]models.Mutation, int, error)
	Update(ctx context.Context, m *models.Mutation) error
}

type pgxMutationRepository struct {
	db *pgxpool.Pool
}

func NewMutationRepository(db *pgxpool.Pool) MutationRepository {
	return &pgxMutationRepository{db: db}
}

func (r *pgxMutationRepository) Create(ctx context.Context, m *models.Mutation) error {
	query := `
		INSERT INTO mutations (
			id, school_id, academic_period_id, student_id, mutation_type,
			calon_nama, calon_nisn, asal_sekolah, tujuan_sekolah, alasan,
			dokumen_url, status, catatan_verifikasi, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`
	_, err := r.db.Exec(ctx, query,
		m.ID, m.SchoolID, m.AcademicPeriodID, m.StudentID, m.MutationType.String(),
		m.CalonNama, m.CalonNISN, m.AsalSekolah, m.TujuanSekolah, m.Alasan,
		m.DokumenURL, m.Status.String(), m.CatatanVerifikasi, m.CreatedAt, m.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxMutationRepository) GetByID(ctx context.Context, schoolID, academicPeriodID, mutationID string) (*models.Mutation, error) {
	query := `
		SELECT id, school_id, academic_period_id, student_id, mutation_type,
		       calon_nama, calon_nisn, asal_sekolah, tujuan_sekolah, alasan,
		       dokumen_url, status, catatan_verifikasi, created_at, updated_at
		FROM mutations
		WHERE id = $1 AND school_id = $2 AND academic_period_id = $3
	`

	var m models.Mutation
	var mutationTypeStr, statusStr string

	err := r.db.QueryRow(ctx, query, mutationID, schoolID, academicPeriodID).Scan(
		&m.ID, &m.SchoolID, &m.AcademicPeriodID, &m.StudentID, &mutationTypeStr,
		&m.CalonNama, &m.CalonNISN, &m.AsalSekolah, &m.TujuanSekolah, &m.Alasan,
		&m.DokumenURL, &statusStr, &m.CatatanVerifikasi, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMutationNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	m.MutationType, _ = valueobjects.NewMutationType(mutationTypeStr)
	m.Status, _ = valueobjects.NewMutationStatus(statusStr)

	return &m, nil
}

func (r *pgxMutationRepository) List(ctx context.Context, filter MutationFilter) ([]models.Mutation, int, error) {
	whereClause := "WHERE school_id = $1 AND academic_period_id = $2"
	args := []interface{}{filter.SchoolID, filter.AcademicPeriodID}
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

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM mutations %s", whereClause)
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
		SELECT id, school_id, academic_period_id, student_id, mutation_type,
		       calon_nama, calon_nisn, asal_sekolah, tujuan_sekolah, alasan,
		       dokumen_url, status, catatan_verifikasi, created_at, updated_at
		FROM mutations %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var mutations []models.Mutation
	for rows.Next() {
		var m models.Mutation
		var mutationTypeStr, statusStr string

		if err := rows.Scan(
			&m.ID, &m.SchoolID, &m.AcademicPeriodID, &m.StudentID, &mutationTypeStr,
			&m.CalonNama, &m.CalonNISN, &m.AsalSekolah, &m.TujuanSekolah, &m.Alasan,
			&m.DokumenURL, &statusStr, &m.CatatanVerifikasi, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}

		m.MutationType, _ = valueobjects.NewMutationType(mutationTypeStr)
		m.Status, _ = valueobjects.NewMutationStatus(statusStr)

		mutations = append(mutations, m)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return mutations, total, nil
}

func (r *pgxMutationRepository) Update(ctx context.Context, m *models.Mutation) error {
	query := `
		UPDATE mutations
		SET student_id = $1, mutation_type = $2, calon_nama = $3, calon_nisn = $4,
		    asal_sekolah = $5, tujuan_sekolah = $6, alasan = $7, dokumen_url = $8,
		    status = $9, catatan_verifikasi = $10, updated_at = $11
		WHERE id = $12 AND school_id = $13 AND academic_period_id = $14
	`
	result, err := r.db.Exec(ctx, query,
		m.StudentID, m.MutationType.String(), m.CalonNama, m.CalonNISN,
		m.AsalSekolah, m.TujuanSekolah, m.Alasan, m.DokumenURL,
		m.Status.String(), m.CatatanVerifikasi, m.UpdatedAt,
		m.ID, m.SchoolID, m.AcademicPeriodID,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrMutationNotFound
	}
	return nil
}
