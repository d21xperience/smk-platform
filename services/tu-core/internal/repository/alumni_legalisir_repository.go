package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AlumniLegalisir struct {
	ID               string    `json:"id"`
	SchoolID         string    `json:"schoolId"`
	AcademicPeriodID string    `json:"academicPeriodId"`
	Nama             string    `json:"nama"`
	TahunLulus       string    `json:"tahunLulus"`
	Jurusan          string    `json:"jurusan"`
	NoIjazah         string    `json:"noIjazah"`
	JumlahLembar     int       `json:"jumlahLembar"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type AlumniLegalisirRepository interface {
	List(ctx context.Context, schoolID, academicPeriodID string) ([]AlumniLegalisir, error)
	Create(ctx context.Context, data *AlumniLegalisir) error
	UpdateStatus(ctx context.Context, id, schoolID, academicPeriodID, status string) error
	Delete(ctx context.Context, id, schoolID, academicPeriodID string) error
}

type pgxAlumniLegalisirRepository struct {
	db *pgxpool.Pool
}

func NewAlumniLegalisirRepository(db *pgxpool.Pool) AlumniLegalisirRepository {
	return &pgxAlumniLegalisirRepository{db: db}
}

func (r *pgxAlumniLegalisirRepository) List(ctx context.Context, schoolID, academicPeriodID string) ([]AlumniLegalisir, error) {
	query := `
		SELECT id::text, school_id, academic_period_id, nama, tahun_lulus, jurusan, no_ijazah, jumlah_lembar, status, created_at, updated_at
		FROM alumni_legalisir
		WHERE school_id = $1 AND academic_period_id = $2
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, schoolID, academicPeriodID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []AlumniLegalisir
	for rows.Next() {
		var d AlumniLegalisir
		if err := rows.Scan(
			&d.ID, &d.SchoolID, &d.AcademicPeriodID, &d.Nama, &d.TahunLulus, &d.Jurusan,
			&d.NoIjazah, &d.JumlahLembar, &d.Status, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	return list, nil
}

func (r *pgxAlumniLegalisirRepository) Create(ctx context.Context, data *AlumniLegalisir) error {
	query := `
		INSERT INTO alumni_legalisir (
			id, school_id, academic_period_id, nama, tahun_lulus, jurusan, no_ijazah, jumlah_lembar, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
	`
	_, err := r.db.Exec(ctx, query,
		data.ID, data.SchoolID, data.AcademicPeriodID, data.Nama, data.TahunLulus,
		data.Jurusan, data.NoIjazah, data.JumlahLembar, data.Status,
	)
	return err
}

func (r *pgxAlumniLegalisirRepository) UpdateStatus(ctx context.Context, id, schoolID, academicPeriodID, status string) error {
	query := `
		UPDATE alumni_legalisir
		SET status = $1, updated_at = NOW()
		WHERE id = $2 AND school_id = $3 AND academic_period_id = $4
	`
	_, err := r.db.Exec(ctx, query, status, id, schoolID, academicPeriodID)
	return err
}

func (r *pgxAlumniLegalisirRepository) Delete(ctx context.Context, id, schoolID, academicPeriodID string) error {
	query := `
		DELETE FROM alumni_legalisir
		WHERE id = $1 AND school_id = $2 AND academic_period_id = $3
	`
	_, err := r.db.Exec(ctx, query, id, schoolID, academicPeriodID)
	return err
}
