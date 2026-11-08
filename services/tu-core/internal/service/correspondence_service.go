package service

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"sekolah-platform/services/tu-core/internal/domain/correspondence/engine"
	"sekolah-platform/services/tu-core/internal/domain/correspondence/models"
	"sekolah-platform/services/tu-core/internal/domain/correspondence/valueobjects"
	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

type CorrespondenceService struct {
	repo       repository.CorrespondenceRepository
	sigService *SignatureService
	corrEngine *engine.CorrespondenceEngine
	publisher  EventPublisher
}

func NewCorrespondenceService(
	repo repository.CorrespondenceRepository,
	sigService *SignatureService,
	corrEngine *engine.CorrespondenceEngine,
	publisher EventPublisher,
) *CorrespondenceService {
	if publisher == nil {
		publisher = NewNoopEventPublisher()
	}
	return &CorrespondenceService{
		repo:       repo,
		sigService: sigService,
		corrEngine: corrEngine,
		publisher:  publisher,
	}
}

func (s *CorrespondenceService) CreateCorrespondence(ctx context.Context, req CreateCorrespondenceRequest) (*CorrespondenceResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	corrType, err := valueobjects.NewCorrespondenceType(req.Type)
	if err != nil {
		return nil, NewServiceError("CreateCorrespondence", "Correspondence", "", "invalid type", err)
	}

	if req.Number == "" {
		seq, err := s.repo.GenerateNumber(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, req.JenisSurat)
		if err != nil {
			return nil, NewServiceError("CreateCorrespondence", "Correspondence", "", "failed to generate number", err)
		}
		monthRoman := []string{"", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X", "XI", "XII"}[int(req.Date.Month())]
		
		if req.KodeKlasifikasi != "" {
			req.Number = fmt.Sprintf("%s/%s", req.KodeKlasifikasi, generateNumberFormat(seq, req.JenisSurat, monthRoman, req.Date.Year()))
		} else {
			req.Number = generateNumberFormat(seq, req.JenisSurat, monthRoman, req.Date.Year())
		}
	}

	data := models.CorrespondenceData{
		Type:            corrType,
		JenisSurat:      req.JenisSurat,
		Kategori:        req.Kategori,
		KodeKlasifikasi: req.KodeKlasifikasi,
		Number:          req.Number,
		Date:            req.Date,
		Subject:         req.Subject,
		FromParty:       req.From,
		ToParty:         req.To,
		Description:     req.Description,
		AttachmentURL:   req.AttachmentURL,
		FileURL:         req.FileURL,
		TTDURL:          req.TTDURL,
		Konten:          req.Konten,
	}

	result := s.corrEngine.CreateCorrespondence(data, opCtx)
	if result.Error != nil {
		return nil, NewServiceError("CreateCorrespondence", "Correspondence", "", "invalid correspondence data", result.Error)
	}

	if err := s.repo.Create(ctx, result.Correspondence); err != nil {
		return nil, NewServiceError("CreateCorrespondence", "Correspondence", result.Correspondence.ID, "failed to save to database", err)
	}

	event := DomainEvent{
		EventType:     "CorrespondenceCreated",
		AggregateType: "Correspondence",
		AggregateID:   result.Correspondence.ID,
		Payload: map[string]interface{}{
			"correspondenceId": result.Correspondence.ID,
			"type":             result.Correspondence.Type.String(),
			"number":           result.Correspondence.Number,
		},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return toCorrespondenceResponse(result.Correspondence), nil
}

func (s *CorrespondenceService) UpdateCorrespondence(ctx context.Context, id string, req UpdateCorrespondenceRequest) (*CorrespondenceResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	c, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, id)
	if err != nil {
		return nil, NewServiceError("UpdateCorrespondence", "Correspondence", id, "failed to get correspondence", err)
	}

	if err := c.Update(req.Kategori, req.KodeKlasifikasi, req.Subject, req.Description, req.AttachmentURL, req.FileURL, req.TTDURL, req.Konten); err != nil {
		return nil, NewServiceError("UpdateCorrespondence", "Correspondence", id, "failed to update", err)
	}

	if err := s.repo.Update(ctx, c); err != nil {
		return nil, NewServiceError("UpdateCorrespondence", "Correspondence", id, "failed to save to database", err)
	}

	return toCorrespondenceResponse(c), nil
}

func (s *CorrespondenceService) DeleteCorrespondence(ctx context.Context, id string) error {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return ErrMissingOperationalContext
	}

	err := s.repo.Delete(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, id)
	if err != nil {
		return NewServiceError("DeleteCorrespondence", "Correspondence", id, "failed to delete", err)
	}

	return nil
}

