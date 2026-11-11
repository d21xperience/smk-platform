package service

import (
	"context"

	"sekolah-platform/services/tu-core/internal/domain/correspondence/models"
	"sekolah-platform/services/tu-core/internal/repository"
)

type LetterTypeService struct {
	repo repository.LetterTypeRepository
}

func NewLetterTypeService(repo repository.LetterTypeRepository) *LetterTypeService {
	return &LetterTypeService{repo: repo}
}

func (s *LetterTypeService) CreateLetterType(ctx context.Context, req LetterTypeRequest) (*LetterTypeResponse, error) {
	lt := &models.LetterType{
		Kategori:        req.Kategori,
		Peruntukan:      req.Peruntukan,
		Value:           req.Value,
		Label:           req.Label,
		KodeKlasifikasi: req.KodeKlasifikasi,
	}
	if err := s.repo.Create(ctx, lt); err != nil {
		return nil, NewServiceError("CreateLetterType", "LetterType", "", "failed to save", err)
	}
	return toLetterTypeResponse(lt), nil
}

func (s *LetterTypeService) GetLetterType(ctx context.Context, id string) (*LetterTypeResponse, error) {
	lt, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, NewServiceError("GetLetterType", "LetterType", id, "failed to get", err)
	}
	return toLetterTypeResponse(lt), nil
}

func (s *LetterTypeService) ListLetterTypes(ctx context.Context) ([]*LetterTypeResponse, error) {
	list, err := s.repo.List(ctx)
	if err != nil {
		return nil, NewServiceError("ListLetterTypes", "LetterType", "", "failed to list", err)
	}
	var res []*LetterTypeResponse
	for _, lt := range list {
		res = append(res, toLetterTypeResponse(&lt))
	}
	return res, nil
}

func (s *LetterTypeService) UpdateLetterType(ctx context.Context, id string, req LetterTypeRequest) (*LetterTypeResponse, error) {
	lt, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, NewServiceError("UpdateLetterType", "LetterType", id, "failed to get", err)
	}

	lt.Kategori = req.Kategori
	lt.Peruntukan = req.Peruntukan
	lt.Value = req.Value
	lt.Label = req.Label
	lt.KodeKlasifikasi = req.KodeKlasifikasi

	if err := s.repo.Update(ctx, lt); err != nil {
		return nil, NewServiceError("UpdateLetterType", "LetterType", id, "failed to update", err)
	}
	return toLetterTypeResponse(lt), nil
}

func (s *LetterTypeService) DeleteLetterType(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return NewServiceError("DeleteLetterType", "LetterType", id, "failed to delete", err)
	}
	return nil
}

func toLetterTypeResponse(lt *models.LetterType) *LetterTypeResponse {
	return &LetterTypeResponse{
		ID:              lt.ID,
		Kategori:        lt.Kategori,
		Peruntukan:      lt.Peruntukan,
		Value:           lt.Value,
		Label:           lt.Label,
		KodeKlasifikasi: lt.KodeKlasifikasi,
		CreatedAt:       lt.CreatedAt,
	}
}
