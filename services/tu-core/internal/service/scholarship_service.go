package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
)

type ScholarshipService struct {
	repo        repository.ScholarshipRepository
	studentRepo repository.StudentRepository
	publisher   EventPublisher
}

func NewScholarshipService(
	repo repository.ScholarshipRepository,
	studentRepo repository.StudentRepository,
	publisher EventPublisher,
) *ScholarshipService {
	if publisher == nil {
		publisher = NewNoopEventPublisher()
	}
	return &ScholarshipService{repo: repo, studentRepo: studentRepo, publisher: publisher}
}

type CreateScholarshipRequest struct {
	StudentID       string  `json:"studentId"`
	Category        string  `json:"category"`
	CardNumber      string  `json:"cardNumber"`
	Reason          string  `json:"reason"`
	IsDocumentValid bool    `json:"isDocumentValid"`
	BankName        string  `json:"bankName"`
	BankAccount     string  `json:"bankAccount"`
	DisbursedAmount float64 `json:"disbursedAmount"`
}

type UpdateScholarshipRequest struct {
	Category        string  `json:"category"`
	CardNumber      string  `json:"cardNumber"`
	Reason          string  `json:"reason"`
	IsDocumentValid bool    `json:"isDocumentValid"`
	Status          string  `json:"status"`
	BankName        string  `json:"bankName"`
	BankAccount     string  `json:"bankAccount"`
	DisbursedAmount float64 `json:"disbursedAmount"`
}

func (s *ScholarshipService) CreateScholarship(ctx context.Context, req CreateScholarshipRequest) (*repository.StudentScholarship, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	// Validasi: pastikan siswa ada
	student, err := s.studentRepo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, req.StudentID)
	if err != nil {
		return nil, NewServiceError("CreateScholarship", "Scholarship", req.StudentID, "siswa tidak ditemukan", err)
	}

	scholarship := &repository.StudentScholarship{
		ID:               uuid.New().String(),
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		StudentID:        req.StudentID,
		Category:         req.Category,
		CardNumber:       req.CardNumber,
		Reason:           req.Reason,
		IsDocumentValid:  req.IsDocumentValid,
		Status:           "SK NOMINASI",
		BankName:         req.BankName,
		BankAccount:      req.BankAccount,
		DisbursedAmount:  req.DisbursedAmount,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, scholarship); err != nil {
		return nil, NewServiceError("CreateScholarship", "Scholarship", scholarship.ID, "gagal simpan ke database", err)
	}

	scholarship.StudentName = student.Name
	scholarship.ClassName = student.ClassName
	scholarship.NISN = student.NISN
	return scholarship, nil
}

func (s *ScholarshipService) ListScholarships(ctx context.Context, status, search string) ([]repository.StudentScholarship, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	filter := repository.ScholarshipFilter{
		SchoolID: opCtx.SchoolID, AcademicPeriodID: opCtx.AcademicPeriodID,
		Status: status, Search: search, Limit: 100, Offset: 0,
	}

	items, _, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, NewServiceError("ListScholarships", "Scholarship", "", "gagal ambil daftar", err)
	}
	return items, nil
}

func (s *ScholarshipService) GetScholarshipByID(ctx context.Context, scholarshipID string) (*repository.StudentScholarship, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}
	return s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, scholarshipID)
}

func (s *ScholarshipService) UpdateScholarship(ctx context.Context, scholarshipID string, req UpdateScholarshipRequest) (*repository.StudentScholarship, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	existing, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, scholarshipID)
	if err != nil {
		return nil, NewServiceError("UpdateScholarship", "Scholarship", scholarshipID, "data tidak ditemukan", err)
	}

	if req.Category != "" { existing.Category = req.Category }
	if req.CardNumber != "" { existing.CardNumber = req.CardNumber }
	if req.Reason != "" { existing.Reason = req.Reason }
	if req.Status != "" { existing.Status = req.Status }
	if req.BankName != "" { existing.BankName = req.BankName }
	if req.BankAccount != "" { existing.BankAccount = req.BankAccount }
	existing.IsDocumentValid = req.IsDocumentValid
	if req.DisbursedAmount > 0 { existing.DisbursedAmount = req.DisbursedAmount }
	existing.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, NewServiceError("UpdateScholarship", "Scholarship", scholarshipID, "gagal update data", err)
	}
	return existing, nil
}

func (s *ScholarshipService) ToggleDocumentValid(ctx context.Context, scholarshipID string, isValid bool) error {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return ErrMissingOperationalContext
	}

	existing, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, scholarshipID)
	if err != nil {
		return NewServiceError("ToggleDocumentValid", "Scholarship", scholarshipID, "data tidak ditemukan", err)
	}

	existing.IsDocumentValid = isValid
	existing.UpdatedAt = time.Now().UTC()
	return s.repo.Update(ctx, existing)
}