func (s *CorrespondenceService) ArchiveCorrespondence(ctx context.Context, id string) (*CorrespondenceResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	c, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, id)
	if err != nil {
		return nil, NewServiceError("ArchiveCorrespondence", "Correspondence", id, "failed to get correspondence", err)
	}

	result := s.corrEngine.ArchiveCorrespondence(c, opCtx)
	if result.Error != nil {
		return nil, NewServiceError("ArchiveCorrespondence", "Correspondence", id, "failed to archive", result.Error)
	}

	if err := s.repo.Update(ctx, result.Correspondence); err != nil {
		return nil, NewServiceError("ArchiveCorrespondence", "Correspondence", id, "failed to save to database", err)
	}

	return toCorrespondenceResponse(result.Correspondence), nil
}

func (s *CorrespondenceService) GetCorrespondence(ctx context.Context, id string) (*CorrespondenceResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	c, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, id)
	if err != nil {
		return nil, NewServiceError("GetCorrespondence", "Correspondence", id, "failed to get correspondence", err)
	}

	return toCorrespondenceResponse(c), nil
}

func (s *CorrespondenceService) ListCorrespondences(ctx context.Context, filter repository.CorrespondenceFilter) ([]*CorrespondenceResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	filter.SchoolID = opCtx.SchoolID
	filter.AcademicPeriodID = opCtx.AcademicPeriodID

	if filter.Limit <= 0 {
		filter.Limit = 100
	}

	correspondences, _, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, NewServiceError("ListCorrespondences", "Correspondence", "", "failed to list correspondences", err)
	}

	var responses []*CorrespondenceResponse
	for _, c := range correspondences {
		responses = append(responses, toCorrespondenceResponse(&c))
	}

	return responses, nil
}

