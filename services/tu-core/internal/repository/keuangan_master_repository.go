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
// KeuanganMasterRepository — CRUD master data keuangan
// =====================================================================
type KeuanganMasterRepository interface {
    // Jenis Iuran
    CreateJenisIuran(ctx context.Context, j *models.JenisIuran) error
    UpdateJenisIuran(ctx context.Context, j *models.JenisIuran) error
    DeleteJenisIuran(ctx context.Context, schoolID, id string) error
    GetJenisIuranByID(ctx context.Context, schoolID, id string) (*models.JenisIuran, error)
    ListJenisIuran(ctx context.Context, filter models.JenisIuranFilter) ([]models.JenisIuran, int, error)

    // Tarif Iuran
    CreateTarifIuran(ctx context.Context, t *models.TarifIuran) error
    UpdateTarifIuran(ctx context.Context, t *models.TarifIuran) error
    DeleteTarifIuran(ctx context.Context, schoolID, id string) error
    GetTarifIuranByID(ctx context.Context, schoolID, id string) (*models.TarifIuran, error)
    ListTarifIuran(ctx context.Context, filter models.TarifIuranFilter) ([]models.TarifIuran, int, error)

    // Kode Rekening
    CreateKodeRekening(ctx context.Context, k *models.KodeRekening) error
    UpdateKodeRekening(ctx context.Context, k *models.KodeRekening) error
    DeleteKodeRekening(ctx context.Context, schoolID, id string) error
    GetKodeRekeningByID(ctx context.Context, schoolID, id string) (*models.KodeRekening, error)
    ListKodeRekening(ctx context.Context, filter models.KodeRekeningFilter) ([]models.KodeRekening, int, error)

    // Rekening Bank
    CreateRekeningBank(ctx context.Context, r *models.RekeningBank) error
    UpdateRekeningBank(ctx context.Context, r *models.RekeningBank) error
    DeleteRekeningBank(ctx context.Context, schoolID, id string) error
    GetRekeningBankByID(ctx context.Context, schoolID, id string) (*models.RekeningBank, error)
    ListRekeningBank(ctx context.Context, filter models.RekeningBankFilter) ([]models.RekeningBank, int, error)
}

type pgxKeuanganMasterRepository struct {
    db *pgxpool.Pool
}

func NewKeuanganMasterRepository(db *pgxpool.Pool) KeuanganMasterRepository {
    return &pgxKeuanganMasterRepository{db: db}
}

// =====================================================================
// Helpers
// =====================================================================
var ErrKeuanganMasterNotFound = errors.New("data master keuangan tidak ditemukan")

func normalizePagination(limit, offset int) (int, int) {
    if limit <= 0 || limit > 200 {
        limit = 50
    }
    if offset < 0 {
        offset = 0
    }
    return limit, offset
}

// =====================================================================
// JENIS IURAN
// =====================================================================
func (r *pgxKeuanganMasterRepository) CreateJenisIuran(ctx context.Context, j *models.JenisIuran) error {
    query := `
        INSERT INTO jenis_iuran
            (id, school_id, kode, nama, kategori, is_wajib, is_aktif, keterangan)
        VALUES
            (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at
    `
    row := r.db.QueryRow(ctx, query,
        j.SchoolID, j.Kode, j.Nama, j.Kategori, j.IsWajib, j.IsAktif, j.Keterangan,
    )
    if err := row.Scan(&j.ID, &j.CreatedAt, &j.UpdatedAt); err != nil {
        return fmt.Errorf("gagal insert jenis_iuran: %w", err)
    }
    return nil
}

func (r *pgxKeuanganMasterRepository) UpdateJenisIuran(ctx context.Context, j *models.JenisIuran) error {
    query := `
        UPDATE jenis_iuran
        SET kode=$1, nama=$2, kategori=$3, is_wajib=$4, is_aktif=$5, keterangan=$6, updated_at=NOW()
        WHERE school_id=$7 AND id=$8
        RETURNING updated_at
    `
    row := r.db.QueryRow(ctx, query,
        j.Kode, j.Nama, j.Kategori, j.IsWajib, j.IsAktif, j.Keterangan,
        j.SchoolID, j.ID,
    )
    if err := row.Scan(&j.UpdatedAt); err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return ErrKeuanganMasterNotFound
        }
        return fmt.Errorf("gagal update jenis_iuran: %w", err)
    }
    return nil
}

