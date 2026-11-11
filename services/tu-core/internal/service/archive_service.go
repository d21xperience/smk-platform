package service

import (
	"context"
	"fmt"
	"time"

	archiveModels "sekolah-platform/services/tu-core/internal/domain/archive/models"
	"sekolah-platform/services/tu-core/internal/domain/correspondence/valueobjects"
	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
)

type ArchiveService struct {
	repo     repository.ArchiveRepository
	corrRepo repository.CorrespondenceRepository
}

func NewArchiveService(repo repository.ArchiveRepository, corrRepo repository.CorrespondenceRepository) *ArchiveService {
	return &ArchiveService{repo: repo, corrRepo: corrRepo}
}

func (s *ArchiveService) ArsipkanSurat(ctx context.Context, correspondenceID string, req ArchiveCorrespondenceRequest) (*ArchiveResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	// 1. Validasi Surat
	corr, err := s.corrRepo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, correspondenceID)
	if err != nil {
		return nil, NewServiceError("ArsipkanSurat", "Archive", "", "surat tidak ditemukan", err)
	}

	// 2. Cek apakah sudah diarsipkan
	existing, _ := s.repo.GetArchiveByCorrespondenceID(ctx, opCtx.SchoolID, correspondenceID)
	if existing != nil {
		return nil, NewServiceError("ArsipkanSurat", "Archive", "", "surat sudah diarsipkan sebelumnya", nil)
	}

	// 3. Hitung Tanggal Musnah
	tglArsip := time.Now()
	tglMusnah := tglArsip.AddDate(req.RetensiTahun, 0, 0)
	lokasiKode := fmt.Sprintf("%s/%s/%s", req.LokasiRak, req.LokasiBox, req.LokasiFolder)

	archive := &archiveModels.Archive{
		SchoolID:         opCtx.SchoolID,
		CorrespondenceID: correspondenceID,
		LokasiKode:       lokasiKode,
		Kategori:         req.Kategori,
		TanggalArsip:     tglArsip,
		TanggalMusnah:    tglMusnah,
		RetensiTahun:     req.RetensiTahun,
		Status:           "TERARSIP",
		Keterangan:       req.Keterangan,
		ArchivedBy:       opCtx.UserID,
	}

	if err := s.repo.CreateArchive(ctx, archive); err != nil {
		return nil, NewServiceError("ArsipkanSurat", "Archive", "", "failed to save archive", err)
	}

	// 4. Update Status Surat
	corr.Status, _ = valueobjects.NewCorrespondenceStatus("ARCHIVED")
	_ = s.corrRepo.Update(ctx, corr)

	// 5. Log Movement
	_ = s.repo.CreateMovement(ctx, &archiveModels.ArchiveMovement{
		ArchiveID:  archive.ID,
		Type:       "ARCHIVE",
		PetugasID:  opCtx.UserID,
		Keterangan: "Pengarsipan awal ke " + lokasiKode,
		Status:     "SELESAI",
	})

	return toArchiveResponse(archive), nil
}

func (s *ArchiveService) PinjamArsip(ctx context.Context, archiveID string, req BorrowArchiveRequest) error {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return ErrMissingOperationalContext
	}

	archive, err := s.repo.GetArchiveByID(ctx, opCtx.SchoolID, archiveID)
	if err != nil {
		return NewServiceError("PinjamArsip", "Archive", archiveID, "arsip tidak ditemukan", err)
	}
	if archive.Status == "DIPINJAM" {
		return NewServiceError("PinjamArsip", "Archive", archiveID, "arsip sedang dipinjam", nil)
	}
	if archive.Status == "DIMUSNAHKAN" {
		return NewServiceError("PinjamArsip", "Archive", archiveID, "arsip sudah dimusnahkan", nil)
	}

	tglKembali, err := time.Parse("2006-01-02", req.TanggalKembaliEstimasi)
	if err != nil {
		return NewServiceError("PinjamArsip", "Archive", archiveID, "format tanggal salah", err)
	}

	if err := s.repo.UpdateArchiveStatus(ctx, archiveID, "DIPINJAM", req.PeminjamID); err != nil {
		return NewServiceError("PinjamArsip", "Archive", archiveID, "gagal update status", err)
	}

	_ = s.repo.CreateMovement(ctx, &archiveModels.ArchiveMovement{
		ArchiveID:              archiveID,
		Type:                   "BORROW",
		PetugasID:              opCtx.UserID,
		PeminjamID:             req.PeminjamID,
		TanggalKembaliEstimasi: &tglKembali,
		Tujuan:                 req.Tujuan,
		Status:                 "AKTIF",
	})

	return nil
}

