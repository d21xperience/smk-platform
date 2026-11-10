package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"sekolah-platform/services/tu-core/internal/repository"
	"sekolah-platform/services/tu-core/internal/storage"
)

type SKPembagianTugasService struct {
	repo           repository.SKPembagianTugasRepository
	schoolRepo     repository.SchoolRepository
	minioClient    *storage.MinioClient
	pdfGenerator   *SKPDFGenerator
	excelGenerator *SKExcelGenerator
}

func NewSKPembagianTugasService(repo repository.SKPembagianTugasRepository, schoolRepo repository.SchoolRepository, minioClient *storage.MinioClient) *SKPembagianTugasService {
	return &SKPembagianTugasService{
		repo:           repo,
		schoolRepo:     schoolRepo,
		minioClient:    minioClient,
		pdfGenerator:   NewSKPDFGenerator(),
		excelGenerator: NewSKExcelGenerator(),
	}
}

func (s *SKPembagianTugasService) CreateDraft(ctx context.Context, req repository.SKCreateRequest, schoolID, userID string) (string, string, error) {
	// Validasi: Semester boleh 0 (Semua Semester), 1 (Ganjil), atau 2 (Genap)
	if req.TahunAjaran == "" || (req.Semester != 0 && req.Semester != 1 && req.Semester != 2) || len(req.Assignments) == 0 {
		return "", "", fmt.Errorf("%w: Data input tidak valid. Semester harus 0, 1, atau 2", ErrInvalidInput)
	}

	teacherHours := make(map[string]int)
	for _, assign := range req.Assignments {
		teacherHours[assign.PtkID] += (assign.JamPerMinggu * len(assign.RombonganBelajarIDs))
	}
	for ptkID, total := range teacherHours {
		if total < 24 || total > 40 {
			return "", "", fmt.Errorf("%w: Guru ID %s mengajar %d jam. Harus 24-40 jam", ErrInvalidInput, ptkID, total)
		}
	}

	nomorSK, err := s.repo.GenerateNextSKNumber(ctx, schoolID)
	if err != nil {
		return "", "", fmt.Errorf("gagal generate nomor SK: %v", err)
	}

	skID, err := s.repo.CreateDraftSK(ctx, req, nomorSK, userID)
	if err != nil {
		return "", "", fmt.Errorf("gagal menyimpan draft SK: %v", err)
	}

	return skID, nomorSK, nil
}

func (s *SKPembagianTugasService) ListSK(ctx context.Context) ([]repository.SKHeader, error) {
	return s.repo.ListSK(ctx)
}

func (s *SKPembagianTugasService) SubmitSK(ctx context.Context, id string) error {
	return s.repo.UpdateStatus(ctx, id, "SUBMITTED", "")
}

func (s *SKPembagianTugasService) RejectSK(ctx context.Context, id, reason string) error {
	if reason == "" {
		return fmt.Errorf("%w: Alasan penolakan wajib diisi", ErrInvalidInput)
	}
	return s.repo.UpdateStatus(ctx, id, "REJECTED", reason)
}

func (s *SKPembagianTugasService) RequestRevision(ctx context.Context, id, notes string) error {
	if notes == "" {
		return fmt.Errorf("%w: Catatan revisi wajib diisi", ErrInvalidInput)
	}
	return s.repo.UpdateStatus(ctx, id, "REVISION", notes)
}

func (s *SKPembagianTugasService) ApproveAndIssueSK(ctx context.Context, id, schoolID string) error {
	school, err := s.schoolRepo.GetByID(ctx, schoolID)
	if err != nil {
		return fmt.Errorf("gagal mengambil profil sekolah: %v", err)
	}
	sk, err := s.repo.GetDetailByID(ctx, id)
	if err != nil {
		return fmt.Errorf("gagal mengambil detail SK: %v", err)
	}

	qrData := fmt.Sprintf("SK:%s|Kepsek:%s|Tanggal:%s", sk.NomorSK, school.PrincipalName, time.Now().Format("2006-01-02"))
	qrPath, err := s.pdfGenerator.GenerateQR(qrData)
	if err != nil {
		return fmt.Errorf("gagal generate QR: %v", err)
	}
	defer os.Remove(qrPath)

	skPdfPath, err := s.pdfGenerator.GenerateSKPDF(school, sk, qrPath)
	if err != nil {
		return fmt.Errorf("gagal generate PDF SK: %v", err)
	}
	defer os.Remove(skPdfPath)

	lampiranPdfPath, err := s.pdfGenerator.GenerateLampiranPDF(school, sk)
	if err != nil {
		return fmt.Errorf("gagal generate PDF Lampiran: %v", err)
	}
	defer os.Remove(lampiranPdfPath)

	uploadFile := func(localPath, folder string) (string, error) {
		file, err := os.Open(localPath)
		if err != nil { return "", err }
		defer file.Close()
		
		objectName := fmt.Sprintf("sk_%s_%s_%d.pdf", id, folder, time.Now().Unix())
		return s.minioClient.UploadRawFile(ctx, "documents", objectName, file, "application/pdf")
	}

	finalPdfUrl, err := uploadFile(skPdfPath, "utama")
	if err != nil { return fmt.Errorf("gagal upload PDF utama: %v", err) }
	
	lampiranUrl, err := uploadFile(lampiranPdfPath, "lampiran")
	if err != nil { return fmt.Errorf("gagal upload Lampiran: %v", err) }

	if err := s.repo.UpdateDocumentURLs(ctx, id, "", finalPdfUrl, lampiranUrl, qrData); err != nil {
		return fmt.Errorf("gagal update URL dokumen: %v", err)
	}

	return s.repo.UpdateStatus(ctx, id, "ISSUED", "")
}

func (s *SKPembagianTugasService) ExportExcel(ctx context.Context, id string) (string, error) {
	sk, err := s.repo.GetDetailByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("gagal mengambil detail SK: %v", err)
	}

	excelPath, err := s.excelGenerator.GenerateRekapExcel(sk)
	if err != nil {
		return "", fmt.Errorf("gagal generate Excel: %v", err)
	}
	defer os.Remove(excelPath)

	file, err := os.Open(excelPath)
	if err != nil {
		return "", fmt.Errorf("gagal membuka file Excel: %v", err)
	}
	defer file.Close()

	objectName := fmt.Sprintf("rekap_sk_%s_%d.xlsx", id, time.Now().Unix())
	url, err := s.minioClient.UploadRawFile(ctx, "documents", objectName, file, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	if err != nil {
		return "", fmt.Errorf("gagal upload Excel ke MinIO: %v", err)
	}

	return url, nil
}
