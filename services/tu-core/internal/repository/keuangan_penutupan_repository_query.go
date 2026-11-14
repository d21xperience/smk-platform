package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
)

// =====================================================================
// GetByID — header penutupan + join user_login (closed_by_nama)
// =====================================================================
func (r *pgxPenutupanRepository) GetByID(
	ctx context.Context, schoolID, id string,
) (*models.PenutupanBukuKas, error) {
	var p models.PenutupanBukuKas
	err := r.db.QueryRow(ctx, `
		SELECT p.id, p.school_id, p.academic_period_id, p.periode_bulan, p.periode_tahun,
		       p.saldo_awal, p.total_debet, p.total_kredit, p.saldo_akhir_buku,
		       p.saldo_fisik_kas, p.selisih_kas,
		       p.total_aset, p.total_kewajiban, p.total_ekuitas,
		       p.status, COALESCE(p.catatan,''),
		       COALESCE(p.closed_by::text,''), COALESCE(u.username,''),
		       p.closed_at,
		       COALESCE(p.ttd_kepsek_hash,''), COALESCE(p.ttd_kepsek_by::text,''), p.ttd_kepsek_at,
		       COALESCE(p.bku_pdf_url,''), COALESCE(p.neraca_pdf_url,''), COALESCE(p.lra_pdf_url,''),
		       p.created_at, p.updated_at
		FROM penutupan_buku_kas p
		LEFT JOIN user_login u ON u.user_id = p.closed_by
		WHERE p.school_id = $1 AND p.id = $2
	`, schoolID, id).Scan(
		&p.ID, &p.SchoolID, &p.AcademicPeriodID, &p.PeriodeBulan, &p.PeriodeTahun,
		&p.SaldoAwal, &p.TotalDebet, &p.TotalKredit, &p.SaldoAkhirBuku,
		&p.SaldoFisikKas, &p.SelisihKas,
		&p.TotalAset, &p.TotalKewajiban, &p.TotalEkuitas,
		&p.Status, &p.Catatan,
		&p.ClosedBy, &p.ClosedByNama,
		&p.ClosedAt,
		&p.TTDKepsekHash, &p.TTDKepsekBy, &p.TTDKepsekAt,
		&p.BkuPdfURL, &p.NeracaPdfURL, &p.LraPdfURL,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTransaksiKeuanganNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("gagal get penutupan: %w", err)
	}
	return &p, nil
}

// =====================================================================
// List — list header tanpa detail
// =====================================================================
func (r *pgxPenutupanRepository) List(
	ctx context.Context, filter models.PenutupanFilter,
) ([]models.PenutupanBukuKas, int, error) {
	limit, offset := normalizePagination(filter.Limit, filter.Offset)

	where := []string{"school_id = $1"}
	args := []interface{}{filter.SchoolID}
	pos := 2

	if filter.AcademicPeriodID != "" {
		where = append(where, fmt.Sprintf("academic_period_id = $%d", pos))
		args = append(args, filter.AcademicPeriodID)
		pos++
	}
	if filter.PeriodeTahun > 0 {
		where = append(where, fmt.Sprintf("periode_tahun = $%d", pos))
		args = append(args, filter.PeriodeTahun)
		pos++
	}
	if filter.Status != "" {
		where = append(where, fmt.Sprintf("status = $%d", pos))
		args = append(args, filter.Status)
		pos++
	}
	whereClause := strings.Join(where, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM penutupan_buku_kas WHERE "+whereClause, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("gagal count penutupan: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT id, school_id, academic_period_id, periode_bulan, periode_tahun,
		       saldo_awal, total_debet, total_kredit, saldo_akhir_buku,
		       saldo_fisik_kas, selisih_kas,
		       total_aset, total_kewajiban, total_ekuitas,
		       status, COALESCE(catatan,''),
		       COALESCE(closed_by::text,''), closed_at, created_at, updated_at
		FROM penutupan_buku_kas
		WHERE %s
		ORDER BY periode_tahun DESC, periode_bulan DESC, created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, pos, pos+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal list penutupan: %w", err)
	}
	defer rows.Close()

	out := []models.PenutupanBukuKas{}
	for rows.Next() {
		var p models.PenutupanBukuKas
		if err := rows.Scan(
			&p.ID, &p.SchoolID, &p.AcademicPeriodID, &p.PeriodeBulan, &p.PeriodeTahun,
			&p.SaldoAwal, &p.TotalDebet, &p.TotalKredit, &p.SaldoAkhirBuku,
			&p.SaldoFisikKas, &p.SelisihKas,
			&p.TotalAset, &p.TotalKewajiban, &p.TotalEkuitas,
			&p.Status, &p.Catatan,
			&p.ClosedBy, &p.ClosedAt, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("gagal scan penutupan: %w", err)
		}
		out = append(out, p)
	}
	return out, total, rows.Err()
}

