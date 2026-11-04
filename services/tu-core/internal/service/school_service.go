package service

import (
	"context"

	"sekolah-platform/services/tu-core/internal/repository"
)

type SchoolService struct {
	schoolRepo repository.SchoolRepository
}

func NewSchoolService(schoolRepo repository.SchoolRepository) *SchoolService {
	return &SchoolService{schoolRepo: schoolRepo}
}

func (s *SchoolService) GetSchool(ctx context.Context, schoolID string) (*repository.School, error) {
	return s.schoolRepo.GetByID(ctx, schoolID)
}

func (s *SchoolService) UpdateSchool(ctx context.Context, school *repository.School) error {
	return s.schoolRepo.Update(ctx, school)
}
