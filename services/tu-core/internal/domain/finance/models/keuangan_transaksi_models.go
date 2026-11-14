package models

import "time"

// =====================================================================
// Constants — Status & Enum (self-documenting, dipakai sebagai CHECK
// constraint di migration dan sebagai validasi di service)
// =====================================================================

// Status Tagihan
const (
	TagihanBelumLunas = "BELUM_LUNAS"
	TagihanLunas      = "LUNAS"
	TagihanSebagian   = "SEBAGIAN"
	TagihanDibatalkan = "DIBATALKAN"
)

// Metode Bayar
const (
	MetodeTunai    = "TUNAI"
	MetodeTransfer = "TRANSFER"
	MetodeQRIS     = "QRIS"
)

// Status Pengajuan Pengeluaran (workflow 3 tingkat Opsi A)
const (
	PengajuanMenungguBendahara = "MENUNGGU_VERIFIKASI_BENDAHARA"
	PengajuanDitolakBendahara  = "DITOLAK_BENDAHARA"
	PengajuanPerluRevisi       = "PERLU_REVISI"
	PengajuanMenungguKepala    = "MENUNGGU_PERSETUJUAN_KEPALA"
	PengajuanDitolakKepala     = "DITOLAK_KEPALA"
	PengajuanDisetujui         = "DISETUJUI"
	PengajuanCair              = "CAIR"
	PengajuanDibatalkan        = "DIBATALKAN"
)

// Status Pengeluaran
const (
	PengeluaranSiapCair = "SIAP_CAIR"
	PengeluaranCair     = "CAIR"
	PengeluaranBatal    = "DIBATALKAN"
)

// Prioritas Pengajuan
const (
	PrioritasRendah  = "RENDAH"
	PrioritasSedang  = "SEDANG"
	PrioritasTinggi  = "TINGGI"
	PrioritasDarurat = "DARURAT"
)

// =====================================================================
// Structs — 5 tabel transaksi
// =====================================================================

// TagihanSiswa — tagihan per siswa × jenis iuran × periode
type TagihanSiswa struct {
	ID                string    `json:"id"`
	SchoolID          string    `json:"school_id"`
	AcademicPeriodID  string    `json:"academic_period_id"`
	PesertaDidikID    string    `json:"peserta_didik_id"`
	NIS               string    `json:"nis"`
	KelasNama         string    `json:"kelas_nama"`
	JenisIuranID      string    `json:"jenis_iuran_id"`
	JenisIuranNama    string    `json:"jenis_iuran_nama,omitempty"`
	PeriodeBulan      int       `json:"periode_bulan"`
	PeriodeTahun      int       `json:"periode_tahun"`
	Nominal           float64   `json:"nominal"`
	TotalDibayar      float64   `json:"total_dibayar"`
	Status            string    `json:"status"`
	TanggalJatuhTempo string    `json:"tanggal_jatuh_tempo,omitempty"`
	TanggalBayar      string    `json:"tanggal_bayar,omitempty"`
	CreatedBy         string    `json:"created_by,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// PenerimaanIuran — transaksi pembayaran (auto BKM)
type PenerimaanIuran struct {
	ID                string    `json:"id"`
	SchoolID          string    `json:"school_id"`
	AcademicPeriodID  string    `json:"academic_period_id"`
	NoBKM             string    `json:"no_bkm"`
	Tanggal           string    `json:"tanggal"`
	TagihanSiswaID    string    `json:"tagihan_siswa_id"`
	PesertaDidikID    string    `json:"peserta_didik_id"`
	NIS               string    `json:"nis"`
	JenisIuranID      string    `json:"jenis_iuran_id"`
	JenisIuranNama    string    `json:"jenis_iuran_nama,omitempty"`
	PeriodeBulan      int       `json:"periode_bulan"`
	PeriodeTahun      int       `json:"periode_tahun"`
	Jumlah            float64   `json:"jumlah"`
	MetodeBayar       string    `json:"metode_bayar"`
	NoReferensi       string    `json:"no_referensi,omitempty"`
	Keterangan        string    `json:"keterangan,omitempty"`
	Status            string    `json:"status"`
	CreatedBy         string    `json:"created_by,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// BukuKasUmum — BKU dengan running balance
type BukuKasUmum struct {
	ID               string    `json:"id"`
	SchoolID         string    `json:"school_id"`
	AcademicPeriodID string    `json:"academic_period_id"`
	Tanggal          string    `json:"tanggal"`
	NoBukti          string    `json:"no_bukti"`
	Uraian           string    `json:"uraian"`
	KodeRekeningID   string    `json:"kode_rekening_id"`
	KodeRekeningKode string    `json:"kode_rekening_kode,omitempty"`
	RekeningBankID   string    `json:"rekening_bank_id,omitempty"`
	Debet            float64   `json:"debet"`
	Kredit           float64   `json:"kredit"`
	SaldoBerjalan    float64   `json:"saldo_berjalan"`
	ReferensiTabel   string    `json:"referensi_tabel,omitempty"`
	ReferensiID      string    `json:"referensi_id,omitempty"`
	IsLocked         bool      `json:"is_locked"`
	IsRekonsiliasi   bool      `json:"is_rekonsiliasi"`
	CreatedBy        string    `json:"created_by,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

// PengajuanPengeluaran — workflow 3 tingkat
type PengajuanPengeluaran struct {
	ID                string     `json:"id"`
	SchoolID          string     `json:"school_id"`
	AcademicPeriodID  string     `json:"academic_period_id"`
	NomorPengajuan    string     `json:"nomor_pengajuan"`
	Tanggal           string     `json:"tanggal"`
	PemohonUserID     string     `json:"pemohon_user_id"`
	PemohonNama       string     `json:"pemohon_nama,omitempty"`
	Keperluan         string     `json:"keperluan"`
	Jumlah            float64    `json:"jumlah"`
	Prioritas         string     `json:"prioritas"`
	Status            string     `json:"status"`
	CatatanPemohon    string     `json:"catatan_pemohon,omitempty"`
	CatatanBendahara  string     `json:"catatan_bendahara,omitempty"`
	CatatanKepala     string     `json:"catatan_kepala,omitempty"`
	VerifiedBy        string     `json:"verified_by,omitempty"`
	VerifiedAt        *time.Time `json:"verified_at,omitempty"`
	ApprovedBy        string     `json:"approved_by,omitempty"`
	ApprovedAt        *time.Time `json:"approved_at,omitempty"`
	NoBKK             string     `json:"no_bkk,omitempty"`
	PosAnggaranID     string     `json:"pos_anggaran_id,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// Pengeluaran — transaksi kas keluar (auto BKK saat DISETUJUI)
type Pengeluaran struct {
	ID               string    `json:"id"`
	SchoolID         string    `json:"school_id"`
	AcademicPeriodID string    `json:"academic_period_id"`
	NoBKK            string    `json:"no_bkk"`
	Tanggal          string    `json:"tanggal"`
	PengajuanID      string    `json:"pengajuan_id"`
	Jumlah           float64   `json:"jumlah"`
	Keperluan        string    `json:"keperluan"`
	Status           string    `json:"status"`
	CreatedBy        string    `json:"created_by,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}


// SiswaInfo — subset tabel_siswa yang diperlukan untuk transaksi keuangan
type SiswaInfo struct {
	PesertaDidikID string `json:"peserta_didik_id"`
	NIS            string `json:"nis"`
	NamaSiswa      string `json:"nama_siswa"`
	KelasNama      string `json:"kelas_nama"`
}
