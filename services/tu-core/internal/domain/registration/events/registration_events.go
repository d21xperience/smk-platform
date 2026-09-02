package events

import (
	"time"

	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"

	"github.com/google/uuid"
)

const (
	EventRegistrationCreated   = "RegistrationCreated"
	EventRegistrationSubmitted = "RegistrationSubmitted"
	EventRegistrationVerified  = "RegistrationVerified"
	EventRegistrationApproved  = "RegistrationApproved"
	EventRegistrationRejected  = "RegistrationRejected"
)

type RegistrationCreatedPayload struct {
	RegistrationID string
	CalonNISN      string
	CalonNama      string
}

func NewRegistrationCreated(p RegistrationCreatedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	builder := types.NewBaseEventBuilder(EventRegistrationCreated, "Registration", p.RegistrationID).
		WithEventID(uuid.New().String()).
		WithVersion(1)
	if ctx != nil {
		builder.WithContext(ctx)
	}
	return builder.WithPayload(map[string]interface{}{
		"registrationId": p.RegistrationID,
		"calonNisn":      p.CalonNISN,
		"calonNama":      p.CalonNama,
		"createdAt":      time.Now().UTC().Format(time.RFC3339),
	}).Build()
}

type RegistrationSubmittedPayload struct {
	RegistrationID string
	CalonNISN      string
	CalonNama      string
}

func NewRegistrationSubmitted(p RegistrationSubmittedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	builder := types.NewBaseEventBuilder(EventRegistrationSubmitted, "Registration", p.RegistrationID).
		WithEventID(uuid.New().String()).
		WithVersion(1)
	if ctx != nil {
		builder.WithContext(ctx)
	}
	return builder.WithPayload(map[string]interface{}{
		"registrationId": p.RegistrationID,
		"calonNisn":      p.CalonNISN,
		"calonNama":      p.CalonNama,
		"submittedAt":    time.Now().UTC().Format(time.RFC3339),
	}).Build()
}

type RegistrationVerifiedPayload struct {
	RegistrationID string
	VerifiedBy     string
	Catatan        string
}

func NewRegistrationVerified(p RegistrationVerifiedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	builder := types.NewBaseEventBuilder(EventRegistrationVerified, "Registration", p.RegistrationID).
		WithEventID(uuid.New().String()).
		WithVersion(1)
	if ctx != nil {
		builder.WithContext(ctx)
	}
	return builder.WithPayload(map[string]interface{}{
		"registrationId": p.RegistrationID,
		"verifiedBy":     p.VerifiedBy,
		"catatan":        p.Catatan,
		"verifiedAt":     time.Now().UTC().Format(time.RFC3339),
	}).Build()
}

type RegistrationApprovedPayload struct {
	RegistrationID string
	CalonNISN      string
	CalonNama      string
	ApprovedBy     string
}

func NewRegistrationApproved(p RegistrationApprovedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	builder := types.NewBaseEventBuilder(EventRegistrationApproved, "Registration", p.RegistrationID).
		WithEventID(uuid.New().String()).
		WithVersion(1)
	if ctx != nil {
		builder.WithContext(ctx)
	}
	return builder.WithPayload(map[string]interface{}{
		"registrationId": p.RegistrationID,
		"calonNisn":      p.CalonNISN,
		"calonNama":      p.CalonNama,
		"approvedBy":     p.ApprovedBy,
		"approvedAt":     time.Now().UTC().Format(time.RFC3339),
	}).Build()
}

type RegistrationRejectedPayload struct {
	RegistrationID string
	CalonNISN      string
	Alasan         string
	RejectedBy     string
}

func NewRegistrationRejected(p RegistrationRejectedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	builder := types.NewBaseEventBuilder(EventRegistrationRejected, "Registration", p.RegistrationID).
		WithEventID(uuid.New().String()).
		WithVersion(1)
	if ctx != nil {
		builder.WithContext(ctx)
	}
	return builder.WithPayload(map[string]interface{}{
		"registrationId": p.RegistrationID,
		"calonNisn":      p.CalonNISN,
		"alasan":         p.Alasan,
		"rejectedBy":     p.RejectedBy,
		"rejectedAt":     time.Now().UTC().Format(time.RFC3339),
	}).Build()
}
