package service

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	corrModels "sekolah-platform/services/tu-core/internal/domain/correspondence/models"
	expModels "sekolah-platform/services/tu-core/internal/domain/expedition/models"
	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

type ExpeditionService struct {
	repo       repository.ExpeditionRepository
	corrRepo   repository.CorrespondenceRepository
	schoolRepo repository.SchoolRepository // <-- TAMBAHAN
}

func NewExpeditionService(
	repo repository.ExpeditionRepository,
	corrRepo repository.CorrespondenceRepository,
	schoolRepo repository.SchoolRepository, // <-- TAMBAHAN
) *ExpeditionService {
	return &ExpeditionService{repo: repo, corrRepo: corrRepo, schoolRepo: schoolRepo}
}

func (s *ExpeditionService) CreateExpedition(ctx context.Context, correspondenceID string, req CreateExpeditionRequest) (*ExpeditionResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	// 1. Validasi Surat (menggunakan corrModels)
	corr, err := s.corrRepo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, correspondenceID)
	if err != nil {
		return nil, NewServiceError("CreateExpedition", "Expedition", "", "surat tidak ditemukan", err)
	}
	if corr.Type.String() != "OUTGOING" {
		return nil, NewServiceError("CreateExpedition", "Expedition", "", "ekspedisi hanya untuk surat keluar", nil)
	}

	// 2. Parse Tanggal
	tglKirim, err := time.Parse("2006-01-02", req.TanggalKirim)
	if err != nil {
		tglKirim = time.Now()
	}

	exp := &expModels.Expedition{
		SchoolID:         opCtx.SchoolID,
		CorrespondenceID: correspondenceID,
		PenerimaNama:     req.PenerimaNama,
		PenerimaInstansi: req.PenerimaInstansi,
		PenerimaAlamat:   req.PenerimaAlamat,
		MetodeKirim:      req.MetodeKirim,
		NomorResi:        req.NomorResi,
		TanggalKirim:     tglKirim,
		JumlahLembar:     req.JumlahLembar,
		Status:           "TERKIRIM",
		VerifyToken:      uuid.New().String()[:8],
		CreatedBy:        opCtx.UserID,
	}

	if err := s.repo.Create(ctx, exp); err != nil {
		return nil, NewServiceError("CreateExpedition", "Expedition", "", "failed to save to db", err)
	}

	// 3. Generate PDF Bukti Ekspedisi
	pdfURL, err := s.generatePDFBukti(ctx, opCtx.SchoolID, exp, corr)
	if err == nil && pdfURL != "" {
		exp.BuktiPdfURL = pdfURL
		_ = s.repo.UpdateBuktiURL(ctx, exp.ID, pdfURL)
	}

	return toExpeditionResponse(exp), nil
}

func (s *ExpeditionService) GetByCorrespondenceID(ctx context.Context, correspondenceID string) ([]*ExpeditionResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}
	list, err := s.repo.GetByCorrespondenceID(ctx, opCtx.SchoolID, correspondenceID)
	if err != nil {
		return nil, err
	}
	var res []*ExpeditionResponse
	for _, e := range list {
		res = append(res, toExpeditionResponse(&e))
	}
	return res, nil
}

func (s *ExpeditionService) generatePDFBukti(ctx context.Context, schoolID string, exp *expModels.Expedition, corr *corrModels.Correspondence) (string, error) {
	// Ambil data sekolah dari database
	school, err := s.schoolRepo.GetByID(ctx, schoolID)
	if err != nil {
		return "", fmt.Errorf("gagal ambil data sekolah: %v", err)
	}

	saveDir := "./storage/expeditions"
	_ = os.MkdirAll(saveDir, os.ModePerm)
	filePath := filepath.Join(saveDir, fmt.Sprintf("%s.pdf", exp.ID))

	// A5 Landscape
	pdf := gofpdf.New("L", "mm", "A5", "")
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(190, 7, "BUKTI EKSPEDISI SURAT", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(190, 5, school.Name, "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(190, 5, fmt.Sprintf("%s, %s, %s", school.Address, school.Kecamatan, school.KabKota), "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Detail
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(40, 6, "Tanggal Kirim", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "C", false, 0, "")
	pdf.CellFormat(100, 6, exp.TanggalKirim.Format("02 January 2006"), "", 1, "L", false, 0, "")

	pdf.CellFormat(40, 6, "Nomor Surat", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "C", false, 0, "")
	pdf.CellFormat(100, 6, corr.Number, "", 1, "L", false, 0, "")

	pdf.CellFormat(40, 6, "Perihal", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "C", false, 0, "")
	pdf.CellFormat(100, 6, corr.Subject, "", 1, "L", false, 0, "")

	pdf.CellFormat(40, 6, "Kepada", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "C", false, 0, "")
	pdf.CellFormat(100, 6, fmt.Sprintf("%s (%s)", exp.PenerimaNama, exp.PenerimaInstansi), "", 1, "L", false, 0, "")

	pdf.CellFormat(40, 6, "Metode", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "C", false, 0, "")
	pdf.CellFormat(100, 6, exp.MetodeKirim, "", 1, "L", false, 0, "")

	// QR Code
	qrContent := fmt.Sprintf("http://localhost:8081/api/v1/verify/expedition/%s", exp.VerifyToken)
	qr, _ := qrcode.New(qrContent, qrcode.Medium)
	qrPng, _ := qr.PNG(128)
	opt := gofpdf.ImageOptions{ImageType: "PNG"}
	pdf.RegisterImageOptionsReader("qr_"+exp.VerifyToken, opt, bytes.NewReader(qrPng))
	pdf.ImageOptions("qr_"+exp.VerifyToken, 160, 35, 25, 25, false, opt, 0, "")

	// Tanda Tangan
	pdf.SetY(80)
	pdf.SetX(30)
	pdf.CellFormat(60, 6, "Petugas Pengirim,", "", 0, "C", false, 0, "")
	pdf.SetX(120)
	pdf.CellFormat(60, 6, "Penerima,", "", 0, "C", false, 0, "")
	pdf.Ln(20)
	pdf.SetX(30)
	pdf.CellFormat(60, 6, "(_______________________)", "", 0, "C", false, 0, "")
	pdf.SetX(120)
	pdf.CellFormat(60, 6, "(_______________________)", "", 1, "C", false, 0, "")

	err = pdf.OutputFileAndClose(filePath)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("/storage/expeditions/%s.pdf", exp.ID), nil
}

func toExpeditionResponse(e *expModels.Expedition) *ExpeditionResponse {
	return &ExpeditionResponse{
		ID:               e.ID,
		CorrespondenceID: e.CorrespondenceID,
		PenerimaNama:     e.PenerimaNama,
		PenerimaInstansi: e.PenerimaInstansi,
		MetodeKirim:      e.MetodeKirim,
		NomorResi:        e.NomorResi,
		TanggalKirim:     e.TanggalKirim,
		JumlahLembar:     e.JumlahLembar,
		Status:           e.Status,
		BuktiPdfURL:      e.BuktiPdfURL,
		VerifyToken:      e.VerifyToken,
	}
}
