package models

// =====================================================================
// Constants — Aging Group
// =====================================================================
const (
	PiutangAgingBelumJatuhTempo = "BELUM_JATUH_TEMPO"
	PiutangAgingLessThan30      = "<30"
	PiutangAging30to60          = "30-60"
	PiutangAging60to90          = "60-90"
	PiutangAgingMoreThan90      = ">90"
)

// =====================================================================
// PiutangItem — 1 baris piutang (tagihan BELUM_LUNAS/SEBAGIAN)
// =====================================================================
type PiutangItem struct {
	ID               string  `json:"id"`
	SchoolID         string  `json:"school_id"`
	AcademicPeriodID string  `json:"academic_period_id"`
	NIS              string  `json:"nis"`
	NamaSiswa        string  `json:"nama_siswa,omitempty"`
	KelasNama        string  `json:"kelas_nama,omitempty"`
	JenisIuranID     string  `json:"jenis_iuran_id"`
	JenisIuranNama   string  `json:"jenis_iuran_nama,omitempty"`
	PeriodeBulan     int     `json:"periode_bulan"`
	PeriodeTahun     int     `json:"periode_tahun"`
	Nominal          float64 `json:"nominal"`
	TotalDibayar     float64 `json:"total_dibayar"`
	Sisa             float64 `json:"sisa"`
	Status           string  `json:"status"`
	TanggalJatuhTempo string `json:"tanggal_jatuh_tempo,omitempty"`
	HariTerlambat    int     `json:"hari_terlambat"`
	AgingGroup       string  `json:"aging_group"`
}

// =====================================================================
// PiutangPerSiswa — grouping per NIS dengan items
// =====================================================================
type PiutangPerSiswa struct {
	NIS              string        `json:"nis"`
	NamaSiswa        string        `json:"nama_siswa,omitempty"`
	KelasNama        string        `json:"kelas_nama,omitempty"`
	TotalSisa        float64       `json:"total_sisa"`
	JumlahTagihan    int           `json:"jumlah_tagihan"`
	OldestAgingGroup string        `json:"oldest_aging_group,omitempty"`
	Items            []PiutangItem `json:"items"`
}

// =====================================================================
// AgingBucket — 1 bucket ringkasan aging
// =====================================================================
type AgingBucket struct {
	Group      string        `json:"group"`
	ItemCount  int           `json:"item_count"`
	TotalSisa  float64       `json:"total_sisa"`
	Items      []PiutangItem `json:"items,omitempty"`
}

// =====================================================================
// AgingReport — full aging report
// =====================================================================
type AgingReport struct {
	AsOfDate       string        `json:"as_of_date"`
	TotalItemCount int           `json:"total_item_count"`
	TotalSisa      float64       `json:"total_sisa"`
	Buckets        []AgingBucket `json:"buckets"`
}

// =====================================================================
// PiutangFilter — filter untuk list
// =====================================================================
type PiutangFilter struct {
	SchoolID         string
	AcademicPeriodID string
	NIS              string
	KelasNama        string
	AgingGroup       string
	JenisIuranID     string
	Limit            int
	Offset           int
}

// =====================================================================
// Validator
// =====================================================================
func IsValidAgingGroup(s string) bool {
	switch s {
	case PiutangAgingBelumJatuhTempo, PiutangAgingLessThan30,
		PiutangAging30to60, PiutangAging60to90, PiutangAgingMoreThan90:
		return true
	}
	return false
}

// AllAgingGroups — urutan bucket untuk konsistensi output
func AllAgingGroups() []string {
	return []string{
		PiutangAgingBelumJatuhTempo,
		PiutangAgingLessThan30,
		PiutangAging30to60,
		PiutangAging60to90,
		PiutangAgingMoreThan90,
	}
}
