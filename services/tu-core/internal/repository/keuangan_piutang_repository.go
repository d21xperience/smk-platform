package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
)

// =====================================================================
// PiutangRepository — interface read-only untuk aging & detail piutang
// =====================================================================
type PiutangRepository interface {
	ListPiutang(ctx context.Context, filter models.PiutangFilter) ([]models.PiutangItem, int, error)
	GetPiutangByNIS(ctx context.Context, schoolID, periodID, nis string) ([]models.PiutangItem, error)
	GetAgingSummary(ctx context.Context, schoolID, periodID string) ([]models.AgingBucket, float64, int, error)
	GetAllPiutangForPeriod(ctx context.Context, schoolID, periodID string) ([]models.PiutangItem, error)
}

type pgxPiutangRepository struct {
	db *pgxpool.Pool
}

func NewPiutangRepository(db *pgxpool.Pool) PiutangRepository {
	return &pgxPiutangRepository{db: db}
}

// =====================================================================
// Base CTE — semua piutang untuk periode + aging computed
//
// Dipakai ulang di semua query supaya filter aging_group bisa WHERE di CTE
// (menghindari duplikasi CASE di subquery).
// =====================================================================
const piutangBaseCTE = `
	WITH piutang AS (
		SELECT 
			t.id::text AS id,
			t.school_id,
			t.academic_period_id,
			t.nis,
			COALESCE(s.nm_siswa, '') AS nama_siswa,
			COALESCE(t.kelas_nama, '') AS kelas_nama,
			t.jenis_iuran_id::text AS jenis_iuran_id,
			COALESCE(j.nama, '') AS jenis_iuran_nama,
			t.periode_bulan,
			t.periode_tahun,
			t.nominal,
			t.total_dibayar,
			(t.nominal - t.total_dibayar) AS sisa,
			t.status,
			COALESCE(t.tanggal_jatuh_tempo::text, '') AS tanggal_jatuh_tempo,
			COALESCE((CURRENT_DATE - t.tanggal_jatuh_tempo)::int, 0) AS hari_terlambat,
			(CASE 
				WHEN t.tanggal_jatuh_tempo IS NULL THEN 'BELUM_JATUH_TEMPO'
				WHEN CURRENT_DATE < t.tanggal_jatuh_tempo THEN 'BELUM_JATUH_TEMPO'
				WHEN (CURRENT_DATE - t.tanggal_jatuh_tempo) < 30 THEN '<30'
				WHEN (CURRENT_DATE - t.tanggal_jatuh_tempo) < 60 THEN '30-60'
				WHEN (CURRENT_DATE - t.tanggal_jatuh_tempo) < 90 THEN '60-90'
				ELSE '>90'
			END) AS aging_group
		FROM tagihan_siswa t
		LEFT JOIN tabel_siswa s ON s.nis = t.nis
		LEFT JOIN jenis_iuran j ON j.id = t.jenis_iuran_id
		WHERE t.school_id = $1
		  AND t.academic_period_id = $2
		  AND t.status IN ('BELUM_LUNAS', 'SEBAGIAN')
		  AND (t.nominal - t.total_dibayar) > 0
	)
`

// =====================================================================
// ListPiutang — dengan filter (NIS, kelas, aging, jenis) + pagination
// =====================================================================
func (r *pgxPiutangRepository) ListPiutang(
	ctx context.Context, filter models.PiutangFilter,
) ([]models.PiutangItem, int, error) {
	limit, offset := normalizePagination(filter.Limit, filter.Offset)

	where := []string{"1=1"}
	args := []interface{}{filter.SchoolID, filter.AcademicPeriodID}
	pos := 3

	if filter.NIS != "" {
		where = append(where, fmt.Sprintf("nis = $%d", pos))
		args = append(args, filter.NIS)
		pos++
	}
	if filter.KelasNama != "" {
		where = append(where, fmt.Sprintf("kelas_nama = $%d", pos))
		args = append(args, filter.KelasNama)
		pos++
	}
	if filter.AgingGroup != "" {
		where = append(where, fmt.Sprintf("aging_group = $%d", pos))
		args = append(args, filter.AgingGroup)
		pos++
	}
	if filter.JenisIuranID != "" {
		where = append(where, fmt.Sprintf("jenis_iuran_id = $%d", pos))
		args = append(args, filter.JenisIuranID)
		pos++
	}
	whereClause := strings.Join(where, " AND ")

	// Total
	countQuery := piutangBaseCTE + " SELECT COUNT(*) FROM piutang WHERE " + whereClause
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("gagal count piutang: %w", err)
	}

	// List
	listQuery := fmt.Sprintf(`
		%s
		SELECT id, school_id, academic_period_id, nis, nama_siswa, kelas_nama,
		       jenis_iuran_id, jenis_iuran_nama, periode_bulan, periode_tahun,
		       nominal, total_dibayar, sisa, status, tanggal_jatuh_tempo,
		       hari_terlambat, aging_group
		FROM piutang
		WHERE %s
		ORDER BY tanggal_jatuh_tempo ASC NULLS LAST, id ASC
		LIMIT $%d OFFSET $%d
	`, piutangBaseCTE, whereClause, pos, pos+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal list piutang: %w", err)
	}
	defer rows.Close()

	out := []models.PiutangItem{}
	for rows.Next() {
		var p models.PiutangItem
		if err := rows.Scan(
			&p.ID, &p.SchoolID, &p.AcademicPeriodID, &p.NIS, &p.NamaSiswa, &p.KelasNama,
			&p.JenisIuranID, &p.JenisIuranNama, &p.PeriodeBulan, &p.PeriodeTahun,
			&p.Nominal, &p.TotalDibayar, &p.Sisa, &p.Status, &p.TanggalJatuhTempo,
			&p.HariTerlambat, &p.AgingGroup,
		); err != nil {
			return nil, 0, fmt.Errorf("gagal scan piutang: %w", err)
		}
		out = append(out, p)
	}
	return out, total, rows.Err()
}

