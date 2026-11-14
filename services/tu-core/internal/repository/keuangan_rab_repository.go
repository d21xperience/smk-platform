package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
)

// =====================================================================
// RABRepository — interface CRUD + realisasi
// =====================================================================
type RABRepository interface {
	// Header & Detail (create atomik)
	CreateRab(ctx context.Context, h *models.RabHeader, items []models.RabDetail) error
	GetRabByID(ctx context.Context, schoolID, id string) (*models.RabHeader, error)
	GetRabByTahun(ctx context.Context, schoolID, periodID string, tahun int) (*models.RabHeader, error)
	ListRab(ctx context.Context, filter models.RabFilter) ([]models.RabHeader, int, error)
	UpdateRab(ctx context.Context, h *models.RabHeader, items []models.RabDetail) error
	SahkanRab(ctx context.Context, schoolID, id, disahkanOleh, catatan string) error
	DeleteRab(ctx context.Context, schoolID, id string) error

	// Realisasi — compute on-the-fly dari BKU
	GetRealisasiByTahun(ctx context.Context, schoolID, periodID string, tahun, sampaiBulan int) (*models.RabRealisasiReport, error)

	// Helper — cek sisa anggaran untuk 1 pos (dipakai PengeluaranService)
	GetSisaAnggaran(ctx context.Context, rabDetailID string) (pagu, realisasi, sisa float64, err error)
}

type pgxRABRepository struct {
	db *pgxpool.Pool
}

func NewRABRepository(db *pgxpool.Pool) RABRepository {
	return &pgxRABRepository{db: db}
}

