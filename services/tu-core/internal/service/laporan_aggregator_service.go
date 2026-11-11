package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"sekolah-platform/services/tu-core/internal/repository"
)

type LaporanData struct {
	TotalPegawai      int     `json:"totalPegawai"`
	RataKehadiran     float64 `json:"rataKehadiran"`
	TotalAlpha        int     `json:"totalAlpha"`
	RataSkp           float64 `json:"rataSkp"`
	TotalCuti         int     `json:"totalCuti"`
	TotalDokumen      int     `json:"totalDokumen"`
	BebanKurang24     int     `json:"bebanKurang24"`
	CatatanTambahan   string  `json:"catatanTambahan"`
}

type LaporanAggregatorService struct {
	statRepo repository.PtkStatisticRepository
	schoolID string
}

func NewLaporanAggregatorService(statRepo repository.PtkStatisticRepository) *LaporanAggregatorService {
	return &LaporanAggregatorService{statRepo: statRepo}
}

// AggregateData memanggil semua repository yang relevan dan menggabungkannya
func (s *LaporanAggregatorService) AggregateData(ctx context.Context, schoolID string, startDate, endDate time.Time) (json.RawMessage, error) {
	// 1. Ambil Statistik Utama (Total Pegawai, Kehadiran, SKP, dll)
	statData, err := s.statRepo.GetDashboardData(ctx, schoolID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("gagal ambil data statistik: %v", err)
	}

	// 2. Susun struktur data laporan
	laporanData := LaporanData{
		TotalPegawai:      statData.KPIs.TotalPegawai,
		RataKehadiran:     statData.KPIs.RataKehadiran,
		TotalAlpha:        statData.KPIs.TotalAlpha,
		RataSkp:           statData.KPIs.RataSkp,
		TotalCuti:         statData.KPIs.CutiPending,
		TotalDokumen:      statData.KPIs.DokumenKadaluarsa,
		BebanKurang24:     statData.KPIs.BebanKurang24,
		CatatanTambahan:   "Data digenerate otomatis oleh sistem.",
	}

	// 3. Konversi ke JSON
	jsonBytes, err := json.Marshal(laporanData)
	if err != nil {
		return nil, fmt.Errorf("gagal marshal data: %v", err)
	}

	return jsonBytes, nil
}