// GeneratePDF membuat file PDF surat, menyisipkan TTD & QR Code, dan mengembalikan path relatif URL-nya
func (s *CorrespondenceService) GeneratePDF(ctx context.Context, id string) (string, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return "", ErrMissingOperationalContext
	}

	c, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, id)
	if err != nil {
		return "", NewServiceError("GeneratePDF", "Correspondence", id, "failed to get correspondence", err)
	}

	// Buat direktori penyimpanan jika belum ada
	saveDir := "./storage/pdfs"
	_ = os.MkdirAll(saveDir, os.ModePerm)
	filePath := filepath.Join(saveDir, fmt.Sprintf("%s.pdf", c.ID))

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Header / Kop Surat Sederhana
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(190, 10, "PEMERINTAH PROVINSI JAWA TIMUR", "", 1, "C", false, 0, "")
	pdf.CellFormat(190, 10, "SMK PASUNDAN JATINANGOR", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(190, 5, "Jl. Raya Jatinangor No. 1, Kabupaten Sumedang, Jawa Barat", "", 1, "C", false, 0, "")
	pdf.Line(10, 38, 200, 38) // Garis Kop Surat
	pdf.Ln(10)

	// Judul Surat
	pdf.SetFont("Arial", "BU", 14)
	pdf.CellFormat(190, 10, c.JenisSurat, "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(190, 10, fmt.Sprintf("Nomor: %s", c.Number), "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Isi Surat
	pdf.SetFont("Arial", "", 11)
	pdf.MultiCell(190, 6, fmt.Sprintf("Yang bertanda tangan di bawah ini, Kepala SMK Pasundan Jatinangor, dengan ini menerangkan bahwa surat ini dibuat untuk keperluan: %s.", c.Description), "", "L", false)
	pdf.Ln(5)

	pdf.MultiCell(190, 6, fmt.Sprintf("Perihal: %s", c.Subject), "", "L", false)
	pdf.MultiCell(190, 6, fmt.Sprintf("Kepada: %s", c.ToParty), "", "L", false)
	pdf.Ln(10)

	// Cek apakah ada TTD Digital
	sig, _ := s.sigService.GetByCorrespondenceID(ctx, opCtx.SchoolID, c.ID)
	
	// Posisi TTD
	pdf.SetY(190) // Set posisi Y di bawah (kira-kira 190mm dari atas)
	pdf.SetX(130)
	
	if sig != nil {
		// Render QR Code
		qrContent := fmt.Sprintf("http://localhost:9000/verify/%s", sig.VerificationCode) // Ganti URL sesuai domain frontend Anda
		qr, _ := qrcode.New(qrContent, qrcode.Medium)
		qrPng, _ := qr.PNG(256)

		opt := gofpdf.ImageOptions{ImageType: "PNG"}
		pdf.RegisterImageOptionsReader("qr_"+sig.VerificationCode, opt, bytes.NewReader(qrPng))
		pdf.ImageOptions("qr_"+sig.VerificationCode, 155, 190, 25, 25, false, opt, 0, "")

		// Render Gambar TTD
		// TTD ditaruh di sebelah kiri QR Code
		optTtd := gofpdf.ImageOptions{ImageType: "PNG"}
		pdf.RegisterImageOptions(sig.ImageURL, optTtd)
		pdf.ImageOptions(sig.ImageURL, 120, 195, 30, 0, false, optTtd, 0, "")

		pdf.SetY(200)
		pdf.SetX(90)
		pdf.SetFont("Arial", "I", 8)
		pdf.CellFormat(90, 5, "Ditandatangani secara elektronik", "", 1, "L", false, 0, "")
		pdf.SetX(90)
		pdf.CellFormat(90, 5, fmt.Sprintf("Kode Verifikasi: %s", sig.VerificationCode), "", 1, "L", false, 0, "")
	} else {
		// TTD Manual (Kosong)
		pdf.CellFormat(90, 6, fmt.Sprintf("Jatinangor, %s", c.Date.Format("02 January 2006")), "", 1, "L", false, 0, "")
		pdf.CellFormat(90, 6, "Kepala Sekolah,", "", 1, "L", false, 0, "")
		pdf.Ln(20)
		pdf.CellFormat(90, 6, "________________________", "", 1, "L", false, 0, "")
	}

	// Simpan PDF
	err = pdf.OutputFileAndClose(filePath)
	if err != nil {
		return "", NewServiceError("GeneratePDF", "Correspondence", id, "failed to generate pdf file", err)
	}

	// Update URL file di database
	fileURL := fmt.Sprintf("/api/v1/correspondences/%s/pdf", c.ID)
	c.FileURL = fileURL
	_ = s.repo.Update(ctx, c)

	return fileURL, nil
}

func toCorrespondenceResponse(c *models.Correspondence) *CorrespondenceResponse {
	if c == nil {
		return nil
	}
	return &CorrespondenceResponse{
		ID:              c.ID,
		Type:            c.Type.String(),
		JenisSurat:      c.JenisSurat,
		Kategori:        c.Kategori,
		KodeKlasifikasi: c.KodeKlasifikasi,
		Number:          c.Number,
		Date:            c.Date,
		Subject:         c.Subject,
		From:            c.FromParty,
		To:              c.ToParty,
		Description:     c.Description,
		AttachmentURL:   c.AttachmentURL,
		FileURL:         c.FileURL,
		TTDURL:          c.TTDURL,
		Konten:          c.Konten,
		Status:          c.Status.String(),
		CreatedAt:       c.CreatedAt,
	}
}

func generateNumberFormat(seq int, jenis, romanMonth string, year int) string {
	return fmt.Sprintf("%03d/%s/SMK.1/%s/%d", seq, jenis, romanMonth, year)
}
