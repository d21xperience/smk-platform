package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// =====================================================================
// LaporanRepository — query aggregat untuk laporan keuangan
// =====================================================================
type LaporanRepository interface {
	// Neraca (compute live dari BKU + piutang)
	GetNeracaData(ctx context.Context, schoolID, periodID string, asOfDate time.Time) (*NeracaData, error)

	// Saldo kas (untuk dashboard)
	GetSaldoKasData(ctx context.Context, schoolID, periodID string, bulan, tahun int) (*SaldoKasData, error)

	// Ringkasan bulanan
	GetRingkasanBulanData(ctx context.Context, schoolID, periodID string, bulan, tahun int) (*RingkasanData, error)

	// Dashboard helpers
	GetPiutangForDashboard(ctx context.Context, schoolID, periodID string) (*PiutangDashboardData, error)
	GetPengajuanForDashboard(ctx context.Context, schoolID, periodID string) (map[string]int, float64, error)
	GetRekonForDashboard(ctx context.Context, schoolID, periodID string, bulan, tahun int) (total, balance, unbalance int, err error)
}

type pgxLaporanRepository struct {
	db *pgxpool.Pool
}

func NewLaporanRepository(db *pgxpool.Pool) LaporanRepository {
	return &pgxLaporanRepository{db: db}
}

// =====================================================================
// Internal DTO — intermediate data untuk orkestrasi service
// =====================================================================

// NeracaDataItem — 1 baris saldo COA (level 3) + piutang virtual
type NeracaDataItem struct {
	KodeRekeningID string  `json:"kode_rekening_id"`
	Kode           string  `json:"kode"`
	Nama           string  `json:"nama"`
	Tipe           string  `json:"tipe"`
	Saldo          float64 `json:"saldo"`
}

type NeracaData struct {
	Items []NeracaDataItem
}

type SaldoKasData struct {
	SaldoTerkini        float64
	TotalDebetBulanIni  float64
	TotalKreditBulanIni float64
}

type RingkasanKategoriData struct {
	KodeRekening    string
	Nama            string
	Total           float64
	JumlahTransaksi int
}

type RingkasanData struct {
	SaldoAwal       float64
	TotalDebet      float64
	TotalKredit     float64
	SaldoAkhir      float64
	JumlahTransaksi int
	TopDebet        []RingkasanKategoriData
	TopKredit       []RingkasanKategoriData
}

type PiutangDashboardData struct {
	TotalSisa   float64
	JumlahItem  int
	JumlahSiswa int
}

