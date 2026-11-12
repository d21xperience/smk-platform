package service

import (
	"context"

	"sekolah-platform/services/tu-core/internal/repository"

	"github.com/google/uuid"
)

type PklPlottingService struct {
	repo repository.PklPlottingRepository
}

func NewPklPlottingService(repo repository.PklPlottingRepository) *PklPlottingService {
	return &PklPlottingService{repo: repo}
}

type PklPlottingRequest struct {
	Nama       string `json:"nama"`
	Kelas      string `json:"kelas"`
	Jurusan    string `json:"jurusan"`
	Perusahaan string `json:"perusahaan"`
	Kota       string `json:"kota"`
	Pembimbing string `json:"pembimbing"`
}

func (s *PklPlottingService) List(ctx context.Context, schoolID, academicPeriodID string) ([]repository.PklPlotting, error) {
	return s.repo.List(ctx, schoolID, academicPeriodID)
}

func (s *PklPlottingService) Create(ctx context.Context, schoolID, academicPeriodID string, req PklPlottingRequest) error {
	data := &repository.PklPlotting{
		ID:               uuid.NewString(),
		SchoolID:         schoolID,
		AcademicPeriodID: academicPeriodID,
		Nama:             req.Nama,
		Kelas:            req.Kelas,
		Jurusan:          req.Jurusan,
		Perusahaan:       req.Perusahaan,
		Kota:             req.Kota,
		Pembimbing:       req.Pembimbing,
	}
	return s.repo.Create(ctx, data)
}

func (s *PklPlottingService) Update(ctx context.Context, id, schoolID, academicPeriodID string, req PklPlottingRequest) error {
	data := &repository.PklPlotting{
		ID:               id,
		SchoolID:         schoolID,
		AcademicPeriodID: academicPeriodID,
		Perusahaan:       req.Perusahaan,
		Kota:             req.Kota,
		Pembimbing:       req.Pembimbing,
	}
	return s.repo.Update(ctx, data)
}

func (s *PklPlottingService) Delete(ctx context.Context, id, schoolID, academicPeriodID string) error {
	return s.repo.Delete(ctx, id, schoolID, academicPeriodID)
}
