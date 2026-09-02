package contract

import (
	"context"
)

// NotificationRequest membawa semua data yang dibutuhkan untuk mengirim notifikasi.
type NotificationRequest struct {
	Channel      string                 // "email", "whatsapp", "firebase"
	Recipient    string                 // Email, No HP, atau FCM Token
	TemplateCode string                 // e.g., "invoice_created", "attendance_absent"
	Data         map[string]interface{} // Data untuk render template
	// Operational Context (Prinsip 5)
	SchoolID         string
	AcademicPeriodID string
	CorrelationID    string
}

// Notifier adalah kontrak untuk semua channel notifikasi.
type Notifier interface {
	Send(ctx context.Context, req NotificationRequest) error
}
