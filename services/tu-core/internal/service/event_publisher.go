package service

import "context"

// DomainEvent merepresentasikan event yang dipublish setelah operasi sukses.
type DomainEvent struct {
	EventType        string                 `json:"eventType"`
	AggregateType    string                 `json:"aggregateType"`
	AggregateID      string                 `json:"aggregateId"`
	Payload          map[string]interface{} `json:"payload"`
	SchoolID         string                 `json:"schoolId"`
	AcademicPeriodID string                 `json:"academicPeriodId"`
	CorrelationID    string                 `json:"correlationId"`
	TraceID          string                 `json:"traceId"`
	Version          string                 `json:"version"`
}

// EventPublisher adalah kontrak untuk publish Domain Event ke message broker.
type EventPublisher interface {
	Publish(ctx context.Context, event DomainEvent) error
	PublishBatch(ctx context.Context, events []DomainEvent) error
}

// NoopEventPublisher adalah implementation no-op untuk testing atau ketika event tidak diperlukan.
type NoopEventPublisher struct{}

// NewNoopEventPublisher membuat instance NoopEventPublisher.
func NewNoopEventPublisher() *NoopEventPublisher {
	return &NoopEventPublisher{}
}

// Publish tidak melakukan apa-apa (no-op).
func (p *NoopEventPublisher) Publish(ctx context.Context, event DomainEvent) error {
	return nil
}

// PublishBatch tidak melakukan apa-apa (no-op).
func (p *NoopEventPublisher) PublishBatch(ctx context.Context, events []DomainEvent) error {
	return nil
}
