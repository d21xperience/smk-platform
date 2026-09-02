package events

import (
	"time"

	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"

	"github.com/google/uuid"
)

const (
	EventAssessmentCreated = "AssessmentCreated"
	EventAssessmentUpdated = "AssessmentUpdated"
	EventAssessmentDeleted = "AssessmentDeleted"
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

type AssessmentCreatedPayload struct {
	AssessmentID   string
	StudentID      string
	SubjectID      string
	AssessmentType string
	Score          float64
	MaxScore       float64
	Semester       int
}

func NewAssessmentCreated(p AssessmentCreatedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventAssessmentCreated, "Assessment", p.AssessmentID, ctx).
		WithPayload(map[string]interface{}{
			"assessmentId":   p.AssessmentID,
			"studentId":      p.StudentID,
			"subjectId":      p.SubjectID,
			"assessmentType": p.AssessmentType,
			"score":          p.Score,
			"maxScore":       p.MaxScore,
			"semester":       p.Semester,
			"createdAt":      time.Now().UTC().Format(time.RFC3339),
		}).Build()
}

type AssessmentUpdatedPayload struct {
	AssessmentID string
	StudentID    string
	SubjectID    string
	OldScore     float64
	NewScore     float64
}

func NewAssessmentUpdated(p AssessmentUpdatedPayload) *types.BaseEvent {
	return baseBuilder(EventAssessmentUpdated, "Assessment", p.AssessmentID, nil).
		WithPayload(map[string]interface{}{
			"assessmentId": p.AssessmentID,
			"studentId":    p.StudentID,
			"subjectId":    p.SubjectID,
			"oldScore":     p.OldScore,
			"newScore":     p.NewScore,
			"updatedAt":    time.Now().UTC().Format(time.RFC3339),
		}).Build()
}

type AssessmentDeletedPayload struct {
	AssessmentID string
	StudentID    string
	SubjectID    string
}

func NewAssessmentDeleted(p AssessmentDeletedPayload) *types.BaseEvent {
	return baseBuilder(EventAssessmentDeleted, "Assessment", p.AssessmentID, nil).
		WithPayload(map[string]interface{}{
			"assessmentId": p.AssessmentID,
			"studentId":    p.StudentID,
			"subjectId":    p.SubjectID,
			"deletedAt":    time.Now().UTC().Format(time.RFC3339),
		}).Build()
}
