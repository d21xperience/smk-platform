package service

import (
	"context"
	"fmt"

	"sekolah-platform/services/tu-core/internal/repository"
)

type TeacherHistoryService struct {
	repo repository.TeacherHistoryRepository
}

func NewTeacherHistoryService(repo repository.TeacherHistoryRepository) *TeacherHistoryService {
	return &TeacherHistoryService{repo: repo}
}

func (s *TeacherHistoryService) ListEducations(ctx context.Context, ptkID string) ([]repository.EducationHistory, error) {
	return s.repo.ListEducations(ctx, ptkID)
}

func (s *TeacherHistoryService) AddEducation(ctx context.Context, ptkID string, edu repository.EducationHistory) error {
	if ptkID == "" { return fmt.Errorf("%w: ID PTK wajib diisi", ErrInvalidInput) }
	edu.PtkID = ptkID
	return s.repo.CreateEducation(ctx, &edu)
}

func (s *TeacherHistoryService) DeleteEducation(ctx context.Context, eduID string) error {
	return s.repo.DeleteEducation(ctx, eduID)
}

func (s *TeacherHistoryService) ListCertifications(ctx context.Context, ptkID string) ([]repository.CertificationHistory, error) {
	return s.repo.ListCertifications(ctx, ptkID)
}

func (s *TeacherHistoryService) AddCertification(ctx context.Context, ptkID string, cert repository.CertificationHistory) error {
	if ptkID == "" { return fmt.Errorf("%w: ID PTK wajib diisi", ErrInvalidInput) }
	cert.PtkID = ptkID
	return s.repo.CreateCertification(ctx, &cert)
}

func (s *TeacherHistoryService) DeleteCertification(ctx context.Context, certID string) error {
	return s.repo.DeleteCertification(ctx, certID)
}
