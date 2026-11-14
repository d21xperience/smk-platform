package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
)

// =====================================================================
// GetRekonsiliasiByID — header + details + unmatched
// =====================================================================
func (r *pgxRekonsiliasiRepository) GetRekonsiliasiByID(
	ctx context.Context, schoolID, id string,
) (*models.RekonsiliasiBank, error) {
	var h models.RekonsiliasiBank
	err := r.db.QueryRow(ctx, `
		SELECT r.id, r.school_id, r.academic_period_id, r.rekening_bank_id,
		       COALESCE(b.nama_bank || ' - ' || b.nomor_rekening, ''),
		       r.periode_bulan, r.periode_tahun,
		       r.saldo_buku, r.saldo_bank, r.selisih,
		       r.jumlah_matched, r.jumlah_unmatched_bank, r.jumlah_unmatched_buku,
		       r.saldo_rekonsiliasi_bank, r.saldo_rekonsiliasi_buku,
		       r.status, COALESCE(r.file_mutasi_url,''),
		       COALESCE(r.catatan,''), COALESCE(r.created_by::text,''), r.created_at
		FROM rekonsiliasi_bank r
		LEFT JOIN rekening_bank b ON b.id = r.rekening_bank_id
		WHERE r.school_id = $1 AND r.id = $2
	`, schoolID, id).Scan(
		&h.ID, &h.SchoolID, &h.AcademicPeriodID, &h.RekeningBankID,
		&h.RekeningBankNama,
		&h.PeriodeBulan, &h.PeriodeTahun,
		&h.SaldoBuku, &h.SaldoBank, &h.Selisih,
		&h.JumlahMatched, &h.JumlahUnmatchedBank, &h.JumlahUnmatchedBuku,
		&h.SaldoRekonsiliasiBank, &h.SaldoRekonsiliasiBuku,
		&h.Status, &h.FileMutasiURL,
		&h.Catatan, &h.CreatedBy, &h.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTransaksiKeuanganNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("gagal get rekonsiliasi: %w", err)
	}

	details, err := r.loadDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	h.Details = details

	unmatched, err := r.loadUnmatched(ctx, id)
	if err != nil {
		return nil, err
	}
	h.Unmatched = unmatched

	return &h, nil
}

