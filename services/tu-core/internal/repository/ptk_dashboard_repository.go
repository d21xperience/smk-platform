package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardStats struct {
	TotalTeachers    int     `json:"totalTeachers"`
	PendingLeaves    int     `json:"pendingLeaves"`
	AverageSkp       float64 `json:"averageSkp"`
	ExpiringArchives int     `json:"expiringArchives"`
}

type PtkDashboardRepository interface {
	GetDashboardStats(ctx context.Context) (*DashboardStats, error)
}

type pgxPtkDashboardRepository struct {
	db *pgxpool.Pool
}

func NewPtkDashboardRepository(db *pgxpool.Pool) PtkDashboardRepository {
	return &pgxPtkDashboardRepository{db: db}
}

func (r *pgxPtkDashboardRepository) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	query := `
		SELECT
			(SELECT COUNT(*) FROM tabel_ptk WHERE soft_delete = 0) AS total_teachers,
			(SELECT COUNT(*) FROM ptk_leaves WHERE status = 'pending') AS pending_leaves,
			(SELECT COALESCE(AVG(nilai_akhir), 0) FROM penilaian_hasil) AS average_skp,
			(SELECT COUNT(*) FROM ptk_archives WHERE tanggal_kadaluarsa IS NOT NULL AND tanggal_kadaluarsa <= NOW() + INTERVAL '30 days') AS expiring_archives
	`
	var stats DashboardStats
	err := r.db.QueryRow(ctx, query).Scan(
		&stats.TotalTeachers,
		&stats.PendingLeaves,
		&stats.AverageSkp,
		&stats.ExpiringArchives,
	)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}
