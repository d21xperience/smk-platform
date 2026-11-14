package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
)

// =====================================================================
// RekonsiliasiRepository — interface CRUD + query untuk auto-match
// =====================================================================
type RekonsiliasiRepository interface {
	CreateRekonsiliasi(
		ctx context.Context,
		h *models.RekonsiliasiBank,
		details []models.RekonsiliasiDetail,
		unmatched []models.RekonsiliasiUnmatched,
	) error
	GetRekonsiliasiByID(ctx context.Context, schoolID, id string) (*models.RekonsiliasiBank, error)
	ListRekonsiliasi(ctx context.Context, filter models.RekonsiliasiFilter) ([]models.RekonsiliasiBank, int, error)
	DeleteRekonsiliasi(ctx context.Context, schoolID, id string) error
	ExistsRekonsiliasi(ctx context.Context, schoolID, periodID, rekeningID string, bulan, tahun int) (bool, error)
	GetBkuForPeriod(ctx context.Context, schoolID, periodID, rekeningID string, bulan, tahun int) ([]models.BkuForMatch, error)
	GetSaldoBkuByPeriode(ctx context.Context, schoolID, periodID, rekeningID string, bulan, tahun int) (float64, error)
	// HasBalanceRekonsiliasi — cek apakah ada minimal 1 rekon BALANCE di bulan tsb (FASE 4B)
	HasBalanceRekonsiliasi(ctx context.Context, schoolID, periodID string, bulan, tahun int) (bool, error)
}

type pgxRekonsiliasiRepository struct {
	db *pgxpool.Pool
}

func NewRekonsiliasiRepository(db *pgxpool.Pool) RekonsiliasiRepository {
	return &pgxRekonsiliasiRepository{db: db}
}