func (r *pgxRekonsiliasiRepository) loadDetails(
	ctx context.Context, rekonID string,
) ([]models.RekonsiliasiDetail, error) {
	rows, err := r.db.Query(ctx, `
		SELECT d.id, d.rekonsiliasi_id, d.mutasi_tanggal::text,
		       COALESCE(d.mutasi_keterangan,''),
		       d.mutasi_debet, d.mutasi_kredit, COALESCE(d.mutasi_saldo,0),
		       COALESCE(d.mutasi_referensi,''),
		       d.buku_kas_umum_id,
		       COALESCE(b.no_bukti,''), COALESCE(b.uraian,''),
		       d.selisih_hari, d.created_at
		FROM rekonsiliasi_detail d
		LEFT JOIN buku_kas_umum b ON b.id = d.buku_kas_umum_id
		WHERE d.rekonsiliasi_id = $1
		ORDER BY d.mutasi_tanggal ASC
	`, rekonID)
	if err != nil {
		return nil, fmt.Errorf("gagal load details: %w", err)
	}
	defer rows.Close()

	out := []models.RekonsiliasiDetail{}
	for rows.Next() {
		var d models.RekonsiliasiDetail
		if err := rows.Scan(
			&d.ID, &d.RekonsiliasiID, &d.MutasiTanggal,
			&d.MutasiKeterangan, &d.MutasiDebet, &d.MutasiKredit, &d.MutasiSaldo,
			&d.MutasiReferensi, &d.BukuKasUmumID,
			&d.NoBuktiBKU, &d.UraianBKU, &d.SelisihHari, &d.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("gagal scan detail: %w", err)
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (r *pgxRekonsiliasiRepository) loadUnmatched(
	ctx context.Context, rekonID string,
) ([]models.RekonsiliasiUnmatched, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, rekonsiliasi_id, tipe, data_json::text,
		       status_tindak_lanjut, created_at
		FROM rekonsiliasi_unmatched
		WHERE rekonsiliasi_id = $1
		ORDER BY tipe ASC, created_at ASC
	`, rekonID)
	if err != nil {
		return nil, fmt.Errorf("gagal load unmatched: %w", err)
	}
	defer rows.Close()

	out := []models.RekonsiliasiUnmatched{}
	for rows.Next() {
		var u models.RekonsiliasiUnmatched
		var jsonStr *string
		if err := rows.Scan(
			&u.ID, &u.RekonsiliasiID, &u.Tipe, &jsonStr,
			&u.StatusTindakLanjut, &u.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("gagal scan unmatched: %w", err)
		}
		if jsonStr != nil {
			u.DataJSON = *jsonStr
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

// =====================================================================
// ListRekonsiliasi — list header (tanpa detail/unmatched, untuk performa)
// =====================================================================
func (r *pgxRekonsiliasiRepository) ListRekonsiliasi(
	ctx context.Context, filter models.RekonsiliasiFilter,
) ([]models.RekonsiliasiBank, int, error) {
	limit, offset := normalizePagination(filter.Limit, filter.Offset)
	whereClause, args, pos := buildRekonsiliasiWhere(filter)

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM rekonsiliasi_bank r WHERE "+whereClause, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("gagal count rekonsiliasi: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT r.id, r.school_id, r.academic_period_id, r.rekening_bank_id,
		       COALESCE(b.nama_bank || ' - ' || b.nomor_rekening, ''),
		       r.periode_bulan, r.periode_tahun,
		       r.saldo_buku, r.saldo_bank, r.selisih,
		       r.jumlah_matched, r.jumlah_unmatched_bank, r.jumlah_unmatched_buku,
		       r.saldo_rekonsiliasi_bank, r.saldo_rekonsiliasi_buku,
		       r.status, COALESCE(r.file_mutasi_url,''),
		       COALESCE(r.catatan,''), COALESCE(r.created_by::text,''), r.created_at
		FROM rekonsiliasi_bank r
		LEFT JOIN rekening_bank b ON b.id = r.rekening_bank_id
		WHERE %s
		ORDER BY r.periode_tahun DESC, r.periode_bulan DESC, r.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, pos, pos+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal list rekonsiliasi: %w", err)
	}
	defer rows.Close()

	out := []models.RekonsiliasiBank{}
	for rows.Next() {
		var h models.RekonsiliasiBank
		if err := rows.Scan(
			&h.ID, &h.SchoolID, &h.AcademicPeriodID, &h.RekeningBankID,
			&h.RekeningBankNama, &h.PeriodeBulan, &h.PeriodeTahun,
			&h.SaldoBuku, &h.SaldoBank, &h.Selisih,
			&h.JumlahMatched, &h.JumlahUnmatchedBank, &h.JumlahUnmatchedBuku,
			&h.SaldoRekonsiliasiBank, &h.SaldoRekonsiliasiBuku,
			&h.Status, &h.FileMutasiURL,
			&h.Catatan, &h.CreatedBy, &h.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("gagal scan rekonsiliasi: %w", err)
		}
		out = append(out, h)
	}
	return out, total, rows.Err()
}

// =====================================================================
// GetBkuForPeriod — list BKU untuk auto-match source
// Hanya BKU yang rekening_bank_id = target (yang terkait mutasi bank)
// =====================================================================
func (r *pgxRekonsiliasiRepository) GetBkuForPeriod(
	ctx context.Context, schoolID, periodID, rekeningID string, bulan, tahun int,
) ([]models.BkuForMatch, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, tanggal, no_bukti, uraian,
		       debet, kredit, saldo_berjalan, kode_rekening_id
		FROM buku_kas_umum
		WHERE school_id = $1
		  AND academic_period_id = $2
		  AND rekening_bank_id = $3
		  AND EXTRACT(MONTH FROM tanggal) = $4
		  AND EXTRACT(YEAR FROM tanggal) = $5
		ORDER BY tanggal ASC, created_at ASC
	`, schoolID, periodID, rekeningID, bulan, tahun)
	if err != nil {
		return nil, fmt.Errorf("gagal get BKU for period: %w", err)
	}
	defer rows.Close()

	out := []models.BkuForMatch{}
	for rows.Next() {
		var b models.BkuForMatch
		if err := rows.Scan(
			&b.ID, &b.Tanggal, &b.NoBukti, &b.Uraian,
			&b.Debet, &b.Kredit, &b.SaldoBerjalan, &b.KodeRekeningID,
		); err != nil {
			return nil, fmt.Errorf("gagal scan BKU: %w", err)
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// =====================================================================
// GetSaldoBkuByPeriode — saldo terakhir s.d. akhir bulan (untuk header)
// =====================================================================
func (r *pgxRekonsiliasiRepository) GetSaldoBkuByPeriode(
	ctx context.Context, schoolID, periodID, rekeningID string, bulan, tahun int,
) (float64, error) {
	var saldo float64
	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(saldo_berjalan, 0)
		FROM buku_kas_umum
		WHERE school_id = $1
		  AND academic_period_id = $2
		  AND rekening_bank_id = $3
		  AND EXTRACT(MONTH FROM tanggal) <= $4
		  AND EXTRACT(YEAR FROM tanggal) = $5
		ORDER BY tanggal DESC, created_at DESC
		LIMIT 1
	`, schoolID, periodID, rekeningID, bulan, tahun).Scan(&saldo)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("gagal get saldo BKU: %w", err)
	}
	return saldo, nil
}

// =====================================================================
// HasBalanceRekonsiliasi — cek apakah ada minimal 1 rekon BALANCE di bulan tsb
// Dipakai PenutupanService untuk validasi sebelum lock periode
// =====================================================================
func (r *pgxRekonsiliasiRepository) HasBalanceRekonsiliasi(
	ctx context.Context, schoolID, periodID string, bulan, tahun int,
) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM rekonsiliasi_bank
			WHERE school_id = $1 AND academic_period_id = $2
			  AND periode_bulan = $3 AND periode_tahun = $4
			  AND status = 'BALANCE'
		)
	`, schoolID, periodID, bulan, tahun).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("gagal check balance rekon: %w", err)
	}
	return exists, nil
}
