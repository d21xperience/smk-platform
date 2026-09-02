// Package events menyediakan Domain Events untuk Teacher aggregate.
// PENTING: Package ini TIDAK boleh mengimport models package!
package events

import (
	"time"

	"github.com/google/uuid"
	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
)

const (
	EventTeacherCreated            = "TeacherCreated"
	EventTeacherProfileUpdated     = "TeacherProfileUpdated"
	EventTeacherDeactivated        = "TeacherDeactivated"
	EventTeacherCertificationAdded = "TeacherCertificationAdded"
)

// === Payload Structs ===

type TeacherCreatedPayload struct {
	TeacherID string
	NIP       string
	NUPTK     string
	FullName  string
	Subject   string
	Gender    string
}

type TeacherProfileUpdatedPayload struct {
	TeacherID string
	UpdatedAt time.Time
}

type TeacherDeactivatedPayload struct {
	TeacherID     string
	Reason        string
	DeactivatedAt time.Time
}

type TeacherCertificationAddedPayload struct {
	TeacherID         string
	CertificationID   string
	CertificationName string
	Issuer            string
}

// === Helper ===

func baseBuilder(eventName, aggregateID string, ctx *context.OperationalContext) *types.BaseEventBuilder {
	builder := types.NewBaseEventBuilder(eventName, "Teacher", aggregateID).
		WithEventID(uuid.New().String()).
		WithVersion(1)
	if ctx != nil {
		builder.WithContext(ctx)
	}
	return builder
}

// === Event Factories ===

func NewTeacherCreated(p TeacherCreatedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventTeacherCreated, p.TeacherID, ctx).
		WithPayload(map[string]any{
			"teacherId": p.TeacherID,
			"nip":       p.NIP,
			"nuptk":     p.NUPTK,
			"fullName":  p.FullName,
			"subject":   p.Subject,
			"gender":    p.Gender,
		}).
		Build()
}

func NewTeacherProfileUpdated(p TeacherProfileUpdatedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventTeacherProfileUpdated, p.TeacherID, ctx).
		WithPayload(map[string]interface{}{
			"teacherId": p.TeacherID,
			"updatedAt": p.UpdatedAt.Format(time.RFC3339),
		}).
		Build()
}

func NewTeacherDeactivated(p TeacherDeactivatedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventTeacherDeactivated, p.TeacherID, ctx).
		WithPayload(map[string]interface{}{
			"teacherId":     p.TeacherID,
			"reason":        p.Reason,
			"deactivatedAt": p.DeactivatedAt.Format(time.RFC3339),
		}).
		Build()
}

func NewTeacherCertificationAdded(p TeacherCertificationAddedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventTeacherCertificationAdded, p.TeacherID, ctx).
		WithPayload(map[string]interface{}{
			"teacherId":         p.TeacherID,
			"certificationId":   p.CertificationID,
			"certificationName": p.CertificationName,
			"issuer":            p.Issuer,
		}).
		Build()
}