func (r *pgxKeuanganMasterRepository) DeleteJenisIuran(ctx context.Context, schoolID, id string) error {
    ct, err := r.db.Exec(ctx, `DELETE FROM jenis_iuran WHERE school_id=$1 AND id=$2`, schoolID, id)
    if err != nil {
        return fmt.Errorf("gagal delete jenis_iuran: %w", err)
    }
    if ct.RowsAffected() == 0 {
        return ErrKeuanganMasterNotFound
    }
    return nil
}

func (r *pgxKeuanganMasterRepository) GetJenisIuranByID(ctx context.Context, schoolID, id string) (*models.JenisIuran, error) {
    query := `
        SELECT id, school_id, kode, nama, kategori, is_wajib, is_aktif,
               COALESCE(keterangan, ''), created_at, updated_at
        FROM jenis_iuran
        WHERE school_id=$1 AND id=$2
    `
    var j models.JenisIuran
    err := r.db.QueryRow(ctx, query, schoolID, id).Scan(
        &j.ID, &j.SchoolID, &j.Kode, &j.Nama, &j.Kategori,
        &j.IsWajib, &j.IsAktif, &j.Keterangan, &j.CreatedAt, &j.UpdatedAt,
    )
    if errors.Is(err, pgx.ErrNoRows) {
        return nil, ErrKeuanganMasterNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("gagal get jenis_iuran: %w", err)
    }
    return &j, nil
}

func (r *pgxKeuanganMasterRepository) ListJenisIuran(ctx context.Context, filter models.JenisIuranFilter) ([]models.JenisIuran, int, error) {
    limit, offset := normalizePagination(filter.Limit, filter.Offset)

    where := []string{"school_id = $1"}
    args := []interface{}{filter.SchoolID}
    argPos := 2

    if filter.Kategori != "" {
        where = append(where, fmt.Sprintf("kategori = $%d", argPos))
        args = append(args, filter.Kategori)
        argPos++
    }
    if filter.IsAktif != nil {
        where = append(where, fmt.Sprintf("is_aktif = $%d", argPos))
        args = append(args, *filter.IsAktif)
        argPos++
    }
    if filter.Search != "" {
        where = append(where, fmt.Sprintf("(nama ILIKE $%d OR kode ILIKE $%d)", argPos, argPos))
        args = append(args, "%"+filter.Search+"%")
        argPos++
    }

    whereClause := strings.Join(where, " AND ")

    var total int
    countQuery := "SELECT COUNT(*) FROM jenis_iuran WHERE " + whereClause
    if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
        return nil, 0, fmt.Errorf("gagal count jenis_iuran: %w", err)
    }

    listQuery := fmt.Sprintf(`
        SELECT id, school_id, kode, nama, kategori, is_wajib, is_aktif,
               COALESCE(keterangan, ''), created_at, updated_at
        FROM jenis_iuran
        WHERE %s
        ORDER BY kode ASC
        LIMIT $%d OFFSET $%d
    `, whereClause, argPos, argPos+1)
    args = append(args, limit, offset)

    rows, err := r.db.Query(ctx, listQuery, args...)
    if err != nil {
        return nil, 0, fmt.Errorf("gagal list jenis_iuran: %w", err)
    }
    defer rows.Close()

    var result []models.JenisIuran
    for rows.Next() {
        var j models.JenisIuran
        if err := rows.Scan(
            &j.ID, &j.SchoolID, &j.Kode, &j.Nama, &j.Kategori,
            &j.IsWajib, &j.IsAktif, &j.Keterangan, &j.CreatedAt, &j.UpdatedAt,
        ); err != nil {
            return nil, 0, fmt.Errorf("gagal scan jenis_iuran: %w", err)
        }
        result = append(result, j)
    }
    if err := rows.Err(); err != nil {
        return nil, 0, err
    }
    if result == nil {
        result = []models.JenisIuran{}
    }
    return result, total, nil
}

// =====================================================================
// TARIF IURAN
// =====================================================================
func (r *pgxKeuanganMasterRepository) CreateTarifIuran(ctx context.Context, t *models.TarifIuran) error {
    query := `
        INSERT INTO tarif_iuran
            (id, school_id, academic_period_id, jenis_iuran_id, kelas, nominal, tahun_berlaku)
        VALUES
            (gen_random_uuid(), $1, $2, $3, NULLIF($4, ''), $5, $6)
        RETURNING id, created_at, updated_at
    `
    row := r.db.QueryRow(ctx, query,
        t.SchoolID, t.AcademicPeriodID, t.JenisIuranID, t.Kelas, t.Nominal, t.TahunBerlaku,
    )
    if err := row.Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt); err != nil {
        return fmt.Errorf("gagal insert tarif_iuran: %w", err)
    }
    return nil
}

