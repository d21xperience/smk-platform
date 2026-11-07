package service

import (
	"context"

	"sekolah-platform/services/tu-core/internal/repository"
)

type PtkLeaveService struct {
	repo repository.PtkLeaveRepository
}

func NewPtkLeaveService(repo repository.PtkLeaveRepository) *PtkLeaveService {
	return &PtkLeaveService{repo: repo}
}

func (s *PtkLeaveService) ListLeaves(ctx context.Context) ([]repository.LeaveRequest, error) {
	return s.repo.ListLeaves(ctx)
}

func (s *PtkLeaveService) CreateLeave(ctx context.Context, req repository.LeaveRequest) error {
	return s.repo.CreateLeave(ctx, &req)
}

func (s *PtkLeaveService) UpdateStatus(ctx context.Context, id string, status string, note string) error {
	return s.repo.UpdateStatus(ctx, id, status, note)
}

func (s *PtkLeaveService) DeleteLeave(ctx context.Context, id string) error {
	return s.repo.DeleteLeave(ctx, id)
}
