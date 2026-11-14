package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
)

// =====================================================================
// GetRealisasiByTahun — hitung realisasi per pos RAB
//
// Sumber data:
//   - PENDAPATAN: SUM(debet) dari BKU (referensi_tabel='penerimaan_iuran')
//   - BELANJA   : SUM(kredit) dari BKU (referensi_tabel='pengeluaran',
//     status pengeluaran IN SIAP_CAIR/CAIR)
//
// Filter periode: berdasarkan `tanggal` BKU (bukan `created_at`)
// karena tanggal transaksi adalah tanggal efektif.
//
// sampaiBulan=0 → fallback ke 12 (full year)
// =====================================================================
func (r *pgxRABRepository) GetRealisasiByTahun(
	ctx context.Context, schoolID, periodID string, tahun, sampaiBulan int,
) (*models.RabRealisasiReport, error) {
	if sampaiBulan <= 0 || sampaiBulan > 12 {
		sampaiBulan = 12
	}

	// Cek RAB header ada
	var headerID, headerStatus string
	err := r.db.QueryRow(ctx, `
		SELECT id::text, status FROM rab_header
		WHERE school_id = $1 AND academic_period_id = $2 AND tahun_anggaran = $3
		LIMIT 1
	`, schoolID, periodID, tahun).Scan(&headerID, &headerStatus)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("RAB tahun %d belum dibuat", tahun)
	}
	if err != nil {
		return nil, fmt.Errorf("gagal lookup RAB: %w", err)
	}

	// Query realisasi per pos
	rows, err := r.db.Query(ctx, `
		SELECT
			d.id::text,
			d.kode_rekening_id::text,
			COALESCE(k.kode, ''),
			d.uraian,
			COALESCE(k.tipe, ''),
			d.subtotal,
			COALESCE((
				SELECT CASE
					WHEN k.tipe = 'PENDAPATAN' THEN
						(SELECT COALESCE(SUM(b.debet), 0) FROM buku_kas_umum b
						 WHERE b.kode_rekening_id = d.kode_rekening_id
						   AND b.academic_period_id = $2
						   AND EXTRACT(YEAR FROM b.tanggal) = $3
						   AND EXTRACT(MONTH FROM b.tanggal) <= $4
						   AND b.referensi_tabel = 'penerimaan_iuran')
					WHEN k.tipe = 'BELANJA' THEN
						(SELECT COALESCE(SUM(b.kredit), 0) FROM buku_kas_umum b
						 JOIN pengeluaran p ON p.id = b.referensi_id AND b.referensi_tabel = 'pengeluaran'
						 WHERE b.kode_rekening_id = d.kode_rekening_id
						   AND b.academic_period_id = $2
						   AND EXTRACT(YEAR FROM b.tanggal) = $3
						   AND EXTRACT(MONTH FROM b.tanggal) <= $4
						   AND p.status IN ('SIAP_CAIR', 'CAIR'))
					ELSE 0
				END
			), 0) AS realisasi
		FROM rab_detail d
		JOIN rab_header h ON h.id = d.rab_header_id
		LEFT JOIN kode_rekening k ON k.id = d.kode_rekening_id
		WHERE h.school_id = $1
		  AND h.academic_period_id = $2
		  AND h.tahun_anggaran = $3
		ORDER BY k.kode ASC
	`, schoolID, periodID, tahun, sampaiBulan)
	if err != nil {
		return nil, fmt.Errorf("gagal query realisasi: %w", err)
	}
	defer rows.Close()

	report := &models.RabRealisasiReport{
		TahunAnggaran: tahun,
		Items:         []models.RabRealisasiItem{},
	}

	for rows.Next() {
		var item models.RabRealisasiItem
		if err := rows.Scan(
			&item.RabDetailID, &item.KodeRekeningID,
			&item.KodeRekeningKode, &item.Uraian, &item.Tipe,
			&item.Pagu, &item.Realisasi,
		); err != nil {
			return nil, fmt.Errorf("gagal scan realisasi: %w", err)
		}

		// Hitung sisa & persen
		item.Sisa = item.Pagu - item.Realisasi
		if item.Pagu > 0 {
			item.PersenRealisasi = (item.Realisasi / item.Pagu) * 100
		}

		// Status warning
		switch {
		case item.PersenRealisasi > 95:
			item.StatusWarning = "OVER_BUDGET_WARNING"
		case item.PersenRealisasi < 10 && sampaiBulan > 6:
			item.StatusWarning = "UNDER_BUDGET_WARNING"
		default:
			item.StatusWarning = "NORMAL"
		}

		report.Items = append(report.Items, item)

		// Akumulasi total
		switch item.Tipe {
		case models.TipePendapatan:
			report.TotalPaguPendapatan += item.Pagu
			report.TotalRealisasiPendapatan += item.Realisasi
		case models.TipeBelanja:
			report.TotalPaguBelanja += item.Pagu
			report.TotalRealisasiBelanja += item.Realisasi
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	report.TotalSisaPendapatan = report.TotalPaguPendapatan - report.TotalRealisasiPendapatan
	report.TotalSisaBelanja = report.TotalPaguBelanja - report.TotalRealisasiBelanja

	return report, nil
}