func (r *pgxKeuanganMasterRepository) UpdateTarifIuran(ctx context.Context, t *models.TarifIuran) error {
    query := `
        UPDATE tarif_iuran
        SET jenis_iuran_id=$1, kelas=NULLIF($2, ''), nominal=$3, tahun_berlaku=$4, updated_at=NOW()
        WHERE school_id=$5 AND id=$6
        RETURNING updated_at
    `
    row := r.db.QueryRow(ctx, query,
        t.JenisIuranID, t.Kelas, t.Nominal, t.TahunBerlaku,
        t.SchoolID, t.ID,
    )
    if err := row.Scan(&t.UpdatedAt); err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return ErrKeuanganMasterNotFound
        }
        return fmt.Errorf("gagal update tarif_iuran: %w", err)
    }
    return nil
}

func (r *pgxKeuanganMasterRepository) DeleteTarifIuran(ctx context.Context, schoolID, id string) error {
    ct, err := r.db.Exec(ctx, `DELETE FROM tarif_iuran WHERE school_id=$1 AND id=$2`, schoolID, id)
    if err != nil {
        return fmt.Errorf("gagal delete tarif_iuran: %w", err)
    }
    if ct.RowsAffected() == 0 {
        return ErrKeuanganMasterNotFound
    }
    return nil
}

func (r *pgxKeuanganMasterRepository) GetTarifIuranByID(ctx context.Context, schoolID, id string) (*models.TarifIuran, error) {
    query := `
        SELECT t.id, t.school_id, t.academic_period_id, t.jenis_iuran_id,
               COALESCE(j.nama, ''), COALESCE(t.kelas, ''), t.nominal, t.tahun_berlaku,
               t.created_at, t.updated_at
        FROM tarif_iuran t
        LEFT JOIN jenis_iuran j ON j.id = t.jenis_iuran_id
        WHERE t.school_id=$1 AND t.id=$2
    `
    var t models.TarifIuran
    err := r.db.QueryRow(ctx, query, schoolID, id).Scan(
        &t.ID, &t.SchoolID, &t.AcademicPeriodID, &t.JenisIuranID,
        &t.JenisIuranNama, &t.Kelas, &t.Nominal, &t.TahunBerlaku,
        &t.CreatedAt, &t.UpdatedAt,
    )
    if errors.Is(err, pgx.ErrNoRows) {
        return nil, ErrKeuanganMasterNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("gagal get tarif_iuran: %w", err)
    }
    return &t, nil
}

func (r *pgxKeuanganMasterRepository) ListTarifIuran(ctx context.Context, filter models.TarifIuranFilter) ([]models.TarifIuran, int, error) {
    limit, offset := normalizePagination(filter.Limit, filter.Offset)

    where := []string{"t.school_id = $1"}
    args := []interface{}{filter.SchoolID}
    argPos := 2

    if filter.AcademicPeriodID != "" {
        where = append(where, fmt.Sprintf("t.academic_period_id = $%d", argPos))
        args = append(args, filter.AcademicPeriodID)
        argPos++
    }
    if filter.JenisIuranID != "" {
        where = append(where, fmt.Sprintf("t.jenis_iuran_id = $%d", argPos))
        args = append(args, filter.JenisIuranID)
        argPos++
    }
    if filter.Kelas != "" {
        where = append(where, fmt.Sprintf("COALESCE(t.kelas, '') = $%d", argPos))
        args = append(args, filter.Kelas)
        argPos++
    }

    whereClause := strings.Join(where, " AND ")

    var total int
    countQuery := "SELECT COUNT(*) FROM tarif_iuran t WHERE " + whereClause
    if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
        return nil, 0, fmt.Errorf("gagal count tarif_iuran: %w", err)
    }

    listQuery := fmt.Sprintf(`
        SELECT t.id, t.school_id, t.academic_period_id, t.jenis_iuran_id,
               COALESCE(j.nama, ''), COALESCE(t.kelas, ''), t.nominal, t.tahun_berlaku,
               t.created_at, t.updated_at
        FROM tarif_iuran t
        LEFT JOIN jenis_iuran j ON j.id = t.jenis_iuran_id
        WHERE %s
        ORDER BY j.kode ASC, t.kelas ASC NULLS FIRST
        LIMIT $%d OFFSET $%d
    `, whereClause, argPos, argPos+1)
    args = append(args, limit, offset)

    rows, err := r.db.Query(ctx, listQuery, args...)
    if err != nil {
        return nil, 0, fmt.Errorf("gagal list tarif_iuran: %w", err)
    }
    defer rows.Close()

    var result []models.TarifIuran
    for rows.Next() {
        var t models.TarifIuran
        if err := rows.Scan(
            &t.ID, &t.SchoolID, &t.AcademicPeriodID, &t.JenisIuranID,
            &t.JenisIuranNama, &t.Kelas, &t.Nominal, &t.TahunBerlaku,
            &t.CreatedAt, &t.UpdatedAt,
        ); err != nil {
            return nil, 0, fmt.Errorf("gagal scan tarif_iuran: %w", err)
        }
        result = append(result, t)
    }
    if err := rows.Err(); err != nil {
        return nil, 0, err
    }
    if result == nil {
        result = []models.TarifIuran{}
    }
    return result, total, nil
}

