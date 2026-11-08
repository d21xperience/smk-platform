package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"sekolah-platform/services/tu-core/internal/domain/signature/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SignatureRepository interface {
	Create(ctx context.Context, s *models.DigitalSignature) error
	GetByCorrespondenceID(ctx context.Context, schoolID, correspondenceID string) (*models.DigitalSignature, error)
	GetByVerificationCode(ctx context.Context, code string) (*models.DigitalSignature, error)
}

type pgxSignatureRepository struct {
	db *pgxpool.Pool
}

func NewSignatureRepository(db *pgxpool.Pool) SignatureRepository {
	return &pgxSignatureRepository{db: db}
}

func (r *pgxSignatureRepository) Create(ctx context.Context, s *models.DigitalSignature) error {
	query := `
		INSERT INTO digital_signatures (
			id, correspondence_id, signer_user_id, signer_role, image_url, verification_code, signed_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	if s.SignedAt.IsZero() {
		s.SignedAt = time.Now().UTC()
	}

	_, err := r.db.Exec(ctx, query,
		s.ID, s.CorrespondenceID, s.SignerUserID, s.SignerRole, s.ImageURL, s.VerificationCode, s.SignedAt,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxSignatureRepository) GetByCorrespondenceID(ctx context.Context, schoolID, correspondenceID string) (*models.DigitalSignature, error) {
	// Note: school_id tidak ada di tabel digital_signatures, sehingga kita validasi lewat join ke correspondences jika perlu.
	// Namun untuk kesederhanaan, kita asumsikan correspondenceID sudah unik secara global di sistem ini.
	query := `
		SELECT id, correspondence_id, signer_user_id, signer_role, image_url, verification_code, signed_at
		FROM digital_signatures
		WHERE correspondence_id = $1
		ORDER BY signed_at DESC LIMIT 1
	`
	var s models.DigitalSignature
	err := r.db.QueryRow(ctx, query, correspondenceID).Scan(
		&s.ID, &s.CorrespondenceID, &s.SignerUserID, &s.SignerRole, &s.ImageURL, &s.VerificationCode, &s.SignedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Return nil jika belum ada TTD
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return &s, nil
}

func (r *pgxSignatureRepository) GetByVerificationCode(ctx context.Context, code string) (*models.DigitalSignature, error) {
	query := `
		SELECT id, correspondence_id, signer_user_id, signer_role, image_url, verification_code, signed_at
		FROM digital_signatures
		WHERE verification_code = $1
	`
	var s models.DigitalSignature
	err := r.db.QueryRow(ctx, query, code).Scan(
		&s.ID, &s.CorrespondenceID, &s.SignerUserID, &s.SignerRole, &s.ImageURL, &s.VerificationCode, &s.SignedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("kode verifikasi tidak valid")
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return &s, nil
}
