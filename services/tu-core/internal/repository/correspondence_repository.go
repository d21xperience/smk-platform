package repository

import (
	"context"
	"encoding/json"
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
	JenisSurat       string
	Kategori         string
	KodeKlasifikasi  string
	DateFrom         string
	DateTo           string
	Search           string
	Limit            int
	Offset           int
}

type CorrespondenceRepository interface {
	Create(ctx context.Context, c *models.Correspondence) error
	GetByID(ctx context.Context, schoolID, academicPeriodID, correspondenceID string) (*models.Correspondence, error)
	List(ctx context.Context, filter CorrespondenceFilter) ([]models.Correspondence, int, error)
	Update(ctx context.Context, c *models.Correspondence) error
	Delete(ctx context.Context, schoolID, academicPeriodID, correspondenceID string) error
	GenerateNumber(ctx context.Context, schoolID, academicPeriodID, jenisSurat string) (int, error)
}

type pgxCorrespondenceRepository struct {
	db *pgxpool.Pool
}

func NewCorrespondenceRepository(db *pgxpool.Pool) CorrespondenceRepository {
	return &pgxCorrespondenceRepository{db: db}
}

func (r *pgxCorrespondenceRepository) GenerateNumber(ctx context.Context, schoolID, academicPeriodID, jenisSurat string) (int, error) {
	query := `
		INSERT INTO correspondence_counters (id, school_id, academic_period_id, jenis_surat, last_number)
		VALUES (gen_random_uuid(), $1, $2, $3, 1)
		ON CONFLICT (school_id, academic_period_id, jenis_surat)
		DO UPDATE SET last_number = correspondence_counters.last_number + 1
		RETURNING last_number
	`
	var nextNum int
	err := r.db.QueryRow(ctx, query, schoolID, academicPeriodID, jenisSurat).Scan(&nextNum)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nextNum, nil
}

