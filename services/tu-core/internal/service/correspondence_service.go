package service

import (
	"context"

	"sekolah-platform/services/tu-core/internal/domain/correspondence/engine"
	"sekolah-platform/services/tu-core/internal/domain/correspondence/models"
	"sekolah-platform/services/tu-core/internal/domain/correspondence/valueobjects"
	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"

	"github.com/google/uuid"
)

type CorrespondenceService struct {
	repo       repository.CorrespondenceRepository
	corrEngine *engine.CorrespondenceEngine
	publisher  EventPublisher
}

func NewCorrespondenceService(
	repo repository.CorrespondenceRepository,
	corrEngine *engine.CorrespondenceEngine,
	publisher EventPublisher,
) *CorrespondenceService {
	if publisher == nil {
		publisher = NewNoopEventPublisher()
	}
	return &CorrespondenceService{
		repo:       repo,
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

	data := models.CorrespondenceData{
		Type:          corrType,
		Number:        req.Number,
		Date:          req.Date,
		Subject:       req.Subject,
		FromParty:     req.From,
		ToParty:       req.To,
		Description:   req.Description,
		AttachmentURL: req.AttachmentURL,
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

func (s *CorrespondenceService) UpdateCorrespondence(ctx context.Context, id, subject, description, attachmentURL string) (*CorrespondenceResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	c, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, id)
	if err != nil {
		return nil, NewServiceError("UpdateCorrespondence", "Correspondence", id, "failed to get correspondence", err)
	}

	if err := c.Update(subject, description, attachmentURL); err != nil {
		return nil, NewServiceError("UpdateCorrespondence", "Correspondence", id, "failed to update", err)
	}

	if err := s.repo.Update(ctx, c); err != nil {
		return nil, NewServiceError("UpdateCorrespondence", "Correspondence", id, "failed to save to database", err)
	}

	return toCorrespondenceResponse(c), nil
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

func (s *CorrespondenceService) ListCorrespondences(ctx context.Context, corrType, status, search string) ([]*CorrespondenceResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	filter := repository.CorrespondenceFilter{
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		Type:             corrType,
		Status:           status,
		Search:           search,
		Limit:            100,
		Offset:           0,
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

func toCorrespondenceResponse(c *models.Correspondence) *CorrespondenceResponse {
	if c == nil {
		return nil
	}
	return &CorrespondenceResponse{
		ID:            c.ID,
		Type:          c.Type.String(),
		Number:        c.Number,
		Date:          c.Date,
		Subject:       c.Subject,
		From:          c.FromParty,
		To:            c.ToParty,
		Description:   c.Description,
		AttachmentURL: c.AttachmentURL,
		Status:        c.Status.String(),
		CreatedAt:     c.CreatedAt,
	}
}
