package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
)

// =====================================================================
// Interface Aggregator — semua sub-domain transaksi keuangan
// =====================================================================
type KeuanganTransaksiRepository interface {
	// Meta
	BeginTx(ctx context.Context) (pgx.Tx, error)

	// Nomor auto-generate (dipanggil dalam transaksi)
	GenerateNomorBKM(ctx context.Context, tx pgx.Tx, schoolID, periodID string, tanggal time.Time) (string, error)
	GenerateNomorBKK(ctx context.Context, tx pgx.Tx, schoolID, periodID string, tanggal time.Time) (string, error)
	GenerateNomorPengajuan(ctx context.Context, tx pgx.Tx, schoolID, periodID string, tanggal time.Time) (string, error)

	// Get saldo terakhir (untuk menghitung running balance BKU)
	GetSaldoTerakhir(ctx context.Context, schoolID, periodID string, tanggal time.Time) (float64, error)

	// Tagihan
	CreateTagihan(ctx context.Context, tx pgx.Tx, t *models.TagihanSiswa) error
	GetTagihanByID(ctx context.Context, schoolID, id string) (*models.TagihanSiswa, error)
	ListTagihan(ctx context.Context, filter models.TagihanFilter) ([]models.TagihanSiswa, int, error)
	UpdateTagihanStatus(ctx context.Context, tx pgx.Tx, id string, status string, totalDibayar float64, tanggalBayar *time.Time) error
	FindTagihanByNISPeriode(ctx context.Context, schoolID, periodID, nis, jenisIuranID string, bulan, tahun int) (*models.TagihanSiswa, error)

	// Penerimaan
	CreatePenerimaan(ctx context.Context, tx pgx.Tx, p *models.PenerimaanIuran) error
	GetPenerimaanByNoBKM(ctx context.Context, schoolID, noBKM string) (*models.PenerimaanIuran, error)
	ListPenerimaan(ctx context.Context, filter models.PenerimaanFilter) ([]models.PenerimaanIuran, int, error)

	// Buku Kas Umum
	InsertBKU(ctx context.Context, tx pgx.Tx, b *models.BukuKasUmum) error
	ListBKU(ctx context.Context, filter models.BukuKasFilter) ([]models.BukuKasUmum, int, error)

	// Pengajuan Pengeluaran
	CreatePengajuan(ctx context.Context, tx pgx.Tx, p *models.PengajuanPengeluaran) error
	GetPengajuanByID(ctx context.Context, schoolID, id string) (*models.PengajuanPengeluaran, error)
	ListPengajuan(ctx context.Context, filter models.PengajuanFilter) ([]models.PengajuanPengeluaran, int, error)
	UpdatePengajuanStatus(ctx context.Context, tx pgx.Tx, id string, newStatus string, verifiedBy, approvedBy *string, catatanBendahara, catatanKepala, noBKK *string) error

	// Pengeluaran
	CreatePengeluaran(ctx context.Context, tx pgx.Tx, p *models.Pengeluaran) error
	GetPengeluaranByNoBKK(ctx context.Context, schoolID, noBKK string) (*models.Pengeluaran, error)
	ListPengeluaran(ctx context.Context, filter models.PengeluaranFilter) ([]models.Pengeluaran, int, error)

	// Siswa lookup (untuk validasi NIS saat terima pembayaran)
	FindSiswaByNIS(ctx context.Context, schoolID, nis string) (*models.SiswaInfo, error)
}

// =====================================================================
// Implementation
// =====================================================================
type pgxKeuanganTransaksiRepository struct {
	db *pgxpool.Pool
}

func NewKeuanganTransaksiRepository(db *pgxpool.Pool) KeuanganTransaksiRepository {
	return &pgxKeuanganTransaksiRepository{db: db}
}

var (
	ErrTransaksiKeuanganNotFound = errors.New("data transaksi keuangan tidak ditemukan")
)

