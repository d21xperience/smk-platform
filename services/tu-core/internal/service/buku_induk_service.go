package service

import (
	"context"

	"sekolah-platform/services/tu-core/internal/repository"
)

type BukuIndukService struct {
	repo repository.BukuIndukRepository
}

func NewBukuIndukService(repo repository.BukuIndukRepository) *BukuIndukService {
	return &BukuIndukService{repo: repo}
}

func (s *BukuIndukService) ListBukuInduk(ctx context.Context) ([]repository.BukuIndukDetail, error) {
	list, err := s.repo.ListBukuInduk(ctx)
	if err != nil {
		return nil, err
	}
	for i := range list {
		calculateCompleteness(&list[i])
	}
	return list, nil
}

func (s *BukuIndukService) GetBukuIndukDetail(ctx context.Context, ptkID string) (*repository.BukuIndukDetail, error) {
	detail, err := s.repo.GetBukuIndukDetail(ctx, ptkID)
	if err != nil {
		return nil, err
	}
	calculateCompleteness(detail)
	return detail, nil
}

// Logika sederhana untuk menghitung skor (Total 100)
func calculateCompleteness(d *repository.BukuIndukDetail) {
	skor := 0
	// Identitas Dasar (40 poin)
	if d.Name != "" { skor += 10 }
	if d.NIK != "" { skor += 10 }
	if d.Gender != "" { skor += 5 }
	if d.BirthDate != nil { skor += 5 }
	if d.Address != "" { skor += 5 }
	if d.Email != "" { skor += 5 }

	// Kepegawaian (40 poin)
	if d.NIP != "" || d.NUPTK != "" { skor += 10 }
	if d.StatusPegawai != "" { skor += 10 }
	if d.TMTMasuk != nil { skor += 10 }
	if d.UnitKerja != "" { skor += 10 }

	// Dokumen Pendukung (20 poin)
	if d.FotoURL != "" { skor += 20 }

	d.SkorKelengkapan = skor

	// Label
	if skor >= 81 {
		d.StatusKelengkapan = "Lengkap"
	} else if skor >= 51 {
		d.StatusKelengkapan = "Cukup"
	} else {
		d.StatusKelengkapan = "Belum Lengkap"
	}
}
