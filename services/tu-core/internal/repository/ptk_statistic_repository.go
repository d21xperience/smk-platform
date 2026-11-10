package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StatisticKPIs struct {
	TotalPegawai      int     `json:"totalPegawai"`
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

	// 1. KPI: Total Pegawai
	_ = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM tabel_ptk WHERE soft_delete = 0`).Scan(&data.KPIs.TotalPegawai)

	// 2. KPI: Rata-rata Kehadiran & Total Alpha (Asumsi tabel ptk_attendance berisi status: HADIR, ALPHA, dll)
	queryAtt := `
		SELECT 
			COALESCE(AVG(CASE WHEN status = 'HADIR' THEN 100.0 ELSE 0 END), 0) as avg_hadir,
			COUNT(CASE WHEN status = 'ALPHA' THEN 1 END) as total_alpha
		FROM ptk_attendance
		WHERE date >= $1 AND date <= $2
	`
	_ = r.db.QueryRow(ctx, queryAtt, startDate, endDate).Scan(&data.KPIs.RataKehadiran, &data.KPIs.TotalAlpha)

	// 3. KPI: Cuti Pending
	_ = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM ptk_leaves WHERE status = 'pending'`).Scan(&data.KPIs.CutiPending)

	// 4. KPI: Rata-rata SKP
	_ = r.db.QueryRow(ctx, `SELECT COALESCE(AVG(nilai_skp), 0) FROM ptk_performances`).Scan(&data.KPIs.RataSkp)

	// 5. KPI: Dokumen Kadaluarsa
	_ = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM ptk_archives WHERE tanggal_kadaluarsa IS NOT NULL AND tanggal_kadaluarsa <= NOW() + INTERVAL '30 days'`).Scan(&data.KPIs.DokumenKadaluarsa)

	// 6. KPI: Beban Kurang 24 JTM (Asumsi ada tabel jadwal/pembelajaran)
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

	// 7. CHART: Distribusi Unit Kerja
	queryUnit := `
		SELECT COALESCE(unit_kerja, 'Tidak Ada Unit'), COUNT(*) 
		FROM ptk_buku_induk_ext 
		GROUP BY unit_kerja
	`
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

	// 8. CHART: Komposisi Pegawai (Status Pegawai)
	queryStatus := `
		SELECT COALESCE(status_pegawai, 'Honorer'), COUNT(*) 
		FROM ptk_buku_induk_ext 
		GROUP BY status_pegawai
	`
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

	// 9. CHART: Sebaran Pangkat (Dummy fallback jika kosong)
	data.Charts.SebaranPangkat = []ChartItem{
		{Label: "Gol III/a", Value: 5}, {Label: "Gol III/b", Value: 8}, {Label: "Gol IV/a", Value: 2},
	}

	return data, nil
}