func (r *pgxCorrespondenceRepository) Create(ctx context.Context, c *models.Correspondence) error {
	query := `
		INSERT INTO correspondences (
			id, school_id, academic_period_id, type, jenis_surat, kategori, kode_klasifikasi, number, date, subject,
			from_party, to_party, description, attachment_url, file_url, ttd_url, konten, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
	`
	kontenJSON, _ := json.Marshal(c.Konten)
	_, err := r.db.Exec(ctx, query,
		c.ID, c.SchoolID, c.AcademicPeriodID, c.Type.String(), c.JenisSurat, c.Kategori, c.KodeKlasifikasi, c.Number, c.Date, c.Subject,
		c.FromParty, c.ToParty, c.Description, c.AttachmentURL, c.FileURL, c.TTDURL, kontenJSON, c.Status.String(), c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxCorrespondenceRepository) GetByID(ctx context.Context, schoolID, academicPeriodID, correspondenceID string) (*models.Correspondence, error) {
	query := `
		SELECT id, school_id, academic_period_id, type, jenis_surat, COALESCE(kategori,'') as kategori, COALESCE(kode_klasifikasi,'') as kode_klasifikasi, number, date, subject,
		       from_party, to_party, COALESCE(description, '') as description,
		       COALESCE(attachment_url, '') as attachment_url, COALESCE(file_url, '') as file_url,
		       COALESCE(ttd_url, '') as ttd_url, konten, status, created_at, updated_at
		FROM correspondences
		WHERE id = $1 AND school_id = $2 AND academic_period_id = $3
	`

	var c models.Correspondence
	var typeStr, statusStr string
	var kontenRaw []byte

	err := r.db.QueryRow(ctx, query, correspondenceID, schoolID, academicPeriodID).Scan(
		&c.ID, &c.SchoolID, &c.AcademicPeriodID, &typeStr, &c.JenisSurat, &c.Kategori, &c.KodeKlasifikasi, &c.Number, &c.Date, &c.Subject,
		&c.FromParty, &c.ToParty, &c.Description, &c.AttachmentURL, &c.FileURL, &c.TTDURL, &kontenRaw, &statusStr, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCorrespondenceNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	c.Type, _ = valueobjects.NewCorrespondenceType(typeStr)
	c.Status, _ = valueobjects.NewCorrespondenceStatus(statusStr)
	if len(kontenRaw) > 0 {
		_ = json.Unmarshal(kontenRaw, &c.Konten)
	}

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
	if filter.JenisSurat != "" {
		whereClause += fmt.Sprintf(" AND jenis_surat = $%d", argIndex)
		args = append(args, filter.JenisSurat)
		argIndex++
	}
	if filter.Kategori != "" {
		whereClause += fmt.Sprintf(" AND kategori = $%d", argIndex)
		args = append(args, filter.Kategori)
		argIndex++
	}
	if filter.KodeKlasifikasi != "" {
		whereClause += fmt.Sprintf(" AND kode_klasifikasi = $%d", argIndex)
		args = append(args, filter.KodeKlasifikasi)
		argIndex++
	}
	if filter.DateFrom != "" {
		whereClause += fmt.Sprintf(" AND date >= $%d", argIndex)
		args = append(args, filter.DateFrom)
		argIndex++
	}
	if filter.DateTo != "" {
		whereClause += fmt.Sprintf(" AND date <= $%d", argIndex)
		args = append(args, filter.DateTo)
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
		SELECT id, school_id, academic_period_id, type, jenis_surat, COALESCE(kategori,'') as kategori, COALESCE(kode_klasifikasi,'') as kode_klasifikasi, number, date, subject,
		       from_party, to_party, COALESCE(description, '') as description,
		       COALESCE(attachment_url, '') as attachment_url, COALESCE(file_url, '') as file_url,
		       COALESCE(ttd_url, '') as ttd_url, konten, status, created_at, updated_at
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
		var kontenRaw []byte

		if err := rows.Scan(
			&c.ID, &c.SchoolID, &c.AcademicPeriodID, &typeStr, &c.JenisSurat, &c.Kategori, &c.KodeKlasifikasi, &c.Number, &c.Date, &c.Subject,
			&c.FromParty, &c.ToParty, &c.Description, &c.AttachmentURL, &c.FileURL, &c.TTDURL, &kontenRaw, &statusStr, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}

		c.Type, _ = valueobjects.NewCorrespondenceType(typeStr)
		c.Status, _ = valueobjects.NewCorrespondenceStatus(statusStr)
		if len(kontenRaw) > 0 {
			_ = json.Unmarshal(kontenRaw, &c.Konten)
		}

		correspondences = append(correspondences, c)
	}

	return correspondences, total, rows.Err()
}

func (r *pgxCorrespondenceRepository) Update(ctx context.Context, c *models.Correspondence) error {
	query := `
		UPDATE correspondences
		SET kategori = $1, kode_klasifikasi = $2, subject = $3, description = $4, attachment_url = $5, file_url = $6, ttd_url = $7,
		    konten = $8, status = $9, updated_at = $10
		WHERE id = $11 AND school_id = $12 AND academic_period_id = $13
	`
	kontenJSON, _ := json.Marshal(c.Konten)
	result, err := r.db.Exec(ctx, query,
		c.Kategori, c.KodeKlasifikasi, c.Subject, c.Description, c.AttachmentURL, c.FileURL, c.TTDURL, kontenJSON, c.Status.String(), c.UpdatedAt,
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

func (r *pgxCorrespondenceRepository) Delete(ctx context.Context, schoolID, academicPeriodID, correspondenceID string) error {
	query := `
		DELETE FROM correspondences 
		WHERE id = $1 AND school_id = $2 AND academic_period_id = $3
	`
	result, err := r.db.Exec(ctx, query, correspondenceID, schoolID, academicPeriodID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrCorrespondenceNotFound
	}
	return nil
}
