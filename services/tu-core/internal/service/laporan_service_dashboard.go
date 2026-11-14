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
// GetDashboard — combine 3 query + build alerts
// =====================================================================
func (s *LaporanService) GetDashboard(
	ctx context.Context, filter models.LaporanFilter,
) (*models.DashboardSummary, error) {
	if filter.SchoolID == "" || filter.AcademicPeriodID == "" {
		return nil, errors.New("school_id dan academic_period_id wajib diisi")
	}

	now := time.Now()
	bulan := int(now.Month())
	tahun := now.Year()

	// 1. Kas
	kasData, err := s.repo.GetSaldoKasData(ctx, filter.SchoolID, filter.AcademicPeriodID, bulan, tahun)
	if err != nil {
		return nil, err
	}

	// 2. Piutang
	piutangData, err := s.repo.GetPiutangForDashboard(ctx, filter.SchoolID, filter.AcademicPeriodID)
	if err != nil {
		return nil, err
	}

	// 3. Piutang aging summary (reuse piutangService)
	agingReport, err := s.piutangService.GetAgingReport(ctx, filter.SchoolID, filter.AcademicPeriodID)
	if err != nil {
		return nil, err
	}

	// 4. Pengajuan
	pengajuanCounts, totalPending, err := s.repo.GetPengajuanForDashboard(ctx, filter.SchoolID, filter.AcademicPeriodID)
	if err != nil {
		return nil, err
	}

	// 5. Rekon bulan ini
	rekonTotal, rekonBalance, rekonUnbalance, err := s.repo.GetRekonForDashboard(ctx, filter.SchoolID, filter.AcademicPeriodID, bulan, tahun)
	if err != nil {
		return nil, err
	}

	// 6. RAB tahun ini (opsional — kalau ada)
	var rabInfo *models.DashboardRabInfo
	rabHeader, err := s.rabService.GetRabByTahun(ctx, filter.SchoolID, filter.AcademicPeriodID, tahun)
	if err == nil && rabHeader != nil {
		persenReal := 0.0
		if rabHeader.TotalBelanja > 0 {
			// Ambil realisasi s.d. bulan ini
			realisasi, rErr := s.rabService.GetRealisasi(ctx, filter.SchoolID, filter.AcademicPeriodID, tahun, bulan)
			if rErr == nil && realisasi != nil {
				persenReal = 0.0
				if realisasi.TotalPaguBelanja > 0 {
					persenReal = (realisasi.TotalRealisasiBelanja / realisasi.TotalPaguBelanja) * 100
				}
				rabInfo = &models.DashboardRabInfo{
					TahunAnggaran:    tahun,
					Status:           rabHeader.Status,
					TotalPaguBelanja: realisasi.TotalPaguBelanja,
					RealisasiBelanja: realisasi.TotalRealisasiBelanja,
					PersenRealisasi:  persenReal,
				}
			}
		}
		if rabInfo == nil {
			rabInfo = &models.DashboardRabInfo{
				TahunAnggaran:    tahun,
				Status:           rabHeader.Status,
				TotalPaguBelanja: rabHeader.TotalBelanja,
			}
		}
	}

	// 7. Build alerts
	alerts := buildDashboardAlerts(
		agingReport, pengajuanCounts, rekonUnbalance, rabInfo, kasData,
	)

	return &models.DashboardSummary{
		AsOfDate: now.Format("2006-01-02"),
		Kas: models.DashboardKasInfo{
			SaldoTerkini:        kasData.SaldoTerkini,
			TotalDebetBulanIni:  kasData.TotalDebetBulanIni,
			TotalKreditBulanIni: kasData.TotalKreditBulanIni,
		},
		Piutang: models.DashboardPiutangInfo{
			TotalSisa:    piutangData.TotalSisa,
			JumlahItem:   piutangData.JumlahItem,
			JumlahSiswa:  piutangData.JumlahSiswa,
			AgingSummary: agingReport.Buckets,
		},
		Rab: rabInfo,
		Pengajuan: models.DashboardPengajuanInfo{
			MenungguVerifikasi:  pengajuanCounts["MENUNGGU_VERIFIKASI_BENDAHARA"],
			MenungguPersetujuan: pengajuanCounts["MENUNGGU_PERSETUJUAN_KEPALA"],
			Disetujui:           pengajuanCounts["DISETUJUI"],
			TotalNominalPending: totalPending,
		},
		Rekonsiliasi: models.DashboardRekonsiliasiInfo{
			TotalRekonBulanIni: rekonTotal,
			Balance:            rekonBalance,
			Unbalance:          rekonUnbalance,
		},
		Alerts: alerts,
	}, nil
}

