package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
	"sekolah-platform/services/tu-core/internal/repository"
)

// =====================================================================
// LaporanService — orkestrasi laporan keuangan
// =====================================================================
type LaporanService struct {
	repo          repository.LaporanRepository
	rabService    *RABService
	piutangService *PiutangService
	penutupanRepo repository.PenutupanRepository
}

func NewLaporanService(
	repo repository.LaporanRepository,
	rabService *RABService,
	piutangService *PiutangService,
	penutupanRepo repository.PenutupanRepository,
) *LaporanService {
	return &LaporanService{
		repo:           repo,
		rabService:     rabService,
		piutangService: piutangService,
		penutupanRepo:  penutupanRepo,
	}
}

// =====================================================================
// GetNeraca — mode live atau snapshot
//
//   - Live: compute dari BKU + piutang
//   - Snapshot: baca dari penutupan_buku_kas yang sudah di-lock
// =====================================================================
func (s *LaporanService) GetNeraca(
	ctx context.Context, filter models.LaporanFilter,
) (*models.NeracaReport, error) {
	if filter.SchoolID == "" || filter.AcademicPeriodID == "" {
		return nil, errors.New("school_id dan academic_period_id wajib diisi")
	}

	mode := filter.Mode
	if mode == "" {
		mode = models.LaporanModeLive
	}
	if !models.IsValidLaporanMode(mode) {
		return nil, fmt.Errorf("mode tidak valid: %s (LIVE/SNAPSHOT)", mode)
	}

	if mode == models.LaporanModeSnapshot {
		return s.getNeracaSnapshot(ctx, filter)
	}
	return s.getNeracaLive(ctx, filter)
}

// =====================================================================
// getNeracaLive — compute dari BKU + piutang
// =====================================================================
func (s *LaporanService) getNeracaLive(
	ctx context.Context, filter models.LaporanFilter,
) (*models.NeracaReport, error) {
	asOf := time.Now()
	if filter.Tanggal != "" {
		t, err := time.Parse("2006-01-02", filter.Tanggal)
		if err != nil {
			return nil, fmt.Errorf("tanggal harus YYYY-MM-DD: %w", err)
		}
		asOf = t
	}

	data, err := s.repo.GetNeracaData(ctx, filter.SchoolID, filter.AcademicPeriodID, asOf)
	if err != nil {
		return nil, err
	}

	return buildNeracaReport(data, models.LaporanModeLive, asOf.Format("2006-01-02")), nil
}

// =====================================================================
// getNeracaSnapshot — baca dari penutupan yang sudah di-lock
// =====================================================================
func (s *LaporanService) getNeracaSnapshot(
	ctx context.Context, filter models.LaporanFilter,
) (*models.NeracaReport, error) {
	if filter.PenutupanID == "" {
		return nil, errors.New("penutupan_id wajib diisi untuk mode SNAPSHOT")
	}
	p, err := s.penutupanRepo.GetByID(ctx, filter.SchoolID, filter.PenutupanID)
	if err != nil {
		return nil, err
	}

	// Snapshot hanya punya total, tidak ada detail per COA.
	// Bangun NeracaReport minimal dengan 3 baris (Kas, Kewajiban, Ekuitas).
	kas := models.NeracaItem{
		KodeRekeningID: "SNAPSHOT_KAS",
		Kode:           "1.1.1",
		Nama:           "Kas (snapshot)",
		Tipe:           models.TipeAset,
		Saldo:          p.SaldoAkhirBuku,
	}
	kewajiban := models.NeracaItem{
		KodeRekeningID: "SNAPSHOT_KEWAJIBAN",
		Kode:           "2.1",
		Nama:           "Kewajiban (snapshot)",
		Tipe:           models.TipeKewajiban,
		Saldo:          p.TotalKewajiban,
	}
	ekuitas := models.NeracaItem{
		KodeRekeningID: "SNAPSHOT_EKUITAS",
		Kode:           "3.1",
		Nama:           "Ekuitas (snapshot)",
		Tipe:           models.TipeEkuitas,
		Saldo:          p.TotalEkuitas,
	}

	report := &models.NeracaReport{
		AsOfDate:       p.ClosedAt.Format("2006-01-02"),
		Mode:           models.LaporanModeSnapshot,
		PenutupanID:    p.ID,
		PeriodeBulan:   p.PeriodeBulan,
		PeriodeTahun:   p.PeriodeTahun,
		Aset:           []models.NeracaItem{kas},
		Kewajiban:      []models.NeracaItem{kewajiban},
		Ekuitas:        []models.NeracaItem{ekuitas},
		Pendapatan:     []models.NeracaItem{},
		Belanja:        []models.NeracaItem{},
		TotalAset:      p.TotalAset,
		TotalKewajiban: p.TotalKewajiban,
		TotalEkuitas:   p.TotalEkuitas,
	}
	report.Selisih = report.TotalAset - (report.TotalKewajiban + report.TotalEkuitas)
	report.IsBalance = absFloat(report.Selisih) < 0.01

	return report, nil
}

