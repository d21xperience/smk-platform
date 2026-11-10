package service

import (
	"context"
	"fmt"
	"math"
	"os"

	"sekolah-platform/services/tu-core/internal/repository"
	"sekolah-platform/services/tu-core/internal/storage"
)

type SKPService struct {
	repo         repository.SKPRepository
	schoolRepo   repository.SchoolRepository
	minioClient  *storage.MinioClient
	pdfGenerator *SKPPDFGenerator
}

func NewSKPService(repo repository.SKPRepository, schoolRepo repository.SchoolRepository, minioClient *storage.MinioClient) *SKPService {
	return &SKPService{
		repo:         repo,
		schoolRepo:   schoolRepo,
		minioClient:  minioClient,
		pdfGenerator: NewSKPPDFGenerator(),
	}
}

func (s *SKPService) ListTemplates(ctx context.Context) ([]repository.SKPTemplate, error) {
	return s.repo.ListTemplates(ctx)
}

func (s *SKPService) ListPeriodes(ctx context.Context) ([]repository.SKPPeriode, error) {
	return s.repo.ListPeriodes(ctx)
}

func (s *SKPService) CreatePeriode(ctx context.Context, req repository.SKPPeriode) (string, error) {
	if req.NamaPeriode == "" || req.JenisPeriode == "" {
		return "", fmt.Errorf("%w: Nama dan Jenis Periode wajib diisi", ErrInvalidInput)
	}
	return s.repo.CreatePeriode(ctx, req)
}

func (s *SKPService) ListPenilaianByPeriode(ctx context.Context, periodeID string) ([]repository.SKPPenilaian, error) {
	return s.repo.ListPenilaianByPeriode(ctx, periodeID)
}

func (s *SKPService) GetHasilByPenilaianID(ctx context.Context, penilaianID string) (*repository.SKPHasil, []repository.SKPDetailItem, error) {
	hasil, err := s.repo.GetHasilByPenilaianID(ctx, penilaianID)
	if err != nil {
		return nil, nil, err
	}
	detail, err := s.repo.GetDetailByPenilaianID(ctx, penilaianID)
	if err != nil {
		return nil, nil, err
	}
	return hasil, detail, nil
}

// Fungsi BARU untuk mengambil detail penilaian untuk Form
func (s *SKPService) GetPenilaianDetail(ctx context.Context, penilaianID string) (*repository.SKPPenilaian, *repository.SKPHasil, []repository.SKPDetailItem, error) {
	penilaian, err := s.repo.GetPenilaianByID(ctx, penilaianID)
	if err != nil {
		return nil, nil, nil, err
	}

	hasil, err := s.repo.GetHasilByPenilaianID(ctx, penilaianID)
	if err != nil && err.Error() != "no rows in result set" {
		return nil, nil, nil, err
	}

	detail, err := s.repo.GetDetailByPenilaianID(ctx, penilaianID)
	if err != nil {
		return nil, nil, nil, err
	}

	return penilaian, hasil, detail, nil
}

