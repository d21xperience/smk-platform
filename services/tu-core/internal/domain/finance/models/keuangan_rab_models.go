package models

import (
	"errors"
	"fmt"
	"time"
)

// =====================================================================
// Constants — Status RAB & Sumber Dana
// =====================================================================
const (
	RabDraft    = "DRAFT"
	RabDiajukan = "DIAJUKAN"
	RabDisahkan = "DISAHKAN"
	RabDirevisi = "DIREVISI"
)

const (
	SumberDanaBOS     = "BOS"
	SumberDanaBOP     = "BOP"
	SumberDanaHibah   = "HIBAH"
	SumberDanaMandiri = "MANDIRI"
)

// =====================================================================
// Structs — RAB Header, Detail, Realisasi
// =====================================================================

// RabHeader — 1 RAB per school × period × tahun anggaran
type RabHeader struct {
	ID                string     `json:"id"`
	SchoolID          string     `json:"school_id"`
	AcademicPeriodID  string     `json:"academic_period_id"`
	TahunAnggaran     int        `json:"tahun_anggaran"`
	TotalPendapatan   float64    `json:"total_pendapatan"`
	TotalBelanja      float64    `json:"total_belanja"`
	SurplusDefisit    float64    `json:"surplus_defisit"`
	Status            string     `json:"status"`
	DisusunOleh       string     `json:"disusun_oleh,omitempty"`
	DisusunOlehNama   string     `json:"disusun_oleh_nama,omitempty"`
	TanggalPenyusunan string     `json:"tanggal_penyusunan,omitempty"`
	DisahkanOleh      string     `json:"disahkan_oleh,omitempty"`
	DisahkanOlehNama  string     `json:"disahkan_oleh_nama,omitempty"`
	TanggalPengesahan string     `json:"tanggal_pengesahan,omitempty"`
	Catatan           string     `json:"catatan,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	Details           []RabDetail `json:"details,omitempty"`
}

// RabDetail — 1 item RAB per kode_rekening
type RabDetail struct {
	ID               string    `json:"id"`
	RabHeaderID      string    `json:"rab_header_id"`
	KodeRekeningID   string    `json:"kode_rekening_id"`
	KodeRekeningKode string    `json:"kode_rekening_kode,omitempty"`
	KodeRekeningNama string    `json:"kode_rekening_nama,omitempty"`
	Tipe             string    `json:"tipe,omitempty"`
	Uraian           string    `json:"uraian"`
	Volume           float64   `json:"volume"`
	Satuan           string    `json:"satuan"`
	HargaSatuan      float64   `json:"harga_satuan"`
	Subtotal         float64   `json:"subtotal"`
	SumberDana       string    `json:"sumber_dana,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// RabDetailInput — payload input item RAB dari client
type RabDetailInput struct {
	KodeRekeningID string  `json:"kode_rekening_id"`
	Uraian         string  `json:"uraian"`
	Volume         float64 `json:"volume"`
	Satuan         string  `json:"satuan"`
	HargaSatuan    float64 `json:"harga_satuan"`
	SumberDana     string  `json:"sumber_dana,omitempty"`
}

// HitungSubtotal — volume × harga_satuan
func (d RabDetailInput) HitungSubtotal() float64 {
	return d.Volume * d.HargaSatuan
}

// CreateRabRequest — payload create/update RAB
type CreateRabRequest struct {
	SchoolID         string           `json:"-"`
	AcademicPeriodID string           `json:"academic_period_id"`
	TahunAnggaran    int              `json:"tahun_anggaran"`
	Catatan          string           `json:"catatan,omitempty"`
	Items            []RabDetailInput `json:"items"`
	DisusunOleh      string           `json:"-"`
}

// SahkanRabRequest — payload sahkan RAB
type SahkanRabRequest struct {
	SchoolID    string `json:"-"`
	RabID       string `json:"-"`
	DisahkanOleh string `json:"disahkan_oleh"`
	Catatan     string `json:"catatan,omitempty"`
}

// RabFilter — filter list RAB
type RabFilter struct {
	SchoolID         string
	AcademicPeriodID string
	TahunAnggaran    int
	Status           string
	Limit            int
	Offset           int
}

// =====================================================================
// Structs — Realisasi Report
// =====================================================================

// RabRealisasiItem — 1 baris realisasi per pos
type RabRealisasiItem struct {
	RabDetailID      string  `json:"rab_detail_id"`
	KodeRekeningID   string  `json:"kode_rekening_id"`
	KodeRekeningKode string  `json:"kode_rekening_kode"`
	Uraian           string  `json:"uraian"`
	Tipe             string  `json:"tipe"`
	Pagu             float64 `json:"pagu"`
	Realisasi        float64 `json:"realisasi"`
	Sisa             float64 `json:"sisa"`
	PersenRealisasi  float64 `json:"persen_realisasi"`
	StatusWarning    string  `json:"status_warning"`
}

// RabRealisasiReport — laporan realisasi per tahun
type RabRealisasiReport struct {
	TahunAnggaran              int                `json:"tahun_anggaran"`
	Items                      []RabRealisasiItem `json:"items"`
	TotalPaguPendapatan        float64            `json:"total_pagu_pendapatan"`
	TotalRealisasiPendapatan   float64            `json:"total_realisasi_pendapatan"`
	TotalPaguBelanja           float64            `json:"total_pagu_belanja"`
	TotalRealisasiBelanja      float64            `json:"total_realisasi_belanja"`
	TotalSisaPendapatan        float64            `json:"total_sisa_pendapatan"`
	TotalSisaBelanja           float64            `json:"total_sisa_belanja"`
}

// =====================================================================
// Validators
// =====================================================================
func IsValidStatusRab(s string) bool {
	switch s {
	case RabDraft, RabDiajukan, RabDisahkan, RabDirevisi:
		return true
	}
	return false
}

func IsValidSumberDana(s string) bool {
	if s == "" {
		return true
	}
	switch s {
	case SumberDanaBOS, SumberDanaBOP, SumberDanaHibah, SumberDanaMandiri:
		return true
	}
	return false
}

// ValidateRabDetailInput — validasi per item RAB sebelum hitung
func ValidateRabDetailInput(d RabDetailInput) error {
	if d.KodeRekeningID == "" {
		return errors.New("kode_rekening_id wajib diisi")
	}
	if d.Uraian == "" {
		return errors.New("uraian wajib diisi")
	}
	if d.Volume <= 0 {
		return fmt.Errorf("volume harus > 0 untuk: %s", d.Uraian)
	}
	if d.HargaSatuan <= 0 {
		return fmt.Errorf("harga_satuan harus > 0 untuk: %s", d.Uraian)
	}
	if d.Satuan == "" {
		return fmt.Errorf("satuan wajib diisi untuk: %s", d.Uraian)
	}
	if !IsValidSumberDana(d.SumberDana) {
		return fmt.Errorf("sumber_dana tidak valid (BOS/BOP/HIBAH/MANDIRI): %s", d.SumberDana)
	}
	return nil
}
