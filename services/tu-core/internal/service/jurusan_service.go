package service

import (
	"context"

	"sekolah-platform/services/tu-core/internal/repository"
)

type JurusanService struct {
	repo repository.JurusanRepository
}

func NewJurusanService(repo repository.JurusanRepository) *JurusanService {
	return &JurusanService{repo: repo}
}

func (s *JurusanService) List(ctx context.Context) ([]repository.Jurusan, error) {
	return s.repo.List(ctx)
}

func (s *JurusanService) Create(ctx context.Context, id, nama, level string) error {
	j := &repository.Jurusan{
		JurusanID:   id,
		NamaJurusan: nama,
		LevelBidang: level,
	}
	return s.repo.Create(ctx, j)
}
