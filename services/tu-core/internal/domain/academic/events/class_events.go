package events

import (
	"time"

	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
)

// =============================================================================
// CLASS CREATED
// =============================================================================
const (
	EventClassCreated         = "ClassCreated"
	EventClassHomeroomChanged = "ClassHomeroomChanged"
	EventClassUpdated         = "ClassUpdated"
	EventClassStatusChanged   = "ClassChanged"
)
const ClassCreatedEventType = "academic.class.created"

type ClassCreatedPayload struct {
	ClassID    string `json:"classId"`
	PeriodID   string `json:"periodId"`
	Name       string `json:"name"`
	Grade      int    `json:"grade"`
	Level      string `json:"level"`
	Major      string `json:"major"`
	HomeroomID string `json:"homeroomId"`
	Capacity   int    `json:"capacity"`
}

type ClassCreated struct {
	payload  ClassCreatedPayload
	metadata types.EventMetadata
}

func NewClassCreated(p ClassCreatedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventClassCreated, "Class", p.ClassID, ctx).
		WithPayload(map[string]any{
			"classId":    p.ClassID,
			"periodId":   p.PeriodID,
			"name":       p.Name,
			"grade":      p.Grade,
			"level":      p.Level,
			"major":      p.Major,
			"homeroomId": p.HomeroomID,
			"capacity":   p.Capacity,
		}).
		Build()
}

// =============================================================================
// CLASS UPDATED
// =============================================================================

const ClassUpdatedEventType = "academic.class.updated"

type ClassUpdatedPayload struct {
	ClassID    string `json:"classId"`
	Name       string `json:"name"`
	Grade      int    `json:"grade"`
	Level      string `json:"level"`
	Major      string `json:"major"`
	HomeroomID string `json:"homeroomId"`
	Capacity   int    `json:"capacity"`
	Status     string `json:"status"`
}

func NewClassUpdated(p ClassUpdatedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventClassUpdated, "Class", p.ClassID, ctx).
		WithPayload(map[string]any{
			"classId":    p.ClassID,
			"name":       p.Name,
			"grade":      p.Grade,
			"level":      p.Level,
			"major":      p.Major,
			"homeroomId": p.HomeroomID,
			"capacity":   p.Capacity,
			"status":     p.Status,
		}).
		Build()
}

// =============================================================================
// CLASS HOMEROOM CHANGED
// =============================================================================

const ClassHomeroomChangedEventType = "academic.class.homeroom_changed"

type ClassHomeroomChangedPayload struct {
	ClassID    string    `json:"classId"`
	HomeroomID string    `json:"homeroomId"`
	ChangedAt  time.Time `json:"changedAt"`
}

func NewClassHomeroomChanged(p ClassHomeroomChangedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventClassHomeroomChanged, "Class", p.ClassID, ctx).
		WithPayload(map[string]any{
			"classId":    p.ClassID,
			"homeroomId": p.HomeroomID,
		}).
		Build()
}

// =============================================================================
// CLASS STATUS CHANGED (opsional)
// =============================================================================

const ClassStatusChangedEventType = "academic.class.status_changed"

type ClassStatusChangedPayload struct {
	ClassID   string    `json:"classId"`
	OldStatus string    `json:"oldStatus"`
	NewStatus string    `json:"newStatus"`
	ChangedAt time.Time `json:"changedAt"`
}

func NewClassStatusChanged(p ClassStatusChangedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventClassStatusChanged, "Class", p.ClassID, ctx).
		WithPayload(map[string]any{
			"classId":   p.ClassID,
			"oldStatus": p.OldStatus,
			"newStatus": p.NewStatus,
		}).
		Build()
}
