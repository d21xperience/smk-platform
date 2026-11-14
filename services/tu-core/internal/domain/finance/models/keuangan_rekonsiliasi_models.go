package models

import "time"

// =====================================================================
// Constants — Status & Tipe
// =====================================================================
const (
	RekonBalance   = "BALANCE"
	RekonUnbalance = "UNBALANCE"
)

const (
	UnmatchedTipeBank = "BANK"
	UnmatchedTipeBuku = "BUKU"
)

const (
	UnmatchedPending   = "PENDING"
	UnmatchedDiabaikan = "DIABAIKAN"
	UnmatchedDikoreksi = "DIKOREKSI"
)

// Tolerance — batas hari untuk auto-match (|tgl_mutasi - tgl_bku|)
const RekonMatchToleranceHari = 2

// =====================================================================
// Structs — Table representations
// =====================================================================

// RekonsiliasiBank — header 1 rekonsiliasi per rekening × bulan × tahun
type RekonsiliasiBank struct {
	ID                     string                  `json:"id"`
	SchoolID               string                  `json:"school_id"`
	AcademicPeriodID       string                  `json:"academic_period_id"`
	RekeningBankID         string                  `json:"rekening_bank_id"`
	RekeningBankNama       string                  `json:"rekening_bank_nama,omitempty"`
	PeriodeBulan           int                     `json:"periode_bulan"`
	PeriodeTahun           int                     `json:"periode_tahun"`
	SaldoBuku              float64                 `json:"saldo_buku"`
	SaldoBank              float64                 `json:"saldo_bank"`
	Selisih                float64                 `json:"selisih"`
	JumlahMatched          int                     `json:"jumlah_matched"`
	JumlahUnmatchedBank    int                     `json:"jumlah_unmatched_bank"`
	JumlahUnmatchedBuku    int                     `json:"jumlah_unmatched_buku"`
	SaldoRekonsiliasiBank  float64                 `json:"saldo_rekonsiliasi_bank"`
	SaldoRekonsiliasiBuku  float64                 `json:"saldo_rekonsiliasi_buku"`
	Status                 string                  `json:"status"`
	FileMutasiURL          string                  `json:"file_mutasi_url,omitempty"`
	Catatan                string                  `json:"catatan,omitempty"`
	CreatedBy              string                  `json:"created_by,omitempty"`
	CreatedAt              time.Time               `json:"created_at"`
	Details                []RekonsiliasiDetail    `json:"details,omitempty"`
	Unmatched              []RekonsiliasiUnmatched `json:"unmatched,omitempty"`
}

// RekonsiliasiDetail — 1 matched pair (mutasi bank ↔ BKU)
type RekonsiliasiDetail struct {
	ID              string    `json:"id"`
	RekonsiliasiID  string    `json:"rekonsiliasi_id"`
	MutasiTanggal   string    `json:"mutasi_tanggal"`
	MutasiKeterangan string   `json:"mutasi_keterangan,omitempty"`
	MutasiDebet     float64   `json:"mutasi_debet"`
	MutasiKredit    float64   `json:"mutasi_kredit"`
	MutasiSaldo     float64   `json:"mutasi_saldo,omitempty"`
	MutasiReferensi string    `json:"mutasi_referensi,omitempty"`
	BukuKasUmumID   string    `json:"buku_kas_umum_id"`
	NoBuktiBKU      string    `json:"no_bukti_bku,omitempty"`
	UraianBKU       string    `json:"uraian_bku,omitempty"`
	SelisihHari     int       `json:"selisih_hari"`
	CreatedAt       time.Time `json:"created_at"`
}

// RekonsiliasiUnmatched — item belum match (BANK atau BUKU)
type RekonsiliasiUnmatched struct {
	ID                  string    `json:"id"`
	RekonsiliasiID      string    `json:"rekonsiliasi_id"`
	Tipe                string    `json:"tipe"`
	DataJSON            string    `json:"data_json"` // raw JSON string
	StatusTindakLanjut  string    `json:"status_tindak_lanjut"`
	CreatedAt           time.Time `json:"created_at"`
}

// =====================================================================
// Input structs — dari client
// =====================================================================

// MutasiInput — 1 baris mutasi bank yang dikirim client
type MutasiInput struct {
	Tanggal    string  `json:"tanggal"`               // YYYY-MM-DD
	Keterangan string  `json:"keterangan,omitempty"`
	Debet      float64 `json:"debet"`                 // masuk
	Kredit     float64 `json:"kredit"`                // keluar
	Saldo      float64 `json:"saldo,omitempty"`
	Referensi  string  `json:"referensi,omitempty"`
}

// CreateRekonsiliasiRequest — POST payload
type CreateRekonsiliasiRequest struct {
	SchoolID         string        `json:"-"`
	AcademicPeriodID string        `json:"academic_period_id"`
	RekeningBankID   string        `json:"rekening_bank_id"`
	PeriodeBulan     int           `json:"periode_bulan"`
	PeriodeTahun     int           `json:"periode_tahun"`
	MutasiItems      []MutasiInput `json:"mutasi_items"`
	Catatan          string        `json:"catatan,omitempty"`
	CreatedBy        string        `json:"-"`
}

// RekonsiliasiFilter — list filter
type RekonsiliasiFilter struct {
	SchoolID         string
	AcademicPeriodID string
	RekeningBankID   string
	PeriodeTahun     int
	Status           string
	Limit            int
	Offset           int
}

// =====================================================================
// Internal structs — untuk auto-match algorithm
// =====================================================================

// bkuForMatch — snapshot BKU minimal untuk matching
type BkuForMatch struct {
	ID          string
	Tanggal     time.Time
	NoBukti     string
	Uraian      string
	Debet       float64
	Kredit      float64
	SaldoBerjalan float64
	KodeRekeningID string
}

// mutasiParsed — mutasi dengan tanggal sudah di-parse
type MutasiParsed struct {
	Tanggal    time.Time
	Keterangan string
	Debet      float64
	Kredit     float64
	Saldo      float64
	Referensi  string
}

// MatchedPair — hasil auto-match
type MatchedPair struct {
	Mutasi     MutasiParsed
	BKU        BkuForMatch
	SelisihHari int
}

// UnmatchedBankItem — mutasi bank tidak match
type UnmatchedBankItem struct {
	Mutasi MutasiParsed
}

// UnmatchedBukuItem — BKU tidak match
type UnmatchedBukuItem struct {
	BKU BkuForMatch
}

// =====================================================================
// Validators
// =====================================================================
func IsValidStatusRekonsiliasi(s string) bool {
	return s == RekonBalance || s == RekonUnbalance
}

func IsValidTipeUnmatched(s string) bool {
	return s == UnmatchedTipeBank || s == UnmatchedTipeBuku
}
