package service

import (
	"context"

	"sekolah-platform/services/tu-core/internal/repository"
)

type PtkArchiveService struct {
	repo repository.PtkArchiveRepository
}

func NewPtkArchiveService(repo repository.PtkArchiveRepository) *PtkArchiveService {
	return &PtkArchiveService{repo: repo}
}

func (s *PtkArchiveService) ListArchives(ctx context.Context) ([]repository.Archive, error) {
	return s.repo.ListArchives(ctx)
}

func (s *PtkArchiveService) CreateArchive(ctx context.Context, req repository.Archive) error {
	return s.repo.CreateArchive(ctx, &req)
}

func (s *PtkArchiveService) UpdateArchive(ctx context.Context, req repository.Archive) error {
	return s.repo.UpdateArchive(ctx, &req)
}

func (s *PtkArchiveService) DeleteArchive(ctx context.Context, id string) error {
	return s.repo.DeleteArchive(ctx, id)
}
