package events

import (
	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"

	"github.com/google/uuid"
)

const (
	EventWorkflowStarted   = "WorkflowStarted"
	EventStepCompleted     = "StepCompleted"
	EventWorkflowApproved  = "WorkflowApproved"
	EventWorkflowRejected  = "WorkflowRejected"
	EventWorkflowCancelled = "WorkflowCancelled"
)

// === Payload Structs ===

type WorkflowStartedPayload struct {
	InstanceID    string
	DefinitionID  string
	ReferenceID   string
	ReferenceType string
	RequestedBy   string
	TotalSteps    int
}

type StepCompletedPayload struct {
	InstanceID  string
	StepID      string
	StepName    string
	StepOrder   int
	Action      string
	Comment     string
	CompletedBy string
	CompletedAt string
}

type WorkflowApprovedPayload struct {
	InstanceID    string
	ReferenceID   string
	ReferenceType string
	ApprovedAt    string
}

type WorkflowRejectedPayload struct {
	InstanceID  string
	ReferenceID string
	Reason      string
	RejectedAt  string
}

type WorkflowCancelledPayload struct {
	InstanceID  string
	ReferenceID string
	Reason      string
	CancelledAt string
}

// === Helper ===

func baseBuilder(eventName, aggregateID string, ctx *context.OperationalContext) *types.BaseEventBuilder {
	builder := types.NewBaseEventBuilder(eventName, "WorkflowInstance", aggregateID).
		WithEventID(uuid.New().String()).
		WithVersion(1)
	if ctx != nil {
		builder.WithContext(ctx)
	}
	return builder
}

// === Event Factories ===

func NewWorkflowStarted(p WorkflowStartedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventWorkflowStarted, p.InstanceID, ctx).
		WithPayload(map[string]interface{}{
			"instanceId":    p.InstanceID,
			"definitionId":  p.DefinitionID,
			"referenceId":   p.ReferenceID,
			"referenceType": p.ReferenceType,
			"requestedBy":   p.RequestedBy,
			"totalSteps":    p.TotalSteps,
		}).
		Build()
}

func NewStepCompleted(p StepCompletedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventStepCompleted, p.InstanceID, ctx).
		WithPayload(map[string]interface{}{
			"instanceId":  p.InstanceID,
			"stepId":      p.StepID,
			"stepName":    p.StepName,
			"stepOrder":   p.StepOrder,
			"action":      p.Action,
			"comment":     p.Comment,
			"completedBy": p.CompletedBy,
			"completedAt": p.CompletedAt,
		}).
		Build()
}

func NewWorkflowApproved(p WorkflowApprovedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventWorkflowApproved, p.InstanceID, ctx).
		WithPayload(map[string]interface{}{
			"instanceId":    p.InstanceID,
			"referenceId":   p.ReferenceID,
			"referenceType": p.ReferenceType,
			"approvedAt":    p.ApprovedAt,
		}).
		Build()
}

func NewWorkflowRejected(p WorkflowRejectedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventWorkflowRejected, p.InstanceID, ctx).
		WithPayload(map[string]interface{}{
			"instanceId":  p.InstanceID,
			"referenceId": p.ReferenceID,
			"reason":      p.Reason,
			"rejectedAt":  p.RejectedAt,
		}).
		Build()
}

func NewWorkflowCancelled(p WorkflowCancelledPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventWorkflowCancelled, p.InstanceID, ctx).
		WithPayload(map[string]interface{}{
			"instanceId":  p.InstanceID,
			"referenceId": p.ReferenceID,
			"reason":      p.Reason,
			"cancelledAt": p.CancelledAt,
		}).
		Build()
}
