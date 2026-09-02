package events

import (
	"time"

	"github.com/google/uuid"
	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
)

const (
	EventCorrespondenceCreated  = "CorrespondenceCreated"
	EventCorrespondenceProcessed = "CorrespondenceProcessed"
	EventCorrespondenceArchived = "CorrespondenceArchived"
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

type CorrespondenceCreatedPayload struct {
	CorrespondenceID string
	Type             string
	Number           string
	Subject          string
}

func NewCorrespondenceCreated(p CorrespondenceCreatedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventCorrespondenceCreated, "Correspondence", p.CorrespondenceID, ctx).
		WithPayload(map[string]interface{}{
			"correspondenceId": p.CorrespondenceID,
			"type":             p.Type,
			"number":           p.Number,
			"subject":          p.Subject,
			"createdAt":        time.Now().UTC().Format(time.RFC3339),
		}).Build()
}

type CorrespondenceProcessedPayload struct {
	CorrespondenceID string
	Number           string
}

func NewCorrespondenceProcessed(p CorrespondenceProcessedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventCorrespondenceProcessed, "Correspondence", p.CorrespondenceID, ctx).
		WithPayload(map[string]interface{}{
			"correspondenceId": p.CorrespondenceID,
			"number":           p.Number,
			"processedAt":      time.Now().UTC().Format(time.RFC3339),
		}).Build()
}

type CorrespondenceArchivedPayload struct {
	CorrespondenceID string
	Number           string
}

func NewCorrespondenceArchived(p CorrespondenceArchivedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventCorrespondenceArchived, "Correspondence", p.CorrespondenceID, ctx).
		WithPayload(map[string]interface{}{
			"correspondenceId": p.CorrespondenceID,
			"number":           p.Number,
			"archivedAt":       time.Now().UTC().Format(time.RFC3339),
		}).Build()
}
