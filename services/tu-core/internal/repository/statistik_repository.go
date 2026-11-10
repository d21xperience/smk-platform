package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type KpiData struct {
	TotalPegawai    int     `json:"totalPegawai"`
	TotalHadir      int     `json:"totalHadir"`
	TotalCuti       int     `json:"totalCuti"`
	TotalAlpha      int     `json:"totalAlpha"`
	AvgSkp          float64 `json:"avgSkp"`
	GuruBebanKurang int     `json:"guruBebanKurang"`
}

type KomposisiData struct {
	Label string `json:"label"`
	Value int    `json:"value"`
}

type ChartData struct {
	Komposisi    []KomposisiData `json:"komposisi"`
	DistribusiUnit []KomposisiData `json:"distribusiUnit"`
}

type StatistikData struct {
	Kpi    KpiData   `json:"kpi"`
	Charts ChartData `json:"charts"`
}

type StatistikRepository interface {
	GetDashboardStats(ctx context.Context) (*StatistikData, error)
}

type pgxStatistikRepository struct {
	db *pgxpool.Pool
}

func NewStatistikRepository(db *pgxpool.Pool) StatistikRepository {
	return &pgxStatistikRepository{db: db}
}

func (r *pgxStatistikRepository) GetDashboardStats(ctx context.Context) (*StatistikData, error) {
	var result StatistikData

	// 1. Query KPI
	kpiQuery := `
		SELECT 
			(SELECT COUNT(*) FROM tabel_ptk WHERE soft_delete = 0) AS total_pegawai,
			(SELECT COUNT(*) FROM ptk_attendance WHERE status = 'HADIR' AND date >= DATE_TRUNC('month', CURRENT_DATE)) AS total_hadir,
			(SELECT COUNT(*) FROM ptk_leaves WHERE status = 'disetujui' AND tanggal_mulai >= DATE_TRUNC('month', CURRENT_DATE)) AS total_cuti,
			(SELECT COUNT(*) FROM ptk_attendance WHERE status = 'ALPHA' AND date >= DATE_TRUNC('month', CURRENT_DATE)) AS total_alpha,
			(SELECT COALESCE(AVG(nilai_akhir), 0) FROM penilaian_hasil) AS avg_skp,
			(SELECT COUNT(*) FROM (
				SELECT ptk_id, SUM(jam_per_minggu) as total 
				FROM sk_pembagian_tugas_kelas 
				JOIN sk_pembagian_tugas_detail ON sk_pembagian_tugas_kelas.detail_id = sk_pembagian_tugas_detail.id 
				GROUP BY ptk_id 
				HAVING SUM(jam_per_minggu) < 24
			) sub) AS guru_beban_kurang
	`
	err := r.db.QueryRow(ctx, kpiQuery).Scan(
		&result.Kpi.TotalPegawai, &result.Kpi.TotalHadir, &result.Kpi.TotalCuti,
		&result.Kpi.TotalAlpha, &result.Kpi.AvgSkp, &result.Kpi.GuruBebanKurang,
	)
	if err != nil {
		return nil, err
	}

	// 2. Query Chart Komposisi
	komposisiQuery := `
		SELECT COALESCE(ext.status_pegawai, 'Honorer') AS label, COUNT(*) AS value
		FROM tabel_ptk tp
		LEFT JOIN ptk_buku_induk_ext ext ON tp.ptk_id = ext.ptk_id
		WHERE tp.soft_delete = 0
		GROUP BY ext.status_pegawai
	`
	rows, err := r.db.Query(ctx, komposisiQuery)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var k KomposisiData
			if err := rows.Scan(&k.Label, &k.Value); err == nil {
				result.Charts.Komposisi = append(result.Charts.Komposisi, k)
			}
		}
	}

	// 3. Query Chart Distribusi
	distribusiQuery := `
		SELECT COALESCE(ext.unit_kerja, 'Belum Diatur') AS label, COUNT(*) AS value
		FROM tabel_ptk tp
		LEFT JOIN ptk_buku_induk_ext ext ON tp.ptk_id = ext.ptk_id
		WHERE tp.soft_delete = 0
		GROUP BY ext.unit_kerja
	`
	rows2, err := dbQueryRows(ctx, r.db, distribusiQuery)
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var k KomposisiData
			if err := rows2.Scan(&k.Label, &k.Value); err == nil {
				result.Charts.DistribusiUnit = append(result.Charts.DistribusiUnit, k)
			}
		}
	}

	return &result, nil
}

// Helper kecil agar query konsisten
func dbQueryRows(ctx context.Context, db *pgxpool.Pool, query string) (interface{ Close(); Next() bool; Scan(dest ...interface{}) error }, error) {
	return db.Query(ctx, query)
}
