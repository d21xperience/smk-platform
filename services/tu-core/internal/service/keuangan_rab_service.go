package service

import (
	"context"
	"errors"
	"fmt"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
	"sekolah-platform/services/tu-core/internal/repository"
)

// =====================================================================
// RABService — RAB header/detail CRUD + realisasi
// =====================================================================
type RABService struct {
	repo       repository.RABRepository
	masterRepo repository.KeuanganMasterRepository
}

func NewRABService(repo repository.RABRepository, masterRepo repository.KeuanganMasterRepository) *RABService {
	return &RABService{repo: repo, masterRepo: masterRepo}
}

// ---------------------------------------------------------------------
// CreateRab — validasi + enrich items (lookup tipe dari COA) + repo
// ---------------------------------------------------------------------
func (s *RABService) CreateRab(
	ctx context.Context, req models.CreateRabRequest,
) (*models.RabHeader, error) {
	if err := validateCreateRab(&req); err != nil {
		return nil, err
	}

	// Cek duplikasi: 1 RAB per school×period×tahun
	_, err := s.repo.GetRabByTahun(ctx, req.SchoolID, req.AcademicPeriodID, req.TahunAnggaran)
	if err == nil {
		return nil, fmt.Errorf("RAB tahun %d untuk periode %s sudah ada (gunakan update)", req.TahunAnggaran, req.AcademicPeriodID)
	}
	if !errors.Is(err, repository.ErrTransaksiKeuanganNotFound) {
		return nil, err
	}

	// Enrich items — lookup kode_rekening (untuk tipe & nama)
	details, err := s.enrichRabDetails(ctx, req.SchoolID, req.Items)
	if err != nil {
		return nil, err
	}

	h := &models.RabHeader{
		SchoolID:         req.SchoolID,
		AcademicPeriodID: req.AcademicPeriodID,
		TahunAnggaran:    req.TahunAnggaran,
		Status:           models.RabDraft,
		DisusunOleh:      req.DisusunOleh,
		Catatan:          req.Catatan,
	}

	if err := s.repo.CreateRab(ctx, h, details); err != nil {
		return nil, err
	}

	return s.repo.GetRabByID(ctx, req.SchoolID, h.ID)
}

// ---------------------------------------------------------------------
// Read methods
// ---------------------------------------------------------------------
func (s *RABService) GetRab(ctx context.Context, schoolID, id string) (*models.RabHeader, error) {
	return s.repo.GetRabByID(ctx, schoolID, id)
}

func (s *RABService) GetRabByTahun(
	ctx context.Context, schoolID, periodID string, tahun int,
) (*models.RabHeader, error) {
	return s.repo.GetRabByTahun(ctx, schoolID, periodID, tahun)
}

func (s *RABService) ListRab(
	ctx context.Context, filter models.RabFilter,
) ([]models.RabHeader, int, error) {
	return s.repo.ListRab(ctx, filter)
}

// ---------------------------------------------------------------------
// UpdateRab — replace header + details (status harus DRAFT/DIREVISI)
// ---------------------------------------------------------------------
func (s *RABService) UpdateRab(
	ctx context.Context, id string, req models.CreateRabRequest,
) (*models.RabHeader, error) {
	if id == "" {
		return nil, errors.New("id wajib diisi")
	}
	if err := validateCreateRab(&req); err != nil {
		return nil, err
	}

	existing, err := s.repo.GetRabByID(ctx, req.SchoolID, id)
	if err != nil {
		return nil, err
	}

	details, err := s.enrichRabDetails(ctx, req.SchoolID, req.Items)
	if err != nil {
		return nil, err
	}

	existing.Catatan = req.Catatan
	if err := s.repo.UpdateRab(ctx, existing, details); err != nil {
		return nil, err
	}
	return s.repo.GetRabByID(ctx, req.SchoolID, id)
}

// ---------------------------------------------------------------------
// SahkanRab — transisi ke DISAHKAN (hanya dari DRAFT/DIAJUKAN)
// ---------------------------------------------------------------------
func (s *RABService) SahkanRab(
	ctx context.Context, req models.SahkanRabRequest,
) (*models.RabHeader, error) {
	if req.SchoolID == "" || req.RabID == "" {
		return nil, errors.New("school_id dan rab_id wajib diisi")
	}
	if err := s.repo.SahkanRab(ctx, req.SchoolID, req.RabID, req.DisahkanOleh, req.Catatan); err != nil {
		return nil, err
	}
	return s.repo.GetRabByID(ctx, req.SchoolID, req.RabID)
}

// ---------------------------------------------------------------------
// DeleteRab — hapus (hanya DRAFT/DIREVISI)
// ---------------------------------------------------------------------
func (s *RABService) DeleteRab(ctx context.Context, schoolID, id string) error {
	return s.repo.DeleteRab(ctx, schoolID, id)
}

// ---------------------------------------------------------------------
// GetRealisasi — compute on-the-fly dari BKU
// ---------------------------------------------------------------------
func (s *RABService) GetRealisasi(
	ctx context.Context, schoolID, periodID string, tahun, sampaiBulan int,
) (*models.RabRealisasiReport, error) {
	return s.repo.GetRealisasiByTahun(ctx, schoolID, periodID, tahun, sampaiBulan)
}

// ---------------------------------------------------------------------
// GetSisaAnggaran — exposed untuk budget check di PengeluaranService
// ---------------------------------------------------------------------
func (s *RABService) GetSisaAnggaran(
	ctx context.Context, rabDetailID string,
) (pagu, realisasi, sisa float64, err error) {
	return s.repo.GetSisaAnggaran(ctx, rabDetailID)
}

// =====================================================================
// Private helpers
// =====================================================================

// enrichRabDetails — lookup kode_rekening per item, isi Tipe & KodeRekeningNama
// dari master COA, lalu hitung subtotal = volume × harga_satuan
func (s *RABService) enrichRabDetails(
	ctx context.Context, schoolID string, items []models.RabDetailInput,
) ([]models.RabDetail, error) {
	details := make([]models.RabDetail, 0, len(items))
	for i, in := range items {
		k, err := s.masterRepo.GetKodeRekeningByID(ctx, schoolID, in.KodeRekeningID)
		if err != nil {
			return nil, fmt.Errorf("item #%d: kode_rekening %s tidak ditemukan: %w", i+1, in.KodeRekeningID, err)
		}
		details = append(details, models.RabDetail{
			KodeRekeningID:   in.KodeRekeningID,
			KodeRekeningKode: k.Kode,
			KodeRekeningNama: k.Nama,
			Tipe:             k.Tipe,
			Uraian:           in.Uraian,
			Volume:           in.Volume,
			Satuan:           in.Satuan,
			HargaSatuan:      in.HargaSatuan,
			Subtotal:         in.HitungSubtotal(),
			SumberDana:       in.SumberDana,
		})
	}
	return details, nil
}

func validateCreateRab(r *models.CreateRabRequest) error {
	if r.SchoolID == "" {
		return errors.New("school_id wajib diisi")
	}
	if r.AcademicPeriodID == "" {
		return errors.New("academic_period_id wajib diisi")
	}
	if r.TahunAnggaran < 2000 || r.TahunAnggaran > 2100 {
		return errors.New("tahun_anggaran tidak valid (2000-2100)")
	}
	if len(r.Items) == 0 {
		return errors.New("items RAB tidak boleh kosong")
	}
	for i, item := range r.Items {
		if err := models.ValidateRabDetailInput(item); err != nil {
			return fmt.Errorf("item #%d: %w", i+1, err)
		}
	}
	return nil
}
