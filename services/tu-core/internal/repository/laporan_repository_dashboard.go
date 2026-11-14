package repository

import (
	"context"
	"fmt"
)

// =====================================================================
// GetPiutangForDashboard — total piutang + jumlah item + jumlah siswa
// =====================================================================
func (r *pgxLaporanRepository) GetPiutangForDashboard(
	ctx context.Context, schoolID, periodID string,
) (*PiutangDashboardData, error) {
	out := &PiutangDashboardData{}
	err := r.db.QueryRow(ctx, `
		SELECT 
			COALESCE(SUM(nominal - total_dibayar), 0) AS total_sisa,
			COUNT(*) AS jumlah_item,
			COUNT(DISTINCT nis) AS jumlah_siswa
		FROM tagihan_siswa
		WHERE school_id = $1 AND academic_period_id = $2
		  AND status IN ('BELUM_LUNAS', 'SEBAGIAN')
		  AND (nominal - total_dibayar) > 0
	`, schoolID, periodID).Scan(&out.TotalSisa, &out.JumlahItem, &out.JumlahSiswa)
	if err != nil {
		return nil, fmt.Errorf("gagal query piutang dashboard: %w", err)
	}
	return out, nil
}

// =====================================================================
// GetPengajuanForDashboard — count per status + total nominal pending
//
// "Pending" = MENUNGGU_VERIFIKASI_BENDAHARA + MENUNGGU_PERSETUJUAN_KEPALA
// =====================================================================
func (r *pgxLaporanRepository) GetPengajuanForDashboard(
	ctx context.Context, schoolID, periodID string,
) (map[string]int, float64, error) {
	rows, err := r.db.Query(ctx, `
		SELECT status, COUNT(*) AS cnt, COALESCE(SUM(jumlah), 0) AS total
		FROM pengajuan_pengeluaran
		WHERE school_id = $1 AND academic_period_id = $2
		GROUP BY status
	`, schoolID, periodID)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal query pengajuan dashboard: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int)
	var totalPending float64

	for rows.Next() {
		var status string
		var cnt int
		var total float64
		if err := rows.Scan(&status, &cnt, &total); err != nil {
			return nil, 0, fmt.Errorf("gagal scan pengajuan: %w", err)
		}
		counts[status] = cnt
		// Akumulasi total nominal untuk status pending
		if status == "MENUNGGU_VERIFIKASI_BENDAHARA" || status == "MENUNGGU_PERSETUJUAN_KEPALA" {
			totalPending += total
		}
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return counts, totalPending, nil
}

// =====================================================================
// GetRekonForDashboard — total, balance, unbalance untuk bulan tertentu
// =====================================================================
func (r *pgxLaporanRepository) GetRekonForDashboard(
	ctx context.Context, schoolID, periodID string, bulan, tahun int,
) (total, balance, unbalance int, err error) {
	err = r.db.QueryRow(ctx, `
		SELECT 
			COUNT(*)::int AS total,
			COUNT(*) FILTER (WHERE status = 'BALANCE')::int AS balance,
			COUNT(*) FILTER (WHERE status = 'UNBALANCE')::int AS unbalance
		FROM rekonsiliasi_bank
		WHERE school_id = $1 AND academic_period_id = $2
		  AND periode_bulan = $3 AND periode_tahun = $4
	`, schoolID, periodID, bulan, tahun).Scan(&total, &balance, &unbalance)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("gagal query rekon dashboard: %w", err)
	}
	return total, balance, unbalance, nil
}
