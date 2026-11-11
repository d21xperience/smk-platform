package service

import (
	"context"
	"fmt"
	"time"

	"sekolah-platform/services/tu-core/internal/repository"
)

type UsulanKenaikanService struct {
	repo repository.UsulanKenaikanRepository
}

func NewUsulanKenaikanService(repo repository.UsulanKenaikanRepository) *UsulanKenaikanService {
	return &UsulanKenaikanService{repo: repo}
}

func (s *UsulanKenaikanService) ListUsulan(ctx context.Context) ([]repository.UsulanKenaikan, error) {
	return s.repo.ListUsulan(ctx)
}

func (s *UsulanKenaikanService) GetSaranUsulan(ctx context.Context) ([]repository.SaranUsulan, error) {
	return s.repo.GetSaranUsulan(ctx)
}

func (s *UsulanKenaikanService) CreateUsulan(ctx context.Context, ptkID, jenis, catatan, userID string, golSekarangID, golBaruID int, tmtBaru time.Time) (string, error) {
	// Generate Nomor Usulan Sederhana
	prefix := "USUL"
	switch jenis {
	case "KGB":
		prefix = "KGB"
	case "SATYALANCANA":
		prefix = "SAT"
	}

	tahun := time.Now().Year()
	bulan := time.Now().Month()
	nomorUsulan := fmt.Sprintf("%s/%d/%02d/%d", prefix, tahun, bulan, time.Now().Unix()%1000)

	id, err := s.repo.CreateUsulan(ctx, ptkID, jenis, nomorUsulan, userID, golSekarangID, golBaruID, tmtBaru, catatan)
	if err != nil {
		return "", fmt.Errorf("gagal membuat usulan: %v", err)
	}
	return id, nil
}

func (s *UsulanKenaikanService) UpdateStatus(ctx context.Context, id, status, catatan string) error {
	return s.repo.UpdateStatus(ctx, id, status, catatan)
}
