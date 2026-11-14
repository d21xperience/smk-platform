package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
)

// =====================================================================
// PENGAJUAN PENGELUARAN (workflow 3 tingkat)
// =====================================================================
func (r *pgxKeuanganTransaksiRepository) CreatePengajuan(
	ctx context.Context, tx pgx.Tx, p *models.PengajuanPengeluaran,
) error {
	query := `
		INSERT INTO pengajuan_pengeluaran (
			id, school_id, academic_period_id, nomor_pengajuan, tanggal,
			pemohon_user_id, keperluan, jumlah, prioritas, status, catatan_pemohon
		) VALUES (
			gen_random_uuid(), $1, $2, $3, $4::date, $5, $6, $7, $8, $9, NULLIF($10,'')
		)
		RETURNING id, created_at, updated_at
	`
	row := tx.QueryRow(ctx, query,
		p.SchoolID, p.AcademicPeriodID, p.NomorPengajuan, p.Tanggal,
		p.PemohonUserID, p.Keperluan, p.Jumlah, p.Prioritas, p.Status, p.CatatanPemohon,
	)
	if err := row.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return fmt.Errorf("gagal insert pengajuan: %w", err)
	}
	return nil
}

func (r *pgxKeuanganTransaksiRepository) GetPengajuanByID(
	ctx context.Context, schoolID, id string,
) (*models.PengajuanPengeluaran, error) {
	query := `
		SELECT p.id, p.school_id, p.academic_period_id, p.nomor_pengajuan, p.tanggal::text,
		       p.pemohon_user_id, COALESCE(u.username,''), p.keperluan, p.jumlah,
		       p.prioritas, p.status,
		       COALESCE(p.catatan_pemohon,''), COALESCE(p.catatan_bendahara,''), COALESCE(p.catatan_kepala,''),
		       COALESCE(p.verified_by::text,''), p.verified_at,
		       COALESCE(p.approved_by::text,''), p.approved_at,
		       COALESCE(p.no_bkk,''), COALESCE(p.pos_anggaran_id::text,''),
		       p.created_at, p.updated_at
		FROM pengajuan_pengeluaran p
		LEFT JOIN user_login u ON u.user_id = p.pemohon_user_id
		WHERE p.school_id = $1 AND p.id = $2
	`
	var p models.PengajuanPengeluaran
	err := r.db.QueryRow(ctx, query, schoolID, id).Scan(
		&p.ID, &p.SchoolID, &p.AcademicPeriodID, &p.NomorPengajuan, &p.Tanggal,
		&p.PemohonUserID, &p.PemohonNama, &p.Keperluan, &p.Jumlah,
		&p.Prioritas, &p.Status,
		&p.CatatanPemohon, &p.CatatanBendahara, &p.CatatanKepala,
		&p.VerifiedBy, &p.VerifiedAt,
		&p.ApprovedBy, &p.ApprovedAt,
		&p.NoBKK, &p.PosAnggaranID,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTransaksiKeuanganNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("gagal get pengajuan: %w", err)
	}
	return &p, nil
}

