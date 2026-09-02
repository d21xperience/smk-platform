package pagination

import "testing"

func TestNormalize(t *testing.T) {
	tests := []struct {
		name           string
		page, limit    int
		expectedPage   int
		expectedLimit  int
	}{
		{"valid", 2, 10, 2, 10},
		{"page < 1", 0, 10, 1, 10},
		{"limit < 1", 1, 0, 1, DefaultLimit},
		{"limit > max", 1, 500, 1, MaxLimit},
		{"both zero", 0, 0, 1, DefaultLimit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Normalize(tt.page, tt.limit)
			if p.Page != tt.expectedPage {
				t.Errorf("Expected page %d, got %d", tt.expectedPage, p.Page)
			}
			if p.Limit != tt.expectedLimit {
				t.Errorf("Expected limit %d, got %d", tt.expectedLimit, p.Limit)
			}
		})
	}
}

func TestOffset(t *testing.T) {
	p := Params{Page: 3, Limit: 20}
	if p.Offset() != 40 {
		t.Errorf("Expected offset 40, got %d", p.Offset())
	}
}

func TestBuildMeta(t *testing.T) {
	params := Params{Page: 1, Limit: 20}
	meta := BuildMeta(params, 50)

	if meta.Total != 50 {
		t.Errorf("Expected total 50, got %d", meta.Total)
	}
	if meta.TotalPages != 3 {
		t.Errorf("Expected totalPages 3, got %d", meta.TotalPages)
	}
	if !meta.HasMore {
		t.Error("Expected hasMore=true")
	}
}

func TestBuildMetaLastPage(t *testing.T) {
	params := Params{Page: 3, Limit: 20}
	meta := BuildMeta(params, 50)

	if meta.HasMore {
		t.Error("Expected hasMore=false on last page")
	}
}

func TestBuildMetaEmpty(t *testing.T) {
	params := Params{Page: 1, Limit: 20}
	meta := BuildMeta(params, 0)

	if meta.TotalPages != 1 {
		t.Errorf("Expected totalPages 1 for empty, got %d", meta.TotalPages)
	}
	if meta.HasMore {
		t.Error("Expected hasMore=false for empty")
	}
}