// =====================================================================
// KODE REKENING
// =====================================================================
func (r *pgxKeuanganMasterRepository) CreateKodeRekening(ctx context.Context, k *models.KodeRekening) error {
    parent := nullableStr(k.ParentID)
    query := `
        INSERT INTO kode_rekening
            (id, school_id, kode, nama, tipe, parent_id, level, is_aktif, keterangan)
        VALUES
            (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, created_at, updated_at
    `
    row := r.db.QueryRow(ctx, query,
        k.SchoolID, k.Kode, k.Nama, k.Tipe, parent, k.Level, k.IsAktif, k.Keterangan,
    )
    if err := row.Scan(&k.ID, &k.CreatedAt, &k.UpdatedAt); err != nil {
        return fmt.Errorf("gagal insert kode_rekening: %w", err)
    }
    return nil
}

func (r *pgxKeuanganMasterRepository) UpdateKodeRekening(ctx context.Context, k *models.KodeRekening) error {
    parent := nullableStr(k.ParentID)
    query := `
        UPDATE kode_rekening
        SET kode=$1, nama=$2, tipe=$3, parent_id=$4, level=$5,
            is_aktif=$6, keterangan=$7, updated_at=NOW()
        WHERE school_id=$8 AND id=$9
        RETURNING updated_at
    `
    row := r.db.QueryRow(ctx, query,
        k.Kode, k.Nama, k.Tipe, parent, k.Level, k.IsAktif, k.Keterangan,
        k.SchoolID, k.ID,
    )
    if err := row.Scan(&k.UpdatedAt); err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return ErrKeuanganMasterNotFound
        }
        return fmt.Errorf("gagal update kode_rekening: %w", err)
    }
    return nil
}

func (r *pgxKeuanganMasterRepository) DeleteKodeRekening(ctx context.Context, schoolID, id string) error {
    ct, err := r.db.Exec(ctx, `DELETE FROM kode_rekening WHERE school_id=$1 AND id=$2`, schoolID, id)
    if err != nil {
        return fmt.Errorf("gagal delete kode_rekening: %w", err)
    }
    if ct.RowsAffected() == 0 {
        return ErrKeuanganMasterNotFound
    }
    return nil
}

func (r *pgxKeuanganMasterRepository) GetKodeRekeningByID(ctx context.Context, schoolID, id string) (*models.KodeRekening, error) {
    query := `
        SELECT id, school_id, kode, nama, tipe, COALESCE(parent_id::text, ''),
               level, is_aktif, COALESCE(keterangan, ''), created_at, updated_at
        FROM kode_rekening
        WHERE school_id=$1 AND id=$2
    `
    var k models.KodeRekening
    err := r.db.QueryRow(ctx, query, schoolID, id).Scan(
        &k.ID, &k.SchoolID, &k.Kode, &k.Nama, &k.Tipe, &k.ParentID,
        &k.Level, &k.IsAktif, &k.Keterangan, &k.CreatedAt, &k.UpdatedAt,
    )
    if errors.Is(err, pgx.ErrNoRows) {
        return nil, ErrKeuanganMasterNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("gagal get kode_rekening: %w", err)
    }
    return &k, nil
}