func (s *ArchiveService) KembalikanArsip(ctx context.Context, archiveID string) error {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return ErrMissingOperationalContext
	}

	archive, err := s.repo.GetArchiveByID(ctx, opCtx.SchoolID, archiveID)
	if err != nil || archive.Status != "DIPINJAM" {
		return NewServiceError("KembalikanArsip", "Archive", archiveID, "arsip tidak dalam status dipinjam", err)
	}

	if err := s.repo.UpdateArchiveStatus(ctx, archiveID, "TERARSIP", ""); err != nil {
		return NewServiceError("KembalikanArsip", "Archive", archiveID, "gagal update status", err)
	}

	tglAktual := time.Now()
	_ = s.repo.CreateMovement(ctx, &archiveModels.ArchiveMovement{
		ArchiveID:            archiveID,
		Type:                 "RETURN",
		PetugasID:            opCtx.UserID,
		TanggalKembaliAktual: &tglAktual,
		Status:               "SELESAI",
	})

	return nil
}

func (s *ArchiveService) PindahArsip(ctx context.Context, archiveID string, req MoveArchiveRequest) error {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return ErrMissingOperationalContext
	}

	lokasiBaru := fmt.Sprintf("%s/%s/%s", req.LokasiRak, req.LokasiBox, req.LokasiFolder)
	if err := s.repo.UpdateArchiveLocation(ctx, archiveID, lokasiBaru); err != nil {
		return NewServiceError("PindahArsip", "Archive", archiveID, "gagal pindah lokasi", err)
	}

	_ = s.repo.CreateMovement(ctx, &archiveModels.ArchiveMovement{
		ArchiveID:  archiveID,
		Type:       "MOVE",
		PetugasID:  opCtx.UserID,
		Keterangan: fmt.Sprintf("Pindah ke %s: %s", lokasiBaru, req.Keterangan),
		Status:     "SELESAI",
	})

	return nil
}

func (s *ArchiveService) MusnahkanArsip(ctx context.Context, archiveID string) error {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return ErrMissingOperationalContext
	}

	if err := s.repo.UpdateArchiveStatus(ctx, archiveID, "DIMUSNAHKAN", ""); err != nil {
		return NewServiceError("MusnahkanArsip", "Archive", archiveID, "gagal memusnahkan", err)
	}

	_ = s.repo.CreateMovement(ctx, &archiveModels.ArchiveMovement{
		ArchiveID:  archiveID,
		Type:       "DESTROY",
		PetugasID:  opCtx.UserID,
		Keterangan: "Pemusnahan arsip",
		Status:     "SELESAI",
	})

	return nil
}

func (s *ArchiveService) ListMovements(ctx context.Context, archiveID string) ([]*ArchiveMovementResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	list, err := s.repo.ListMovements(ctx, archiveID)
	if err != nil {
		return nil, err
	}

	var res []*ArchiveMovementResponse
	for _, m := range list {
		res = append(res, &ArchiveMovementResponse{
			ID:                     m.ID,
			ArchiveID:              m.ArchiveID,
			Type:                   m.Type,
			PetugasID:              m.PetugasID,
			PeminjamID:             m.PeminjamID,
			Tanggal:                m.Tanggal,
			TanggalKembaliEstimasi: m.TanggalKembaliEstimasi,
			TanggalKembaliAktual:   m.TanggalKembaliAktual,
			Tujuan:                 m.Tujuan,
			Keterangan:             m.Keterangan,
			Status:                 m.Status,
		})
	}
	return res, nil
}

func toArchiveResponse(a *archiveModels.Archive) *ArchiveResponse {
	return &ArchiveResponse{
		ID:               a.ID,
		CorrespondenceID: a.CorrespondenceID,
		LokasiKode:       a.LokasiKode,
		Kategori:         a.Kategori,
		TanggalArsip:     a.TanggalArsip,
		TanggalMusnah:    a.TanggalMusnah,
		RetensiTahun:     a.RetensiTahun,
		Status:           a.Status,
		CurrentHolderID:  a.CurrentHolderID,
		Keterangan:       a.Keterangan,
		ArchivedBy:       a.ArchivedBy,
	}
}
