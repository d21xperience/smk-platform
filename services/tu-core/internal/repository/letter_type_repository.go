package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"sekolah-platform/services/tu-core/internal/domain/correspondence/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LetterTypeRepository interface {
	Create(ctx context.Context, lt *models.LetterType) error
	GetByID(ctx context.Context, id string) (*models.LetterType, error)
	List(ctx context.Context) ([]models.LetterType, error)
	Update(ctx context.Context, lt *models.LetterType) error
	Delete(ctx context.Context, id string) error
}

type pgxLetterTypeRepository struct {
	db *pgxpool.Pool
}

func NewLetterTypeRepository(db *pgxpool.Pool) LetterTypeRepository {
	return &pgxLetterTypeRepository{db: db}
}

func (r *pgxLetterTypeRepository) Create(ctx context.Context, lt *models.LetterType) error {
	query := `
		INSERT INTO letter_types (id, kategori, peruntukan, value, label, kode_klasifikasi, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	if lt.ID == "" {
		lt.ID = uuid.New().String()
	}
	lt.CreatedAt = time.Now().UTC()
	lt.UpdatedAt = time.Now().UTC()

	_, err := r.db.Exec(ctx, query, lt.ID, lt.Kategori, lt.Peruntukan, lt.Value, lt.Label, lt.KodeKlasifikasi, lt.CreatedAt, lt.UpdatedAt)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxLetterTypeRepository) GetByID(ctx context.Context, id string) (*models.LetterType, error) {
	query := `SELECT id, kategori, peruntukan, value, label, kode_klasifikasi, created_at, updated_at FROM letter_types WHERE id = $1`
	var lt models.LetterType
	err := r.db.QueryRow(ctx, query, id).Scan(&lt.ID, &lt.Kategori, &lt.Peruntukan, &lt.Value, &lt.Label, &lt.KodeKlasifikasi, &lt.CreatedAt, &lt.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("letter type not found")
		}
		return nil, err
	}
	return &lt, nil
}

func (r *pgxLetterTypeRepository) List(ctx context.Context) ([]models.LetterType, error) {
	query := `SELECT id, kategori, peruntukan, value, label, kode_klasifikasi, created_at, updated_at FROM letter_types ORDER BY kategori, label ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.LetterType
	for rows.Next() {
		var lt models.LetterType
		if err := rows.Scan(&lt.ID, &lt.Kategori, &lt.Peruntukan, &lt.Value, &lt.Label, &lt.KodeKlasifikasi, &lt.CreatedAt, &lt.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, lt)
	}
	return list, rows.Err()
}

func (r *pgxLetterTypeRepository) Update(ctx context.Context, lt *models.LetterType) error {
	query := `
		UPDATE letter_types
		SET kategori = $1, peruntukan = $2, value = $3, label = $4, kode_klasifikasi = $5, updated_at = $6
		WHERE id = $7
	`
	lt.UpdatedAt = time.Now().UTC()
	result, err := r.db.Exec(ctx, query, lt.Kategori, lt.Peruntukan, lt.Value, lt.Label, lt.KodeKlasifikasi, lt.UpdatedAt, lt.ID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return errors.New("letter type not found")
	}
	return nil
}

func (r *pgxLetterTypeRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM letter_types WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return errors.New("letter type not found")
	}
	return nil
}
