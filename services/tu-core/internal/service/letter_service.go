package service

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
	"sekolah-platform/services/tu-core/internal/storage"
)

type LetterService struct {
	repo        repository.LetterRepository
	studentRepo repository.StudentRepository
	minio       *storage.MinioClient
	publisher   EventPublisher
}

func NewLetterService(
	repo repository.LetterRepository,
	studentRepo repository.StudentRepository,
	minio *storage.MinioClient,
	publisher EventPublisher,
) *LetterService {
	if publisher == nil {
		publisher = NewNoopEventPublisher()
	}
	return &LetterService{repo: repo, studentRepo: studentRepo, minio: minio, publisher: publisher}
}

type CreateLetterRequest struct {
	StudentID    string `json:"studentId"`
	LetterNumber string `json:"letterNumber"`
	ParentName   string `json:"parentName"`
	Purpose      string `json:"purpose"`
}

func (s *LetterService) CreateLetter(ctx context.Context, req CreateLetterRequest) (*repository.StudentLetter, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	// 1. Ambil data siswa untuk PDF
	student, err := s.studentRepo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, req.StudentID)
	if err != nil {
		return nil, NewServiceError("CreateLetter", "Letter", req.StudentID, "siswa tidak ditemukan", err)
	}

	// 2. Generate PDF
	pdfBytes, err := s.generatePDF(student, req)
	if err != nil {
		return nil, NewServiceError("CreateLetter", "Letter", "", "gagal generate PDF", err)
	}

	// 3. Upload ke MinIO menggunakan UploadRawFile (karena method yang tersedia menerima io.Reader)
	fileName := fmt.Sprintf("surat-aktif-%s-%d.pdf", student.NIS, time.Now().Unix())
	url, err := s.minio.UploadRawFile(ctx, "documents", fileName, bytes.NewReader(pdfBytes), "application/pdf")
	if err != nil {
		return nil, NewServiceError("CreateLetter", "Letter", "", "gagal upload PDF ke MinIO", err)
	}

	// 4. Simpan metadata ke database
	letter := &repository.StudentLetter{
		ID:               uuid.New().String(),
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		StudentID:        req.StudentID,
		LetterNumber:     req.LetterNumber,
		ParentName:       req.ParentName,
		Purpose:          req.Purpose,
		Status:           "PENDING",
		FileURL:          url,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, letter); err != nil {
		return nil, NewServiceError("CreateLetter", "Letter", letter.ID, "gagal simpan metadata", err)
	}

	// Set nama untuk response
	letter.StudentName = student.Name
	return letter, nil
}

func (s *LetterService) ListLetters(ctx context.Context, status, search string) ([]repository.StudentLetter, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	filter := repository.LetterFilter{
		SchoolID: opCtx.SchoolID, AcademicPeriodID: opCtx.AcademicPeriodID,
		Status: status, Search: search, Limit: 100, Offset: 0,
	}
	letters, _, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, NewServiceError("ListLetters", "Letter", "", "gagal ambil daftar", err)
	}
	return letters, nil
}

func (s *LetterService) CompleteLetter(ctx context.Context, letterID string) error {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return ErrMissingOperationalContext
	}
	return s.repo.UpdateStatus(ctx, letterID, "COMPLETED")
}

func (s *LetterService) GetLetterByID(ctx context.Context, letterID string) (*repository.StudentLetter, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}
	return s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, letterID)
}

// generatePDF membuat PDF Surat Keterangan Aktif
func (s *LetterService) generatePDF(student *repository.Student, req CreateLetterRequest) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	
	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "PEMERINTAH PROVINSI JAWA BARAT")
	pdf.Ln(6)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "SMK SWASTA INDONESIA")
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(190, 5, "Jl. Raya Pendidikan No. 123, Bandung")
	pdf.Ln(15)
	
	// Garis
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(10)
	
	// Judul
	pdf.SetFont("Arial", "BU", 14)
	pdf.Cell(190, 10, "SURAT KETERANGAN AKTIF BELAJAR")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(190, 5, fmt.Sprintf("Nomor: %s", req.LetterNumber))
	pdf.Ln(15)
	
	// Isi
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(190, 6, "Yang bertanda tangan di bawah ini, Kepala SMK Swasta Indonesia, menerangkan bahwa:")
	pdf.Ln(8)
	
	pdf.Cell(30, 6, "Nama")
	pdf.Cell(5, 6, ":")
	pdf.Cell(100, 6, student.Name)
	pdf.Ln(6)
	
	pdf.Cell(30, 6, "NIS / NISN")
	pdf.Cell(5, 6, ":")
	pdf.Cell(100, 6, fmt.Sprintf("%s / %s", student.NIS, student.NISN))
	pdf.Ln(6)
	
	pdf.Cell(30, 6, "Kelas")
	pdf.Cell(5, 6, ":")
	pdf.Cell(100, 6, student.ClassName)
	pdf.Ln(6)
	
	pdf.Cell(30, 6, "Nama Orang Tua")
	pdf.Cell(5, 6, ":")
	pdf.Cell(100, 6, req.ParentName)
	pdf.Ln(10)
	
	pdf.MultiCell(190, 6, fmt.Sprintf("Adalah benar siswa aktif di SMK Swasta Indonesia pada Tahun Ajaran 2026/2027. Surat keterangan ini dibuat untuk keperluan: %s.", req.Purpose), "", "J", false)
	pdf.Ln(15)
	
	// Tanda tangan
	pdf.Cell(120, 6, "")
	pdf.Cell(70, 6, "Bandung, "+time.Now().Format("2 January 2006"))
	pdf.Ln(6)
	pdf.Cell(120, 6, "")
	pdf.Cell(70, 6, "Kepala Sekolah")
	pdf.Ln(20)
	pdf.Cell(120, 6, "")
	pdf.SetFont("Arial", "BU", 11)
	pdf.Cell(70, 6, "Drs. H. Ahmad Subarjo, M.Pd.")
	pdf.Ln(6)
	pdf.Cell(120, 6, "")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(70, 6, "NIP. 19700101 199001 1 001")
	
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
