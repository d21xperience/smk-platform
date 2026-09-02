package events

import (
	"time"

	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"

	"github.com/google/uuid"
)

const (
	EventInvoiceCreated   = "InvoiceCreated"
	EventInvoiceCancelled = "InvoiceCancelled"
	EventPaymentApplied   = "PaymentApplied"
)

func baseBuilder(eventName, aggregateType, aggregateID string, ctx *context.OperationalContext) *types.BaseEventBuilder {
	builder := types.NewBaseEventBuilder(eventName, aggregateType, aggregateID).
		WithEventID(uuid.New().String()).
		WithVersion(1)
	if ctx != nil {
		builder.WithContext(ctx)
	}
	return builder
}

type InvoiceCreatedPayload struct {
	InvoiceID     string
	StudentID     string
	ComponentType string
	Amount        float64
	Month         int
	Year          int
}

func NewInvoiceCreated(p InvoiceCreatedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventInvoiceCreated, "Invoice", p.InvoiceID, ctx).
		WithPayload(map[string]interface{}{
			"invoiceId":     p.InvoiceID,
			"studentId":     p.StudentID,
			"componentType": p.ComponentType,
			"amount":        p.Amount,
			"month":         p.Month,
			"year":          p.Year,
			"createdAt":     time.Now().UTC().Format(time.RFC3339),
		}).Build()
}

type InvoiceCancelledPayload struct {
	InvoiceID string
	Reason    string
}

func NewInvoiceCancelled(p InvoiceCancelledPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventInvoiceCancelled, "Invoice", p.InvoiceID, ctx).
		WithPayload(map[string]interface{}{
			"invoiceId":   p.InvoiceID,
			"reason":      p.Reason,
			"cancelledAt": time.Now().UTC().Format(time.RFC3339),
		}).Build()
}

type PaymentAppliedPayload struct {
	InvoiceID string
	PaymentID string
	Amount    float64
	Status    string
}

func NewPaymentApplied(p PaymentAppliedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventPaymentApplied, "Invoice", p.InvoiceID, ctx).
		WithPayload(map[string]interface{}{
			"invoiceId": p.InvoiceID,
			"paymentId": p.PaymentID,
			"amount":    p.Amount,
			"status":    p.Status,
			"appliedAt": time.Now().UTC().Format(time.RFC3339),
		}).Build()
}
