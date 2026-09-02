package repository

import (
	"context"
	"errors"
	"fmt"

	"sekolah-platform/services/tu-core/internal/domain/correspondence/models"
	"sekolah-platform/services/tu-core/internal/domain/correspondence/valueobjects"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrCorrespondenceNotFound = errors.New("correspondence not found")

type CorrespondenceFilter struct {
	SchoolID         string
	AcademicPeriodID string
	Type             string
	Status           string
	Search           string
	Limit            int
	Offset           int
}

type CorrespondenceRepository interface {
	Create(ctx context.Context, c *models.Correspondence) error
	GetByID(ctx context.Context, schoolID, academicPeriodID, correspondenceID string) (*models.Correspondence, error)
	List(ctx context.Context, filter CorrespondenceFilter) ([]models.Correspondence, int, error)
	Update(ctx context.Context, c *models.Correspondence) error
}

type pgxCorrespondenceRepository struct {
	db *pgxpool.Pool
}

func NewCorrespondenceRepository(db *pgxpool.Pool) CorrespondenceRepository {
	return &pgxCorrespondenceRepository{db: db}
}

func (r *pgxCorrespondenceRepository) Create(ctx context.Context, c *models.Correspondence) error {
	query := `
		INSERT INTO correspondences (
			id, school_id, academic_period_id, type, number, date, subject,
			from_party, to_party, description, attachment_url, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	_, err := r.db.Exec(ctx, query,
		c.ID, c.SchoolID, c.AcademicPeriodID, c.Type.String(), c.Number, c.Date, c.Subject,
		c.FromParty, c.ToParty, c.Description, c.AttachmentURL, c.Status.String(), c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxCorrespondenceRepository) GetByID(ctx context.Context, schoolID, academicPeriodID, correspondenceID string) (*models.Correspondence, error) {
	query := `
		SELECT id, school_id, academic_period_id, type, number, date, subject,
		       from_party, to_party, COALESCE(description, '') as description,
		       COALESCE(attachment_url, '') as attachment_url, status, created_at, updated_at
		FROM correspondences
		WHERE id = $1 AND school_id = $2 AND academic_period_id = $3
	`

	var c models.Correspondence
	var typeStr, statusStr string

	err := r.db.QueryRow(ctx, query, correspondenceID, schoolID, academicPeriodID).Scan(
		&c.ID, &c.SchoolID, &c.AcademicPeriodID, &typeStr, &c.Number, &c.Date, &c.Subject,
		&c.FromParty, &c.ToParty, &c.Description, &c.AttachmentURL, &statusStr, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCorrespondenceNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	c.Type, _ = valueobjects.NewCorrespondenceType(typeStr)
	c.Status, _ = valueobjects.NewCorrespondenceStatus(statusStr)

	return &c, nil
}

func (r *pgxCorrespondenceRepository) List(ctx context.Context, filter CorrespondenceFilter) ([]models.Correspondence, int, error) {
	whereClause := "WHERE school_id = $1 AND academic_period_id = $2"
	args := []interface{}{filter.SchoolID, filter.AcademicPeriodID}
	argIndex := 3

	if filter.Type != "" {
		whereClause += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, filter.Type)
		argIndex++
	}
	if filter.Status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}
	if filter.Search != "" {
		whereClause += fmt.Sprintf(" AND (subject ILIKE $%d OR number ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filter.Search+"%")
		argIndex++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM correspondences %s", whereClause)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	dataQuery := fmt.Sprintf(`
		SELECT id, school_id, academic_period_id, type, number, date, subject,
		       from_party, to_party, COALESCE(description, '') as description,
		       COALESCE(attachment_url, '') as attachment_url, status, created_at, updated_at
		FROM correspondences %s
		ORDER BY date DESC, created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var correspondences []models.Correspondence
	for rows.Next() {
		var c models.Correspondence
		var typeStr, statusStr string

		if err := rows.Scan(
			&c.ID, &c.SchoolID, &c.AcademicPeriodID, &typeStr, &c.Number, &c.Date, &c.Subject,
			&c.FromParty, &c.ToParty, &c.Description, &c.AttachmentURL, &statusStr, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}

		c.Type, _ = valueobjects.NewCorrespondenceType(typeStr)
		c.Status, _ = valueobjects.NewCorrespondenceStatus(statusStr)

		correspondences = append(correspondences, c)
	}

	return correspondences, total, rows.Err()
}

func (r *pgxCorrespondenceRepository) Update(ctx context.Context, c *models.Correspondence) error {
	query := `
		UPDATE correspondences
		SET subject = $1, description = $2, attachment_url = $3, status = $4, updated_at = $5
		WHERE id = $6 AND school_id = $7 AND academic_period_id = $8
	`
	result, err := r.db.Exec(ctx, query,
		c.Subject, c.Description, c.AttachmentURL, c.Status.String(), c.UpdatedAt,
		c.ID, c.SchoolID, c.AcademicPeriodID,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrCorrespondenceNotFound
	}
	return nil
}
