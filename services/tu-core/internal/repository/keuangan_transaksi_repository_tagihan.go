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
// TAGIHAN SISWA
// =====================================================================
func (r *pgxKeuanganTransaksiRepository) CreateTagihan(
	ctx context.Context, tx pgx.Tx, t *models.TagihanSiswa,
) error {
	query := `
		INSERT INTO tagihan_siswa (
			id, school_id, academic_period_id, peserta_didik_id, nis, kelas_nama,
			jenis_iuran_id, periode_bulan, periode_tahun, nominal, total_dibayar,
			status, tanggal_jatuh_tempo, created_by
		) VALUES (
			gen_random_uuid(), $1, $2, $3, $4, NULLIF($5,''), $6, $7, $8, $9, $10,
			$11, NULLIF($12,'')::date, NULLIF($13,'')
		)
		RETURNING id, created_at, updated_at
	`
	row := tx.QueryRow(ctx, query,
		t.SchoolID, t.AcademicPeriodID, t.PesertaDidikID, t.NIS, t.KelasNama,
		t.JenisIuranID, t.PeriodeBulan, t.PeriodeTahun, t.Nominal, t.TotalDibayar,
		t.Status, t.TanggalJatuhTempo, t.CreatedBy,
	)
	if err := row.Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return fmt.Errorf("gagal insert tagihan_siswa: %w", err)
	}
	return nil
}

func (r *pgxKeuanganTransaksiRepository) GetTagihanByID(
	ctx context.Context, schoolID, id string,
) (*models.TagihanSiswa, error) {
	query := `
		SELECT t.id, t.school_id, t.academic_period_id, t.peserta_didik_id, t.nis,
		       COALESCE(t.kelas_nama,''), t.jenis_iuran_id, COALESCE(j.nama,''),
		       t.periode_bulan, t.periode_tahun, t.nominal, t.total_dibayar, t.status,
		       COALESCE(t.tanggal_jatuh_tempo::text,''), COALESCE(t.tanggal_bayar::text,''),
		       COALESCE(t.created_by::text,''), t.created_at, t.updated_at
		FROM tagihan_siswa t
		LEFT JOIN jenis_iuran j ON j.id = t.jenis_iuran_id
		WHERE t.school_id = $1 AND t.id = $2
	`
	var t models.TagihanSiswa
	err := r.db.QueryRow(ctx, query, schoolID, id).Scan(
		&t.ID, &t.SchoolID, &t.AcademicPeriodID, &t.PesertaDidikID, &t.NIS,
		&t.KelasNama, &t.JenisIuranID, &t.JenisIuranNama,
		&t.PeriodeBulan, &t.PeriodeTahun, &t.Nominal, &t.TotalDibayar, &t.Status,
		&t.TanggalJatuhTempo, &t.TanggalBayar, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTransaksiKeuanganNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("gagal get tagihan: %w", err)
	}
	return &t, nil
}

