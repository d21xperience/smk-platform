package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
)

// =====================================================================
// ListRab — list header RAB (tanpa details, untuk performa)
// =====================================================================
func (r *pgxRABRepository) ListRab(
	ctx context.Context, filter models.RabFilter,
) ([]models.RabHeader, int, error) {
	limit, offset := normalizePagination(filter.Limit, filter.Offset)
	whereClause, args, pos := buildRabWhereClause(filter)

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM rab_header WHERE "+whereClause, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("gagal count rab: %w", err)
	}

	listQuery := fmt.Sprintf(`
		SELECT id, school_id, academic_period_id, tahun_anggaran,
		       total_pendapatan, total_belanja, surplus_defisit, status,
		       created_at, updated_at
		FROM rab_header
		WHERE %s
		ORDER BY tahun_anggaran DESC, created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, pos, pos+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal list rab: %w", err)
	}
	defer rows.Close()

	out := []models.RabHeader{}
	for rows.Next() {
		var h models.RabHeader
		if err := rows.Scan(
			&h.ID, &h.SchoolID, &h.AcademicPeriodID, &h.TahunAnggaran,
			&h.TotalPendapatan, &h.TotalBelanja, &h.SurplusDefisit, &h.Status,
			&h.CreatedAt, &h.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("gagal scan rab: %w", err)
		}
		out = append(out, h)
	}
	return out, total, rows.Err()
}

// =====================================================================
// UpdateRab — replace header + semua details (atomik)
// Hanya boleh update jika status DRAFT atau DIREVISI
// =====================================================================
func (r *pgxRABRepository) UpdateRab(
	ctx context.Context, h *models.RabHeader, items []models.RabDetail,
) error {
	if len(items) == 0 {
		return errors.New("item RAB tidak boleh kosong")
	}

	// Cek status existing
	var currentStatus string
	err := r.db.QueryRow(ctx,
		`SELECT status FROM rab_header WHERE school_id = $1 AND id = $2`,
		h.SchoolID, h.ID,
	).Scan(&currentStatus)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrTransaksiKeuanganNotFound
	}
	if err != nil {
		return fmt.Errorf("gagal cek status: %w", err)
	}
	if currentStatus != models.RabDraft && currentStatus != models.RabDirevisi {
		return fmt.Errorf("RAB dengan status %s tidak dapat diubah (hanya DRAFT/DIREVISI)", currentStatus)
	}

	// Recompute totals
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

	// Update header
	if _, err := tx.Exec(ctx, `
		UPDATE rab_header
		SET total_pendapatan = $1, total_belanja = $2, surplus_defisit = $3,
		    catatan = NULLIF($4,''), updated_at = NOW()
		WHERE school_id = $5 AND id = $6
	`, h.TotalPendapatan, h.TotalBelanja, h.SurplusDefisit, h.Catatan,
		h.SchoolID, h.ID,
	); err != nil {
		return fmt.Errorf("gagal update header: %w", err)
	}

	// Delete existing details
	if _, err := tx.Exec(ctx, `DELETE FROM rab_detail WHERE rab_header_id = $1`, h.ID); err != nil {
		return fmt.Errorf("gagal hapus detail lama: %w", err)
	}

	// Insert new details
	for i := range items {
		d := &items[i]
		d.RabHeaderID = h.ID
		err := tx.QueryRow(ctx, `
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
			return fmt.Errorf("gagal insert detail baru (kode %s): %w", d.KodeRekeningKode, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("gagal commit: %w", err)
	}
	committed = true
	return nil
}

// =====================================================================
// SahkanRab — transisi status ke DISAHKAN (hanya dari DRAFT/DIAJUKAN)
// =====================================================================
func (r *pgxRABRepository) SahkanRab(
	ctx context.Context, schoolID, id, disahkanOleh, catatan string,
) error {
	ct, err := r.db.Exec(ctx, `
		UPDATE rab_header
		SET status = 'DISAHKAN',
		    disahkan_oleh = NULLIF($1,''),
		    tanggal_pengesahan = CURRENT_DATE,
		    catatan = COALESCE(NULLIF($2,''), catatan),
		    updated_at = NOW()
		WHERE school_id = $3 AND id = $4
		  AND status IN ('DRAFT', 'DIAJUKAN')
	`, disahkanOleh, catatan, schoolID, id)
	if err != nil {
		return fmt.Errorf("gagal sahkan RAB: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("RAB tidak ditemukan atau sudah disahkan (hanya DRAFT/DIAJUKAN yang bisa disahkan)")
	}
	return nil
}

// =====================================================================
// DeleteRab — hapus header + details (CASCADE)
// Hanya boleh hapus jika status DRAFT atau DIREVISI
// =====================================================================
func (r *pgxRABRepository) DeleteRab(
	ctx context.Context, schoolID, id string,
) error {
	ct, err := r.db.Exec(ctx, `
		DELETE FROM rab_header
		WHERE school_id = $1 AND id = $2
		  AND status IN ('DRAFT', 'DIREVISI')
	`, schoolID, id)
	if err != nil {
		return fmt.Errorf("gagal delete RAB: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("RAB tidak ditemukan atau berstatus DISAHKAN (tidak dapat dihapus)")
	}
	return nil
}
