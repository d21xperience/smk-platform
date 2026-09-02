package retry

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrMaxRetriesExceeded adalah error ketika retry habis.
var ErrMaxRetriesExceeded = errors.New("max retries exceeded")

// Operation adalah function yang akan di-retry.
type Operation func() error

// Executor mengeksekusi operation dengan retry policy.
type Executor struct {
	policy Policy
}

// NewExecutor membuat Executor baru.
func NewExecutor(policy Policy) *Executor {
	return &Executor{policy: policy}
}

// Result adalah hasil eksekusi dengan metadata retry.
type Result struct {
	Success    bool
	Attempts   int
	TotalDelay time.Duration
	LastError  error
	Delays     []time.Duration
}

// Execute menjalankan operation dengan retry.
func (e *Executor) Execute(ctx context.Context, op Operation) Result {
	result := Result{}

	for attempt := 0; attempt <= e.policy.MaxRetries; attempt++ {
		result.Attempts = attempt + 1

		// Check context
		if err := ctx.Err(); err != nil {
			result.LastError = err
			return result
		}

		// Execute operation
		err := op()
		if err == nil {
			result.Success = true
			return result
		}

		result.LastError = err

		// Jika masih ada retry, tunggu sesuai delay
		if attempt < e.policy.MaxRetries {
			delay := e.policy.DelayForAttempt(attempt)
			result.Delays = append(result.Delays, delay)
			result.TotalDelay += delay

			select {
			case <-ctx.Done():
				result.LastError = ctx.Err()
				return result
			case <-time.After(delay):
				// Lanjut ke attempt berikutnya
			}
		}
	}

	return result
}

// ExecuteWithCallback sama seperti Execute, tapi memanggil callback setiap retry.
func (e *Executor) ExecuteWithCallback(
	ctx context.Context,
	op Operation,
	onRetry func(attempt int, err error, delay time.Duration),
) Result {
	result := Result{}

	for attempt := 0; attempt <= e.policy.MaxRetries; attempt++ {
		result.Attempts = attempt + 1

		if err := ctx.Err(); err != nil {
			result.LastError = err
			return result
		}

		err := op()
		if err == nil {
			result.Success = true
			return result
		}

		result.LastError = err

		if attempt < e.policy.MaxRetries {
			delay := e.policy.DelayForAttempt(attempt)
			result.Delays = append(result.Delays, delay)
			result.TotalDelay += delay

			if onRetry != nil {
				onRetry(attempt+1, err, delay)
			}

			select {
			case <-ctx.Done():
				result.LastError = ctx.Err()
				return result
			case <-time.After(delay):
			}
		}
	}

	return result
}

// WrapError membungkus error dengan informasi retry.
func WrapError(result Result, originalErr error) error {
	return fmt.Errorf("%w after %d attempts (total delay: %v): %v",
		ErrMaxRetriesExceeded,
		result.Attempts,
		result.TotalDelay,
		originalErr,
	)
}