// =====================================================================
// CreateRab — atomik insert header + details + hitung total
// =====================================================================
func (r *pgxRABRepository) CreateRab(
	ctx context.Context, h *models.RabHeader, items []models.RabDetail,
) error {
	if len(items) == 0 {
		return errors.New("item RAB tidak boleh kosong")
	}

	// Hitung total dari items
	var totalPendapatan, totalBelanja float64
	for _, d := range items {
		switch d.Tipe {
		case models.TipePendapatan:
			totalPendapatan += d.Subtotal
		case models.TipeBelanja:
			totalBelanja += d.Subtotal
		}
	}
	h.TotalPendapatan = totalPendapatan
	h.TotalBelanja = totalBelanja
	h.SurplusDefisit = totalPendapatan - totalBelanja

	if totalBelanja > totalPendapatan {
		return fmt.Errorf("total belanja (%.2f) melebihi total pendapatan (%.2f)",
			totalBelanja, totalPendapatan)
	}

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
		INSERT INTO rab_header (
			id, school_id, academic_period_id, tahun_anggaran,
			total_pendapatan, total_belanja, surplus_defisit, status,
			disusun_oleh, tanggal_penyusunan, catatan
		) VALUES (
			gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7,
			NULLIF($8,''), CURRENT_DATE, NULLIF($9,'')
		)
		RETURNING id, created_at, updated_at
	`, h.SchoolID, h.AcademicPeriodID, h.TahunAnggaran,
		h.TotalPendapatan, h.TotalBelanja, h.SurplusDefisit, h.Status,
		h.DisusunOleh, h.Catatan,
	).Scan(&h.ID, &h.CreatedAt, &h.UpdatedAt)
	if err != nil {
		return fmt.Errorf("gagal insert rab_header: %w", err)
	}

	// 2. Insert details
	for i := range items {
		d := &items[i]
		d.RabHeaderID = h.ID
		err = tx.QueryRow(ctx, `
			INSERT INTO rab_detail (
				id, rab_header_id, kode_rekening_id, uraian,
				volume, satuan, harga_satuan, subtotal, sumber_dana
			) VALUES (
				gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, NULLIF($8,'')
			)
			RETURNING id, created_at, updated_at
		`, h.ID, d.KodeRekeningID, d.Uraian,
			d.Volume, d.Satuan, d.HargaSatuan, d.Subtotal, d.SumberDana,
		).Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt)
		if err != nil {
			return fmt.Errorf("gagal insert rab_detail (kode %s): %w", d.KodeRekeningKode, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("gagal commit: %w", err)
	}
	committed = true
	return nil
}

// =====================================================================
// GetRabByID — header + details
// =====================================================================
func (r *pgxRABRepository) GetRabByID(
	ctx context.Context, schoolID, id string,
) (*models.RabHeader, error) {
	row := r.db.QueryRow(ctx, `
		SELECT h.id, h.school_id, h.academic_period_id, h.tahun_anggaran,
		       h.total_pendapatan, h.total_belanja, h.surplus_defisit, h.status,
		       COALESCE(h.disusun_oleh::text,''), COALESCE(u1.username,''),
		       COALESCE(h.tanggal_penyusunan::text,''),
		       COALESCE(h.disahkan_oleh::text,''), COALESCE(u2.username,''),
		       COALESCE(h.tanggal_pengesahan::text,''),
		       COALESCE(h.catatan,''), h.created_at, h.updated_at
		FROM rab_header h
		LEFT JOIN user_login u1 ON u1.user_id = h.disusun_oleh
		LEFT JOIN user_login u2 ON u2.user_id = h.disahkan_oleh
		WHERE h.school_id = $1 AND h.id = $2
	`, schoolID, id)

	var h models.RabHeader
	err := row.Scan(
		&h.ID, &h.SchoolID, &h.AcademicPeriodID, &h.TahunAnggaran,
		&h.TotalPendapatan, &h.TotalBelanja, &h.SurplusDefisit, &h.Status,
		&h.DisusunOleh, &h.DisusunOlehNama, &h.TanggalPenyusunan,
		&h.DisahkanOleh, &h.DisahkanOlehNama, &h.TanggalPengesahan,
		&h.Catatan, &h.CreatedAt, &h.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTransaksiKeuanganNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("gagal get rab_header: %w", err)
	}

	// Load details
	details, err := r.loadDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	h.Details = details

	return &h, nil
}

// GetRabByTahun — shortcut lookup 1 RAB per school×period×tahun
func (r *pgxRABRepository) GetRabByTahun(
	ctx context.Context, schoolID, periodID string, tahun int,
) (*models.RabHeader, error) {
	var id string
	err := r.db.QueryRow(ctx, `
		SELECT id::text FROM rab_header
		WHERE school_id = $1 AND academic_period_id = $2 AND tahun_anggaran = $3
		LIMIT 1
	`, schoolID, periodID, tahun).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTransaksiKeuanganNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("gagal lookup rab by tahun: %w", err)
	}
	return r.GetRabByID(ctx, schoolID, id)
}

// =====================================================================
// loadDetails — private helper
// =====================================================================
func (r *pgxRABRepository) loadDetails(
	ctx context.Context, headerID string,
) ([]models.RabDetail, error) {
	rows, err := r.db.Query(ctx, `
		SELECT d.id, d.rab_header_id, d.kode_rekening_id,
		       COALESCE(k.kode,''), COALESCE(k.nama,''), COALESCE(k.tipe,''),
		       d.uraian, d.volume, d.satuan, d.harga_satuan, d.subtotal,
		       COALESCE(d.sumber_dana,''), d.created_at, d.updated_at
		FROM rab_detail d
		LEFT JOIN kode_rekening k ON k.id = d.kode_rekening_id
		WHERE d.rab_header_id = $1
		ORDER BY k.kode ASC
	`, headerID)
	if err != nil {
		return nil, fmt.Errorf("gagal load details: %w", err)
	}
	defer rows.Close()

	out := []models.RabDetail{}
	for rows.Next() {
		var d models.RabDetail
		if err := rows.Scan(
			&d.ID, &d.RabHeaderID, &d.KodeRekeningID,
			&d.KodeRekeningKode, &d.KodeRekeningNama, &d.Tipe,
			&d.Uraian, &d.Volume, &d.Satuan, &d.HargaSatuan, &d.Subtotal,
			&d.SumberDana, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("gagal scan detail: %w", err)
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// =====================================================================
// GetSisaAnggaran — untuk budget check di PengeluaranService
// =====================================================================
func (r *pgxRABRepository) GetSisaAnggaran(
	ctx context.Context, rabDetailID string,
) (pagu, realisasi, sisa float64, err error) {
	// 1. Pagu dari rab_detail
	err = r.db.QueryRow(ctx,
		`SELECT subtotal FROM rab_detail WHERE id = $1`, rabDetailID,
	).Scan(&pagu)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, 0, 0, ErrTransaksiKeuanganNotFound
	}
	if err != nil {
		return 0, 0, 0, fmt.Errorf("gagal get pagu: %w", err)
	}

	// 2. Realisasi = SUM(kredit) dari BKU dengan kode_rekening yang sama
	//    yang referensi_tabel='pengeluaran' di tahun yang sama
	//    (kita asumsikan tahun_anggaran dari tanggal BKU)
	var rekeningID string
	if err := r.db.QueryRow(ctx,
		`SELECT kode_rekening_id::text FROM rab_detail WHERE id = $1`, rabDetailID,
	).Scan(&rekeningID); err != nil {
		return 0, 0, 0, fmt.Errorf("gagal get rekening_id: %w", err)
	}

	err = r.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(b.kredit), 0)
		FROM buku_kas_umum b
		JOIN pengeluaran p ON p.id = b.referensi_id AND b.referensi_tabel = 'pengeluaran'
		WHERE b.kode_rekening_id = $1
		  AND p.status IN ('SIAP_CAIR', 'CAIR')
	`, rekeningID).Scan(&realisasi)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("gagal hitung realisasi: %w", err)
	}

	sisa = pagu - realisasi
	return pagu, realisasi, sisa, nil
}

// =====================================================================
// normalizeListQuery — helper pagination (shared)
// =====================================================================
func buildRabWhereClause(filter models.RabFilter) (string, []interface{}, int) {
	where := []string{"school_id = $1"}
	args := []interface{}{filter.SchoolID}
	pos := 2

	if filter.AcademicPeriodID != "" {
		where = append(where, fmt.Sprintf("academic_period_id = $%d", pos))
		args = append(args, filter.AcademicPeriodID)
		pos++
	}
	if filter.TahunAnggaran > 0 {
		where = append(where, fmt.Sprintf("tahun_anggaran = $%d", pos))
		args = append(args, filter.TahunAnggaran)
		pos++
	}
	if filter.Status != "" {
		where = append(where, fmt.Sprintf("status = $%d", pos))
		args = append(args, filter.Status)
		pos++
	}
	return strings.Join(where, " AND "), args, pos
}
