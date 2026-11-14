package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
	"sekolah-platform/services/tu-core/internal/repository"
)

// =====================================================================
// UserPinTTDService — setup, verify, change, unlock PIN TTD
// =====================================================================
type UserPinTTDService struct {
	repo repository.UserPinTTDRepository
}

func NewUserPinTTDService(repo repository.UserPinTTDRepository) *UserPinTTDService {
	return &UserPinTTDService{repo: repo}
}

const pinBcryptCost = 10
const pinLockDuration = 15 * time.Minute

// =====================================================================
// SetupPin — daftarkan PIN pertama kali (fail jika sudah ada)
// =====================================================================
func (s *UserPinTTDService) SetupPin(
	ctx context.Context, req models.SetupPinRequest,
) (*models.UserPinTTD, error) {
	if req.UserID == "" {
		return nil, errors.New("user_id wajib diisi")
	}
	if !models.IsValidPinFormat(req.NewPin) {
		return nil, fmt.Errorf("PIN harus %d digit angka", models.PinLengthDefault)
	}

	// Cek belum ada
	existing, err := s.repo.GetByUserID(ctx, req.UserID)
	if err == nil && existing != nil {
		return nil, errors.New("PIN sudah disetup sebelumnya, gunakan endpoint change untuk ganti")
	}
	if err != nil && !errors.Is(err, repository.ErrPinTTDNotFound) {
		return nil, err
	}

	// Hash
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPin), pinBcryptCost)
	if err != nil {
		return nil, fmt.Errorf("gagal hash PIN: %w", err)
	}

	p := &models.UserPinTTD{
		UserID:     req.UserID,
		PinHash:    string(hash),
		MaxAttempt: models.MaxAttemptDefault,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// =====================================================================
// VerifyPin — verifikasi PIN, kelola attempt & lock
//
// Behavior:
//   1. PIN belum disetup → error
//   2. locked_until > NOW() → return valid=false, message "locked"
//   3. locked_until <= NOW() (auto-unlock) → reset attempt dulu
//   4. bcrypt compare:
//        - valid   → reset attempt=0, lockedUntil=nil
//        - invalid → increment attempt, lock jika >= max
// =====================================================================
func (s *UserPinTTDService) VerifyPin(
	ctx context.Context, req models.VerifyPinRequest,
) (*models.VerifyPinResponse, error) {
	if req.UserID == "" {
		return nil, errors.New("user_id wajib diisi")
	}
	if !models.IsValidPinFormat(req.Pin) {
		return nil, fmt.Errorf("PIN harus %d digit angka", models.PinLengthDefault)
	}

	p, err := s.repo.GetByUserID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	// Cek lock aktif
	if p.LockedUntil != nil && p.LockedUntil.After(now) {
		return &models.VerifyPinResponse{
			Valid:            false,
			CurrentAttempt:   p.CurrentAttempt,
			RemainingAttempt: 0,
			LockedUntil:      p.LockedUntil.Format(time.RFC3339),
			Message:          fmt.Sprintf("PIN terkunci sampai %s", p.LockedUntil.Format("2006-01-02 15:04:05")),
		}, nil
	}

	// Auto-unlock: lock sudah lewat tapi masih ada state lock
	// (bisa juga karena is_locked stale dari sebelum NOW)
	if p.IsLocked || (p.LockedUntil != nil && !p.LockedUntil.After(now)) {
		// Reset attempt sebelum verifikasi (clean state)
		if err := s.repo.UpdateAttempt(ctx, req.UserID, 0, now, nil); err != nil {
			return nil, err
		}
		p.CurrentAttempt = 0
		p.IsLocked = false
		p.LockedUntil = nil
	}

	// Verifikasi PIN
	cmpErr := bcrypt.CompareHashAndPassword([]byte(p.PinHash), []byte(req.Pin))
	if cmpErr == nil {
		// Valid → reset
		if err := s.repo.UpdateAttempt(ctx, req.UserID, 0, now, nil); err != nil {
			return nil, err
		}
		return &models.VerifyPinResponse{
			Valid:            true,
			CurrentAttempt:   0,
			RemainingAttempt: p.MaxAttempt,
			Message:          "PIN valid",
		}, nil
	}

	// Invalid → increment
	newAttempt := p.CurrentAttempt + 1
	var lockedUntil *time.Time
	remaining := p.MaxAttempt - newAttempt

	if newAttempt >= p.MaxAttempt {
		lu := now.Add(pinLockDuration)
		lockedUntil = &lu
		remaining = 0
	}

	if err := s.repo.UpdateAttempt(ctx, req.UserID, newAttempt, now, lockedUntil); err != nil {
		return nil, err
	}

	resp := &models.VerifyPinResponse{
		Valid:            false,
		CurrentAttempt:   newAttempt,
		RemainingAttempt: remaining,
		Message:          "PIN salah",
	}
	if lockedUntil != nil {
		resp.LockedUntil = lockedUntil.Format(time.RFC3339)
		resp.Message = fmt.Sprintf("PIN salah %d kali, terkunci sampai %s",
			newAttempt, lockedUntil.Format("2006-01-02 15:04:05"))
	}
	return resp, nil
}

// =====================================================================
// ChangePin — ganti PIN (verify old dulu, lalu update hash)
// =====================================================================
func (s *UserPinTTDService) ChangePin(
	ctx context.Context, req models.ChangePinRequest,
) error {
	if req.UserID == "" {
		return errors.New("user_id wajib diisi")
	}
	if !models.IsValidPinFormat(req.OldPin) {
		return errors.New("old_pin harus 6 digit angka")
	}
	if !models.IsValidPinFormat(req.NewPin) {
		return errors.New("new_pin harus 6 digit angka")
	}
	if req.OldPin == req.NewPin {
		return errors.New("new_pin harus berbeda dari old_pin")
	}

	// Verify old PIN
	verifyResp, err := s.VerifyPin(ctx, models.VerifyPinRequest{
		UserID: req.UserID,
		Pin:    req.OldPin,
	})
	if err != nil {
		return err
	}
	if !verifyResp.Valid {
		return fmt.Errorf("old_pin tidak valid: %s", verifyResp.Message)
	}

	// Hash new PIN
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPin), pinBcryptCost)
	if err != nil {
		return fmt.Errorf("gagal hash new_pin: %w", err)
	}

	return s.repo.UpdatePin(ctx, req.UserID, string(hash))
}

