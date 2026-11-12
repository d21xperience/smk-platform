package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AlumniTracking struct {
	ID               string    `json:"id"`
	SchoolID         string    `json:"schoolId"`
	AcademicPeriodID string    `json:"academicPeriodId"`
	Nama             string    `json:"nama"`
	TahunLulus       string    `json:"tahunLulus"`
	Jurusan          string    `json:"jurusan"`
	Status           string    `json:"status"`
	Instansi         string    `json:"instansi"`
	PosisiUsaha      string    `json:"posisiUsaha"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type AlumniTrackingRepository interface {
	List(ctx context.Context, schoolID, academicPeriodID string) ([]AlumniTracking, error)
	Update(ctx context.Context, data *AlumniTracking) error
}

type pgxAlumniTrackingRepository struct {
	db *pgxpool.Pool
}

func NewAlumniTrackingRepository(db *pgxpool.Pool) AlumniTrackingRepository {
	return &pgxAlumniTrackingRepository{db: db}
}

func (r *pgxAlumniTrackingRepository) List(ctx context.Context, schoolID, academicPeriodID string) ([]AlumniTracking, error) {
	query := `
		SELECT id::text, school_id, academic_period_id, nama, tahun_lulus, jurusan, 
		       status, COALESCE(instansi, ''), COALESCE(posisi_usaha, ''), created_at, updated_at
		FROM alumni_tracking
		WHERE school_id = $1 AND academic_period_id = $2
		ORDER BY nama ASC
	`
	rows, err := r.db.Query(ctx, query, schoolID, academicPeriodID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []AlumniTracking
	for rows.Next() {
		var d AlumniTracking
		if err := rows.Scan(
			&d.ID, &d.SchoolID, &d.AcademicPeriodID, &d.Nama, &d.TahunLulus, &d.Jurusan,
			&d.Status, &d.Instansi, &d.PosisiUsaha, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	return list, nil
}

func (r *pgxAlumniTrackingRepository) Update(ctx context.Context, data *AlumniTracking) error {
	query := `
		UPDATE alumni_tracking
		SET status = $1, instansi = $2, posisi_usaha = $3, updated_at = NOW()
		WHERE id = $4 AND school_id = $5 AND academic_period_id = $6
	`
	_, err := r.db.Exec(ctx, query,
		data.Status, data.Instansi, data.PosisiUsaha,
		data.ID, data.SchoolID, data.AcademicPeriodID,
	)
	return err
}
