package search

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// pgxSearchRepository mengimplementasikan SearchRepository menggunakan PostgreSQL Full Text Search.
type pgxSearchRepository struct {
	db *pgxpool.Pool
}

// NewPostgresSearchRepository membuat instance baru dari repository.
func NewPostgresSearchRepository(db *pgxpool.Pool) SearchRepository {
	return &pgxSearchRepository{db: db}
}

// FullTextSearch melakukan pencarian menggunakan tsvector dan tsquery.
func (r *pgxSearchRepository) FullTextSearch(ctx context.Context, query SearchQuery) ([]SearchResult, error) {
	// Tentukan tabel dan kolom berdasarkan entity
	var tableName, titleColumn, subtitleColumn string
	switch query.Entity {
	case "student":
		tableName = "students_search_index"
		titleColumn = "name"
		subtitleColumn = "class_name"
	case "teacher":
		tableName = "teachers_search_index"
		titleColumn = "name"
		subtitleColumn = "subject"
	case "document":
		tableName = "documents_search_index"
		titleColumn = "title"
		subtitleColumn = "document_number"
	default:
		// Default ke student jika entity tidak dikenali
		tableName = "students_search_index"
		titleColumn = "name"
		subtitleColumn = "class_name"
	}

	// Query dengan Full Text Search dan ranking
	sqlQuery := fmt.Sprintf(`
		SELECT
			id,
			%s as title,
			%s as subtitle,
			ts_rank(search_vector, websearch_to_tsquery('indonesian', $1)) as score
		FROM %s
		WHERE school_id = $2
		  AND academic_period_id = $3
		  AND search_vector @@ websearch_to_tsquery('indonesian', $1)
		ORDER BY score DESC
		LIMIT $4 OFFSET $5
	`, titleColumn, subtitleColumn, tableName)

	rows, err := r.db.Query(ctx, sqlQuery,
		query.Keyword,
		query.SchoolID,
		query.AcademicPeriodID,
		query.Limit,
		query.Offset,
	)
	if err != nil {
		return nil, fmt.Errorf("query search gagal: %w", err)
	}
	defer rows.Close()

	var results []SearchResult
	for rows.Next() {
		var res SearchResult
		if err := rows.Scan(&res.ID, &res.Title, &res.Subtitle, &res.Score); err != nil {
			return nil, fmt.Errorf("gagal scan search result: %w", err)
		}
		res.Entity = query.Entity
		results = append(results, res)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating search results: %w", err)
	}

	return results, nil
}
