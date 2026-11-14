package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
)

// =====================================================================
// UserPinTTDRepository — CRUD PIN TTD
// =====================================================================
type UserPinTTDRepository interface {
	Create(ctx context.Context, p *models.UserPinTTD) error
	GetByUserID(ctx context.Context, userID string) (*models.UserPinTTD, error)
	UpdateAttempt(ctx context.Context, userID string, newAttempt int, lastAt time.Time, lockedUntil *time.Time) error
	UpdatePin(ctx context.Context, userID, newHash string) error
	Unlock(ctx context.Context, userID string) error
}

type pgxUserPinTTDRepository struct {
	db *pgxpool.Pool
}

func NewUserPinTTDRepository(db *pgxpool.Pool) UserPinTTDRepository {
	return &pgxUserPinTTDRepository{db: db}
}

// =====================================================================
// Error sentinel
// =====================================================================
var ErrPinTTDNotFound = errors.New("PIN TTD belum disetup")

// =====================================================================
// Create — insert PIN baru
// =====================================================================
func (r *pgxUserPinTTDRepository) Create(
	ctx context.Context, p *models.UserPinTTD,
) error {
	err := r.db.QueryRow(ctx, `
		INSERT INTO user_pin_ttd (
			user_id, pin_hash, max_attempt, current_attempt, is_locked
		) VALUES (
			$1, $2, $3, 0, FALSE
		)
		RETURNING id, created_at, updated_at
	`, p.UserID, p.PinHash, p.MaxAttempt,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("gagal insert user_pin_ttd: %w", err)
	}
	return nil
}

// =====================================================================
// GetByUserID — read 1 row PIN dengan username (optional join)
// =====================================================================
func (r *pgxUserPinTTDRepository) GetByUserID(
	ctx context.Context, userID string,
) (*models.UserPinTTD, error) {
	var p models.UserPinTTD
	err := r.db.QueryRow(ctx, `
		SELECT p.id, p.user_id, COALESCE(u.username,''), p.pin_hash,
		       p.max_attempt, p.current_attempt, p.last_attempt_at,
		       p.is_locked, p.locked_until, p.created_at, p.updated_at
		FROM user_pin_ttd p
		LEFT JOIN user_login u ON u.user_id = p.user_id
		WHERE p.user_id = $1
	`, userID).Scan(
		&p.ID, &p.UserID, &p.Username, &p.PinHash,
		&p.MaxAttempt, &p.CurrentAttempt, &p.LastAttemptAt,
		&p.IsLocked, &p.LockedUntil, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrPinTTDNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("gagal get PIN TTD: %w", err)
	}
	return &p, nil
}

// =====================================================================
// UpdateAttempt — update attempt counter + lock state
//
// Semantik:
//   - newAttempt=0, lockedUntil=nil → reset (sukses verifikasi)
//   - newAttempt=N, lockedUntil=nil → gagal tapi belum lock
//   - newAttempt=N, lockedUntil=t   → gagal & terkunci
//
// is_locked di-set otomatis dari locked_until > NOW() (auto-unlock saat
// locked_until lewat).
// =====================================================================
func (r *pgxUserPinTTDRepository) UpdateAttempt(
	ctx context.Context, userID string,
	newAttempt int, lastAt time.Time, lockedUntil *time.Time,
) error {
	ct, err := r.db.Exec(ctx, `
		UPDATE user_pin_ttd
		SET current_attempt = $2,
		    last_attempt_at = $3,
		    is_locked = ($4::timestamptz IS NOT NULL AND $4::timestamptz > NOW()),
		    locked_until = $4,
		    updated_at = NOW()
		WHERE user_id = $1
	`, userID, newAttempt, lastAt, lockedUntil)
	if err != nil {
		return fmt.Errorf("gagal update attempt: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrPinTTDNotFound
	}
	return nil
}

// =====================================================================
// UpdatePin — ganti hash PIN (untuk ChangePin)
// =====================================================================
func (r *pgxUserPinTTDRepository) UpdatePin(
	ctx context.Context, userID, newHash string,
) error {
	ct, err := r.db.Exec(ctx, `
		UPDATE user_pin_ttd
		SET pin_hash = $2, updated_at = NOW()
		WHERE user_id = $1
	`, userID, newHash)
	if err != nil {
		return fmt.Errorf("gagal update PIN: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrPinTTDNotFound
	}
	return nil
}

// =====================================================================
// Unlock — reset attempt & lock (admin action)
// =====================================================================
func (r *pgxUserPinTTDRepository) Unlock(
	ctx context.Context, userID string,
) error {
	ct, err := r.db.Exec(ctx, `
		UPDATE user_pin_ttd
		SET current_attempt = 0, is_locked = FALSE, locked_until = NULL,
		    updated_at = NOW()
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return fmt.Errorf("gagal unlock PIN: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrPinTTDNotFound
	}
	return nil
}
