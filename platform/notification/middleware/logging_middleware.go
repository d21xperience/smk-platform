package middleware

import (
	"context"
	"log"
	"time"

	"sekolah-platform/platform/notification/contract"
)

// LoggingMiddleware menambahkan logging sebelum dan sesudah pengiriman notifikasi.
type LoggingMiddleware struct {
	next contract.Notifier
}

// NewLoggingMiddleware membuat middleware logging baru.
func NewLoggingMiddleware(next contract.Notifier) *LoggingMiddleware {
	return &LoggingMiddleware{next: next}
}

// Send mengirimkan notifikasi dengan logging.
func (m *LoggingMiddleware) Send(ctx context.Context, req contract.NotificationRequest) error {
	start := time.Now()

	log.Printf("[Notification] START channel=%s recipient=%s template=%s correlation=%s",
		req.Channel, req.Recipient, req.TemplateCode, req.CorrelationID)

	err := m.next.Send(ctx, req)
	duration := time.Since(start)

	if err != nil {
		log.Printf("[Notification] FAILED channel=%s recipient=%s duration=%v error=%v correlation=%s",
			req.Channel, req.Recipient, duration, err, req.CorrelationID)
	} else {
		log.Printf("[Notification] SUCCESS channel=%s recipient=%s duration=%v correlation=%s",
			req.Channel, req.Recipient, duration, req.CorrelationID)
	}

	return err
}
