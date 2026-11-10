package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"sekolah-platform/services/tu-core/internal/repository"
	"sekolah-platform/services/tu-core/internal/storage"
)

type LaporanBerkalaService struct {
	repo      repository.LaporanBerkalaRepository
	schoolRepo repository.SchoolRepository
	minio     *storage.MinioClient
}

func NewLaporanBerkalaService(repo repository.LaporanBerkalaRepository, schoolRepo repository.SchoolRepository, minio *storage.MinioClient) *LaporanBerkalaService {
	return &LaporanBerkalaService{repo: repo, schoolRepo: schoolRepo, minio: minio}
}

func (s *LaporanBerkalaService) GenerateDocuments(ctx context.Context, id string, schoolID string) error {
	laporan, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("laporan tidak ditemukan: %v", err)
	}

	school, err := s.schoolRepo.GetByID(ctx, schoolID)
	if err != nil {
		return fmt.Errorf("gagal ambil data sekolah: %v", err)
	}

	// 1. Generate Excel
	excelDinasPath, _ := NewLaporanExcelGenerator().GenerateExcel("Dinas", laporan.Tahun, laporan.Triwulan, laporan.DataJSON)
	defer os.Remove(excelDinasPath)
	
	excelYayasanPath, _ := NewLaporanExcelGenerator().GenerateExcel("Yayasan", laporan.Tahun, laporan.Triwulan, laporan.DataJSON)
	defer os.Remove(excelYayasanPath)

	// 2. Generate PDF
	pdfGen := NewLaporanPDFGenerator()
	pdfPath, _ := pdfGen.GenerateLaporanPDF(school, laporan)
	defer os.Remove(pdfPath)

	ringkasanPath, _ := pdfGen.GenerateRingkasanPDF(school, laporan)
	defer os.Remove(ringkasanPath)

	suratPath, _ := pdfGen.GenerateSuratPengantar(school, laporan)
	defer os.Remove(suratPath)

	// 3. Upload ke MinIO
	upload := func(localPath, folder string) (string, error) {
		file, err := os.Open(localPath)
		if err != nil { return "", err }
		defer file.Close()
		
		ext := ".pdf"
		if folder == "excel_dinas" || folder == "excel_yayasan" { ext = ".xlsx" }
		
		objectName := fmt.Sprintf("laporan_%s_%s_%d%s", id, folder, time.Now().Unix(), ext)
		contentType := "application/pdf"
		if ext == ".xlsx" { contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" }
		
		return s.minio.UploadRawFile(ctx, "documents", objectName, file, contentType)
	}

	excelDinasUrl, err := upload(excelDinasPath, "excel_dinas")
	if err != nil { return fmt.Errorf("gagal upload excel dinas: %v", err) }
	
	excelYayasanUrl, err := upload(excelYayasanPath, "excel_yayasan")
	if err != nil { return fmt.Errorf("gagal upload excel yayasan: %v", err) }

	pdfUrl, err := upload(pdfPath, "pdf_utama")
	if err != nil { return fmt.Errorf("gagal upload pdf: %v", err) }

	ringkasanUrl, err := upload(ringkasanPath, "ringkasan")
	if err != nil { return fmt.Errorf("gagal upload ringkasan: %v", err) }

	suratUrl, err := upload(suratPath, "surat")
	if err != nil { return fmt.Errorf("gagal upload surat: %v", err) }

	// 4. Update Database
	return s.repo.UpdateDocuments(ctx, id, pdfUrl, excelDinasUrl, excelYayasanUrl, ringkasanUrl, suratUrl)
}

func (s *LaporanBerkalaService) ReviewLaporan(ctx context.Context, id, userID, notes string) error {
	if notes == "" { return fmt.Errorf("catatan KTU wajib diisi") }
	err := s.repo.UpdateStatus(ctx, id, "REVIEWED", userID, notes)
	if err != nil { return err }
	return s.repo.LogTimeline(ctx, id, "DRAFT", "REVIEWED", userID, notes)
}

func (s *LaporanBerkalaService) ApproveLaporan(ctx context.Context, id, userID, notes string) error {
	err := s.repo.UpdateStatus(ctx, id, "APPROVED", userID, notes)
	if err != nil { return err }
	return s.repo.LogTimeline(ctx, id, "REVIEWED", "APPROVED", userID, notes)
}

func (s *LaporanBerkalaService) SendLaporan(ctx context.Context, id, userID string, emailService *LaporanEmailService, toEmail string) error {
	laporan, err := s.repo.GetByID(ctx, id)
	if err != nil { return err }
	
	subject := fmt.Sprintf("Laporan Kepegawaian Triwulan %d Tahun %d", laporan.Triwulan, laporan.Tahun)
	body := fmt.Sprintf("<p>Berikut kami sampaikan Laporan Kepegawaian Triwulan %d Tahun %d.</p>", laporan.Triwulan, laporan.Tahun)
	
	if err := emailService.SendLaporan(ctx, laporan, toEmail, subject, body); err != nil {
		return err
	}
	
	return s.repo.LogTimeline(ctx, id, "APPROVED", "SENT", userID, "Dikirim ke "+toEmail)
}
