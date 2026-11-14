package models

import (
	"errors"
	"strings"
	"time"
)

// =====================================================================
// Filter structs untuk list endpoints
// =====================================================================

type TagihanFilter struct {
	SchoolID         string
	AcademicPeriodID string
	NIS              string
	JenisIuranID     string
	Status           string
	Bulan            int
	Tahun            int
	Limit            int
	Offset           int
}

type PenerimaanFilter struct {
	SchoolID         string
	AcademicPeriodID string
	NIS              string
	JenisIuranID     string
	TanggalDari      string
	TanggalSampai    string
	Limit            int
	Offset           int
}

type BukuKasFilter struct {
	SchoolID         string
	AcademicPeriodID string
	KodeRekeningID   string
	TanggalDari      string
	TanggalSampai    string
	Limit            int
	Offset           int
}

type PengajuanFilter struct {
	SchoolID         string
	AcademicPeriodID string
	Status           string
	PemohonUserID    string
	Limit            int
	Offset           int
}

type PengeluaranFilter struct {
	SchoolID         string
	AcademicPeriodID string
	Status           string
	TanggalDari      string
	TanggalSampai    string
	Limit            int
	Offset           int
}

// =====================================================================
// Validators
// =====================================================================

func IsValidStatusTagihan(s string) bool {
	switch s {
	case TagihanBelumLunas, TagihanLunas, TagihanSebagian, TagihanDibatalkan:
		return true
	}
	return false
}

func IsValidMetodeBayar(s string) bool {
	switch s {
	case MetodeTunai, MetodeTransfer, MetodeQRIS:
		return true
	}
	return false
}

func IsValidPrioritas(s string) bool {
	switch s {
	case PrioritasRendah, PrioritasSedang, PrioritasTinggi, PrioritasDarurat:
		return true
	}
	return false
}

// IsValidTransisiPengajuan memeriksa apakah transisi status diizinkan.
// Workflow Opsi A:
//   MENUNGGU_VERIFIKASI_BENDAHARA → MENUNGGU_PERSETUJUAN_KEPALA | DITOLAK_BENDAHARA | PERLU_REVISI
//   MENUNGGU_PERSETUJUAN_KEPALA   → DISETUJUI | DITOLAK_KEPALA
//   DISETUJUI                     → CAIR
func IsValidTransisiPengajuan(dari, ke string) bool {
	switch dari {
	case PengajuanMenungguBendahara:
		return ke == PengajuanMenungguKepala || ke == PengajuanDitolakBendahara || ke == PengajuanPerluRevisi
	case PengajuanMenungguKepala:
		return ke == PengajuanDisetujui || ke == PengajuanDitolakKepala
	case PengajuanDisetujui:
		return ke == PengajuanCair
	}
	return false
}

// ParseDateOrEmpty — helper parse tanggal; kosong = nil
func ParseDateOrEmpty(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return "", errors.New("format tanggal harus YYYY-MM-DD")
	}
	return t.Format("2006-01-02"), nil
}
