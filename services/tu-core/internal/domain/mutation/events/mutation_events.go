package events

import (
	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"

	"github.com/google/uuid"
)

const (
	EventMutationCreated  = "MutationCreated"
	EventMutationApproved = "MutationApproved"
	EventMutationRejected = "MutationRejected"
)

type MutationCreatedPayload struct {
	MutationID    string
	StudentID     *string
	MutationType  string
	CalonNama     string
	CalonNISN     string
	AsalSekolah   string
	TujuanSekolah string
	Alasan        string
}

func NewMutationCreated(p MutationCreatedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	builder := types.NewBaseEventBuilder(EventMutationCreated, "Mutation", p.MutationID).
		WithEventID(uuid.New().String()).
		WithVersion(1)
	if ctx != nil {
		builder.WithContext(ctx)
	}
	return builder.WithPayload(map[string]interface{}{
		"mutationId":    p.MutationID,
		"studentId":     p.StudentID,
		"mutationType":  p.MutationType,
		"calonNama":     p.CalonNama,
		"calonNisn":     p.CalonNISN,
		"asalSekolah":   p.AsalSekolah,
		"tujuanSekolah": p.TujuanSekolah,
		"alasan":        p.Alasan,
	}).Build()
}

type MutationApprovedPayload struct {
	MutationID    string
	StudentID     *string
	MutationType  string
	CalonNama     string
	CalonNISN     string
	AsalSekolah   string
	TujuanSekolah string
	Catatan       string
}

func NewMutationApproved(p MutationApprovedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	builder := types.NewBaseEventBuilder(EventMutationApproved, "Mutation", p.MutationID).
		WithEventID(uuid.New().String()).
		WithVersion(1)
	if ctx != nil {
		builder.WithContext(ctx)
	}
	return builder.WithPayload(map[string]interface{}{
		"mutationId":    p.MutationID,
		"studentId":     p.StudentID,
		"mutationType":  p.MutationType,
		"calonNama":     p.CalonNama,
		"calonNisn":     p.CalonNISN,
		"asalSekolah":   p.AsalSekolah,
		"tujuanSekolah": p.TujuanSekolah,
		"catatan":       p.Catatan,
	}).Build()
}

type MutationRejectedPayload struct {
	MutationID string
	StudentID  *string
	Catatan    string
}

func NewMutationRejected(p MutationRejectedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	builder := types.NewBaseEventBuilder(EventMutationRejected, "Mutation", p.MutationID).
		WithEventID(uuid.New().String()).
		WithVersion(1)
	if ctx != nil {
		builder.WithContext(ctx)
	}
	return builder.WithPayload(map[string]interface{}{
		"mutationId": p.MutationID,
		"studentId":  p.StudentID,
		"catatan":    p.Catatan,
	}).Build()
}