// =====================================================================
// ExistsByPeriode — cek duplikasi
// =====================================================================
func (r *pgxPenutupanRepository) ExistsByPeriode(
	ctx context.Context, schoolID, periodID string, bulan, tahun int,
) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM penutupan_buku_kas
			WHERE school_id = $1 AND academic_period_id = $2
			  AND periode_bulan = $3 AND periode_tahun = $4
		)
	`, schoolID, periodID, bulan, tahun).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("gagal check exists: %w", err)
	}
	return exists, nil
}

// =====================================================================
// GetSnapshotBku — agregat BKU untuk bulan tsb
//
//	saldo_awal: saldo_berjalan terakhir SEBELUM bulan ini
//	total_debet/kredit: SUM di bulan ini
//	saldo_akhir_buku: saldo_awal + debet - kredit
//
// =====================================================================
func (r *pgxPenutupanRepository) GetSnapshotBku(
	ctx context.Context, schoolID, periodID string, bulan, tahun int,
) (*models.SnapshotBku, error) {
	s := &models.SnapshotBku{}

	// Saldo awal: saldo terakhir sebelum tanggal 1 bulan ini
	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(saldo_berjalan, 0)
		FROM buku_kas_umum
		WHERE school_id = $1 AND academic_period_id = $2
		  AND tanggal < make_date($3, $4, 1)
		ORDER BY tanggal DESC, created_at DESC
		LIMIT 1
	`, schoolID, periodID, tahun, bulan).Scan(&s.SaldoAwal)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("gagal get saldo awal: %w", err)
	}

	// Sum debet & kredit bulan ini
	err = r.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(debet), 0), COALESCE(SUM(kredit), 0)
		FROM buku_kas_umum
		WHERE school_id = $1 AND academic_period_id = $2
		  AND EXTRACT(MONTH FROM tanggal) = $3
		  AND EXTRACT(YEAR FROM tanggal) = $4
	`, schoolID, periodID, bulan, tahun).Scan(&s.TotalDebet, &s.TotalKredit)
	if err != nil {
		return nil, fmt.Errorf("gagal sum BKU bulanan: %w", err)
	}

	s.SaldoAkhirBuku = s.SaldoAwal + s.TotalDebet - s.TotalKredit
	return s, nil
}

// =====================================================================
// GetSnapshotPosisi — simplified posisi keuangan
//
//	aset = saldo kas (dari BKU) + saldo rekening (dari master)
//	kewajiban = 0 (belum ada modul utang di FASE 4)
//	ekuitas = aset - kewajiban
//
// =====================================================================
func (r *pgxPenutupanRepository) GetSnapshotPosisi(
	ctx context.Context, schoolID, periodID string,
) (*models.SnapshotPosisi, error) {
	var saldoKas, saldoRekening float64

	// Saldo kas dari BKU terakhir
	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(saldo_berjalan, 0)
		FROM buku_kas_umum
		WHERE school_id = $1 AND academic_period_id = $2
		ORDER BY tanggal DESC, created_at DESC
		LIMIT 1
	`, schoolID, periodID).Scan(&saldoKas)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("gagal get saldo kas: %w", err)
	}

	// Total saldo rekening bank aktif
	err = r.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(saldo_terkini), 0)
		FROM rekening_bank
		WHERE school_id = $1 AND is_aktif = TRUE
	`, schoolID).Scan(&saldoRekening)
	if err != nil {
		return nil, fmt.Errorf("gagal sum saldo rekening: %w", err)
	}

	aset := saldoKas + saldoRekening
	return &models.SnapshotPosisi{
		TotalAset:      aset,
		TotalKewajiban: 0,
		TotalEkuitas:   aset,
	}, nil
}

// =====================================================================
// SignPenutupan — update hash + signer + status ke CLOSED_FULLY_SIGNED
//
// Atomic: hanya update kalau status saat ini CLOSED (belum signed).
// Ini mencegah double-sign / race condition.
// =====================================================================
func (r *pgxPenutupanRepository) SignPenutupan(
	ctx context.Context, schoolID, id, hash, signerID string, signedAt time.Time,
) error {
	ct, err := r.db.Exec(ctx, `
		UPDATE penutupan_buku_kas
		SET ttd_kepsek_hash = $1,
		    ttd_kepsek_by = NULLIF($2,''),
		    ttd_kepsek_at = $3,
		    status = 'CLOSED_FULLY_SIGNED',
		    updated_at = NOW()
		WHERE school_id = $4 AND id = $5
		  AND status = 'CLOSED'
	`, hash, signerID, signedAt, schoolID, id)
	if err != nil {
		return fmt.Errorf("gagal sign penutupan: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("penutupan tidak ditemukan atau sudah di-sign (status bukan CLOSED)")
	}
	return nil
}

// =====================================================================
// GetPenutupanByTTDHash — lookup dokumen via hash TTD (untuk endpoint verify publik)
// =====================================================================
func (r *pgxPenutupanRepository) GetPenutupanByTTDHash(
	ctx context.Context, hash string,
) (*models.PenutupanBukuKas, error) {
	var p models.PenutupanBukuKas
	err := r.db.QueryRow(ctx, `
		SELECT p.id, p.school_id, p.academic_period_id, p.periode_bulan, p.periode_tahun,
		       p.saldo_awal, p.total_debet, p.total_kredit, p.saldo_akhir_buku,
		       p.saldo_fisik_kas, p.selisih_kas,
		       p.total_aset, p.total_kewajiban, p.total_ekuitas,
		       p.status, COALESCE(p.catatan,''),
		       COALESCE(p.closed_by::text,''), COALESCE(u1.username,''),
		       p.closed_at,
		       COALESCE(p.ttd_kepsek_hash,''), COALESCE(p.ttd_kepsek_by::text,''),
		       COALESCE(u2.username,''), p.ttd_kepsek_at,
		       COALESCE(p.bku_pdf_url,''), COALESCE(p.neraca_pdf_url,''), COALESCE(p.lra_pdf_url,''),
		       p.created_at, p.updated_at
		FROM penutupan_buku_kas p
		LEFT JOIN user_login u1 ON u1.user_id = p.closed_by
		LEFT JOIN user_login u2 ON u2.user_id = p.ttd_kepsek_by
		WHERE p.ttd_kepsek_hash = $1
		LIMIT 1
	`, hash).Scan(
		&p.ID, &p.SchoolID, &p.AcademicPeriodID, &p.PeriodeBulan, &p.PeriodeTahun,
		&p.SaldoAwal, &p.TotalDebet, &p.TotalKredit, &p.SaldoAkhirBuku,
		&p.SaldoFisikKas, &p.SelisihKas,
		&p.TotalAset, &p.TotalKewajiban, &p.TotalEkuitas,
		&p.Status, &p.Catatan,
		&p.ClosedBy, &p.ClosedByNama,
		&p.ClosedAt,
		&p.TTDKepsekHash, &p.TTDKepsekBy, &p.TTDKepsekNama, &p.TTDKepsekAt,
		&p.BkuPdfURL, &p.NeracaPdfURL, &p.LraPdfURL,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTransaksiKeuanganNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("gagal get penutupan by hash: %w", err)
	}
	return &p, nil
}
