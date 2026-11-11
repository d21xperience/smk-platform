package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"sekolah-platform/services/tu-core/internal/repository"
	"sekolah-platform/services/tu-core/internal/storage"
)

type PelayananService struct {
	repo        repository.PelayananRepository
	schoolRepo  repository.SchoolRepository
	minioClient *storage.MinioClient
}

func NewPelayananService(repo repository.PelayananRepository, schoolRepo repository.SchoolRepository, minioClient *storage.MinioClient) *PelayananService {
	return &PelayananService{repo: repo, schoolRepo: schoolRepo, minioClient: minioClient}
}

func (s *PelayananService) ListKatalog(ctx context.Context) ([]repository.MasterLayanan, error) {
	return s.repo.ListKatalog(ctx)
}

func (s *PelayananService) CreatePengajuan(ctx context.Context, ptkID, kode, tujuan, kepada, ket, dokURL string) (string, error) {
	katalog, err := s.repo.GetKatalogByKode(ctx, kode)
	if err != nil {
		return "", fmt.Errorf("layanan tidak ditemukan: %v", err)
	}

	nomor := fmt.Sprintf("PL/%d/%03d", time.Now().Year(), time.Now().UnixNano()%1000)
	estimasi := time.Now().AddDate(0, 0, katalog.SLABulan)

	id, err := s.repo.CreatePengajuan(ctx, ptkID, kode, tujuan, kepada, ket, dokURL, nomor, estimasi)
	if err != nil {
		return "", err
	}

	s.repo.AddTimeline(ctx, id, "", "SUBMITTED", ptkID, "Pengajuan dibuat")
	return id, nil
}

func (s *PelayananService) ListPengajuan(ctx context.Context, status string) ([]repository.PengajuanLayanan, error) {
	return s.repo.ListPengajuan(ctx, status)
}

func (s *PelayananService) Verify(ctx context.Context, id, action, catatan, userID string) error {
	pengajuan, err := s.repo.GetPengajuanByID(ctx, id)
	if err != nil { return err }

	statusBaru := "VERIFIED"
	if action == "REJECT" { statusBaru = "REJECTED" }
	if action == "REVISION" { statusBaru = "REVISION" }

	err = s.repo.UpdateStatus(ctx, id, statusBaru, catatan, userID)
	if err == nil {
		s.repo.AddTimeline(ctx, id, pengajuan.Status, statusBaru, userID, catatan)
	}
	return err
}

func (s *PelayananService) Approve(ctx context.Context, id, userID, schoolID string) (string, error) {
	pengajuan, err := s.repo.GetPengajuanByID(ctx, id)
	if err != nil { return "", err }
	if pengajuan.Status != "VERIFIED" {
		return "", fmt.Errorf("pengajuan belum diverifikasi")
	}

	school, err := s.schoolRepo.GetByID(ctx, schoolID)
	if err != nil { return "", err }

	nomorSurat := fmt.Sprintf("421.2/%d/%02d/%d", time.Now().UnixNano()%1000, time.Now().Month(), time.Now().Year())

	// Generate QR
	qrPath, err := GenerateQRPelayanan(nomorSurat)
	if err != nil { return "", err }
	defer os.Remove(qrPath)

	// Generate PDF
	pdfPath, err := GeneratePelayananPDF(school, pengajuan, nomorSurat, qrPath)
	if err != nil { return "", err }
	defer os.Remove(pdfPath)

	// Upload MinIO
	file, err := os.Open(pdfPath)
	if err != nil { return "", err }
	defer file.Close()

	objectName := fmt.Sprintf("surat_pelayanan_%s.pdf", id)
	url, err := s.minioClient.UploadRawFile(ctx, "documents", objectName, file, "application/pdf")
	if err != nil { return "", err }

	// Update Status
	s.repo.UpdateStatus(ctx, id, "ISSUED", "", userID)
	s.repo.AddTimeline(ctx, id, pengajuan.Status, "ISSUED", userID, "Surat diterbitkan: "+nomorSurat)

	// Simpan ke pelayanan_surat (Log)
	// (Akan ditambahkan di Fase 3 jika perlu)

	return url, nil
}
