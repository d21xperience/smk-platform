package service

import (
	"context"

	"sekolah-platform/services/tu-core/internal/repository"

	"github.com/google/uuid"
)

type AlumniLegalisirService struct {
	repo repository.AlumniLegalisirRepository
}

func NewAlumniLegalisirService(repo repository.AlumniLegalisirRepository) *AlumniLegalisirService {
	return &AlumniLegalisirService{repo: repo}
}

type CreateAlumniLegalisirRequest struct {
	Nama         string `json:"nama"`
	TahunLulus   string `json:"tahunLulus"`
	Jurusan      string `json:"jurusan"`
	NoIjazah     string `json:"noIjazah"`
	JumlahLembar int    `json:"jumlahLembar"`
}

func (s *AlumniLegalisirService) List(ctx context.Context, schoolID, academicPeriodID string) ([]repository.AlumniLegalisir, error) {
	return s.repo.List(ctx, schoolID, academicPeriodID)
}

func (s *AlumniLegalisirService) Create(ctx context.Context, schoolID, academicPeriodID string, req CreateAlumniLegalisirRequest) error {
	data := &repository.AlumniLegalisir{
		ID:               uuid.NewString(),
		SchoolID:         schoolID,
		AcademicPeriodID: academicPeriodID,
		Nama:             req.Nama,
		TahunLulus:       req.TahunLulus,
		Jurusan:          req.Jurusan,
		NoIjazah:         req.NoIjazah,
		JumlahLembar:     req.JumlahLembar,
		Status:           "PROSES TTD",
	}
	return s.repo.Create(ctx, data)
}

func (s *AlumniLegalisirService) UpdateStatus(ctx context.Context, id, schoolID, academicPeriodID, status string) error {
	return s.repo.UpdateStatus(ctx, id, schoolID, academicPeriodID, status)
}

func (s *AlumniLegalisirService) Delete(ctx context.Context, id, schoolID, academicPeriodID string) error {
	return s.repo.Delete(ctx, id, schoolID, academicPeriodID)
}