// =====================================================================
// GetRingkasanBulan — ringkasan 1 bulan
// =====================================================================
func (s *LaporanService) GetRingkasanBulan(
	ctx context.Context, filter models.LaporanFilter,
) (*models.RingkasanBulan, error) {
	if filter.SchoolID == "" || filter.AcademicPeriodID == "" {
		return nil, errors.New("school_id dan academic_period_id wajib diisi")
	}
	if filter.Bulan < 1 || filter.Bulan > 12 {
		return nil, errors.New("bulan harus 1-12")
	}
	if filter.Tahun < 2000 || filter.Tahun > 2100 {
		return nil, errors.New("tahun tidak valid")
	}

	data, err := s.repo.GetRingkasanBulanData(ctx, filter.SchoolID, filter.AcademicPeriodID, filter.Bulan, filter.Tahun)
	if err != nil {
		return nil, err
	}

	out := &models.RingkasanBulan{
		PeriodeBulan:    filter.Bulan,
		PeriodeTahun:    filter.Tahun,
		SaldoAwal:       data.SaldoAwal,
		TotalDebet:      data.TotalDebet,
		TotalKredit:     data.TotalKredit,
		SaldoAkhir:      data.SaldoAkhir,
		JumlahTransaksi: data.JumlahTransaksi,
		TopDebet:        []models.RingkasanKategoriItem{},
		TopKredit:       []models.RingkasanKategoriItem{},
	}

	for _, it := range data.TopDebet {
		out.TopDebet = append(out.TopDebet, models.RingkasanKategoriItem{
			KodeRekening:    it.KodeRekening,
			Nama:            it.Nama,
			Total:           it.Total,
			JumlahTransaksi: it.JumlahTransaksi,
		})
	}
	for _, it := range data.TopKredit {
		out.TopKredit = append(out.TopKredit, models.RingkasanKategoriItem{
			KodeRekening:    it.KodeRekening,
			Nama:            it.Nama,
			Total:           it.Total,
			JumlahTransaksi: it.JumlahTransaksi,
		})
	}

	return out, nil
}

// =====================================================================
// buildDashboardAlerts — generate alerts dari data
// =====================================================================
func buildDashboardAlerts(
	agingReport *models.AgingReport,
	pengajuanCounts map[string]int,
	rekonUnbalance int,
	rabInfo *models.DashboardRabInfo,
	kasData *repository.SaldoKasData,
) []models.DashboardAlert {
	alerts := []models.DashboardAlert{}

	// Alert piutang aging >90
	if agingReport != nil {
		for _, b := range agingReport.Buckets {
			if b.Group == ">90" && b.ItemCount > 0 {
				alerts = append(alerts, models.DashboardAlert{
					Severity: models.AlertSeverityDanger,
					Category: models.AlertCategoryPiutang,
					Message:  fmt.Sprintf("Piutang >90 hari: %d tagihan, total %.0f", b.ItemCount, b.TotalSisa),
				})
			}
		}
	}

	// Alert pengajuan pending
	if cnt := pengajuanCounts["MENUNGGU_VERIFIKASI_BENDAHARA"]; cnt > 0 {
		alerts = append(alerts, models.DashboardAlert{
			Severity: models.AlertSeverityWarning,
			Category: models.AlertCategoryPengajuan,
			Message:  fmt.Sprintf("%d pengajuan menunggu verifikasi bendahara", cnt),
		})
	}
	if cnt := pengajuanCounts["MENUNGGU_PERSETUJUAN_KEPALA"]; cnt > 0 {
		alerts = append(alerts, models.DashboardAlert{
			Severity: models.AlertSeverityWarning,
			Category: models.AlertCategoryPengajuan,
			Message:  fmt.Sprintf("%d pengajuan menunggu persetujuan kepala sekolah", cnt),
		})
	}

	// Alert rekonsiliasi unbalance
	if rekonUnbalance > 0 {
		alerts = append(alerts, models.DashboardAlert{
			Severity: models.AlertSeverityDanger,
			Category: models.AlertCategoryRekonsiliasi,
			Message:  fmt.Sprintf("%d rekonsiliasi bulan ini UNBALANCE", rekonUnbalance),
		})
	}

	// Alert RAB over-budget
	if rabInfo != nil && rabInfo.PersenRealisasi > 90 {
		alerts = append(alerts, models.DashboardAlert{
			Severity: models.AlertSeverityWarning,
			Category: models.AlertCategoryRab,
			Message:  fmt.Sprintf("Realisasi RAB %.1f%% dari pagu, mendekati batas", rabInfo.PersenRealisasi),
		})
	}

	// Alert kas negatif
	if kasData != nil && kasData.SaldoTerkini < 0 {
		alerts = append(alerts, models.DashboardAlert{
			Severity: models.AlertSeverityWarning,
			Category: models.AlertCategoryKas,
			Message:  fmt.Sprintf("Saldo kas negatif: %.0f", kasData.SaldoTerkini),
		})
	}

	return alerts
}
