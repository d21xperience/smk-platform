package projection

import (
	"context"
)

// DomainEvent merepresentasikan event yang diproses oleh projection.
type DomainEvent struct {
	EventType        string                 `json:"eventType"`
	AggregateType    string                 `json:"aggregateType"`
	AggregateID      string                 `json:"aggregateId"`
	Payload          map[string]interface{} `json:"payload"`
	SchoolID         string                 `json:"schoolId"`
	AcademicPeriodID string                 `json:"academicPeriodId"`
	CorrelationID    string                 `json:"correlationId"`
	TraceID          string                 `json:"traceId"`
	Version          string                 `json:"version"` // Event versioning (ADR-009)
}

// ProjectionHandler adalah kontrak untuk semua projection handlers.
// Setiap handler HARUS idempotent (Prinsip ADR-009).
type ProjectionHandler interface {
	// HandleEvent memproses domain event dan mengupdate read model.
	HandleEvent(ctx context.Context, event DomainEvent) error

	// SupportedEvents mengembalikan daftar event type yang didukung.
	SupportedEvents() []string

	// Name mengembalikan nama projection (untuk logging dan monitoring).
	Name() string
}

// ProjectionMetadata menyimpan metadata tentang sebuah projection.
type ProjectionMetadata struct {
	Name            string
	SupportedEvents []string
	Handler         ProjectionHandler
}
