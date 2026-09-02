package service

import (
	"context"
	"fmt"

	"sekolah-platform/services/tu-core/internal/repository"
)

// AcademicService mengorchestrasi use case untuk domain Academic.
type AcademicService struct {
	repo repository.AcademicRepository
}

// NewAcademicService membuat instance baru.
func NewAcademicService(repo repository.AcademicRepository) *AcademicService {
	return &AcademicService{repo: repo}
}

// GetActivePeriod mengambil periode akademik aktif.
func (s *AcademicService) GetActivePeriod(ctx context.Context, schoolID string) (*repository.AcademicPeriod, error) {
	if schoolID == "" {
		return nil, ErrMissingOperationalContext
	}

	period, err := s.repo.GetActivePeriod(ctx, schoolID)
	if err != nil {
		if isNotFoundError(err) {
			return nil, fmt.Errorf("%w: active period untuk school %s", ErrNotFound, schoolID)
		}
		return nil, NewServiceError("GetActivePeriod", "AcademicPeriod", "", "gagal ambil data", err)
	}

	return period, nil
}

// GetPeriodByID mengambil periode akademik berdasarkan ID.
func (s *AcademicService) GetPeriodByID(ctx context.Context, schoolID, periodID string) (*repository.AcademicPeriod, error) {
	if schoolID == "" {
		return nil, ErrMissingOperationalContext
	}

	period, err := s.repo.GetPeriodByID(ctx, schoolID, periodID)
	if err != nil {
		if isNotFoundError(err) {
			return nil, fmt.Errorf("%w: period %s", ErrNotFound, periodID)
		}
		return nil, NewServiceError("GetPeriodByID", "AcademicPeriod", periodID, "gagal ambil data", err)
	}

	return period, nil
}

// GetClassByID mengambil kelas berdasarkan ID.
func (s *AcademicService) GetClassByID(ctx context.Context, schoolID, academicPeriodID, classID string) (*repository.Class, error) {
	if schoolID == "" || academicPeriodID == "" {
		return nil, ErrMissingOperationalContext
	}

	class, err := s.repo.GetClassByID(ctx, schoolID, academicPeriodID, classID)
	if err != nil {
		if isNotFoundError(err) {
			return nil, fmt.Errorf("%w: class %s", ErrNotFound, classID)
		}
		return nil, NewServiceError("GetClassByID", "Class", classID, "gagal ambil data", err)
	}

	return class, nil
}

// ListClasses mengambil daftar kelas.
func (s *AcademicService) ListClasses(ctx context.Context, schoolID, academicPeriodID string) ([]repository.Class, error) {
	if schoolID == "" || academicPeriodID == "" {
		return nil, ErrMissingOperationalContext
	}

	classes, err := s.repo.ListClasses(ctx, schoolID, academicPeriodID)
	if err != nil {
		return nil, NewServiceError("ListClasses", "Class", "", "gagal ambil daftar", err)
	}

	return classes, nil
}

// ListSubjects mengambil daftar mata pelajaran.
func (s *AcademicService) ListSubjects(ctx context.Context, schoolID string) ([]repository.Subject, error) {
	if schoolID == "" {
		return nil, ErrMissingOperationalContext
	}

	subjects, err := s.repo.ListSubjects(ctx, schoolID)
	if err != nil {
		return nil, NewServiceError("ListSubjects", "Subject", "", "gagal ambil daftar", err)
	}

	return subjects, nil
}
