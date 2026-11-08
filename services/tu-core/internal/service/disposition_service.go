package service

import (
	"context"
	"errors"

	"sekolah-platform/services/tu-core/internal/domain/disposition/models"
	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
)

type DispositionService struct {
	repo repository.DispositionRepository
}

func NewDispositionService(repo repository.DispositionRepository) *DispositionService {
	return &DispositionService{repo: repo}
}

func (s *DispositionService) CreateDisposition(ctx context.Context, req CreateDispositionRequest) (*DispositionResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	d := &models.Disposition{
		CorrespondenceID: req.CorrespondenceID,
		SchoolID:         opCtx.SchoolID,
		FromUserID:       opCtx.UserID, // Diambil dari JWT/Header
		ToUserID:         req.ToUserID,
		Note:             req.Note,
		Status:           "PENDING",
	}

	if err := s.repo.Create(ctx, d); err != nil {
		return nil, NewServiceError("CreateDisposition", "Disposition", "", "failed to save", err)
	}

	return toDispositionResponse(d), nil
}

func (s *DispositionService) ListByCorrespondenceID(ctx context.Context, correspondenceID string) ([]*DispositionResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	list, err := s.repo.ListByCorrespondenceID(ctx, opCtx.SchoolID, correspondenceID)
	if err != nil {
		return nil, NewServiceError("ListByCorrespondenceID", "Disposition", "", "failed to list", err)
	}

	var res []*DispositionResponse
	for _, d := range list {
		res = append(res, toDispositionResponse(&d))
	}
	return res, nil
}

func (s *DispositionService) GetMyInbox(ctx context.Context) ([]*DispositionResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	list, err := s.repo.GetInboxByUser(ctx, opCtx.SchoolID, opCtx.UserID)
	if err != nil {
		return nil, NewServiceError("GetMyInbox", "Disposition", "", "failed to get inbox", err)
	}

	var res []*DispositionResponse
	for _, d := range list {
		res = append(res, toDispositionResponse(&d))
	}
	return res, nil
}

func (s *DispositionService) UpdateStatus(ctx context.Context, id string, req UpdateDispositionStatusRequest) error {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return ErrMissingOperationalContext
	}

	if req.Status != "READ" && req.Status != "DONE" && req.Status != "PENDING" {
		return errors.New("invalid status")
	}

	return s.repo.UpdateStatus(ctx, opCtx.SchoolID, id, req.Status)
}

func toDispositionResponse(d *models.Disposition) *DispositionResponse {
	return &DispositionResponse{
		ID:               d.ID,
		CorrespondenceID: d.CorrespondenceID,
		FromUserID:       d.FromUserID,
		ToUserID:         d.ToUserID,
		Note:             d.Note,
		Status:           d.Status,
		CreatedAt:        d.CreatedAt,
		UpdatedAt:        d.UpdatedAt,
	}
}