// =====================================================================
// GetNeracaData — aggregate BKU per COA level 3 + piutang virtual
//
// Saldo natural:
//   - ASET, BELANJA     : debet - kredit
//   - KEWAJIBAN, EKUITAS, PENDAPATAN : kredit - debet
//
// Piutang (tagihan BELUM_LUNAS/SEBAGIAN) di-inject sebagai item ASET
// virtual dengan kode "1.1.9" agar neraca lengkap.
// =====================================================================
func (r *pgxLaporanRepository) GetNeracaData(
	ctx context.Context, schoolID, periodID string, asOfDate time.Time,
) (*NeracaData, error) {
	items := []NeracaDataItem{}
	asOfStr := asOfDate.Format("2006-01-02")

	// 1. Kas Tunai — dari running balance terakhir BKU
	var saldoKas float64
	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(saldo_berjalan, 0)
		FROM buku_kas_umum
		WHERE school_id = $1 AND academic_period_id = $2 AND tanggal <= $3::date
		ORDER BY tanggal DESC, created_at DESC
		LIMIT 1
	`, schoolID, periodID, asOfStr).Scan(&saldoKas)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("gagal get saldo kas neraca: %w", err)
	}
	items = append(items, NeracaDataItem{
		KodeRekeningID: "KAS",
		Kode:           "1.1.1",
		Nama:           "Kas Tunai",
		Tipe:           "ASET",
		Saldo:          saldoKas,
	})

	// 2. Piutang Siswa — dari tagihan BELUM_LUNAS/SEBAGIAN
	var totalPiutang float64
	err = r.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(nominal - total_dibayar), 0)
		FROM tagihan_siswa
		WHERE school_id = $1 AND academic_period_id = $2
		  AND status IN ('BELUM_LUNAS', 'SEBAGIAN')
		  AND (nominal - total_dibayar) > 0
	`, schoolID, periodID).Scan(&totalPiutang)
	if err != nil {
		return nil, fmt.Errorf("gagal get piutang neraca: %w", err)
	}
	if totalPiutang > 0 {
		items = append(items, NeracaDataItem{
			KodeRekeningID: "PIUTANG_SISWA",
			Kode:           "1.1.9",
			Nama:           "Piutang Siswa",
			Tipe:           "ASET",
			Saldo:          totalPiutang,
		})
	}

	// 3. Pendapatan & Belanja — dari BKU (kas-oriented sign)
	//    PENDAPATAN: kas masuk → debet (saldo = SUM(debet))
	//    BELANJA   : kas keluar → kredit (saldo = SUM(kredit))
	rows, err := r.db.Query(ctx, `
		SELECT 
			k.id::text AS kode_rekening_id,
			k.kode,
			k.nama,
			k.tipe,
			SUM(CASE 
				WHEN k.tipe = 'PENDAPATAN' THEN b.debet
				WHEN k.tipe = 'BELANJA' THEN b.kredit
				ELSE 0
			END) AS saldo
		FROM buku_kas_umum b
		JOIN kode_rekening k ON k.id = b.kode_rekening_id
		WHERE b.school_id = $1
		  AND b.academic_period_id = $2
		  AND b.tanggal <= $3::date
		  AND k.level = 3
		  AND k.tipe IN ('PENDAPATAN', 'BELANJA')
		GROUP BY k.id, k.kode, k.nama, k.tipe
		ORDER BY k.tipe ASC, k.kode ASC
	`, schoolID, periodID, asOfStr)
	if err != nil {
		return nil, fmt.Errorf("gagal query pendapatan/belanja neraca: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var it NeracaDataItem
		if err := rows.Scan(&it.KodeRekeningID, &it.Kode, &it.Nama, &it.Tipe, &it.Saldo); err != nil {
			return nil, fmt.Errorf("gagal scan pendapatan/belanja: %w", err)
		}
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &NeracaData{Items: items}, nil
}

// =====================================================================
// GetSaldoKasData — saldo terkini + total debet/kredit bulan
// =====================================================================
func (r *pgxLaporanRepository) GetSaldoKasData(
	ctx context.Context, schoolID, periodID string, bulan, tahun int,
) (*SaldoKasData, error) {
	out := &SaldoKasData{}

	// Saldo terkini (running balance terakhir)
	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(saldo_berjalan, 0)
		FROM buku_kas_umum
		WHERE school_id = $1 AND academic_period_id = $2
		ORDER BY tanggal DESC, created_at DESC
		LIMIT 1
	`, schoolID, periodID).Scan(&out.SaldoTerkini)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("gagal get saldo: %w", err)
	}

	// Debet & kredit bulan ini (kalau bulan di-set)
	if bulan >= 1 && bulan <= 12 && tahun >= 2000 {
		err = r.db.QueryRow(ctx, `
			SELECT COALESCE(SUM(debet), 0), COALESCE(SUM(kredit), 0)
			FROM buku_kas_umum
			WHERE school_id = $1 AND academic_period_id = $2
			  AND EXTRACT(MONTH FROM tanggal) = $3
			  AND EXTRACT(YEAR FROM tanggal) = $4
		`, schoolID, periodID, bulan, tahun).Scan(&out.TotalDebetBulanIni, &out.TotalKreditBulanIni)
		if err != nil {
			return nil, fmt.Errorf("gagal sum debet/kredit bulan: %w", err)
		}
	}

	return out, nil
}

// =====================================================================
// GetRingkasanBulanData — aggregate 1 bulan
// =====================================================================
func (r *pgxLaporanRepository) GetRingkasanBulanData(
	ctx context.Context, schoolID, periodID string, bulan, tahun int,
) (*RingkasanData, error) {
	out := &RingkasanData{}

	// Saldo awal (sebelum tanggal 1 bulan ini)
	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(saldo_berjalan, 0)
		FROM buku_kas_umum
		WHERE school_id = $1 AND academic_period_id = $2
		  AND tanggal < make_date($3, $4, 1)
		ORDER BY tanggal DESC, created_at DESC
		LIMIT 1
	`, schoolID, periodID, tahun, bulan).Scan(&out.SaldoAwal)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("gagal get saldo awal: %w", err)
	}

	// Total bulan
	err = r.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(debet), 0), COALESCE(SUM(kredit), 0), COUNT(*)
		FROM buku_kas_umum
		WHERE school_id = $1 AND academic_period_id = $2
		  AND EXTRACT(MONTH FROM tanggal) = $3
		  AND EXTRACT(YEAR FROM tanggal) = $4
	`, schoolID, periodID, bulan, tahun).Scan(&out.TotalDebet, &out.TotalKredit, &out.JumlahTransaksi)
	if err != nil {
		return nil, fmt.Errorf("gagal sum bulan: %w", err)
	}
	out.SaldoAkhir = out.SaldoAwal + out.TotalDebet - out.TotalKredit

	// Top 3 debet
	out.TopDebet, err = r.topKategoriBy(ctx, schoolID, periodID, bulan, tahun, "debet", 3)
	if err != nil {
		return nil, err
	}

	// Top 3 kredit
	out.TopKredit, err = r.topKategoriBy(ctx, schoolID, periodID, bulan, tahun, "kredit", 3)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// helper: top N kategori berdasarkan debet/kredit
func (r *pgxLaporanRepository) topKategoriBy(
	ctx context.Context, schoolID, periodID string, bulan, tahun int,
	kolom string, limit int,
) ([]RingkasanKategoriData, error) {
	// Validasi kolom untuk hindari SQL injection
	if kolom != "debet" && kolom != "kredit" {
		return nil, fmt.Errorf("kolom tidak valid: %s", kolom)
	}

	query := fmt.Sprintf(`
		SELECT k.kode, k.nama, SUM(b.%s) AS total, COUNT(*) AS cnt
		FROM buku_kas_umum b
		JOIN kode_rekening k ON k.id = b.kode_rekening_id
		WHERE b.school_id = $1 AND b.academic_period_id = $2
		  AND EXTRACT(MONTH FROM b.tanggal) = $3
		  AND EXTRACT(YEAR FROM b.tanggal) = $4
		  AND b.%s > 0
		GROUP BY k.kode, k.nama
		ORDER BY total DESC
		LIMIT %d
	`, kolom, kolom, limit)

	rows, err := r.db.Query(ctx, query, schoolID, periodID, bulan, tahun)
	if err != nil {
		return nil, fmt.Errorf("gagal query top %s: %w", kolom, err)
	}
	defer rows.Close()

	out := []RingkasanKategoriData{}
	for rows.Next() {
		var it RingkasanKategoriData
		if err := rows.Scan(&it.KodeRekening, &it.Nama, &it.Total, &it.JumlahTransaksi); err != nil {
			return nil, fmt.Errorf("gagal scan top %s: %w", kolom, err)
		}
		out = append(out, it)
	}
	return out, rows.Err()
}
