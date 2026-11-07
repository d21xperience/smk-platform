package service

import (
	"context"
	"errors"

	"sekolah-platform/services/tu-core/internal/repository"
)

type LogisticService struct {
	repo repository.LogisticRepository
}

func NewLogisticService(repo repository.LogisticRepository) *LogisticService {
	return &LogisticService{repo: repo}
}

func (s *LogisticService) CreateLogistic(ctx context.Context, req *repository.Logistic) error {
	if req.Kode == "" || req.Nama == "" {
		return errors.New("kode dan nama barang wajib diisi")
	}
	return s.repo.Create(ctx, req)
}

func (s *LogisticService) ListLogistics(ctx context.Context) ([]repository.Logistic, error) {
	return s.repo.FindAll(ctx)
}

func (s *LogisticService) RestockLogistic(ctx context.Context, id string, jumlah int, keterangan string) error {
	logistic, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	newStock := logistic.Stok + jumlah
	if err := s.repo.UpdateStock(ctx, id, newStock); err != nil {
		return err
	}

	mutation := &repository.LogisticMutation{
		LogisticID: id,
		Tipe:       "IN",
		Jumlah:     jumlah,
		Keterangan: keterangan,
	}
	return s.repo.RecordMutation(ctx, mutation)
}

func (s *LogisticService) UseLogistic(ctx context.Context, id string, jumlah int, keterangan string) error {
	logistic, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if logistic.Stok < jumlah {
		return errors.New("stok tidak mencukupi untuk digunakan")
	}

	newStock := logistic.Stok - jumlah
	if err := s.repo.UpdateStock(ctx, id, newStock); err != nil {
		return err
	}

	mutation := &repository.LogisticMutation{
		LogisticID: id,
		Tipe:       "OUT",
		Jumlah:     jumlah,
		Keterangan: keterangan,
	}
	return s.repo.RecordMutation(ctx, mutation)
}