func (r *pgxKeuanganMasterRepository) ListKodeRekening(ctx context.Context, filter models.KodeRekeningFilter) ([]models.KodeRekening, int, error) {
    limit, offset := normalizePagination(filter.Limit, filter.Offset)

    where := []string{"school_id = $1"}
    args := []interface{}{filter.SchoolID}
    argPos := 2

    if filter.Tipe != "" {
        where = append(where, fmt.Sprintf("tipe = $%d", argPos))
        args = append(args, filter.Tipe)
        argPos++
    }
    if filter.IsAktif != nil {
        where = append(where, fmt.Sprintf("is_aktif = $%d", argPos))
        args = append(args, *filter.IsAktif)
        argPos++
    }
    if filter.Search != "" {
        where = append(where, fmt.Sprintf("(nama ILIKE $%d OR kode ILIKE $%d)", argPos, argPos))
        args = append(args, "%"+filter.Search+"%")
        argPos++
    }

    whereClause := strings.Join(where, " AND ")

    var total int
    countQuery := "SELECT COUNT(*) FROM kode_rekening WHERE " + whereClause
    if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
        return nil, 0, fmt.Errorf("gagal count kode_rekening: %w", err)
    }

    listQuery := fmt.Sprintf(`
        SELECT id, school_id, kode, nama, tipe, COALESCE(parent_id::text, ''),
               level, is_aktif, COALESCE(keterangan, ''), created_at, updated_at
        FROM kode_rekening
        WHERE %s
        ORDER BY kode ASC
        LIMIT $%d OFFSET $%d
    `, whereClause, argPos, argPos+1)
    args = append(args, limit, offset)

    rows, err := r.db.Query(ctx, listQuery, args...)
    if err != nil {
        return nil, 0, fmt.Errorf("gagal list kode_rekening: %w", err)
    }
    defer rows.Close()

    var result []models.KodeRekening
    for rows.Next() {
        var k models.KodeRekening
        if err := rows.Scan(
            &k.ID, &k.SchoolID, &k.Kode, &k.Nama, &k.Tipe, &k.ParentID,
            &k.Level, &k.IsAktif, &k.Keterangan, &k.CreatedAt, &k.UpdatedAt,
        ); err != nil {
            return nil, 0, fmt.Errorf("gagal scan kode_rekening: %w", err)
        }
        result = append(result, k)
    }
    if err := rows.Err(); err != nil {
        return nil, 0, err
    }
    if result == nil {
        result = []models.KodeRekening{}
    }
    return result, total, nil
}

// =====================================================================
// REKENING BANK
// =====================================================================
func (r *pgxKeuanganMasterRepository) CreateRekeningBank(ctx context.Context, rb *models.RekeningBank) error {
    kodeRek := nullableStr(rb.KodeRekeningID)
    query := `
        INSERT INTO rekening_bank
            (id, school_id, nama_bank, nomor_rekening, atas_nama, cabang,
             saldo_awal, saldo_terkini, kode_rekening_id, format_file_mutasi, is_aktif)
        VALUES
            (gen_random_uuid(), $1, $2, $3, $4, NULLIF($5,''), $6, $7, $8, NULLIF($9,''), $10)
        RETURNING id, created_at, updated_at
    `
    row := r.db.QueryRow(ctx, query,
        rb.SchoolID, rb.NamaBank, rb.NomorRekening, rb.AtasNama, rb.Cabang,
        rb.SaldoAwal, rb.SaldoTerkini, kodeRek, rb.FormatFileMutasi, rb.IsAktif,
    )
    if err := row.Scan(&rb.ID, &rb.CreatedAt, &rb.UpdatedAt); err != nil {
        return fmt.Errorf("gagal insert rekening_bank: %w", err)
    }
    return nil
}

func (r *pgxKeuanganMasterRepository) UpdateRekeningBank(ctx context.Context, rb *models.RekeningBank) error {
    kodeRek := nullableStr(rb.KodeRekeningID)
    query := `
        UPDATE rekening_bank
        SET nama_bank=$1, nomor_rekening=$2, atas_nama=$3, cabang=NULLIF($4,''),
            saldo_awal=$5, kode_rekening_id=$6, format_file_mutasi=NULLIF($7,''),
            is_aktif=$8, updated_at=NOW()
        WHERE school_id=$9 AND id=$10
        RETURNING updated_at
    `
    row := r.db.QueryRow(ctx, query,
        rb.NamaBank, rb.NomorRekening, rb.AtasNama, rb.Cabang,
        rb.SaldoAwal, kodeRek, rb.FormatFileMutasi, rb.IsAktif,
        rb.SchoolID, rb.ID,
    )
    if err := row.Scan(&rb.UpdatedAt); err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return ErrKeuanganMasterNotFound
        }
        return fmt.Errorf("gagal update rekening_bank: %w", err)
    }
    return nil
}

