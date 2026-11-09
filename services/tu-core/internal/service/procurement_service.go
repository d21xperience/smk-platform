package service

import (
	"context"
	"errors"

	"sekolah-platform/services/tu-core/internal/repository"
)

type ProcurementService struct {
	repo repository.ProcurementRepository
}

func NewProcurementService(repo repository.ProcurementRepository) *ProcurementService {
	return &ProcurementService{repo: repo}
}

func (s *ProcurementService) CreateRequest(ctx context.Context, req *repository.ProcurementRequest) (*repository.ProcurementRequest, error) {
	if req.ItemName == "" || req.Quantity <= 0 || req.Unit == "" || req.Justification == "" {
		return nil, errors.New("data pengajuan tidak lengkap")
	}

	reqNo, err := s.repo.GenerateRequestNumber(ctx, req.SemesterID)
	if err != nil {
		return nil, err
	}
	req.RequestNumber = reqNo
	req.Status = "PENDING"
	req.TotalEstimated = float64(req.Quantity) * req.EstimatedPrice

	if err := s.repo.Create(ctx, req); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *ProcurementService) GetAllRequests(ctx context.Context, semesterID string) ([]*repository.ProcurementRequest, error) {
	return s.repo.GetAll(ctx, semesterID)
}

func (s *ProcurementService) UpdateRequest(ctx context.Context, req *repository.ProcurementRequest) error {
	existing, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if existing.Status != "PENDING" {
		return errors.New("hanya pengajuan PENDING yang dapat diedit")
	}
	req.TotalEstimated = float64(req.Quantity) * req.EstimatedPrice
	return s.repo.Update(ctx, req)
}

func (s *ProcurementService) DeleteRequest(ctx context.Context, id string) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing.Status != "PENDING" {
		return errors.New("hanya pengajuan PENDING yang dapat dihapus")
	}
	return s.repo.Delete(ctx, id)
}

func (s *ProcurementService) ApproveRequest(ctx context.Context, id string, approverID string, notes string) error {
	req, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if req.Status != "PENDING" {
		return errors.New("hanya pengajuan PENDING yang dapat disetujui")
	}
	return s.repo.UpdateStatus(ctx, id, "APPROVED", approverID, nil, &notes)
}

func (s *ProcurementService) RejectRequest(ctx context.Context, id string, approverID string, reason string, notes string) error {
	if reason == "" {
		return errors.New("alasan penolakan wajib diisi")
	}
	req, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if req.Status != "PENDING" {
		return errors.New("hanya pengajuan PENDING yang dapat ditolak")
	}
	return s.repo.UpdateStatus(ctx, id, "REJECTED", approverID, &reason, &notes)
}

func (s *ProcurementService) MarkAsProcured(ctx context.Context, id string) error {
	req, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if req.Status != "APPROVED" {
		return errors.New("hanya pengajuan APPROVED yang dapat ditandai PROCURED")
	}
	return s.repo.UpdateStatus(ctx, id, "PROCURED", req.RequesterID, nil, nil)
}
