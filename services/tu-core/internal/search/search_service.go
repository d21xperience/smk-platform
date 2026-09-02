package search

import (
	"context"
	"fmt"
)

// SearchQuery merepresentasikan parameter pencarian.
type SearchQuery struct {
	Keyword          string
	SchoolID         string
	AcademicPeriodID string
	Entity           string // "student", "teacher", "document"
	Limit            int
	Offset           int
}

// SearchResult adalah hasil pencarian yang sudah di-rank.
type SearchResult struct {
	ID       string  `json:"id"`
	Entity   string  `json:"entity"`
	Title    string  `json:"title"`
	Subtitle string  `json:"subtitle"`
	Score    float64 `json:"score"`
}

// SearchRepository adalah kontrak untuk lapisan data (Prinsip 6: Pure logic di engine, repo handle DB).
type SearchRepository interface {
	FullTextSearch(ctx context.Context, query SearchQuery) ([]SearchResult, error)
}

// SearchService adalah orchestrator untuk pencarian (Prinsip 11).
type SearchService struct {
	repo SearchRepository
}

// NewSearchService membuat instance baru dari SearchService.
func NewSearchService(repo SearchRepository) *SearchService {
	return &SearchService{repo: repo}
}

// Search melakukan pencarian full-text dengan validasi operational context.
func (s *SearchService) Search(ctx context.Context, query SearchQuery) ([]SearchResult, error) {
	// Validasi Operational Context (Prinsip 5)
	if query.SchoolID == "" || query.AcademicPeriodID == "" {
		return nil, ErrMissingOperationalContext
	}

	if query.Keyword == "" {
		return nil, ErrInvalidQuery
	}

	// Set default limit dan offset jika tidak diisi
	if query.Limit <= 0 {
		query.Limit = 20
	}
	if query.Offset < 0 {
		query.Offset = 0
	}

	// Delegasikan ke repository (Prinsip 11: Service = Orchestrator)
	results, err := s.repo.FullTextSearch(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSearchFailed, err)
	}

	return results, nil
}