func (r *pgxKeuanganMasterRepository) DeleteRekeningBank(ctx context.Context, schoolID, id string) error {
    ct, err := r.db.Exec(ctx, `DELETE FROM rekening_bank WHERE school_id=$1 AND id=$2`, schoolID, id)
    if err != nil {
        return fmt.Errorf("gagal delete rekening_bank: %w", err)
    }
    if ct.RowsAffected() == 0 {
        return ErrKeuanganMasterNotFound
    }
    return nil
}

func (r *pgxKeuanganMasterRepository) GetRekeningBankByID(ctx context.Context, schoolID, id string) (*models.RekeningBank, error) {
    query := `
        SELECT id, school_id, nama_bank, nomor_rekening, atas_nama, COALESCE(cabang, ''),
               saldo_awal, saldo_terkini, COALESCE(kode_rekening_id::text, ''),
               COALESCE(format_file_mutasi, ''), is_aktif, created_at, updated_at
        FROM rekening_bank
        WHERE school_id=$1 AND id=$2
    `
    var rb models.RekeningBank
    err := r.db.QueryRow(ctx, query, schoolID, id).Scan(
        &rb.ID, &rb.SchoolID, &rb.NamaBank, &rb.NomorRekening, &rb.AtasNama, &rb.Cabang,
        &rb.SaldoAwal, &rb.SaldoTerkini, &rb.KodeRekeningID,
        &rb.FormatFileMutasi, &rb.IsAktif, &rb.CreatedAt, &rb.UpdatedAt,
    )
    if errors.Is(err, pgx.ErrNoRows) {
        return nil, ErrKeuanganMasterNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("gagal get rekening_bank: %w", err)
    }
    return &rb, nil
}

func (r *pgxKeuanganMasterRepository) ListRekeningBank(ctx context.Context, filter models.RekeningBankFilter) ([]models.RekeningBank, int, error) {
    limit, offset := normalizePagination(filter.Limit, filter.Offset)

    where := []string{"school_id = $1"}
    args := []interface{}{filter.SchoolID}
    argPos := 2

    if filter.IsAktif != nil {
        where = append(where, fmt.Sprintf("is_aktif = $%d", argPos))
        args = append(args, *filter.IsAktif)
        argPos++
    }
    if filter.Search != "" {
        where = append(where, fmt.Sprintf("(nama_bank ILIKE $%d OR nomor_rekening ILIKE $%d OR atas_nama ILIKE $%d)", argPos, argPos, argPos))
        args = append(args, "%"+filter.Search+"%")
        argPos++
    }

    whereClause := strings.Join(where, " AND ")

    var total int
    countQuery := "SELECT COUNT(*) FROM rekening_bank WHERE " + whereClause
    if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
        return nil, 0, fmt.Errorf("gagal count rekening_bank: %w", err)
    }

    listQuery := fmt.Sprintf(`
        SELECT id, school_id, nama_bank, nomor_rekening, atas_nama, COALESCE(cabang, ''),
               saldo_awal, saldo_terkini, COALESCE(kode_rekening_id::text, ''),
               COALESCE(format_file_mutasi, ''), is_aktif, created_at, updated_at
        FROM rekening_bank
        WHERE %s
        ORDER BY nama_bank ASC, nomor_rekening ASC
        LIMIT $%d OFFSET $%d
    `, whereClause, argPos, argPos+1)
    args = append(args, limit, offset)

    rows, err := r.db.Query(ctx, listQuery, args...)
    if err != nil {
        return nil, 0, fmt.Errorf("gagal list rekening_bank: %w", err)
    }
    defer rows.Close()

    var result []models.RekeningBank
    for rows.Next() {
        var rb models.RekeningBank
        if err := rows.Scan(
            &rb.ID, &rb.SchoolID, &rb.NamaBank, &rb.NomorRekening, &rb.AtasNama, &rb.Cabang,
            &rb.SaldoAwal, &rb.SaldoTerkini, &rb.KodeRekeningID,
            &rb.FormatFileMutasi, &rb.IsAktif, &rb.CreatedAt, &rb.UpdatedAt,
        ); err != nil {
            return nil, 0, fmt.Errorf("gagal scan rekening_bank: %w", err)
        }
        result = append(result, rb)
    }
    if err := rows.Err(); err != nil {
        return nil, 0, err
    }
    if result == nil {
        result = []models.RekeningBank{}
    }
    return result, total, nil
}

// =====================================================================
// Internal helper
// =====================================================================
func nullableStr(s string) interface{} {
    if s == "" {
        return nil
    }
    return s
}