func (r *pgxKeuanganTransaksiRepository) ListTagihan(
	ctx context.Context, filter models.TagihanFilter,
) ([]models.TagihanSiswa, int, error) {
	limit, offset := normalizePagination(filter.Limit, filter.Offset)

	where := []string{"t.school_id = $1"}
	args := []interface{}{filter.SchoolID}
	pos := 2

	if filter.AcademicPeriodID != "" {
		where = append(where, fmt.Sprintf("t.academic_period_id = $%d", pos))
		args = append(args, filter.AcademicPeriodID)
		pos++
	}
	if filter.NIS != "" {
		where = append(where, fmt.Sprintf("t.nis = $%d", pos))
		args = append(args, filter.NIS)
		pos++
	}
	if filter.JenisIuranID != "" {
		where = append(where, fmt.Sprintf("t.jenis_iuran_id = $%d", pos))
		args = append(args, filter.JenisIuranID)
		pos++
	}
	if filter.Status != "" {
		where = append(where, fmt.Sprintf("t.status = $%d", pos))
		args = append(args, filter.Status)
		pos++
	}
	if filter.Bulan > 0 {
		where = append(where, fmt.Sprintf("t.periode_bulan = $%d", pos))
		args = append(args, filter.Bulan)
		pos++
	}
	if filter.Tahun > 0 {
		where = append(where, fmt.Sprintf("t.periode_tahun = $%d", pos))
		args = append(args, filter.Tahun)
		pos++
	}

	whereClause := strings.Join(where, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM tagihan_siswa t WHERE "+whereClause, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("gagal count tagihan: %w", err)
	}

	listQuery := fmt.Sprintf(`
		SELECT t.id, t.school_id, t.academic_period_id, t.peserta_didik_id, t.nis,
		       COALESCE(t.kelas_nama,''), t.jenis_iuran_id, COALESCE(j.nama,''),
		       t.periode_bulan, t.periode_tahun, t.nominal, t.total_dibayar, t.status,
		       COALESCE(t.tanggal_jatuh_tempo::text,''), COALESCE(t.tanggal_bayar::text,''),
		       COALESCE(t.created_by::text,''), t.created_at, t.updated_at
		FROM tagihan_siswa t
		LEFT JOIN jenis_iuran j ON j.id = t.jenis_iuran_id
		WHERE %s
		ORDER BY t.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, pos, pos+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal list tagihan: %w", err)
	}
	defer rows.Close()

	out := []models.TagihanSiswa{}
	for rows.Next() {
		var t models.TagihanSiswa
		if err := rows.Scan(
			&t.ID, &t.SchoolID, &t.AcademicPeriodID, &t.PesertaDidikID, &t.NIS,
			&t.KelasNama, &t.JenisIuranID, &t.JenisIuranNama,
			&t.PeriodeBulan, &t.PeriodeTahun, &t.Nominal, &t.TotalDibayar, &t.Status,
			&t.TanggalJatuhTempo, &t.TanggalBayar, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("gagal scan tagihan: %w", err)
		}
		out = append(out, t)
	}
	return out, total, rows.Err()
}

func (r *pgxKeuanganTransaksiRepository) UpdateTagihanStatus(
	ctx context.Context, tx pgx.Tx,
	id, status string, totalDibayar float64, tanggalBayar *time.Time,
) error {
	query := `
		UPDATE tagihan_siswa
		SET status = $1, total_dibayar = $2, tanggal_bayar = $3, updated_at = NOW()
		WHERE id = $4
	`
	ct, err := tx.Exec(ctx, query, status, totalDibayar, tanggalBayar, id)
	if err != nil {
		return fmt.Errorf("gagal update status tagihan: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrTransaksiKeuanganNotFound
	}
	return nil
}

func (r *pgxKeuanganTransaksiRepository) FindTagihanByNISPeriode(
	ctx context.Context,
	schoolID, periodID, nis, jenisIuranID string, bulan, tahun int,
) (*models.TagihanSiswa, error) {
	query := `
		SELECT id, school_id, academic_period_id, peserta_didik_id, nis,
		       COALESCE(kelas_nama,''), jenis_iuran_id, '' AS jenis_nama,
		       periode_bulan, periode_tahun, nominal, total_dibayar, status,
		       COALESCE(tanggal_jatuh_tempo::text,''), COALESCE(tanggal_bayar::text,''),
		       COALESCE(created_by::text,''), created_at, updated_at
		FROM tagihan_siswa
		WHERE school_id = $1 AND academic_period_id = $2
		  AND nis = $3 AND jenis_iuran_id = $4
		  AND periode_bulan = $5 AND periode_tahun = $6
		LIMIT 1
	`
	var t models.TagihanSiswa
	err := r.db.QueryRow(ctx, query,
		schoolID, periodID, nis, jenisIuranID, bulan, tahun,
	).Scan(
		&t.ID, &t.SchoolID, &t.AcademicPeriodID, &t.PesertaDidikID, &t.NIS,
		&t.KelasNama, &t.JenisIuranID, &t.JenisIuranNama,
		&t.PeriodeBulan, &t.PeriodeTahun, &t.Nominal, &t.TotalDibayar, &t.Status,
		&t.TanggalJatuhTempo, &t.TanggalBayar, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTransaksiKeuanganNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("gagal find tagihan: %w", err)
	}
	return &t, nil
}

// =====================================================================
// PENERIMAAN IURAN
// =====================================================================
func (r *pgxKeuanganTransaksiRepository) CreatePenerimaan(
	ctx context.Context, tx pgx.Tx, p *models.PenerimaanIuran,
) error {
	query := `
		INSERT INTO penerimaan_iuran (
			id, school_id, academic_period_id, no_bkm, tanggal,
			tagihan_siswa_id, peserta_didik_id, nis, jenis_iuran_id,
			periode_bulan, periode_tahun, jumlah, metode_bayar, no_referensi,
			keterangan, status, created_by
		) VALUES (
			gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8,
			$9, $10, $11, $12, NULLIF($13,''), NULLIF($14,''), $15, NULLIF($16,'')
		)
		RETURNING id, created_at, updated_at
	`
	row := tx.QueryRow(ctx, query,
		p.SchoolID, p.AcademicPeriodID, p.NoBKM, p.Tanggal,
		p.TagihanSiswaID, p.PesertaDidikID, p.NIS, p.JenisIuranID,
		p.PeriodeBulan, p.PeriodeTahun, p.Jumlah, p.MetodeBayar, p.NoReferensi,
		p.Keterangan, p.Status, p.CreatedBy,
	)
	if err := row.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return fmt.Errorf("gagal insert penerimaan: %w", err)
	}
	return nil
}

func (r *pgxKeuanganTransaksiRepository) GetPenerimaanByNoBKM(
	ctx context.Context, schoolID, noBKM string,
) (*models.PenerimaanIuran, error) {
	query := `
		SELECT p.id, p.school_id, p.academic_period_id, p.no_bkm, p.tanggal::text,
		       p.tagihan_siswa_id, p.peserta_didik_id, p.nis, p.jenis_iuran_id,
		       COALESCE(j.nama,''), p.periode_bulan, p.periode_tahun, p.jumlah,
		       p.metode_bayar, COALESCE(p.no_referensi,''), COALESCE(p.keterangan,''),
		       p.status, COALESCE(p.created_by::text,''), p.created_at, p.updated_at
		FROM penerimaan_iuran p
		LEFT JOIN jenis_iuran j ON j.id = p.jenis_iuran_id
		WHERE p.school_id = $1 AND p.no_bkm = $2
	`
	var p models.PenerimaanIuran
	err := r.db.QueryRow(ctx, query, schoolID, noBKM).Scan(
		&p.ID, &p.SchoolID, &p.AcademicPeriodID, &p.NoBKM, &p.Tanggal,
		&p.TagihanSiswaID, &p.PesertaDidikID, &p.NIS, &p.JenisIuranID,
		&p.JenisIuranNama, &p.PeriodeBulan, &p.PeriodeTahun, &p.Jumlah,
		&p.MetodeBayar, &p.NoReferensi, &p.Keterangan,
		&p.Status, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTransaksiKeuanganNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("gagal get penerimaan: %w", err)
	}
	return &p, nil
}

func (r *pgxKeuanganTransaksiRepository) ListPenerimaan(
	ctx context.Context, filter models.PenerimaanFilter,
) ([]models.PenerimaanIuran, int, error) {
	limit, offset := normalizePagination(filter.Limit, filter.Offset)

	where := []string{"p.school_id = $1"}
	args := []interface{}{filter.SchoolID}
	pos := 2

	if filter.AcademicPeriodID != "" {
		where = append(where, fmt.Sprintf("p.academic_period_id = $%d", pos))
		args = append(args, filter.AcademicPeriodID)
		pos++
	}
	if filter.NIS != "" {
		where = append(where, fmt.Sprintf("p.nis = $%d", pos))
		args = append(args, filter.NIS)
		pos++
	}
	if filter.JenisIuranID != "" {
		where = append(where, fmt.Sprintf("p.jenis_iuran_id = $%d", pos))
		args = append(args, filter.JenisIuranID)
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
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM penerimaan_iuran p WHERE "+whereClause, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("gagal count penerimaan: %w", err)
	}

	listQuery := fmt.Sprintf(`
		SELECT p.id, p.school_id, p.academic_period_id, p.no_bkm, p.tanggal::text,
		       p.tagihan_siswa_id, p.peserta_didik_id, p.nis, p.jenis_iuran_id,
		       COALESCE(j.nama,''), p.periode_bulan, p.periode_tahun, p.jumlah,
		       p.metode_bayar, COALESCE(p.no_referensi,''), COALESCE(p.keterangan,''),
		       p.status, COALESCE(p.created_by::text,''), p.created_at, p.updated_at
		FROM penerimaan_iuran p
		LEFT JOIN jenis_iuran j ON j.id = p.jenis_iuran_id
		WHERE %s
		ORDER BY p.tanggal DESC, p.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, pos, pos+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal list penerimaan: %w", err)
	}
	defer rows.Close()

	out := []models.PenerimaanIuran{}
	for rows.Next() {
		var p models.PenerimaanIuran
		if err := rows.Scan(
			&p.ID, &p.SchoolID, &p.AcademicPeriodID, &p.NoBKM, &p.Tanggal,
			&p.TagihanSiswaID, &p.PesertaDidikID, &p.NIS, &p.JenisIuranID,
			&p.JenisIuranNama, &p.PeriodeBulan, &p.PeriodeTahun, &p.Jumlah,
			&p.MetodeBayar, &p.NoReferensi, &p.Keterangan,
			&p.Status, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("gagal scan penerimaan: %w", err)
		}
		out = append(out, p)
	}
	return out, total, rows.Err()
}
