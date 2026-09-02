// Package types menyediakan Domain Event types untuk seluruh platform.
// Struktur DomainEvent SYMMETRIC dengan DomainEvent di frontend JavaScript.
//
// Setiap Domain Event membawa:
// - Metadata: EventID, EventName, Version, OccurredAt
// - Aggregate: AggregateID, AggregateType
// - Operational Context: SchoolID, AcademicYear, Semester, TenantID, UserID
// - Observability: CorrelationID, TraceID, RequestID
// - Payload: data spesifik event
package types

import (
	"time"

	platformctx "sekolah-platform/platform/context"
)

// DomainEvent adalah interface untuk semua domain events.
// Struktur ini SYMMETRIC dengan DomainEvent di frontend.
type DomainEvent interface {
	// Metadata
	GetEventID() string
	GetEventName() string
	GetVersion() int
	GetOccurredAt() time.Time

	// Aggregate
	GetAggregateID() string
	GetAggregateType() string

	// Operational Context
	GetSchoolID() string
	GetAcademicPeriodID() string
	GetUserID() string

	// Observability
	GetCorrelationID() string
	GetTraceID() string
	GetRequestID() string

	// Payload
	GetPayload() interface{}

	// Serialization
	ToMap() map[string]interface{}
}

// BaseEvent adalah implementasi dasar DomainEvent.
// Semua domain event spesifik harus meng-embed struct ini.
type BaseEvent struct {
	// Metadata
	EventID     string    `json:"eventId"`
	EventName   string    `json:"eventName"`
	Version     int       `json:"version"`
	OccurredAt  time.Time `json:"occurredAt"`

	// Aggregate
	AggregateID   string `json:"aggregateId"`
	AggregateType string `json:"aggregateType"`

	// Operational Context (WAJIB)
	SchoolID         string `json:"schoolId"`
	AcademicYear     string `json:"academicYear,omitempty"`
	Semester         int    `json:"semester,omitempty"`
	AcademicPeriodID string `json:"academicPeriodId"`
	TenantID         string `json:"tenantId,omitempty"`
	UserID           string `json:"userId"`

	// Observability (WAJIB)
	CorrelationID string `json:"correlationId"`
	TraceID       string `json:"traceId"`
	RequestID     string `json:"requestId"`

	// Payload
	Payload interface{} `json:"payload"`
}

// Implementasi DomainEvent interface

func (e *BaseEvent) GetEventID() string        { return e.EventID }
func (e *BaseEvent) GetEventName() string      { return e.EventName }
func (e *BaseEvent) GetVersion() int           { return e.Version }
func (e *BaseEvent) GetOccurredAt() time.Time  { return e.OccurredAt }
func (e *BaseEvent) GetAggregateID() string    { return e.AggregateID }
func (e *BaseEvent) GetAggregateType() string  { return e.AggregateType }
func (e *BaseEvent) GetSchoolID() string       { return e.SchoolID }
func (e *BaseEvent) GetAcademicPeriodID() string { return e.AcademicPeriodID }
func (e *BaseEvent) GetUserID() string         { return e.UserID }
func (e *BaseEvent) GetCorrelationID() string  { return e.CorrelationID }
func (e *BaseEvent) GetTraceID() string        { return e.TraceID }
func (e *BaseEvent) GetRequestID() string      { return e.RequestID }
func (e *BaseEvent) GetPayload() interface{}   { return e.Payload }

// ToMap mengkonversi event ke map (untuk serialization ke Message Broker).
func (e *BaseEvent) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"eventId":          e.EventID,
		"eventName":        e.EventName,
		"version":          e.Version,
		"occurredAt":       e.OccurredAt.Format(time.RFC3339Nano),
		"aggregateId":      e.AggregateID,
		"aggregateType":    e.AggregateType,
		"schoolId":         e.SchoolID,
		"academicYear":     e.AcademicYear,
		"semester":         e.Semester,
		"academicPeriodId": e.AcademicPeriodID,
		"tenantId":         e.TenantID,
		"userId":           e.UserID,
		"correlationId":    e.CorrelationID,
		"traceId":          e.TraceID,
		"requestId":        e.RequestID,
		"payload":          e.Payload,
	}
}

// NewBaseEventBuilder membuat builder untuk BaseEvent.
func NewBaseEventBuilder(eventName, aggregateType, aggregateID string) *BaseEventBuilder {
	return &BaseEventBuilder{
		event: &BaseEvent{
			EventName:     eventName,
			AggregateType: aggregateType,
			AggregateID:   aggregateID,
			Version:       1,
			OccurredAt:    time.Now().UTC(),
		},
	}
}

// BaseEventBuilder adalah builder untuk BaseEvent.
type BaseEventBuilder struct {
	event *BaseEvent
}

// WithEventID set event ID.
func (b *BaseEventBuilder) WithEventID(id string) *BaseEventBuilder {
	b.event.EventID = id
	return b
}

// WithVersion set version.
func (b *BaseEventBuilder) WithVersion(v int) *BaseEventBuilder {
	b.event.Version = v
	return b
}

// WithPayload set payload.
func (b *BaseEventBuilder) WithPayload(payload interface{}) *BaseEventBuilder {
	b.event.Payload = payload
	return b
}

// WithContext set operational context.
func (b *BaseEventBuilder) WithContext(ctx *platformctx.OperationalContext) *BaseEventBuilder {
	if ctx == nil {
		return b
	}
	b.event.SchoolID = ctx.SchoolID
	b.event.AcademicYear = ctx.AcademicYear
	b.event.Semester = ctx.Semester
	b.event.AcademicPeriodID = ctx.AcademicPeriodID
	b.event.TenantID = ctx.TenantID
	b.event.UserID = ctx.UserID
	return b
}

// WithObservability set observability IDs.
func (b *BaseEventBuilder) WithObservability(correlationID, traceID, requestID string) *BaseEventBuilder {
	b.event.CorrelationID = correlationID
	b.event.TraceID = traceID
	b.event.RequestID = requestID
	return b
}

// Build membangun BaseEvent.
func (b *BaseEventBuilder) Build() *BaseEvent {
	return b.event
}
