package log

import (
	"time"
)

// NotificationLog merepresentasikan log pengiriman notifikasi untuk audit trail.
type NotificationLog struct {
	ID               string
	SchoolID         string
	AcademicPeriodID string
	CorrelationID    string
	Channel          string
	Recipient        string
	TemplateCode     string
	Status           string // PENDING, SENT, FAILED
	ErrorMessage     string
	AttemptCount     int
	CreatedAt        time.Time
	SentAt           *time.Time
	UpdatedAt        time.Time
}

// NotificationLogRepository adalah kontrak untuk menyimpan log notifikasi.
type NotificationLogRepository interface {
	Save(log *NotificationLog) error
	UpdateStatus(id, status, errorMessage string) error
	FindByCorrelationID(correlationID string) ([]NotificationLog, error)
	FindBySchoolID(schoolID string, limit, offset int) ([]NotificationLog, error)
}
