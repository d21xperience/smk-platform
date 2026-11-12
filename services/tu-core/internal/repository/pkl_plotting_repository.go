package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PklPlotting struct {
	ID               string    `json:"id"`
	SchoolID         string    `json:"schoolId"`
	AcademicPeriodID string    `json:"academicPeriodId"`
	Nama             string    `json:"nama"`
	Kelas            string    `json:"kelas"`
	Jurusan          string    `json:"jurusan"`
	Perusahaan       string    `json:"perusahaan"`
	Kota             string    `json:"kota"`
	Pembimbing       string    `json:"pembimbing"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type PklPlottingRepository interface {
	List(ctx context.Context, schoolID, academicPeriodID string) ([]PklPlotting, error)
	Create(ctx context.Context, data *PklPlotting) error
	Update(ctx context.Context, data *PklPlotting) error
	Delete(ctx context.Context, id, schoolID, academicPeriodID string) error
}

type pgxPklPlottingRepository struct {
	db *pgxpool.Pool
}

func NewPklPlottingRepository(db *pgxpool.Pool) PklPlottingRepository {
	return &pgxPklPlottingRepository{db: db}
}

func (r *pgxPklPlottingRepository) List(ctx context.Context, schoolID, academicPeriodID string) ([]PklPlotting, error) {
	query := `
		SELECT id::text, school_id, academic_period_id, nama, kelas, jurusan, 
		       perusahaan, COALESCE(kota, ''), COALESCE(pembimbing, ''), created_at, updated_at
		FROM pkl_plotting
		WHERE school_id = $1 AND academic_period_id = $2
		ORDER BY nama ASC
	`
	rows, err := r.db.Query(ctx, query, schoolID, academicPeriodID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []PklPlotting
	for rows.Next() {
		var d PklPlotting
		if err := rows.Scan(
			&d.ID, &d.SchoolID, &d.AcademicPeriodID, &d.Nama, &d.Kelas, &d.Jurusan,
			&d.Perusahaan, &d.Kota, &d.Pembimbing, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	return list, nil
}

func (r *pgxPklPlottingRepository) Create(ctx context.Context, data *PklPlotting) error {
	query := `
		INSERT INTO pkl_plotting (
			id, school_id, academic_period_id, nama, kelas, jurusan, perusahaan, kota, pembimbing, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
	`
	_, err := r.db.Exec(ctx, query,
		data.ID, data.SchoolID, data.AcademicPeriodID, data.Nama, data.Kelas,
		data.Jurusan, data.Perusahaan, data.Kota, data.Pembimbing,
	)
	return err
}

func (r *pgxPklPlottingRepository) Update(ctx context.Context, data *PklPlotting) error {
	query := `
		UPDATE pkl_plotting
		SET perusahaan = $1, kota = $2, pembimbing = $3, updated_at = NOW()
		WHERE id = $4 AND school_id = $5 AND academic_period_id = $6
	`
	_, err := r.db.Exec(ctx, query,
		data.Perusahaan, data.Kota, data.Pembimbing,
		data.ID, data.SchoolID, data.AcademicPeriodID,
	)
	return err
}

func (r *pgxPklPlottingRepository) Delete(ctx context.Context, id, schoolID, academicPeriodID string) error {
	query := `
		DELETE FROM pkl_plotting
		WHERE id = $1 AND school_id = $2 AND academic_period_id = $3
	`
	_, err := r.db.Exec(ctx, query, id, schoolID, academicPeriodID)
	return err
}
