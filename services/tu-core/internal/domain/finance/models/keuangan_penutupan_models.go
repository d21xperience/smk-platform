package models

import "time"

// =====================================================================
// Constants — Status Penutupan
// =====================================================================
const (
	PenutupanClosed             = "CLOSED"
	PenutupanMenungguTTDKepala  = "MENUNGGU_TTD_KEPALA"
	PenutupanFullySigned        = "CLOSED_FULLY_SIGNED"
	PenutupanDibukaKembali      = "DIBUKA_KEMBALI"
)

// =====================================================================
// PenutupanBukuKas — snapshot penutupan 1 bulan per school×period
// =====================================================================
type PenutupanBukuKas struct {
	ID                string     `json:"id"`
	SchoolID          string     `json:"school_id"`
	AcademicPeriodID  string     `json:"academic_period_id"`
	PeriodeBulan      int        `json:"periode_bulan"`
	PeriodeTahun      int        `json:"periode_tahun"`

	// Snapshot buku kas
	SaldoAwal         float64    `json:"saldo_awal"`
	TotalDebet        float64    `json:"total_debet"`
	TotalKredit       float64    `json:"total_kredit"`
	SaldoAkhirBuku    float64    `json:"saldo_akhir_buku"`

	// Cash count
	SaldoFisikKas     float64    `json:"saldo_fisik_kas"`
	SelisihKas        float64    `json:"selisih_kas"`

	// Snapshot posisi keuangan
	TotalAset         float64    `json:"total_aset"`
	TotalKewajiban    float64    `json:"total_kewajiban"`
	TotalEkuitas      float64    `json:"total_ekuitas"`

	// Status & audit
	Status            string     `json:"status"`
	Catatan           string     `json:"catatan,omitempty"`
	ClosedBy          string     `json:"closed_by,omitempty"`
	ClosedByNama      string     `json:"closed_by_nama,omitempty"`
	ClosedAt          *time.Time `json:"closed_at,omitempty"`

	// E-sign (diisi FASE 5)
	PinHash           string     `json:"-"`
	TTDKepsekHash     string     `json:"ttd_kepsek_hash,omitempty"`
	TTDKepsekBy       string     `json:"ttd_kepsek_by,omitempty"`
	TTDKepsekNama     string     `json:"ttd_kepsek_nama,omitempty"`
	TTDKepsekAt       *time.Time `json:"ttd_kepsek_at,omitempty"`

	// Dokumen PDF (diisi FASE 5)
	BkuPdfURL         string     `json:"bku_pdf_url,omitempty"`
	NeracaPdfURL      string     `json:"neraca_pdf_url,omitempty"`
	LraPdfURL         string     `json:"lra_pdf_url,omitempty"`

	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// =====================================================================
// Input structs
// =====================================================================

// CreatePenutupanRequest — payload POST /penutupan
type CreatePenutupanRequest struct {
	SchoolID         string  `json:"-"`
	AcademicPeriodID string  `json:"academic_period_id"`
	PeriodeBulan     int     `json:"periode_bulan"`
	PeriodeTahun     int     `json:"periode_tahun"`
	SaldoFisikKas    float64 `json:"saldo_fisik_kas"`
	Catatan          string  `json:"catatan,omitempty"`
	ClosedBy         string  `json:"-"`
}

// PenutupanFilter — list filter
type PenutupanFilter struct {
	SchoolID         string
	AcademicPeriodID string
	PeriodeTahun     int
	Status           string
	Limit            int
	Offset           int
}

// =====================================================================
// Internal structs — snapshot computation
// =====================================================================

// SnapshotBku — agregat BKU untuk bulan tertentu (dihitung service)
type SnapshotBku struct {
	SaldoAwal      float64
	TotalDebet     float64
	TotalKredit    float64
	SaldoAkhirBuku float64
}

// SnapshotPosisi — snapshot aset/kewajiban/ekuitas (dihitung service)
type SnapshotPosisi struct {
	TotalAset      float64
	TotalKewajiban float64
	TotalEkuitas   float64
}

// =====================================================================
// Validators
// =====================================================================
func IsValidStatusPenutupan(s string) bool {
	switch s {
	case PenutupanClosed, PenutupanMenungguTTDKepala, PenutupanFullySigned, PenutupanDibukaKembali:
		return true
	}
	return false
}
