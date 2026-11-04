package service

import (
	"context"
	"fmt"

	"sekolah-platform/services/tu-core/internal/repository"
)

type AcademicService struct {
	repo repository.AcademicRepository
}

func NewAcademicService(repo repository.AcademicRepository) *AcademicService {
	return &AcademicService{repo: repo}
}

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
