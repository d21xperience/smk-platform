package service

import (
	"context"
	"time"

	"sekolah-platform/services/tu-core/internal/repository"
)

type BukuKontrolService struct {
	repo repository.BukuKontrolRepository
}

func NewBukuKontrolService(repo repository.BukuKontrolRepository) *BukuKontrolService {
	return &BukuKontrolService{repo: repo}
}

func (s *BukuKontrolService) ListGolongan(ctx context.Context) ([]repository.MasterGolongan, error) {
	return s.repo.ListGolongan(ctx)
}

func (s *BukuKontrolService) GetDashboard(ctx context.Context) (map[string]interface{}, error) {
	list, err := s.repo.ListBukuKontrol(ctx)
	if err != nil {
		return nil, err
	}

	// Hitung KPI
	var total, aman, perluPersiapan, segera, urgent, overdue int
	for _, b := range list {
		total++
		switch b.StatusMonitoring {
		case "AMAN": aman++
		case "PERLU_PERSIAPAN": perluPersiapan++
		case "SEGERA": segera++
		case "URGENT": urgent++
		case "OVERDUE": overdue++
		}
	}

	return map[string]interface{}{
		"kpi": map[string]int{
			"total": total, "aman": aman, "perluPersiapan": perluPersiapan,
			"segera": segera, "urgent": urgent, "overdue": overdue,
		},
		"table": list,
	}, nil
}

func (s *BukuKontrolService) InitializeMassal(ctx context.Context) (int, error) {
	pegawaiList, err := s.repo.GetPegawaiForInit(ctx)
	if err != nil { return 0, err }

	count := 0
	for _, p := range pegawaiList {
		ptkID := p["ptkId"].(string)
		golID := p["golonganId"].(int)
		tmtPangkat := p["tmtPangkat"].(time.Time)
		tmtGaji := p["tmtGaji"].(time.Time)
		tglLahir := p["tanggalLahir"].(time.Time)

		// Kalkulasi Otomatis
		tglKGB, tglPangkat, tglPensiun := calculateDates(tmtPangkat, tmtGaji, golID, tglLahir)
		status := determineStatus(tglKGB, tglPangkat)

		err := s.repo.CreateBukuKontrol(ctx, ptkID, golID, tmtPangkat, tmtGaji, tglKGB, tglPangkat, tglPensiun, status)
		if err == nil { count++ }
	}
	return count, nil
}

// Fungsi Helper: Menghitung Tanggal Berikutnya
func calculateDates(tmtPangkat, tmtGaji time.Time, golID int, tglLahir time.Time) (time.Time, time.Time, time.Time) {
	now := time.Now()

	// 1. Hitung KGB Berikutnya (Interval 2 atau 4 tahun)
	intervalKGB := 4
	if golID >= 11 && golID <= 24 { // Golongan I dan II
		intervalKGB = 2
	}

	tglKGB := tmtGaji
	for tglKGB.Before(now) || tglKGB.Equal(now) {
		tglKGB = tglKGB.AddDate(intervalKGB, 0, 0)
	}

	// 2. Hitung Pangkat Berikutnya (Reguler 4 tahun)
	tglPangkat := tmtPangkat
	for tglPangkat.Before(now) || tglPangkat.Equal(now) {
		tglPangkat = tglPangkat.AddDate(4, 0, 0)
	}

	// 3. Hitung Pensiun (60 tahun)
	tglPensiun := tglLahir.AddDate(60, 0, 0)

	return tglKGB, tglPangkat, tglPensiun
}

// Fungsi Helper: Menentukan Status Monitoring
func determineStatus(tglKGB, tglPangkat time.Time) string {
	now := time.Now()
	daysKGB := int(tglKGB.Sub(now).Hours() / 24)
	daysPangkat := int(tglPangkat.Sub(now).Hours() / 24)

	minDays := daysKGB
	if daysPangkat < minDays { minDays = daysPangkat }

	if minDays < 0 { return "OVERDUE" }
	if minDays <= 7 { return "URGENT" }
	if minDays <= 30 { return "SEGERA" }
	if minDays <= 90 { return "PERLU_PERSIAPAN" }
	return "AMAN"
}

func (s *BukuKontrolService) GetByPtkID(ctx context.Context, ptkID string) (*repository.BukuKontrol, error) {
	return s.repo.GetBukuKontrolByPtkID(ctx, ptkID)
}
