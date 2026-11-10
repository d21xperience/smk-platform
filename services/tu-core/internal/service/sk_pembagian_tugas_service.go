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
	repo         repository.SKPembagianTugasRepository
	schoolRepo   repository.SchoolRepository
	minioClient  *storage.MinioClient
	pdfGenerator *SKPDFGenerator
}

func NewSKPembagianTugasService(repo repository.SKPembagianTugasRepository, schoolRepo repository.SchoolRepository, minioClient *storage.MinioClient) *SKPembagianTugasService {
	return &SKPembagianTugasService{
		repo:         repo,
		schoolRepo:   schoolRepo,
		minioClient:  minioClient,
		pdfGenerator: NewSKPDFGenerator(),
	}
}

func (s *SKPembagianTugasService) CreateDraft(ctx context.Context, req repository.SKCreateRequest, schoolID, userID string) (string, string, error) {
	if req.TahunAjaran == "" || (req.Semester != 1 && req.Semester != 2) || len(req.Assignments) == 0 {
		return "", "", fmt.Errorf("%w: Data input tidak valid", ErrInvalidInput)
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
	// 1. Ambil Data
	school, err := s.schoolRepo.GetByID(ctx, schoolID)
	if err != nil {
		return fmt.Errorf("gagal mengambil profil sekolah: %v", err)
	}
	sk, err := s.repo.GetDetailByID(ctx, id)
	if err != nil {
		return fmt.Errorf("gagal mengambil detail SK: %v", err)
	}

	// 2. Generate QR TTE
	qrData := fmt.Sprintf("SK:%s|Kepsek:%s|Tanggal:%s", sk.NomorSK, school.PrincipalName, time.Now().Format("2006-01-02"))
	qrPath, err := s.pdfGenerator.GenerateQR(qrData)
	if err != nil {
		return fmt.Errorf("gagal generate QR: %v", err)
	}
	defer os.Remove(qrPath)

	// 3. Generate PDF SK & Lampiran
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

	// 4. Upload ke MinIO (Gunakan MinIO Client yang menerima io.Reader, bukan multipart.FileHeader)
	uploadFile := func(localPath, folder string) (string, error) {
		file, err := os.Open(localPath)
		if err != nil {
			return "", err
		}
		defer file.Close()

		objectName := fmt.Sprintf("sk_%s_%s_%d.pdf", id, folder, time.Now().Unix())
		// Kita panggil method upload baru di minioClient (lihat step 4 di bawah)
		return s.minioClient.UploadRawFile(ctx, "documents", objectName, file, "application/pdf")
	}

	finalPdfUrl, err := uploadFile(skPdfPath, "utama")
	if err != nil {
		return fmt.Errorf("gagal upload PDF utama: %v", err)
	}

	lampiranUrl, err := uploadFile(lampiranPdfPath, "lampiran")
	if err != nil {
		return fmt.Errorf("gagal upload Lampiran: %v", err)
	}

	// 5. Update Database
	if err := s.repo.UpdateDocumentURLs(ctx, id, "", finalPdfUrl, lampiranUrl, qrData); err != nil {
		return fmt.Errorf("gagal update URL dokumen: %v", err)
	}

	return s.repo.UpdateStatus(ctx, id, "ISSUED", "")
}