func (r *pgxKeuanganTransaksiRepository) ListPengajuan(
	ctx context.Context, filter models.PengajuanFilter,
) ([]models.PengajuanPengeluaran, int, error) {
	limit, offset := normalizePagination(filter.Limit, filter.Offset)

	where := []string{"p.school_id = $1"}
	args := []interface{}{filter.SchoolID}
	pos := 2

	if filter.AcademicPeriodID != "" {
		where = append(where, fmt.Sprintf("p.academic_period_id = $%d", pos))
		args = append(args, filter.AcademicPeriodID)
		pos++
	}
	if filter.Status != "" {
		where = append(where, fmt.Sprintf("p.status = $%d", pos))
		args = append(args, filter.Status)
		pos++
	}
	if filter.PemohonUserID != "" {
		where = append(where, fmt.Sprintf("p.pemohon_user_id = $%d", pos))
		args = append(args, filter.PemohonUserID)
		pos++
	}

	whereClause := strings.Join(where, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM pengajuan_pengeluaran p WHERE "+whereClause, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("gagal count pengajuan: %w", err)
	}

	listQuery := fmt.Sprintf(`
		SELECT p.id, p.school_id, p.academic_period_id, p.nomor_pengajuan, p.tanggal::text,
		       p.pemohon_user_id, COALESCE(u.username,''), p.keperluan, p.jumlah,
		       p.prioritas, p.status,
		       COALESCE(p.catatan_pemohon,''), COALESCE(p.catatan_bendahara,''), COALESCE(p.catatan_kepala,''),
		       COALESCE(p.verified_by::text,''), p.verified_at,
		       COALESCE(p.approved_by::text,''), p.approved_at,
		       COALESCE(p.no_bkk,''), COALESCE(p.pos_anggaran_id::text,''),
		       p.created_at, p.updated_at
		FROM pengajuan_pengeluaran p
		LEFT JOIN user_login u ON u.user_id = p.pemohon_user_id
		WHERE %s
		ORDER BY p.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, pos, pos+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal list pengajuan: %w", err)
	}
	defer rows.Close()

	out := []models.PengajuanPengeluaran{}
	for rows.Next() {
		var p models.PengajuanPengeluaran
		if err := rows.Scan(
			&p.ID, &p.SchoolID, &p.AcademicPeriodID, &p.NomorPengajuan, &p.Tanggal,
			&p.PemohonUserID, &p.PemohonNama, &p.Keperluan, &p.Jumlah,
			&p.Prioritas, &p.Status,
			&p.CatatanPemohon, &p.CatatanBendahara, &p.CatatanKepala,
			&p.VerifiedBy, &p.VerifiedAt,
			&p.ApprovedBy, &p.ApprovedAt,
			&p.NoBKK, &p.PosAnggaranID,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("gagal scan pengajuan: %w", err)
		}
		out = append(out, p)
	}
	return out, total, rows.Err()
}

// UpdatePengajuanStatus — ubah status + optional update catatan/verified/approved
// Service layer bertanggung jawab memvalidasi transisi via models.IsValidTransisiPengajuan
func (r *pgxKeuanganTransaksiRepository) UpdatePengajuanStatus(
	ctx context.Context, tx pgx.Tx,
	id, newStatus string,
	verifiedBy, approvedBy *string,
	catatanBendahara, catatanKepala, noBKK *string,
) error {
	query := `
		UPDATE pengajuan_pengeluaran
		SET status = $1,
		    verified_by = COALESCE($2, verified_by),
		    verified_at = CASE WHEN $2::text IS NOT NULL THEN NOW() ELSE verified_at END,
		    approved_by = COALESCE($3, approved_by),
		    approved_at = CASE WHEN $3::text IS NOT NULL THEN NOW() ELSE approved_at END,
		    catatan_bendahara = COALESCE($4, catatan_bendahara),
		    catatan_kepala = COALESCE($5, catatan_kepala),
		    no_bkk = COALESCE($6, no_bkk),
		    updated_at = NOW()
		WHERE id = $7
	`
	ct, err := tx.Exec(ctx, query,
		newStatus, verifiedBy, approvedBy, catatanBendahara, catatanKepala, noBKK, id,
	)
	if err != nil {
		return fmt.Errorf("gagal update status pengajuan: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrTransaksiKeuanganNotFound
	}
	return nil
}

// =====================================================================
// PENGELUARAN
// =====================================================================
func (r *pgxKeuanganTransaksiRepository) CreatePengeluaran(
	ctx context.Context, tx pgx.Tx, p *models.Pengeluaran,
) error {
	query := `
		INSERT INTO pengeluaran (
			id, school_id, academic_period_id, no_bkk, tanggal,
			pengajuan_id, jumlah, keperluan, status, created_by
		) VALUES (
			gen_random_uuid(), $1, $2, $3, $4::date, $5, $6, $7, $8, NULLIF($9,'')
		)
		RETURNING id, created_at, updated_at
	`
	row := tx.QueryRow(ctx, query,
		p.SchoolID, p.AcademicPeriodID, p.NoBKK, p.Tanggal,
		p.PengajuanID, p.Jumlah, p.Keperluan, p.Status, p.CreatedBy,
	)
	if err := row.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return fmt.Errorf("gagal insert pengeluaran: %w", err)
	}
	return nil
}

func (r *pgxKeuanganTransaksiRepository) GetPengeluaranByNoBKK(
	ctx context.Context, schoolID, noBKK string,
) (*models.Pengeluaran, error) {
	query := `
		SELECT p.id, p.school_id, p.academic_period_id, p.no_bkk, p.tanggal::text,
		       p.pengajuan_id, p.jumlah, p.keperluan, p.status,
		       COALESCE(p.created_by::text,''), p.created_at, p.updated_at
		FROM pengeluaran p
		WHERE p.school_id = $1 AND p.no_bkk = $2
	`
	var p models.Pengeluaran
	err := r.db.QueryRow(ctx, query, schoolID, noBKK).Scan(
		&p.ID, &p.SchoolID, &p.AcademicPeriodID, &p.NoBKK, &p.Tanggal,
		&p.PengajuanID, &p.Jumlah, &p.Keperluan, &p.Status,
		&p.CreatedBy, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTransaksiKeuanganNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("gagal get pengeluaran: %w", err)
	}
	return &p, nil
}

func (r *pgxKeuanganTransaksiRepository) ListPengeluaran(
	ctx context.Context, filter models.PengeluaranFilter,
) ([]models.Pengeluaran, int, error) {
	limit, offset := normalizePagination(filter.Limit, filter.Offset)

	where := []string{"p.school_id = $1"}
	args := []interface{}{filter.SchoolID}
	pos := 2

	if filter.AcademicPeriodID != "" {
		where = append(where, fmt.Sprintf("p.academic_period_id = $%d", pos))
		args = append(args, filter.AcademicPeriodID)
		pos++
	}
	if filter.Status != "" {
		where = append(where, fmt.Sprintf("p.status = $%d", pos))
		args = append(args, filter.Status)
		pos++
	}
	if filter.TanggalDari != "" {
		where = append(where, fmt.Sprintf("p.tanggal >= $%d::date", pos))
		args = append(args, filter.TanggalDari)
		pos++
	}
	if filter.TanggalSampai != "" {
		where = append(where, fmt.Sprintf("p.tanggal <= $%d::date", pos))
		args = append(args, filter.TanggalSampai)
		pos++
	}

	whereClause := strings.Join(where, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM pengeluaran p WHERE "+whereClause, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("gagal count pengeluaran: %w", err)
	}

	listQuery := fmt.Sprintf(`
		SELECT p.id, p.school_id, p.academic_period_id, p.no_bkk, p.tanggal::text,
		       p.pengajuan_id, p.jumlah, p.keperluan, p.status,
		       COALESCE(p.created_by::text,''), p.created_at, p.updated_at
		FROM pengeluaran p
		WHERE %s
		ORDER BY p.tanggal DESC, p.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, pos, pos+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal list pengeluaran: %w", err)
	}
	defer rows.Close()

	out := []models.Pengeluaran{}
	for rows.Next() {
		var p models.Pengeluaran
		if err := rows.Scan(
			&p.ID, &p.SchoolID, &p.AcademicPeriodID, &p.NoBKK, &p.Tanggal,
			&p.PengajuanID, &p.Jumlah, &p.Keperluan, &p.Status,
			&p.CreatedBy, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("gagal scan pengeluaran: %w", err)
		}
		out = append(out, p)
	}
	return out, total, rows.Err()
}
