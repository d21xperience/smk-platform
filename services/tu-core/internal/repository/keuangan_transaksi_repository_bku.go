package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
)

// =====================================================================
// BUKU KAS UMUM (BKU)
// =====================================================================

// InsertBKU — WAJIB dipanggil dalam tx (untuk membaca saldo terakhir
// secara konsisten dan mencegah race condition).
//
// saldo_berjalan dihitung otomatis dari saldo terakhir + (debet - kredit).
// Field SaldoBerjalan pada parameter b akan di-override.
func (r *pgxKeuanganTransaksiRepository) InsertBKU(
	ctx context.Context, tx pgx.Tx, b *models.BukuKasUmum,
) error {
	// Hitung saldo terakhir dalam transaksi yang sama (row lock implisit
	// via ORDER BY created_at DESC)
	var lastSaldo float64
	err := tx.QueryRow(ctx, `
		SELECT COALESCE(saldo_berjalan, 0)
		FROM buku_kas_umum
		WHERE school_id = $1 AND academic_period_id = $2 AND tanggal <= $3::date
		ORDER BY tanggal DESC, created_at DESC
		LIMIT 1
	`, b.SchoolID, b.AcademicPeriodID, b.Tanggal).Scan(&lastSaldo)
	if err != nil && err != pgx.ErrNoRows {
		return fmt.Errorf("gagal hitung saldo terakhir: %w", err)
	}

	newSaldo := lastSaldo + b.Debet - b.Kredit
	b.SaldoBerjalan = newSaldo

	query := `
		INSERT INTO buku_kas_umum (
			id, school_id, academic_period_id, tanggal, no_bukti, uraian,
			kode_rekening_id, rekening_bank_id, debet, kredit, saldo_berjalan,
			referensi_tabel, referensi_id, created_by
		) VALUES (
			gen_random_uuid(), $1, $2, $3::date, $4, $5,
			$6, NULLIF($7,'')::uuid, $8, $9, $10,
			NULLIF($11,''), NULLIF($12,'')::uuid, NULLIF($13,'')
		)
		RETURNING id, created_at
	`
	row := tx.QueryRow(ctx, query,
		b.SchoolID, b.AcademicPeriodID, b.Tanggal, b.NoBukti, b.Uraian,
		b.KodeRekeningID, b.RekeningBankID, b.Debet, b.Kredit, b.SaldoBerjalan,
		b.ReferensiTabel, b.ReferensiID, b.CreatedBy,
	)
	if err := row.Scan(&b.ID, &b.CreatedAt); err != nil {
		return fmt.Errorf("gagal insert BKU: %w", err)
	}
	return nil
}

func (r *pgxKeuanganTransaksiRepository) ListBKU(
	ctx context.Context, filter models.BukuKasFilter,
) ([]models.BukuKasUmum, int, error) {
	limit, offset := normalizePagination(filter.Limit, filter.Offset)

	where := []string{"b.school_id = $1"}
	args := []interface{}{filter.SchoolID}
	pos := 2

	if filter.AcademicPeriodID != "" {
		where = append(where, fmt.Sprintf("b.academic_period_id = $%d", pos))
		args = append(args, filter.AcademicPeriodID)
		pos++
	}
	if filter.KodeRekeningID != "" {
		where = append(where, fmt.Sprintf("b.kode_rekening_id = $%d", pos))
		args = append(args, filter.KodeRekeningID)
		pos++
	}
	if filter.TanggalDari != "" {
		where = append(where, fmt.Sprintf("b.tanggal >= $%d::date", pos))
		args = append(args, filter.TanggalDari)
		pos++
	}
	if filter.TanggalSampai != "" {
		where = append(where, fmt.Sprintf("b.tanggal <= $%d::date", pos))
		args = append(args, filter.TanggalSampai)
		pos++
	}

	whereClause := strings.Join(where, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM buku_kas_umum b WHERE "+whereClause, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("gagal count BKU: %w", err)
	}

	listQuery := fmt.Sprintf(`
		SELECT b.id, b.school_id, b.academic_period_id, b.tanggal::text, b.no_bukti, b.uraian,
		       b.kode_rekening_id, COALESCE(k.kode,''), COALESCE(b.rekening_bank_id::text,''),
		       b.debet, b.kredit, b.saldo_berjalan,
		       COALESCE(b.referensi_tabel,''), COALESCE(b.referensi_id::text,''),
		       b.is_locked, b.is_rekonsiliasi, COALESCE(b.created_by::text,''), b.created_at
		FROM buku_kas_umum b
		LEFT JOIN kode_rekening k ON k.id = b.kode_rekening_id
		WHERE %s
		ORDER BY b.tanggal ASC, b.created_at ASC
		LIMIT $%d OFFSET $%d
	`, whereClause, pos, pos+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal list BKU: %w", err)
	}
	defer rows.Close()

	out := []models.BukuKasUmum{}
	for rows.Next() {
		var b models.BukuKasUmum
		if err := rows.Scan(
			&b.ID, &b.SchoolID, &b.AcademicPeriodID, &b.Tanggal, &b.NoBukti, &b.Uraian,
			&b.KodeRekeningID, &b.KodeRekeningKode, &b.RekeningBankID,
			&b.Debet, &b.Kredit, &b.SaldoBerjalan,
			&b.ReferensiTabel, &b.ReferensiID,
			&b.IsLocked, &b.IsRekonsiliasi, &b.CreatedBy, &b.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("gagal scan BKU: %w", err)
		}
		out = append(out, b)
	}
	return out, total, rows.Err()
}