// =====================================================================
// GetPiutangByNIS — semua piutang 1 siswa
// =====================================================================
func (r *pgxPiutangRepository) GetPiutangByNIS(
	ctx context.Context, schoolID, periodID, nis string,
) ([]models.PiutangItem, error) {
	query := fmt.Sprintf(`
		%s
		SELECT id, school_id, academic_period_id, nis, nama_siswa, kelas_nama,
		       jenis_iuran_id, jenis_iuran_nama, periode_bulan, periode_tahun,
		       nominal, total_dibayar, sisa, status, tanggal_jatuh_tempo,
		       hari_terlambat, aging_group
		FROM piutang
		WHERE nis = $3
		ORDER BY tanggal_jatuh_tempo ASC NULLS LAST
	`, piutangBaseCTE)

	rows, err := r.db.Query(ctx, query, schoolID, periodID, nis)
	if err != nil {
		return nil, fmt.Errorf("gagal get piutang per NIS: %w", err)
	}
	defer rows.Close()

	out := []models.PiutangItem{}
	for rows.Next() {
		var p models.PiutangItem
		if err := rows.Scan(
			&p.ID, &p.SchoolID, &p.AcademicPeriodID, &p.NIS, &p.NamaSiswa, &p.KelasNama,
			&p.JenisIuranID, &p.JenisIuranNama, &p.PeriodeBulan, &p.PeriodeTahun,
			&p.Nominal, &p.TotalDibayar, &p.Sisa, &p.Status, &p.TanggalJatuhTempo,
			&p.HariTerlambat, &p.AgingGroup,
		); err != nil {
			return nil, fmt.Errorf("gagal scan piutang: %w", err)
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// =====================================================================
// GetAgingSummary — GROUP BY aging_group
// Return: buckets (5 selalu muncul), total_sisa, total_item_count
// =====================================================================
func (r *pgxPiutangRepository) GetAgingSummary(
	ctx context.Context, schoolID, periodID string,
) ([]models.AgingBucket, float64, int, error) {
	query := fmt.Sprintf(`
		%s
		SELECT aging_group, COUNT(*)::int AS item_count, COALESCE(SUM(sisa), 0) AS total_sisa
		FROM piutang
		GROUP BY aging_group
	`, piutangBaseCTE)

	rows, err := r.db.Query(ctx, query, schoolID, periodID)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("gagal aging summary: %w", err)
	}
	defer rows.Close()

	// Map dulu untuk lookup cepat
	byGroup := make(map[string]models.AgingBucket)
	var totalSisa float64
	var totalCount int
	for rows.Next() {
		var b models.AgingBucket
		if err := rows.Scan(&b.Group, &b.ItemCount, &b.TotalSisa); err != nil {
			return nil, 0, 0, fmt.Errorf("gagal scan bucket: %w", err)
		}
		byGroup[b.Group] = b
		totalSisa += b.TotalSisa
		totalCount += b.ItemCount
	}
	if err := rows.Err(); err != nil {
		return nil, 0, 0, err
	}

	// Pastikan 5 bucket selalu muncul (0 kalau tidak ada data)
	out := make([]models.AgingBucket, 0, 5)
	for _, g := range models.AllAgingGroups() {
		if b, ok := byGroup[g]; ok {
			out = append(out, b)
		} else {
			out = append(out, models.AgingBucket{Group: g, ItemCount: 0, TotalSisa: 0})
		}
	}
	return out, totalSisa, totalCount, nil
}

// =====================================================================
// GetAllPiutangForPeriod — helper untuk service compose PiutangPerSiswa
// (group by NIS di service layer)
// =====================================================================
func (r *pgxPiutangRepository) GetAllPiutangForPeriod(
	ctx context.Context, schoolID, periodID string,
) ([]models.PiutangItem, error) {
	query := fmt.Sprintf(`
		%s
		SELECT id, school_id, academic_period_id, nis, nama_siswa, kelas_nama,
		       jenis_iuran_id, jenis_iuran_nama, periode_bulan, periode_tahun,
		       nominal, total_dibayar, sisa, status, tanggal_jatuh_tempo,
		       hari_terlambat, aging_group
		FROM piutang
		ORDER BY nis ASC, tanggal_jatuh_tempo ASC NULLS LAST
	`, piutangBaseCTE)

	rows, err := r.db.Query(ctx, query, schoolID, periodID)
	if err != nil {
		return nil, fmt.Errorf("gagal get all piutang: %w", err)
	}
	defer rows.Close()

	out := []models.PiutangItem{}
	for rows.Next() {
		var p models.PiutangItem
		if err := rows.Scan(
			&p.ID, &p.SchoolID, &p.AcademicPeriodID, &p.NIS, &p.NamaSiswa, &p.KelasNama,
			&p.JenisIuranID, &p.JenisIuranNama, &p.PeriodeBulan, &p.PeriodeTahun,
			&p.Nominal, &p.TotalDibayar, &p.Sisa, &p.Status, &p.TanggalJatuhTempo,
			&p.HariTerlambat, &p.AgingGroup,
		); err != nil {
			return nil, fmt.Errorf("gagal scan piutang: %w", err)
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
