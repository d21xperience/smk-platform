package models

import "time"

// =====================================================================
// Constants — PIN
// =====================================================================
const (
	PinLengthDefault        = 6
	MaxAttemptDefault       = 3
	PinLockDurationMinutes  = 15
)

// =====================================================================
// UserPinTTD — model tabel user_pin_ttd
// =====================================================================
type UserPinTTD struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	Username       string     `json:"username,omitempty"`
	PinHash        string     `json:"-"` // NEVER expose ke client
	MaxAttempt     int        `json:"max_attempt"`
	CurrentAttempt int        `json:"current_attempt"`
	LastAttemptAt  *time.Time `json:"last_attempt_at,omitempty"`
	IsLocked       bool       `json:"is_locked"`
	LockedUntil    *time.Time `json:"locked_until,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// =====================================================================
// DTO Request
// =====================================================================
type SetupPinRequest struct {
	UserID string `json:"-"`
	NewPin string `json:"new_pin"`
}

type ChangePinRequest struct {
	UserID string `json:"-"`
	OldPin string `json:"old_pin"`
	NewPin string `json:"new_pin"`
}

type VerifyPinRequest struct {
	UserID string `json:"-"`
	Pin    string `json:"pin"`
}

type UnlockPinRequest struct {
	UserID string `json:"-"`
}

// =====================================================================
// DTO Response
// =====================================================================
type PinStatusResponse struct {
	IsSetup        bool       `json:"is_setup"`
	IsLocked       bool       `json:"is_locked"`
	CurrentAttempt int        `json:"current_attempt"`
	MaxAttempt     int        `json:"max_attempt"`
	LockedUntil    *time.Time `json:"locked_until,omitempty"`
}

type VerifyPinResponse struct {
	Valid            bool   `json:"valid"`
	CurrentAttempt   int    `json:"current_attempt"`
	RemainingAttempt int    `json:"remaining_attempt"`
	LockedUntil      string `json:"locked_until,omitempty"`
	Message          string `json:"message,omitempty"`
}

// =====================================================================
// Validator
// =====================================================================
func IsValidPinFormat(pin string) bool {
	if len(pin) != PinLengthDefault {
		return false
	}
	for _, c := range pin {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
