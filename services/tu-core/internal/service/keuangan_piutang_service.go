package service

import (
	"context"
	"errors"
	"time"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
	"sekolah-platform/services/tu-core/internal/repository"
)

// =====================================================================
// PiutangService — aging report & detail piutang
// Read-only, delegasi ke repo + compose output
// =====================================================================
type PiutangService struct {
	repo repository.PiutangRepository
}

func NewPiutangService(repo repository.PiutangRepository) *PiutangService {
	return &PiutangService{repo: repo}
}

// =====================================================================
// ListPiutang — list flat (setiap tagihan 1 row)
// =====================================================================
func (s *PiutangService) ListPiutang(
	ctx context.Context, filter models.PiutangFilter,
) ([]models.PiutangItem, int, error) {
	if filter.SchoolID == "" {
		return nil, 0, errors.New("school_id wajib diisi")
	}
	if filter.AcademicPeriodID == "" {
		return nil, 0, errors.New("academic_period_id wajib diisi")
	}
	if filter.AgingGroup != "" && !models.IsValidAgingGroup(filter.AgingGroup) {
		return nil, 0, errors.New("aging_group tidak valid")
	}
	return s.repo.ListPiutang(ctx, filter)
}

// =====================================================================
// GetAgingReport — report agregat per bucket
// =====================================================================
func (s *PiutangService) GetAgingReport(
	ctx context.Context, schoolID, periodID string,
) (*models.AgingReport, error) {
	if schoolID == "" || periodID == "" {
		return nil, errors.New("school_id dan academic_period_id wajib diisi")
	}
	buckets, totalSisa, totalCount, err := s.repo.GetAgingSummary(ctx, schoolID, periodID)
	if err != nil {
		return nil, err
	}
	return &models.AgingReport{
		AsOfDate:       time.Now().Format("2006-01-02"),
		TotalItemCount: totalCount,
		TotalSisa:      totalSisa,
		Buckets:        buckets,
	}, nil
}

// =====================================================================
// GetPiutangByNIS — detail piutang 1 siswa (grouped)
// =====================================================================
func (s *PiutangService) GetPiutangByNIS(
	ctx context.Context, schoolID, periodID, nis string,
) (*models.PiutangPerSiswa, error) {
	if schoolID == "" || periodID == "" || nis == "" {
		return nil, errors.New("school_id, academic_period_id, nis wajib diisi")
	}
	items, err := s.repo.GetPiutangByNIS(ctx, schoolID, periodID, nis)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, repository.ErrTransaksiKeuanganNotFound
	}
	grouped := groupPiutangByNIS(items)
	if len(grouped) == 0 {
		return nil, repository.ErrTransaksiKeuanganNotFound
	}
	return &grouped[0], nil
}

// =====================================================================
// ListPiutangPerSiswa — list grouped per NIS (semua siswa dengan piutang)
// =====================================================================
func (s *PiutangService) ListPiutangPerSiswa(
	ctx context.Context, filter models.PiutangFilter,
) ([]models.PiutangPerSiswa, int, error) {
	if filter.SchoolID == "" {
		return nil, 0, errors.New("school_id wajib diisi")
	}
	if filter.AcademicPeriodID == "" {
		return nil, 0, errors.New("academic_period_id wajib diisi")
	}

	// Pakai GetAllPiutangForPeriod supaya grouping konsisten
	// (filter NIS/Kelas tetap bisa di-handle di grouping post-process)
	items, err := s.repo.GetAllPiutangForPeriod(ctx, filter.SchoolID, filter.AcademicPeriodID)
	if err != nil {
		return nil, 0, err
	}

	// Filter in-memory (kalau ada filter NIS/Kelas dari user)
	if filter.NIS != "" || filter.KelasNama != "" {
		filtered := make([]models.PiutangItem, 0, len(items))
		for _, it := range items {
			if filter.NIS != "" && it.NIS != filter.NIS {
				continue
			}
			if filter.KelasNama != "" && it.KelasNama != filter.KelasNama {
				continue
			}
			filtered = append(filtered, it)
		}
		items = filtered
	}

	grouped := groupPiutangByNIS(items)
	total := len(grouped)

	// Pagination in-memory (karena grouping tidak bisa langsung SQL LIMIT)
	limit, offset := filter.Limit, filter.Offset
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	if offset >= len(grouped) {
		return []models.PiutangPerSiswa{}, total, nil
	}
	end := offset + limit
	if end > len(grouped) {
		end = len(grouped)
	}

	return grouped[offset:end], total, nil
}

// =====================================================================
// groupPiutangByNIS — helper untuk grouping items → per siswa
//
// Output di-sort by total_sisa DESC (siswa dengan tunggakan terbesar di atas)
// =====================================================================
func groupPiutangByNIS(items []models.PiutangItem) []models.PiutangPerSiswa {
	type acc struct {
		per    models.PiutangPerSiswa
		worst  int // priority untuk OldestAgingGroup
	}
	byNIS := make(map[string]*acc)

	for _, it := range items {
		a, ok := byNIS[it.NIS]
		if !ok {
			a = &acc{
				per: models.PiutangPerSiswa{
					NIS:       it.NIS,
					NamaSiswa: it.NamaSiswa,
					KelasNama: it.KelasNama,
					Items:     []models.PiutangItem{},
				},
			}
			byNIS[it.NIS] = a
		}
		a.per.TotalSisa += it.Sisa
		a.per.JumlahTagihan++
		a.per.Items = append(a.per.Items, it)

		p := models.AgingGroupPriority(it.AgingGroup)
		if p > a.worst {
			a.worst = p
			a.per.OldestAgingGroup = it.AgingGroup
		}
	}

	out := make([]models.PiutangPerSiswa, 0, len(byNIS))
	for _, a := range byNIS {
		out = append(out, a.per)
	}

	// Sort by total_sisa DESC (tunggakan terbesar dulu) — simple insertion sort
	// (jumlah siswa < 200 biasanya, jadi O(n²) tidak masalah; kalau besar,
	// ganti dengan sort.Slice di masa depan)
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j].TotalSisa > out[j-1].TotalSisa; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}

	return out
}
