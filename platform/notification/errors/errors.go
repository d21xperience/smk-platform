package errors

import (
	"errors"
	"fmt"
)

// Custom error types untuk notification platform
var (
	ErrChannelNotFound      = errors.New("channel notifikasi tidak ditemukan")
	ErrTemplateNotFound     = errors.New("template tidak ditemukan")
	ErrInvalidRecipient     = errors.New("recipient tidak valid")
	ErrInvalidTemplateCode  = errors.New("template code tidak valid")
	ErrTemplateRenderFailed = errors.New("gagal render template")
	ErrSendFailed           = errors.New("gagal mengirim notifikasi")
	ErrRetryExhausted       = errors.New("retry exhausted")
	ErrChannelUnavailable   = errors.New("channel tidak tersedia")
)

// NotificationError adalah error wrapper dengan konteks tambahan.
type NotificationError struct {
	Op            string // Operasi yang gagal (e.g., "Send", "Render")
	Channel       string // Channel yang digunakan
	Recipient     string // Recipient yang dituju
	TemplateCode  string // Template yang digunakan
	CorrelationID string // Untuk tracing
	Err           error  // Underlying error
}

func (e *NotificationError) Error() string {
	return fmt.Sprintf(
		"notification error [op=%s, channel=%s, recipient=%s, template=%s, correlation=%s]: %v",
		e.Op, e.Channel, e.Recipient, e.TemplateCode, e.CorrelationID, e.Err,
	)
}

func (e *NotificationError) Unwrap() error {
	return e.Err
}

// NewNotificationError membuat error baru dengan konteks.
func NewNotificationError(op, channel, recipient, templateCode, correlationID string, err error) *NotificationError {
	return &NotificationError{
		Op:            op,
		Channel:       channel,
		Recipient:     recipient,
		TemplateCode:  templateCode,
		CorrelationID: correlationID,
		Err:           err,
	}
}

// IsChannelNotFound memeriksa apakah error adalah ErrChannelNotFound.
func IsChannelNotFound(err error) bool {
	return errors.Is(err, ErrChannelNotFound)
}

// IsTemplateNotFound memeriksa apakah error adalah ErrTemplateNotFound.
func IsTemplateNotFound(err error) bool {
	return errors.Is(err, ErrTemplateNotFound)
}

// IsRetryExhausted memeriksa apakah error adalah ErrRetryExhausted.
func IsRetryExhausted(err error) bool {
	return errors.Is(err, ErrRetryExhausted)
}
