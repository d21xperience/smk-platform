package search

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// FTSIndexer mendengarkan event perubahan data dan mengupdate search index.
type FTSIndexer struct {
	db *pgxpool.Pool
}

// NewFTSIndexer membuat instance baru dari FTSIndexer.
func NewFTSIndexer(db *pgxpool.Pool) *FTSIndexer {
	return &FTSIndexer{db: db}
}

// studentEvent merepresentasikan payload event StudentCreated/Updated.
type studentEvent struct {
	ID               string `json:"id"`
	SchoolID         string `json:"schoolId"`
	AcademicPeriodID string `json:"academicPeriodId"`
	Name             string `json:"name"`
	ClassName        string `json:"className"`
}

// HandleStudentCreatedOrUpdated dipanggil saat event StudentCreated atau StudentUpdated diterima dari NATS.
func (i *FTSIndexer) HandleStudentCreatedOrUpdated(ctx context.Context, msg []byte) error {
	var event studentEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return fmt.Errorf("gagal parse student event: %w", err)
	}

	// UPSERT ke tabel search index (Idempotent - Prinsip ADR-009)
	query := `
		INSERT INTO students_search_index (id, school_id, academic_period_id, name, class_name, search_vector)
		VALUES ($1, $2, $3, $4, $5, setweight(to_tsvector('indonesian', $4), 'A') || setweight(to_tsvector('indonesian', $5), 'B'))
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			class_name = EXCLUDED.class_name,
			search_vector = setweight(to_tsvector('indonesian', EXCLUDED.name), 'A') || setweight(to_tsvector('indonesian', EXCLUDED.class_name), 'B'),
			updated_at = NOW()
	`

	_, err := i.db.Exec(ctx, query,
		event.ID,
		event.SchoolID,
		event.AcademicPeriodID,
		event.Name,
		event.ClassName,
	)
	if err != nil {
		log.Printf("[FTS Indexer] Gagal update index untuk student %s: %v", event.ID, err)
		return fmt.Errorf("gagal upsert search index: %w", err)
	}

	log.Printf("[FTS Indexer] Search index updated untuk student %s", event.ID)
	return nil
}

// HandleStudentDeleted menghapus student dari search index.
func (i *FTSIndexer) HandleStudentDeleted(ctx context.Context, msg []byte) error {
	var event struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(msg, &event); err != nil {
		return fmt.Errorf("gagal parse student delete event: %w", err)
	}

	query := `DELETE FROM students_search_index WHERE id = $1`
	_, err := i.db.Exec(ctx, query, event.ID)
	if err != nil {
		return fmt.Errorf("gagal delete search index: %w", err)
	}

	log.Printf("[FTS Indexer] Search index deleted untuk student %s", event.ID)
	return nil
}
