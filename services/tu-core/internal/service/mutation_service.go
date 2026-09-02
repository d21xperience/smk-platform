package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"sekolah-platform/services/tu-core/internal/domain/mutation/engine"
	"sekolah-platform/services/tu-core/internal/domain/mutation/models"
	"sekolah-platform/services/tu-core/internal/domain/mutation/valueobjects"
	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
)

type MutationService struct {
	repo        repository.MutationRepository
	studentRepo repository.StudentRepository // <-- TAMBAHKAN INI
	mutEngine   *engine.MutationEngine
	publisher   EventPublisher
}

func NewMutationService(
	repo repository.MutationRepository,
	studentRepo repository.StudentRepository, // <-- TAMBAHKAN INI
	mutEngine *engine.MutationEngine,
	publisher EventPublisher,
) *MutationService {
	if publisher == nil {
		publisher = NewNoopEventPublisher()
	}
	return &MutationService{
		repo:        repo,
		studentRepo: studentRepo,
		mutEngine:   mutEngine,
		publisher:   publisher,
	}
}

type CreateMutationRequest struct {
	StudentID     *string
	MutationType  valueobjects.MutationType
	CalonNama     string
	CalonNISN     string
	AsalSekolah   string
	TujuanSekolah string
	Alasan        string
	DokumenURL    string
}

func (s *MutationService) CreateMutation(ctx context.Context, req CreateMutationRequest) (*MutationResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	data := models.MutationData{
		StudentID:        req.StudentID,
		MutationType:     req.MutationType,
		CalonNama:        req.CalonNama,
		CalonNISN:        req.CalonNISN,
		AsalSekolah:      req.AsalSekolah,
		TujuanSekolah:    req.TujuanSekolah,
		Alasan:           req.Alasan,
		DokumenURL:       req.DokumenURL,
	}

	result := s.mutEngine.CreateMutation(data, opCtx)
	if result.Error != nil {
		return nil, NewServiceError("CreateMutation", "Mutation", "", "invalid mutation data", result.Error)
	}

	if err := s.repo.Create(ctx, result.Mutation); err != nil {
		return nil, NewServiceError("CreateMutation", "Mutation", result.Mutation.ID, "failed to save to database", err)
	}

	event := DomainEvent{
		EventType:        "MutationCreated",
		AggregateType:    "Mutation",
		AggregateID:      result.Mutation.ID,
		Payload: map[string]interface{}{
			"mutationId":   result.Mutation.ID,
			"mutationType": result.Mutation.MutationType.String(),
		},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return toMutationResponse(result.Mutation), nil
}

func (s *MutationService) ApproveMutation(ctx context.Context, mutationID, catatan string) (*MutationResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	m, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, mutationID)
	if err != nil {
		return nil, NewServiceError("ApproveMutation", "Mutation", mutationID, "failed to get mutation", err)
	}

	result := s.mutEngine.ApproveMutation(m, catatan, opCtx)
	if result.Error != nil {
		return nil, NewServiceError("ApproveMutation", "Mutation", mutationID, "failed to approve", result.Error)
	}

	if err := s.repo.Update(ctx, result.Mutation); err != nil {
		return nil, NewServiceError("ApproveMutation", "Mutation", mutationID, "failed to update database", err)
	}

	// ============================================================
	// ✅ SIDE-EFFECT: Jika Mutasi KELUAR disetujui, ubah status siswa
	// ============================================================
	if result.Mutation.MutationType == valueobjects.TypeKeluar && result.Mutation.StudentID != nil {
		student, err := s.studentRepo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, *result.Mutation.StudentID)
		if err == nil && student != nil {
			student.Status = "TRANSFERRED" // Ubah status menjadi TRANSFERRED
			if err := s.studentRepo.Update(ctx, student); err != nil {
				return nil, NewServiceError("ApproveMutation", "Mutation", mutationID, "failed to update student status to TRANSFERRED", err)
			}
		}
	}

	event := DomainEvent{
		EventType:        "MutationApproved",
		AggregateType:    "Mutation",
		AggregateID:      result.Mutation.ID,
		Payload: map[string]interface{}{
			"mutationId":   result.Mutation.ID,
			"mutationType": result.Mutation.MutationType.String(),
			"status":       result.Mutation.Status.String(),
			"catatan":      catatan,
		},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return toMutationResponse(result.Mutation), nil
}

func (s *MutationService) RejectMutation(ctx context.Context, mutationID, catatan string) (*MutationResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	m, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, mutationID)
	if err != nil {
		return nil, NewServiceError("RejectMutation", "Mutation", mutationID, "failed to get mutation", err)
	}

	result := s.mutEngine.RejectMutation(m, catatan, opCtx)
	if result.Error != nil {
		return nil, NewServiceError("RejectMutation", "Mutation", mutationID, "failed to reject", result.Error)
	}

	if err := s.repo.Update(ctx, result.Mutation); err != nil {
		return nil, NewServiceError("RejectMutation", "Mutation", mutationID, "failed to update database", err)
	}

	event := DomainEvent{
		EventType:        "MutationRejected",
		AggregateType:    "Mutation",
		AggregateID:      result.Mutation.ID,
		Payload: map[string]interface{}{
			"mutationId": result.Mutation.ID,
			"status":     result.Mutation.Status.String(),
			"catatan":    catatan,
		},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return toMutationResponse(result.Mutation), nil
}

func (s *MutationService) ListMutations(ctx context.Context, status, mType string) ([]*MutationResponse, error) {
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
		return nil, NewServiceError("ListMutations", "Mutation", "", "failed to list mutations", err)
	}

	var responses []*MutationResponse
	for _, m := range mutations {
		responses = append(responses, toMutationResponse(&m))
	}

	return responses, nil
}

type MutationResponse struct {
	ID                string
	SchoolID          string
	AcademicPeriodID  string
	StudentID         *string
	MutationType      string
	CalonNama         string
	CalonNISN         string
	AsalSekolah       string
	TujuanSekolah     string
	Alasan            string
	DokumenURL        string
	Status            string
	CatatanVerifikasi string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func toMutationResponse(m *models.Mutation) *MutationResponse {
	if m == nil {
		return nil
	}
	return &MutationResponse{
		ID:                m.ID,
		SchoolID:          m.SchoolID,
		AcademicPeriodID:  m.AcademicPeriodID,
		StudentID:         m.StudentID,
		MutationType:      m.MutationType.String(),
		CalonNama:         m.CalonNama,
		CalonNISN:         m.CalonNISN,
		AsalSekolah:       m.AsalSekolah,
		TujuanSekolah:     m.TujuanSekolah,
		Alasan:            m.Alasan,
		DokumenURL:        m.DokumenURL,
		Status:            m.Status.String(),
		CatatanVerifikasi: m.CatatanVerifikasi,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}