// =====================================================================
// GetPosisiKeuangan — alias GetNeraca dengan validasi ketat
// (selalu harus balance, kalau tidak → error)
// =====================================================================
func (s *LaporanService) GetPosisiKeuangan(
	ctx context.Context, filter models.LaporanFilter,
) (*models.NeracaReport, error) {
	report, err := s.GetNeraca(ctx, filter)
	if err != nil {
		return nil, err
	}
	if !report.IsBalance {
		return nil, fmt.Errorf(
			"posisi keuangan tidak balance: Aset %.2f ≠ Kewajiban %.2f + Ekuitas %.2f (selisih %.2f)",
			report.TotalAset, report.TotalKewajiban, report.TotalEkuitas, report.Selisih,
		)
	}
	return report, nil
}

// =====================================================================
// GetLra — wrapper dari RABService.GetRealisasi + info RAB status
// =====================================================================
func (s *LaporanService) GetLra(
	ctx context.Context, filter models.LaporanFilter,
) (*models.LraReport, error) {
	if filter.SchoolID == "" || filter.AcademicPeriodID == "" {
		return nil, errors.New("school_id dan academic_period_id wajib diisi")
	}
	if filter.Tahun < 2000 || filter.Tahun > 2100 {
		return nil, errors.New("tahun tidak valid")
	}
	sampaiBulan := filter.SampaiBulan
	if sampaiBulan <= 0 || sampaiBulan > 12 {
		sampaiBulan = 12
	}

	// Ambil RAB header untuk info status
	rabHeader, err := s.rabService.GetRabByTahun(ctx, filter.SchoolID, filter.AcademicPeriodID, filter.Tahun)
	if err != nil {
		return nil, fmt.Errorf("RAB tahun %d belum dibuat: %w", filter.Tahun, err)
	}

	// Reuse realisasi
	realisasi, err := s.rabService.GetRealisasi(ctx, filter.SchoolID, filter.AcademicPeriodID, filter.Tahun, sampaiBulan)
	if err != nil {
		return nil, err
	}

	return &models.LraReport{
		TahunAnggaran:     filter.Tahun,
		SampaiBulan:       sampaiBulan,
		RabID:             rabHeader.ID,
		RabStatus:         rabHeader.Status,
		TanggalPengesahan: rabHeader.TanggalPengesahan,
		Realisasi:         realisasi,
	}, nil
}

// =====================================================================
// buildNeracaReport — helper: group item per tipe, hitung totals
// =====================================================================
func buildNeracaReport(
	data *repository.NeracaData, mode, asOfDate string,
) *models.NeracaReport {
	report := &models.NeracaReport{
		AsOfDate:   asOfDate,
		Mode:       mode,
		Aset:       []models.NeracaItem{},
		Kewajiban:  []models.NeracaItem{},
		Ekuitas:    []models.NeracaItem{},
		Pendapatan: []models.NeracaItem{},
		Belanja:    []models.NeracaItem{},
	}

	for _, it := range data.Items {
		ni := models.NeracaItem{
			KodeRekeningID: it.KodeRekeningID,
			Kode:           it.Kode,
			Nama:           it.Nama,
			Tipe:           it.Tipe,
			Saldo:          it.Saldo,
		}
		switch it.Tipe {
		case models.TipeAset:
			report.Aset = append(report.Aset, ni)
			report.TotalAset += it.Saldo
		case models.TipeKewajiban:
			report.Kewajiban = append(report.Kewajiban, ni)
			report.TotalKewajiban += it.Saldo
		case models.TipeEkuitas:
			report.Ekuitas = append(report.Ekuitas, ni)
			report.TotalEkuitas += it.Saldo
		case models.TipePendapatan:
			report.Pendapatan = append(report.Pendapatan, ni)
			report.TotalPendapatan += it.Saldo
		case models.TipeBelanja:
			report.Belanja = append(report.Belanja, ni)
			report.TotalBelanja += it.Saldo
		}
	}

	// Ekuitas ALWAYS = aset - kewajiban (karena ekuitas = aset bersih).
	// Ini untuk mode LIVE. Untuk SNAPSHOT, sudah di-set langsung di getNeracaSnapshot.
	if len(report.Ekuitas) == 0 {
		report.TotalEkuitas = report.TotalAset - report.TotalKewajiban
	}

	report.SurplusDefisit = report.TotalPendapatan - report.TotalBelanja
	report.Selisih = report.TotalAset - (report.TotalKewajiban + report.TotalEkuitas)
	report.IsBalance = absFloat(report.Selisih) < 0.01

	return report
}

func absFloat(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
