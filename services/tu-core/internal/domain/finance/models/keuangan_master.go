package models

import "time"

// =====================================================================
// Master Keuangan — DTO & Filter
// =====================================================================

// ---------------------------------------------------------------------
// Jenis Iuran
// ---------------------------------------------------------------------
type JenisIuran struct {
    ID         string    `json:"id"`
    SchoolID   string    `json:"school_id"`
    Kode       string    `json:"kode"`
    Nama       string    `json:"nama"`
    Kategori   string    `json:"kategori"`
    IsWajib    bool      `json:"is_wajib"`
    IsAktif    bool      `json:"is_aktif"`
    Keterangan string    `json:"keterangan"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}

type JenisIuranFilter struct {
    SchoolID string
    Kategori string
    IsAktif  *bool
    Search   string
    Limit    int
    Offset   int
}

// ---------------------------------------------------------------------
// Tarif Iuran
// ---------------------------------------------------------------------
type TarifIuran struct {
    ID               string    `json:"id"`
    SchoolID         string    `json:"school_id"`
    AcademicPeriodID string    `json:"academic_period_id"`
    JenisIuranID     string    `json:"jenis_iuran_id"`
    JenisIuranNama   string    `json:"jenis_iuran_nama,omitempty"`
    Kelas            string    `json:"kelas"`
    Nominal          float64   `json:"nominal"`
    TahunBerlaku     int       `json:"tahun_berlaku"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
}

type TarifIuranFilter struct {
    SchoolID         string
    AcademicPeriodID string
    JenisIuranID     string
    Kelas            string
    Limit            int
    Offset           int
}

// ---------------------------------------------------------------------
// Kode Rekening (COA)
// ---------------------------------------------------------------------
type KodeRekening struct {
    ID         string    `json:"id"`
    SchoolID   string    `json:"school_id"`
    Kode       string    `json:"kode"`
    Nama       string    `json:"nama"`
    Tipe       string    `json:"tipe"`
    ParentID   string    `json:"parent_id,omitempty"`
    Level      int       `json:"level"`
    IsAktif    bool      `json:"is_aktif"`
    Keterangan string    `json:"keterangan,omitempty"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}

type KodeRekeningFilter struct {
    SchoolID string
    Tipe     string
    IsAktif  *bool
    Search   string
    Limit    int
    Offset   int
}

// ---------------------------------------------------------------------
// Rekening Bank
// ---------------------------------------------------------------------
type RekeningBank struct {
    ID               string    `json:"id"`
    SchoolID         string    `json:"school_id"`
    NamaBank         string    `json:"nama_bank"`
    NomorRekening    string    `json:"nomor_rekening"`
    AtasNama         string    `json:"atas_nama"`
    Cabang           string    `json:"cabang,omitempty"`
    SaldoAwal        float64   `json:"saldo_awal"`
    SaldoTerkini     float64   `json:"saldo_terkini"`
    KodeRekeningID   string    `json:"kode_rekening_id,omitempty"`
    FormatFileMutasi string    `json:"format_file_mutasi,omitempty"`
    IsAktif          bool      `json:"is_aktif"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
}

type RekeningBankFilter struct {
    SchoolID string
    IsAktif  *bool
    Search   string
    Limit    int
    Offset   int
}

// =====================================================================
// Constants — Kategori Iuran & Tipe COA
// =====================================================================
const (
    KategoriSPP      = "SPP"
    KategoriGedung   = "GEDUNG"
    KategoriEkstra   = "EKSTRA"
    KategoriSeragam  = "SERAGAM"
    KategoriKegiatan = "KEGIATAN"
    KategoriLainnya  = "LAINNYA"
)

const (
    TipeAset       = "ASET"
    TipeKewajiban  = "KEWAJIBAN"
    TipeEkuitas    = "EKUITAS"
    TipePendapatan = "PENDAPATAN"
    TipeBelanja    = "BELANJA"
)

// =====================================================================
// Helper — Validasi Kategori & Tipe
// =====================================================================
func IsValidKategoriIuran(k string) bool {
    switch k {
    case KategoriSPP, KategoriGedung, KategoriEkstra, KategoriSeragam, KategoriKegiatan, KategoriLainnya:
        return true
    }
    return false
}

func IsValidTipeKodeRekening(t string) bool {
    switch t {
    case TipeAset, TipeKewajiban, TipeEkuitas, TipePendapatan, TipeBelanja:
        return true
    }
    return false
}