// =====================================================================
// Unlock — admin reset attempt & lock
// =====================================================================
func (s *UserPinTTDService) Unlock(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user_id wajib diisi")
	}
	return s.repo.Unlock(ctx, userID)
}

// =====================================================================
// GetStatus — status PIN user
// =====================================================================
func (s *UserPinTTDService) GetStatus(
	ctx context.Context, userID string,
) (*models.PinStatusResponse, error) {
	if userID == "" {
		return nil, errors.New("user_id wajib diisi")
	}

	p, err := s.repo.GetByUserID(ctx, userID)
	if errors.Is(err, repository.ErrPinTTDNotFound) {
		return &models.PinStatusResponse{IsSetup: false}, nil
	}
	if err != nil {
		return nil, err
	}

	// Auto-unlock kalau lock sudah lewat
	now := time.Now()
	if p.LockedUntil != nil && !p.LockedUntil.After(now) {
		// Lock sudah lewat — laporkan sebagai unlocked
		return &models.PinStatusResponse{
			IsSetup:        true,
			IsLocked:       false,
			CurrentAttempt: p.CurrentAttempt,
			MaxAttempt:     p.MaxAttempt,
		}, nil
	}

	return &models.PinStatusResponse{
		IsSetup:        true,
		IsLocked:       p.IsLocked,
		CurrentAttempt: p.CurrentAttempt,
		MaxAttempt:     p.MaxAttempt,
		LockedUntil:    p.LockedUntil,
	}, nil
}
