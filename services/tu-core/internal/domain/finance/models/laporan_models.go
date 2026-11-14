package models

import "time"

// =====================================================================
// Constants
// =====================================================================
const (
	LaporanModeLive     = "LIVE"
	LaporanModeSnapshot = "SNAPSHOT"
)

const (
	AlertSeverityInfo    = "INFO"
	AlertSeverityWarning = "WARNING"
	AlertSeverityDanger  = "DANGER"
)

const (
	AlertCategoryPiutang       = "PIUTANG"
	AlertCategoryPengajuan     = "PENGAJUAN"
	AlertCategoryRekonsiliasi  = "REKONSILIASI"
	AlertCategoryRab           = "RAB"
	AlertCategoryKas           = "KAS"
)

// =====================================================================
// Neraca (Posisi Keuangan)
// =====================================================================

// NeracaItem — 1 baris saldo per COA level 3
//
// Saldo = natural balance:
//   - Aset, Belanja     : debet - kredit
//   - Kewajiban, Ekuitas, Pendapatan : kredit - debet
type NeracaItem struct {
	KodeRekeningID string  `json:"kode_rekening_id"`
	Kode           string  `json:"kode"`
	Nama           string  `json:"nama"`
	Tipe           string  `json:"tipe"`
	Saldo          float64 `json:"saldo"`
}

// NeracaReport — laporan neraca lengkap
type NeracaReport struct {
	AsOfDate       string       `json:"as_of_date"`
	Mode           string       `json:"mode"` // LIVE | SNAPSHOT
	PenutupanID    string       `json:"penutupan_id,omitempty"`
	PeriodeBulan   int          `json:"periode_bulan,omitempty"`
	PeriodeTahun   int          `json:"periode_tahun,omitempty"`

	Aset       []NeracaItem `json:"aset"`
	Kewajiban  []NeracaItem `json:"kewajiban"`
	Ekuitas    []NeracaItem `json:"ekuitas"`
	Pendapatan []NeracaItem `json:"pendapatan"`
	Belanja    []NeracaItem `json:"belanja"`

	TotalAset       float64 `json:"total_aset"`
	TotalKewajiban  float64 `json:"total_kewajiban"`
	TotalEkuitas    float64 `json:"total_ekuitas"`
	TotalPendapatan float64 `json:"total_pendapatan"`
	TotalBelanja    float64 `json:"total_belanja"`

	IsBalance bool    `json:"is_balance"` // Aset == Kewajiban + Ekuitas
	Selisih   float64 `json:"selisih"`    // Aset - (Kewajiban + Ekuitas)

	// Surplus/defisit periode berjalan (Pendapatan - Belanja)
	SurplusDefisit float64 `json:"surplus_defisit"`
}

// =====================================================================
// LRA — Laporan Realisasi Anggaran (wrapper dari RAB)
// =====================================================================
type LraReport struct {
	TahunAnggaran     int     `json:"tahun_anggaran"`
	SampaiBulan       int     `json:"sampai_bulan"`
	RabID             string  `json:"rab_id"`
	RabStatus         string  `json:"rab_status"`
	TanggalPengesahan string  `json:"tanggal_pengesahan,omitempty"`

	// Realisasi di-reuse dari struct RAB existing (RabRealisasiReport)
	Realisasi *RabRealisasiReport `json:"realisasi"`
}

// =====================================================================
// Dashboard
// =====================================================================
type DashboardAlert struct {
	Severity string `json:"severity"` // INFO | WARNING | DANGER
	Category string `json:"category"`
	Message  string `json:"message"`
}

type DashboardKasInfo struct {
	SaldoTerkini        float64 `json:"saldo_terkini"`
	TotalDebetBulanIni  float64 `json:"total_debet_bulan_ini"`
	TotalKreditBulanIni float64 `json:"total_kredit_bulan_ini"`
}

type DashboardPiutangInfo struct {
	TotalSisa    float64       `json:"total_sisa"`
	JumlahItem   int           `json:"jumlah_item"`
	JumlahSiswa  int           `json:"jumlah_siswa"`
	AgingSummary []AgingBucket `json:"aging_summary"`
}

type DashboardRabInfo struct {
	TahunAnggaran    int     `json:"tahun_anggaran"`
	Status           string  `json:"status"`
	TotalPaguBelanja float64 `json:"total_pagu_belanja"`
	RealisasiBelanja float64 `json:"realisasi_belanja"`
	PersenRealisasi  float64 `json:"persen_realisasi"`
}

type DashboardPengajuanInfo struct {
	MenungguVerifikasi  int     `json:"menunggu_verifikasi"`
	MenungguPersetujuan int     `json:"menunggu_persetujuan"`
	Disetujui           int     `json:"disetujui"`
	TotalNominalPending float64 `json:"total_nominal_pending"`
}

type DashboardRekonsiliasiInfo struct {
	TotalRekonBulanIni int `json:"total_rekon_bulan_ini"`
	Balance            int `json:"balance"`
	Unbalance          int `json:"unbalance"`
}

type DashboardSummary struct {
	AsOfDate      string                    `json:"as_of_date"`
	Kas           DashboardKasInfo          `json:"kas"`
	Piutang       DashboardPiutangInfo      `json:"piutang"`
	Rab           *DashboardRabInfo         `json:"rab,omitempty"` // nil kalau belum ada RAB tahun ini
	Pengajuan     DashboardPengajuanInfo    `json:"pengajuan"`
	Rekonsiliasi  DashboardRekonsiliasiInfo `json:"rekonsiliasi"`
	Alerts        []DashboardAlert          `json:"alerts"`
}

// =====================================================================
// Ringkasan Bulanan
// =====================================================================
type RingkasanKategoriItem struct {
	KodeRekening    string  `json:"kode_rekening"`
	Nama            string  `json:"nama"`
	Total           float64 `json:"total"`
	JumlahTransaksi int     `json:"jumlah_transaksi"`
}

type RingkasanBulan struct {
	PeriodeBulan    int                     `json:"periode_bulan"`
	PeriodeTahun    int                     `json:"periode_tahun"`
	SaldoAwal       float64                 `json:"saldo_awal"`
	TotalDebet      float64                 `json:"total_debet"`
	TotalKredit     float64                 `json:"total_kredit"`
	SaldoAkhir      float64                 `json:"saldo_akhir"`
	JumlahTransaksi int                     `json:"jumlah_transaksi"`
	TopDebet        []RingkasanKategoriItem `json:"top_debet"`
	TopKredit       []RingkasanKategoriItem `json:"top_kredit"`
}

// =====================================================================
// Filter
// =====================================================================
type LaporanFilter struct {
	SchoolID         string
	AcademicPeriodID string
	Tanggal          string // YYYY-MM-DD, default today (untuk neraca)
	Tahun            int    // tahun anggaran (untuk LRA)
	SampaiBulan      int    // 1-12 (untuk LRA)
	Bulan            int    // 1-12 (untuk ringkasan)
	Mode             string // "live" | "snapshot" (untuk neraca)
	PenutupanID      string // untuk mode snapshot
}

// =====================================================================
// Validator
// =====================================================================
func IsValidLaporanMode(s string) bool {
	return s == "" || s == LaporanModeLive || s == LaporanModeSnapshot
}

// NotUsedTime — suppress unused-import kalau time belum dipakai
var _ = time.Now