func (s *SKPService) SubmitPenilaian(ctx context.Context, penilaianID string, hasil repository.SKPHasil, detail []repository.SKPDetailItem) error {
	penilaian, err := s.repo.GetPenilaianByID(ctx, penilaianID)
	if err != nil {
		return fmt.Errorf("penilaian tidak ditemukan: %v", err)
	}

	templates, err := s.repo.ListTemplates(ctx)
	if err != nil {
		return err
	}
	var tpl repository.SKPTemplate
	for _, t := range templates {
		if t.ID == penilaian.TemplateID {
			tpl = t
			break
		}
	}

	totalPerilaku := 0
	countPerilaku := 0
	for _, d := range detail {
		if d.Kelompok == "PERILAKU" {
			totalPerilaku += d.Nilai
			countPerilaku++
		}
	}
	if countPerilaku == 0 { return fmt.Errorf("aspek perilaku kerja wajib diisi") }
	rataPerilaku := float64(totalPerilaku) / float64(countPerilaku)

	faktorPerilaku := 1.00
	predikatPerilaku := "BAIK"
	if rataPerilaku >= 91 { faktorPerilaku = 1.20; predikatPerilaku = "SANGAT_BAIK" } else 
	if rataPerilaku >= 76 { faktorPerilaku = 1.00; predikatPerilaku = "BAIK" } else 
	if rataPerilaku >= 61 { faktorPerilaku = 0.90; predikatPerilaku = "BUTUH_BIMBINGAN" } else 
	if rataPerilaku >= 51 { faktorPerilaku = 0.75; predikatPerilaku = "KURANG" } else 
	{ faktorPerilaku = 0.60; predikatPerilaku = "BURUK" }

	totalUtama := 0
	countUtama := 0
	for _, d := range detail {
		if d.Kelompok == "HASIL_UTAMA" {
			totalUtama += d.Nilai
			countUtama++
		}
	}
	if countUtama == 0 { return fmt.Errorf("aspek hasil kerja utama wajib diisi") }
	rataUtama := float64(totalUtama) / float64(countUtama)

	totalTambahan := 0
	countTambahan := 0
	for _, d := range detail {
		if d.Kelompok == "HASIL_TAMBAHAN" {
			totalTambahan += d.Nilai
			countTambahan++
		}
	}
	rataTambahan := 0.0
	if countTambahan > 0 { rataTambahan = float64(totalTambahan) / float64(countTambahan) }

	nilaiUtamaTerkoreksi := rataUtama * faktorPerilaku
	nilaiTambahanTerkoreksi := rataTambahan * faktorPerilaku

	b1 := float64(tpl.BobotPerilaku) / 100.0
	b2 := float64(tpl.BobotHasilUtama) / 100.0
	b3 := float64(tpl.BobotHasilTambahan) / 100.0

	nilaiAkhir := (rataPerilaku * b1) + (nilaiUtamaTerkoreksi * b2) + (nilaiTambahanTerkoreksi * b3)
	nilaiAkhir = math.Round(nilaiAkhir*100) / 100

	predikatAkhir := "BURUK"
	if nilaiAkhir >= 91 { predikatAkhir = "SANGAT_BAIK" } else 
	if nilaiAkhir >= 76 { predikatAkhir = "BAIK" } else 
	if nilaiAkhir >= 61 { predikatAkhir = "BUTUH_BIMBINGAN" } else 
	if nilaiAkhir >= 51 { predikatAkhir = "KURANG" }

	hasil.PenilaianID = penilaianID
	hasil.NilaiPerilaku = rataPerilaku
	hasil.PredikatPerilaku = predikatPerilaku
	hasil.NilaiHasilUtama = rataUtama
	hasil.NilaiHasilTambahan = rataTambahan
	hasil.NilaiAkhir = nilaiAkhir
	hasil.PredikatAkhir = predikatAkhir

	if err := s.repo.SaveHasil(ctx, &hasil, detail); err != nil {
		return err
	}

	return s.repo.UpdateStatus(ctx, penilaianID, "SUBMITTED")
}

func (s *SKPService) ApprovePenilaian(ctx context.Context, penilaianID string, catatan string, schoolID string) error {
	penilaian, err := s.repo.GetPenilaianByID(ctx, penilaianID)
	if err != nil { return err }
	if penilaian.Status != "SUBMITTED" {
		return fmt.Errorf("%w: Hanya penilaian berstatus SUBMITTED yang bisa disetujui", ErrInvalidInput)
	}

	hasil, detail, err := s.GetHasilByPenilaianID(ctx, penilaianID)
	if err != nil { return err }

	school, err := s.schoolRepo.GetByID(ctx, schoolID)
	if err != nil { return fmt.Errorf("gagal mengambil profil sekolah: %v", err) }

	qrData := fmt.Sprintf("SKP:%s|Nilai:%.2f|Predikat:%s", penilaianID, hasil.NilaiAkhir, hasil.PredikatAkhir)
	qrPath, err := s.pdfGenerator.GenerateQR(qrData)
	if err != nil { return fmt.Errorf("gagal generate QR: %v", err) }
	defer os.Remove(qrPath)

	pdfPath, err := s.pdfGenerator.GenerateSKPPDF(school, penilaian, hasil, detail, qrPath)
	if err != nil { return fmt.Errorf("gagal generate PDF: %v", err) }
	defer os.Remove(pdfPath)

	file, err := os.Open(pdfPath)
	if err != nil { return err }
	defer file.Close()

	objectName := fmt.Sprintf("skp_formulir_%s.pdf", penilaianID)
	_, err = s.minioClient.UploadRawFile(ctx, "documents", objectName, file, "application/pdf")
	if err != nil { return fmt.Errorf("gagal upload PDF: %v", err) }

	return s.repo.UpdateStatusWithCatatan(ctx, penilaianID, "FINAL", catatan)
}

func (s *SKPService) ReturnPenilaian(ctx context.Context, penilaianID string, catatan string) error {
	penilaian, err := s.repo.GetPenilaianByID(ctx, penilaianID)
	if err != nil { return err }
	if penilaian.Status != "SUBMITTED" {
		return fmt.Errorf("%w: Hanya penilaian berstatus SUBMITTED yang bisa dikembalikan", ErrInvalidInput)
	}
	return s.repo.UpdateStatusWithCatatan(ctx, penilaianID, "RETURNED", catatan)
}
