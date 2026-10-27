package events

import (
	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"

	"github.com/google/uuid"
)

// === Helper ===

func baseBuilder(eventName, aggregateType, aggregateID string, ctx *context.OperationalContext) *types.BaseEventBuilder {
	builder := types.NewBaseEventBuilder(eventName, aggregateType, aggregateID).
		WithEventID(uuid.New().String()).
		WithVersion(1)
	if ctx != nil {
		builder.WithContext(ctx)
	}
	return builder
}
