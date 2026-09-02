// Package retry menyediakan retry policy dan executor dengan exponential backoff.
package retry

import (
	"math"
	"math/rand"
	"time"
)

// BackoffType adalah tipe backoff strategy.
type BackoffType string

const (
	// BackoffExponential adalah exponential backoff (delay * 2^attempt).
	BackoffExponential BackoffType = "exponential"
	// BackoffLinear adalah linear backoff (delay * attempt).
	BackoffLinear BackoffType = "linear"
	// BackoffConstant adalah constant backoff (delay tetap).
	BackoffConstant BackoffType = "constant"
)

// Policy adalah konfigurasi retry.
type Policy struct {
	// MaxRetries adalah jumlah maksimum retry (tidak termasuk attempt pertama).
	MaxRetries int `json:"maxRetries"`

	// InitialDelay adalah delay awal sebelum retry pertama.
	InitialDelay time.Duration `json:"initialDelay"`

	// MaxDelay adalah delay maksimum (cap).
	MaxDelay time.Duration `json:"maxDelay"`

	// BackoffType adalah tipe backoff strategy.
	BackoffType BackoffType `json:"backoffType"`

	// Jitter menambahkan randomization ke delay untuk menghindari thundering herd.
	Jitter bool `json:"jitter"`
}

// DefaultPolicy adalah policy default untuk event retry.
var DefaultPolicy = Policy{
	MaxRetries:   3,
	InitialDelay: 1 * time.Second,
	MaxDelay:     30 * time.Second,
	BackoffType:  BackoffExponential,
	Jitter:       true,
}

// Validate memvalidasi policy.
func (p Policy) Validate() error {
	if p.MaxRetries < 0 {
		return errInvalidMaxRetries
	}
	if p.InitialDelay < 0 {
		return errInvalidInitialDelay
	}
	if p.MaxDelay < p.InitialDelay {
		return errInvalidMaxDelay
	}
	return nil
}

// DelayForAttempt menghitung delay untuk attempt tertentu (0-indexed).
func (p Policy) DelayForAttempt(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}

	var delay time.Duration
	switch p.BackoffType {
	case BackoffExponential:
		multiplier := math.Pow(2, float64(attempt))
		delay = time.Duration(float64(p.InitialDelay) * multiplier)
	case BackoffLinear:
		delay = p.InitialDelay * time.Duration(attempt+1)
	case BackoffConstant:
		delay = p.InitialDelay
	default:
		delay = p.InitialDelay
	}

	// Cap ke MaxDelay
	if delay > p.MaxDelay {
		delay = p.MaxDelay
	}

	// Add jitter (±25%)
	if p.Jitter && delay > 0 {
		jitterRange := float64(delay) * 0.25
		jitter := (rand.Float64()*2 - 1) * jitterRange
		delay = time.Duration(float64(delay) + jitter)
	}

	return delay
}

var (
	errInvalidMaxRetries   = errString("maxRetries tidak boleh negatif")
	errInvalidInitialDelay = errString("initialDelay tidak boleh negatif")
	errInvalidMaxDelay     = errString("maxDelay harus >= initialDelay")
)

type errString string

func (e errString) Error() string { return string(e) }
