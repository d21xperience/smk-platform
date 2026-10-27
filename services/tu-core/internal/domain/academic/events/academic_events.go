package events

import (
	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
)

const (
	EventAcademicYearCreated = "AcademicYearCreated"
	EventSemesterAdded       = "SemesterAdded"
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
