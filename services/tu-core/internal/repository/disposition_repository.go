package repository

import (
	"context"
	"fmt"
	"time"

	"sekolah-platform/services/tu-core/internal/domain/disposition/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DispositionRepository interface {
	Create(ctx context.Context, d *models.Disposition) error
	ListByCorrespondenceID(ctx context.Context, schoolID, correspondenceID string) ([]models.Disposition, error)
	GetInboxByUser(ctx context.Context, schoolID, userID string) ([]models.Disposition, error)
	UpdateStatus(ctx context.Context, schoolID, id, status string) error
}

type pgxDispositionRepository struct {
	db *pgxpool.Pool
}

func NewDispositionRepository(db *pgxpool.Pool) DispositionRepository {
	return &pgxDispositionRepository{db: db}
}

func (r *pgxDispositionRepository) Create(ctx context.Context, d *models.Disposition) error {
	query := `
		INSERT INTO dispositions (
			id, correspondence_id, school_id, from_user_id, to_user_id, note, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	d.CreatedAt = time.Now().UTC()
	d.UpdatedAt = time.Now().UTC()

	_, err := r.db.Exec(ctx, query,
		d.ID, d.CorrespondenceID, d.SchoolID, d.FromUserID, d.ToUserID, d.Note, d.Status, d.CreatedAt, d.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxDispositionRepository) ListByCorrespondenceID(ctx context.Context, schoolID, correspondenceID string) ([]models.Disposition, error) {
	query := `
		SELECT id, correspondence_id, school_id, from_user_id, to_user_id, COALESCE(note, ''), status, created_at, updated_at
		FROM dispositions
		WHERE school_id = $1 AND correspondence_id = $2
		ORDER BY created_at ASC
	`
	rows, err := r.db.Query(ctx, query, schoolID, correspondenceID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var list []models.Disposition
	for rows.Next() {
		var d models.Disposition
		if err := rows.Scan(
			&d.ID, &d.CorrespondenceID, &d.SchoolID, &d.FromUserID, &d.ToUserID, &d.Note, &d.Status, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		list = append(list, d)
	}
	return list, rows.Err()
}

func (r *pgxDispositionRepository) GetInboxByUser(ctx context.Context, schoolID, userID string) ([]models.Disposition, error) {
	query := `
		SELECT id, correspondence_id, school_id, from_user_id, to_user_id, COALESCE(note, ''), status, created_at, updated_at
		FROM dispositions
		WHERE school_id = $1 AND to_user_id = $2
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, schoolID, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var list []models.Disposition
	for rows.Next() {
		var d models.Disposition
		if err := rows.Scan(
			&d.ID, &d.CorrespondenceID, &d.SchoolID, &d.FromUserID, &d.ToUserID, &d.Note, &d.Status, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		list = append(list, d)
	}
	return list, rows.Err()
}

func (r *pgxDispositionRepository) UpdateStatus(ctx context.Context, schoolID, id, status string) error {
	query := `
		UPDATE dispositions
		SET status = $1, updated_at = NOW()
		WHERE id = $2 AND school_id = $3
	`
	result, err := r.db.Exec(ctx, query, status, id, schoolID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("disposition not found")
	}
	return nil
}