// =====================================================================
// CreateRekonsiliasi — atomik insert header + detail + unmatched
// =====================================================================
func (r *pgxRekonsiliasiRepository) CreateRekonsiliasi(
	ctx context.Context,
	h *models.RekonsiliasiBank,
	details []models.RekonsiliasiDetail,
	unmatched []models.RekonsiliasiUnmatched,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("gagal begin tx: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	// 1. Insert header
	err = tx.QueryRow(ctx, `
		INSERT INTO rekonsiliasi_bank (
			id, school_id, academic_period_id, rekening_bank_id,
			periode_bulan, periode_tahun,
			saldo_buku, saldo_bank, selisih,
			jumlah_matched, jumlah_unmatched_bank, jumlah_unmatched_buku,
			saldo_rekonsiliasi_bank, saldo_rekonsiliasi_buku,
			status, file_mutasi_url, catatan, created_by
		) VALUES (
			gen_random_uuid(), $1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, $11, $12, $13,
			$14, NULLIF($15,''), NULLIF($16,''), NULLIF($17,'')
		)
		RETURNING id, created_at
	`, h.SchoolID, h.AcademicPeriodID, h.RekeningBankID,
		h.PeriodeBulan, h.PeriodeTahun,
		h.SaldoBuku, h.SaldoBank, h.Selisih,
		h.JumlahMatched, h.JumlahUnmatchedBank, h.JumlahUnmatchedBuku,
		h.SaldoRekonsiliasiBank, h.SaldoRekonsiliasiBuku,
		h.Status, h.FileMutasiURL, h.Catatan, h.CreatedBy,
	).Scan(&h.ID, &h.CreatedAt)
	if err != nil {
		return fmt.Errorf("gagal insert rekonsiliasi_bank: %w", err)
	}

	// 2. Insert details (matched pairs)
	for i := range details {
		d := &details[i]
		d.RekonsiliasiID = h.ID
		err = tx.QueryRow(ctx, `
			INSERT INTO rekonsiliasi_detail (
				id, rekonsiliasi_id,
				mutasi_tanggal, mutasi_keterangan, mutasi_debet, mutasi_kredit,
				mutasi_saldo, mutasi_referensi,
				buku_kas_umum_id, selisih_hari
			) VALUES (
				gen_random_uuid(), $1, $2::date, NULLIF($3,''), $4, $5,
				$6, NULLIF($7,''), $8, $9
			)
			RETURNING id, created_at
		`, h.ID, d.MutasiTanggal, d.MutasiKeterangan, d.MutasiDebet, d.MutasiKredit,
			d.MutasiSaldo, d.MutasiReferensi, d.BukuKasUmumID, d.SelisihHari,
		).Scan(&d.ID, &d.CreatedAt)
		if err != nil {
			return fmt.Errorf("gagal insert rekonsiliasi_detail: %w", err)
		}
	}

	// 3. Insert unmatched items
	for i := range unmatched {
		u := &unmatched[i]
		u.RekonsiliasiID = h.ID
		err = tx.QueryRow(ctx, `
			INSERT INTO rekonsiliasi_unmatched (
				id, rekonsiliasi_id, tipe, data_json, status_tindak_lanjut
			) VALUES (
				gen_random_uuid(), $1, $2, $3::jsonb, $4
			)
			RETURNING id, created_at
		`, h.ID, u.Tipe, u.DataJSON, u.StatusTindakLanjut,
		).Scan(&u.ID, &u.CreatedAt)
		if err != nil {
			return fmt.Errorf("gagal insert rekonsiliasi_unmatched: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("gagal commit: %w", err)
	}
	committed = true

	h.Details = details
	h.Unmatched = unmatched
	return nil
}

// =====================================================================
// DeleteRekonsiliasi — hapus header (+CASCADE ke detail & unmatched)
// =====================================================================
func (r *pgxRekonsiliasiRepository) DeleteRekonsiliasi(
	ctx context.Context, schoolID, id string,
) error {
	ct, err := r.db.Exec(ctx,
		`DELETE FROM rekonsiliasi_bank WHERE school_id = $1 AND id = $2`,
		schoolID, id,
	)
	if err != nil {
		return fmt.Errorf("gagal delete rekonsiliasi: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrTransaksiKeuanganNotFound
	}
	return nil
}

// =====================================================================
// ExistsRekonsiliasi — cek duplikasi sebelum create
// =====================================================================
func (r *pgxRekonsiliasiRepository) ExistsRekonsiliasi(
	ctx context.Context, schoolID, periodID, rekeningID string, bulan, tahun int,
) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM rekonsiliasi_bank
			WHERE school_id = $1 AND academic_period_id = $2
			  AND rekening_bank_id = $3
			  AND periode_bulan = $4 AND periode_tahun = $5
		)
	`, schoolID, periodID, rekeningID, bulan, tahun).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("gagal check exists: %w", err)
	}
	return exists, nil
}

// =====================================================================
// Helper — build where clause untuk List
// =====================================================================
func buildRekonsiliasiWhere(filter models.RekonsiliasiFilter) (string, []interface{}, int) {
	where := []string{"r.school_id = $1"}
	args := []interface{}{filter.SchoolID}
	pos := 2

	if filter.AcademicPeriodID != "" {
		where = append(where, fmt.Sprintf("r.academic_period_id = $%d", pos))
		args = append(args, filter.AcademicPeriodID)
		pos++
	}
	if filter.RekeningBankID != "" {
		where = append(where, fmt.Sprintf("r.rekening_bank_id = $%d", pos))
		args = append(args, filter.RekeningBankID)
		pos++
	}
	if filter.PeriodeTahun > 0 {
		where = append(where, fmt.Sprintf("r.periode_tahun = $%d", pos))
		args = append(args, filter.PeriodeTahun)
		pos++
	}
	if filter.Status != "" {
		where = append(where, fmt.Sprintf("r.status = $%d", pos))
		args = append(args, filter.Status)
		pos++
	}
	return strings.Join(where, " AND "), args, pos
}

// =====================================================================
// internal helper — scan unmatched data_json safely
// =====================================================================
func normalizeJSONBString(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "{}"
	}
	return s
}

