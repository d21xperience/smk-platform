package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StatisticKPIs struct {
	TotalPegawai      int     `json:"totalPegawai"`
	TotalGuruTetap    int     `json:"totalGuruTetap"`
	TotalHonorer      int     `json:"totalHonorer"`
	TotalTendik       int     `json:"totalTendik"`
	TotalWaliKelas    int     `json:"totalWaliKelas"`
	RataKehadiran     float64 `json:"rataKehadiran"`
	CutiPending       int     `json:"cutiPending"`
	TotalAlpha        int     `json:"totalAlpha"`
	RataSkp           float64 `json:"rataSkp"`
	DokumenKadaluarsa int     `json:"dokumenKadaluarsa"`
	BebanKurang24     int     `json:"bebanKurang24"`
}

type ChartItem struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

type StatisticCharts struct {
	KomposisiPegawai []ChartItem `json:"komposisiPegawai"`
	DistribusiUnit   []ChartItem `json:"distribusiUnit"`
	SebaranPangkat   []ChartItem `json:"sebaranPangkat"`
}

type StatisticDashboardData struct {
	KPIs   StatisticKPIs   `json:"kpis"`
	Charts StatisticCharts `json:"charts"`
}

type PtkStatisticRepository interface {
	GetDashboardData(ctx context.Context, schoolID string, startDate, endDate time.Time) (*StatisticDashboardData, error)
}

type pgxPtkStatisticRepository struct {
	db *pgxpool.Pool
}

func NewPtkStatisticRepository(db *pgxpool.Pool) PtkStatisticRepository {
	return &pgxPtkStatisticRepository{db: db}
}

func (r *pgxPtkStatisticRepository) GetDashboardData(ctx context.Context, schoolID string, startDate, endDate time.Time) (*StatisticDashboardData, error) {
	data := &StatisticDashboardData{}

	// 1. KPI Dasar (JOIN tabel_ptk dan ext)
	_ = r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM tabel_ptk tp 
		WHERE tp.soft_delete = 0 AND tp.status_keaktifan_id = 1
	`).Scan(&data.KPIs.TotalPegawai)

	// 2. KPI Kategori (Berdasarkan ptk_buku_induk_ext)
	// Guru Tetap: status_pegawai IN ('PNS', 'PPPK', 'GTY')
	_ = r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM ptk_buku_induk_ext ext
		JOIN tabel_ptk tp ON ext.ptk_id = tp.ptk_id
		WHERE tp.soft_delete = 0 AND tp.status_keaktifan_id = 1 
		  AND ext.jenis_pegawai = 'Guru'
		  AND ext.status_pegawai IN ('PNS', 'PPPK', 'GTY')
	`).Scan(&data.KPIs.TotalGuruTetap)

	// Guru Honorer
	_ = r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM ptk_buku_induk_ext ext
		JOIN tabel_ptk tp ON ext.ptk_id = tp.ptk_id
		WHERE tp.soft_delete = 0 AND tp.status_keaktifan_id = 1 
		  AND ext.jenis_pegawai = 'Guru'
		  AND ext.status_pegawai = 'Honorer'
	`).Scan(&data.KPIs.TotalHonorer)

	// Tenaga Kependidikan
	_ = r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM ptk_buku_induk_ext ext
		JOIN tabel_ptk tp ON ext.ptk_id = tp.ptk_id
		WHERE tp.soft_delete = 0 AND tp.status_keaktifan_id = 1 
		  AND ext.jenis_pegawai = 'Tendik'
	`).Scan(&data.KPIs.TotalTendik)

	// 3. KPI Wali Kelas (Berdasarkan tabel_kelas)
	_ = r.db.QueryRow(ctx, `
		SELECT COUNT(DISTINCT ptk_id) FROM tabel_kelas 
		WHERE ptk_id IS NOT NULL
	`).Scan(&data.KPIs.TotalWaliKelas)

	// 4. KPI Kehadiran & Alpha
	queryAtt := `
		SELECT 
			COALESCE(AVG(CASE WHEN status = 'HADIR' THEN 100.0 ELSE 0 END), 0) as avg_hadir,
			COUNT(CASE WHEN status = 'ALPHA' THEN 1 END) as total_alpha
		FROM ptk_attendance
		WHERE date >= $1 AND date <= $2
	`
	_ = r.db.QueryRow(ctx, queryAtt, startDate, endDate).Scan(&data.KPIs.RataKehadiran, &data.KPIs.TotalAlpha)

	// 5. KPI Cuti Pending
	_ = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM ptk_leaves WHERE status = 'pending'`).Scan(&data.KPIs.CutiPending)

	// 6. KPI SKP
	_ = r.db.QueryRow(ctx, `SELECT COALESCE(AVG(nilai_skp), 0) FROM ptk_performances`).Scan(&data.KPIs.RataSkp)

	// 7. KPI Dokumen Kadaluarsa
	_ = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM ptk_archives WHERE tanggal_kadaluarsa IS NOT NULL AND tanggal_kadaluarsa <= NOW() + INTERVAL '30 days'`).Scan(&data.KPIs.DokumenKadaluarsa)

	// 8. KPI Beban Kurang 24 JTM
	queryBeban := `
		SELECT COUNT(*) FROM (
			SELECT tp.ptk_id, COUNT(jp.id) * 2 as total_jam
			FROM tabel_ptk tp
			JOIN tabel_pembelajaran tpm ON tp.ptk_id = tpm.ptk_terdaftar_id
			JOIN tabel_jadwal_pelajaran jp ON tpm.pembelajaran_id = jp.pembelajaran_id
			WHERE tp.soft_delete = 0
			GROUP BY tp.ptk_id
			HAVING COUNT(jp.id) * 2 < 24
		) AS sub
	`
	_ = r.db.QueryRow(ctx, queryBeban).Scan(&data.KPIs.BebanKurang24)

	// --- CHARTS ---

	// Chart: Distribusi Unit Kerja (Dari ext)
	queryUnit := `SELECT COALESCE(unit_kerja, 'Tidak Ada Unit'), COUNT(*) FROM ptk_buku_induk_ext GROUP BY unit_kerja`
	rowsUnit, err := r.db.Query(ctx, queryUnit)
	if err == nil {
		defer rowsUnit.Close()
		for rowsUnit.Next() {
			var label string
			var val float64
			if err := rowsUnit.Scan(&label, &val); err == nil {
				data.Charts.DistribusiUnit = append(data.Charts.DistribusiUnit, ChartItem{Label: label, Value: val})
			}
		}
	}

	// Chart: Komposisi Pegawai (Dari ext)
	queryStatus := `SELECT status_pegawai, COUNT(*) FROM ptk_buku_induk_ext GROUP BY status_pegawai`
	rowsStatus, err := r.db.Query(ctx, queryStatus)
	if err == nil {
		defer rowsStatus.Close()
		for rowsStatus.Next() {
			var label string
			var val float64
			if err := rowsStatus.Scan(&label, &val); err == nil {
				data.Charts.KomposisiPegawai = append(data.Charts.KomposisiPegawai, ChartItem{Label: label, Value: val})
			}
		}
	}

	// Dummy Chart Sebaran Pangkat
	data.Charts.SebaranPangkat = []ChartItem{
		{Label: "Gol III/a", Value: 5}, {Label: "Gol III/b", Value: 8}, {Label: "Gol IV/a", Value: 2},
	}

	return data, nil
}
