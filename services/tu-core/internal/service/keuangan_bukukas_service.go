package service

import (
	"context"
	"time"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
	"sekolah-platform/services/tu-core/internal/repository"
)

// =====================================================================
// BukuKasService — read-only view BKU + saldo
// Write BKU dilakukan oleh IuranService & PengeluaranService saat
// orkestrasi transaksi (bukan di service ini).
// =====================================================================
type BukuKasService struct {
	repo repository.KeuanganTransaksiRepository
}

func NewBukuKasService(repo repository.KeuanganTransaksiRepository) *BukuKasService {
	return &BukuKasService{repo: repo}
}

func (s *BukuKasService) ListBKU(ctx context.Context, filter models.BukuKasFilter) ([]models.BukuKasUmum, int, error) {
	return s.repo.ListBKU(ctx, filter)
}

// GetSaldo — saldo berjalan per tanggal tertentu (default: hari ini)
func (s *BukuKasService) GetSaldo(ctx context.Context, schoolID, periodID, tanggal string) (float64, string, error) {
	tgl := time.Now()
	if tanggal != "" {
		parsed, err := time.Parse("2006-01-02", tanggal)
		if err != nil {
			return 0, "", err
		}
		tgl = parsed
	}
	saldo, err := s.repo.GetSaldoTerakhir(ctx, schoolID, periodID, tgl)
	if err != nil {
		return 0, "", err
	}
	return saldo, tgl.Format("2006-01-02"), nil
}
