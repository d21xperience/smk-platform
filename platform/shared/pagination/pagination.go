// Package pagination menyediakan helper untuk pagination query.
package pagination

// Params adalah parameter pagination dari request.
type Params struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// Meta adalah metadata pagination untuk response.
type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
	HasMore    bool  `json:"hasMore"`
}

// DefaultLimit adalah limit default jika tidak diset.
const DefaultLimit = 20

// MaxLimit adalah limit maksimum yang diizinkan.
const MaxLimit = 100

// Normalize menormalisasi parameter pagination.
func Normalize(page, limit int) Params {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	return Params{Page: page, Limit: limit}
}

// Offset menghitung offset untuk SQL query.
func (p Params) Offset() int {
	return (p.Page - 1) * p.Limit
}

// BuildMeta membangun Meta dari total records.
func BuildMeta(params Params, total int64) Meta {
	totalPages := int(total) / params.Limit
	if int(total)%params.Limit > 0 {
		totalPages++
	}
	if totalPages == 0 {
		totalPages = 1
	}

	return Meta{
		Page:       params.Page,
		Limit:      params.Limit,
		Total:      total,
		TotalPages: totalPages,
		HasMore:    params.Page < totalPages,
	}
}