package service

import (
	"context"

	"sekolah-platform/services/tu-core/internal/repository"
)

type AlumniTrackingService struct {
	repo repository.AlumniTrackingRepository
}

func NewAlumniTrackingService(repo repository.AlumniTrackingRepository) *AlumniTrackingService {
	return &AlumniTrackingService{repo: repo}
}

type UpdateAlumniTrackingRequest struct {
	Status      string `json:"status"`
	Instansi    string `json:"instansi"`
	PosisiUsaha string `json:"posisiUsaha"`
}

func (s *AlumniTrackingService) List(ctx context.Context, schoolID, academicPeriodID string) ([]repository.AlumniTracking, error) {
	return s.repo.List(ctx, schoolID, academicPeriodID)
}

func (s *AlumniTrackingService) Update(ctx context.Context, id, schoolID, academicPeriodID string, req UpdateAlumniTrackingRequest) error {
	data := &repository.AlumniTracking{
		ID:               id,
		SchoolID:         schoolID,
		AcademicPeriodID: academicPeriodID,
		Status:           req.Status,
		Instansi:         req.Instansi,
		PosisiUsaha:      req.PosisiUsaha,
	}
	return s.repo.Update(ctx, data)
}
