package middleware

import (
	"context"
	"log"
	"time"

	"sekolah-platform/platform/notification/contract"
	notifErrors "sekolah-platform/platform/notification/errors"
)

// RetryConfig konfigurasi retry.
type RetryConfig struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
}

// DefaultRetryConfig memberikan konfigurasi default (3x retry, exponential backoff).
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 1 * time.Second,
		MaxBackoff:     30 * time.Second,
	}
}

// RetryMiddleware menambahkan retry logic dengan exponential backoff.
type RetryMiddleware struct {
	next   contract.Notifier
	config RetryConfig
}

// NewRetryMiddleware membuat middleware retry baru.
func NewRetryMiddleware(next contract.Notifier, config RetryConfig) *RetryMiddleware {
	return &RetryMiddleware{
		next:   next,
		config: config,
	}
}

// Send mengirimkan notifikasi dengan retry.
func (m *RetryMiddleware) Send(ctx context.Context, req contract.NotificationRequest) error {
	// var lastErr error

	for attempt := 0; attempt <= m.config.MaxRetries; attempt++ {
		if attempt > 0 {
			backoff := m.calculateBackoff(attempt)
			log.Printf("[Retry] Attempt %d/%d untuk channel=%s, menunggu %v",
				attempt, m.config.MaxRetries, req.Channel, backoff)
			time.Sleep(backoff)
		}

		err := m.next.Send(ctx, req)
		if err == nil {
			return nil // Sukses
		}

		// lastErr := err
		log.Printf("[Retry] Attempt %d gagal untuk channel=%s: %v", attempt+1, req.Channel, err)
	}

	// Semua retry gagal
	return notifErrors.NewNotificationError("Send", req.Channel, req.Recipient, req.TemplateCode, req.CorrelationID,
		notifErrors.ErrRetryExhausted)
}

// calculateBackoff menghitung exponential backoff.
func (m *RetryMiddleware) calculateBackoff(attempt int) time.Duration {
	backoff := m.config.InitialBackoff * time.Duration(1<<(attempt-1))
	if backoff > m.config.MaxBackoff {
		backoff = m.config.MaxBackoff
	}
	return backoff
}
