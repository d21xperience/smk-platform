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

func (s *JurusanService) ListJurusan(ctx context.Context) ([]repository.Jurusan, error) {
	return s.repo.List(ctx)
}
