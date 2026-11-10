package service

import (
	"context"
	"time"

	"sekolah-platform/services/tu-core/internal/repository"
)

type KehadiranService struct {
	repo repository.KehadiranRepository
}

func NewKehadiranService(repo repository.KehadiranRepository) *KehadiranService {
	return &KehadiranService{repo: repo}
}

func (s *KehadiranService) GetAttendanceByDate(ctx context.Context, tanggal time.Time) ([]repository.KehadiranHarian, error) {
	return s.repo.GetAttendanceByDate(ctx, tanggal)
}

func (s *KehadiranService) ManualInput(ctx context.Context, ptkID string, tanggal time.Time, status string, jamMasuk, jamPulang *string) error {
	var durasi, menitTerlambat *int
	if jamMasuk != nil && jamPulang != nil {
		d := 0
		m := 0
		durasi = &d
		menitTerlambat = &m
	}
	return s.repo.UpsertAttendance(ctx, ptkID, tanggal, status, jamMasuk, jamPulang, durasi, menitTerlambat, "MANUAL")
}

func (s *KehadiranService) AutoCollect(ctx context.Context, tanggal time.Time) (int, error) {
	return s.repo.AutoCollectDaily(ctx, tanggal)
}

func (s *KehadiranService) GetRekapBulanan(ctx context.Context) ([]repository.RekapBulananHeader, error) {
	return s.repo.GetRekapBulanan(ctx)
}

func (s *KehadiranService) GenerateRekapBulanan(ctx context.Context, bulan, tahun int) (string, error) {
	return s.repo.GenerateRekapBulanan(ctx, bulan, tahun)
}
