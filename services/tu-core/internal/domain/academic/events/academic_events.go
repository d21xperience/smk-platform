package events

import (
	"time"

	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"

	"github.com/google/uuid"
)

const (
	EventAcademicYearCreated  = "AcademicYearCreated"
	EventSemesterAdded        = "SemesterAdded"
	EventClassCreated         = "ClassCreated"
	EventClassHomeroomChanged = "ClassHomeroomChanged"
)

// === Payload Structs ===

type AcademicYearCreatedPayload struct {
	AcademicYearID string
	Label          string
	StartYear      int
	EndYear        int
	IsActive       bool
}

type SemesterAddedPayload struct {
	AcademicYearID string
	SemesterID     string
	Number         int
	PeriodID       string
	Label          string
}

type ClassCreatedPayload struct {
	ClassID    string
	PeriodID   string
	Name       string
	Grade      int
	Level      string
	Major      string
	HomeroomID string
	Capacity   int
}

type ClassHomeroomChangedPayload struct {
	ClassID    string
	HomeroomID string
	ChangedAt  time.Time
}

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

// === Event Factories ===

func NewAcademicYearCreated(p AcademicYearCreatedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventAcademicYearCreated, "AcademicYear", p.AcademicYearID, ctx).
		WithPayload(map[string]interface{}{
			"academicYearId": p.AcademicYearID,
			"label":          p.Label,
			"startYear":      p.StartYear,
			"endYear":        p.EndYear,
			"isActive":       p.IsActive,
		}).
		Build()
}

func NewSemesterAdded(p SemesterAddedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventSemesterAdded, "AcademicYear", p.AcademicYearID, ctx).
		WithPayload(map[string]interface{}{
			"academicYearId": p.AcademicYearID,
			"semesterId":     p.SemesterID,
			"number":         p.Number,
			"periodId":       p.PeriodID,
			"label":          p.Label,
		}).
		Build()
}

func NewClassCreated(p ClassCreatedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventClassCreated, "Class", p.ClassID, ctx).
		WithPayload(map[string]interface{}{
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

func NewClassHomeroomChanged(p ClassHomeroomChangedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventClassHomeroomChanged, "Class", p.ClassID, ctx).
		WithPayload(map[string]interface{}{
			"classId":    p.ClassID,
			"homeroomId": p.HomeroomID,
			"changedAt":  p.ChangedAt.Format(time.RFC3339),
		}).
		Build()
}
