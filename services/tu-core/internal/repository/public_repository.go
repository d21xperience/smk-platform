package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PublicStats struct {
	TotalSiswa  int
	TotalGuru   int
	TotalRombel int
}

type PublicRepository interface {
	GetDashboardStats(ctx context.Context) (*PublicStats, error)
}

type publicRepository struct {
	db *pgxpool.Pool
}

func NewPublicRepository(db *pgxpool.Pool) PublicRepository {
	return &publicRepository{db: db}
}

// GetDashboardStats mengembalikan statistik dashboard publik.
// Menembak skema DAPODIK: tabel_siswa, tabel_ptk, tabel_kelas.
func (r *publicRepository) GetDashboardStats(ctx context.Context) (*PublicStats, error) {
	stats := &PublicStats{}

	// Siswa aktif (DAPODIK tidak punya flag aktif di tabel_siswa; hitung semua).
	if err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM tabel_siswa`,
	).Scan(&stats.TotalSiswa); err != nil {
		return nil, fmt.Errorf("gagal hitung total siswa: %w", err)
	}

	// Guru = PTK dengan soft_delete = 0 (sesuaikan jenis_ptk_id jika perlu filter "guru").
	// jenis_ptk_id = 1 biasanya "Guru" di DAPODIK; kalau ragu, hitung semua PTK aktif.
	if err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM tabel_ptk WHERE soft_delete = 0`,
	).Scan(&stats.TotalGuru); err != nil {
		return nil, fmt.Errorf("gagal hitung total guru: %w", err)
	}

	// Rombel = jumlah kelas di tabel_kelas (lintas semester — kalau perlu
	// filter semester, tambahkan parameter & WHERE).
	if err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM tabel_kelas`,
	).Scan(&stats.TotalRombel); err != nil {
		return nil, fmt.Errorf("gagal hitung total rombel: %w", err)
	}

	return stats, nil
}
