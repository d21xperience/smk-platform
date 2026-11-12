package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
)

type MutationService struct {
	repo        repository.MutationRepository
	studentRepo repository.StudentRepository
	publisher   EventPublisher
}

func NewMutationService(
	repo repository.MutationRepository,
	studentRepo repository.StudentRepository,
	publisher EventPublisher,
) *MutationService {
	if publisher == nil {
		publisher = NewNoopEventPublisher()
	}
	return &MutationService{
		repo:        repo,
		studentRepo: studentRepo,
		publisher:   publisher,
	}
}

type CreateMutationRequest struct {
	StudentID     *string
	MutationType  string // "MASUK" atau "KELUAR"
	NIS           string
	NISN          string
	Nama          string
	JenisKelamin  string
	TempatLahir   string
	TanggalLahir  *time.Time
	AsalSekolah   string
	TujuanSekolah string
	Alasan        string
	DokumenURL    string
}

func (s *MutationService) CreateMutation(ctx context.Context, req CreateMutationRequest) (*repository.StudentMutation, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	m := &repository.StudentMutation{
		ID:               uuid.New().String(),
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		StudentID:        req.StudentID,
		MutationType:     req.MutationType,
		NIS:              req.NIS,
		NISN:             req.NISN,
		Nama:             req.Nama,
		JenisKelamin:     req.JenisKelamin,
		TempatLahir:      req.TempatLahir,
		TanggalLahir:     req.TanggalLahir,
		AsalSekolah:      req.AsalSekolah,
		TujuanSekolah:    req.TujuanSekolah,
		Alasan:           req.Alasan,
		DokumenURL:       req.DokumenURL,
		Status:           "PENDING",
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, m); err != nil {
		return nil, NewServiceError("CreateMutation", "Mutation", m.ID, "gagal simpan ke database", err)
	}

	return m, nil
}

func (s *MutationService) ApproveMutation(ctx context.Context, mutationID, catatan string) error {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return ErrMissingOperationalContext
	}

	m, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, mutationID)
	if err != nil {
		return NewServiceError("ApproveMutation", "Mutation", mutationID, "mutasi tidak ditemukan", err)
	}

	if m.Status != "PENDING" {
		return NewServiceError("ApproveMutation", "Mutation", mutationID, "hanya mutasi PENDING yang bisa disetujui", nil)
	}

	if m.MutationType == "KELUAR" && m.StudentID != nil {
		student, err := s.studentRepo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, *m.StudentID)
		if err != nil {
			return NewServiceError("ApproveMutation", "Student", *m.StudentID, "siswa tidak ditemukan", err)
		}

		student.Status = "TRANSFERRED"
		if err := s.studentRepo.Update(ctx, student); err != nil {
			return NewServiceError("ApproveMutation", "Student", *m.StudentID, "gagal update status siswa ke TRANSFERRED", err)
		}

		_ = s.studentRepo.DeleteForUI(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, *m.StudentID)
	} else if m.MutationType == "MASUK" {
		fmt.Printf("[MutationService] ⚠️ Mutasi MASUK disetujui, namun siswa belum terdaftar di tabel_siswa.\n")
	}

	if err := s.repo.UpdateStatus(ctx, mutationID, "APPROVED", catatan); err != nil {
		return NewServiceError("ApproveMutation", "Mutation", mutationID, "gagal update status mutasi", err)
	}

	return nil
}

func (s *MutationService) RejectMutation(ctx context.Context, mutationID, catatan string) error {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return ErrMissingOperationalContext
	}

	_, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, mutationID)
	if err != nil {
		return NewServiceError("RejectMutation", "Mutation", mutationID, "mutasi tidak ditemukan", err)
	}

	if err := s.repo.UpdateStatus(ctx, mutationID, "REJECTED", catatan); err != nil {
		return NewServiceError("RejectMutation", "Mutation", mutationID, "gagal update status mutasi", err)
	}

	return nil
}

func (s *MutationService) ListMutations(ctx context.Context, status, mType string) ([]repository.StudentMutation, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	filter := repository.MutationFilter{
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		Status:           status,
		Type:             mType,
		Limit:            100,
		Offset:           0,
	}

	mutations, _, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, NewServiceError("ListMutations", "Mutation", "", "gagal ambil daftar mutasi", err)
	}
	return mutations, nil
}

func (s *MutationService) GetMutationHistory(ctx context.Context, studentID string) ([]repository.StudentMutation, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}
	return s.repo.GetHistoryByStudent(ctx, opCtx.SchoolID, studentID)
}

// GetMutationByID (Tambahan untuk gRPC)
func (s *MutationService) GetMutationByID(ctx context.Context, mutationID string) (*repository.StudentMutation, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}
	return s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, mutationID)
}