// BeginTx — helper untuk service memulai transaksi atomik
func (r *pgxKeuanganTransaksiRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.db.Begin(ctx)
}

// =====================================================================
// Nomor Auto-Generate — format BKM/001/SMK/XI/2026
// =====================================================================

var bulanRomawi = []string{"I", "II", "III", "IV", "V", "VI",
	"VII", "VIII", "IX", "X", "XI", "XII"}

func formatNomorBukti(prefix string, counter int, tanggal time.Time) string {
	bln := bulanRomawi[int(tanggal.Month())-1]
	return fmt.Sprintf("%s/%03d/SMK/%s/%d", prefix, counter, bln, tanggal.Year())
}

// nextCounter — query counter terakhir dari tabel, untuk generate nomor unik
func (r *pgxKeuanganTransaksiRepository) nextCounter(
	ctx context.Context, tx pgx.Tx,
	tableName, kolomNomor, schoolID, periodID string,
	tanggal time.Time,
) (int, error) {
	query := fmt.Sprintf(`
		SELECT COALESCE(MAX(CAST(SPLIT_PART(%s, '/', 2) AS INTEGER)), 0) + 1
		FROM %s
		WHERE school_id = $1
		  AND academic_period_id = $2
		  AND EXTRACT(MONTH FROM tanggal) = $3
		  AND EXTRACT(YEAR FROM tanggal) = $4
	`, kolomNomor, tableName)

	var next int
	err := tx.QueryRow(ctx, query, schoolID, periodID, int(tanggal.Month()), tanggal.Year()).Scan(&next)
	if err != nil {
		return 0, fmt.Errorf("gagal hitung counter %s: %w", tableName, err)
	}
	return next, nil
}

func (r *pgxKeuanganTransaksiRepository) GenerateNomorBKM(
	ctx context.Context, tx pgx.Tx,
	schoolID, periodID string, tanggal time.Time,
) (string, error) {
	n, err := r.nextCounter(ctx, tx, "penerimaan_iuran", "no_bkm", schoolID, periodID, tanggal)
	if err != nil {
		return "", err
	}
	return formatNomorBukti("BKM", n, tanggal), nil
}

func (r *pgxKeuanganTransaksiRepository) GenerateNomorBKK(
	ctx context.Context, tx pgx.Tx,
	schoolID, periodID string, tanggal time.Time,
) (string, error) {
	n, err := r.nextCounter(ctx, tx, "pengeluaran", "no_bkk", schoolID, periodID, tanggal)
	if err != nil {
		return "", err
	}
	return formatNomorBukti("BKK", n, tanggal), nil
}

func (r *pgxKeuanganTransaksiRepository) GenerateNomorPengajuan(
	ctx context.Context, tx pgx.Tx,
	schoolID, periodID string, tanggal time.Time,
) (string, error) {
	n, err := r.nextCounter(ctx, tx, "pengajuan_pengeluaran", "nomor_pengajuan", schoolID, periodID, tanggal)
	if err != nil {
		return "", err
	}
	return formatNomorBukti("PP", n, tanggal), nil
}

// =====================================================================
// Saldo Terakhir BKU — basis untuk running balance
// =====================================================================
func (r *pgxKeuanganTransaksiRepository) GetSaldoTerakhir(
	ctx context.Context, schoolID, periodID string, tanggal time.Time,
) (float64, error) {
	query := `
		SELECT COALESCE(saldo_berjalan, 0)
		FROM buku_kas_umum
		WHERE school_id = $1
		  AND academic_period_id = $2
		  AND tanggal <= $3
		ORDER BY tanggal DESC, created_at DESC
		LIMIT 1
	`
	var saldo float64
	err := r.db.QueryRow(ctx, query, schoolID, periodID, tanggal).Scan(&saldo)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("gagal get saldo terakhir: %w", err)
	}
	return saldo, nil
}

// =====================================================================
// Helpers untuk scan — nullable types
// =====================================================================
func strOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
