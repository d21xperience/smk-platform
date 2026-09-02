package models

import (
	"errors"
	"time"

	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
	workflowevents "sekolah-platform/services/tu-core/internal/platform/workflow/events"

	"github.com/google/uuid"
)

// WorkflowStatus adalah status workflow instance.
type WorkflowStatus string

const (
	WorkflowStatusPending    WorkflowStatus = "PENDING"
	WorkflowStatusInProgress WorkflowStatus = "IN_PROGRESS"
	WorkflowStatusApproved   WorkflowStatus = "APPROVED"
	WorkflowStatusRejected   WorkflowStatus = "REJECTED"
	WorkflowStatusCancelled  WorkflowStatus = "CANCELLED"
	WorkflowStatusTimeout    WorkflowStatus = "TIMEOUT"
)

// WorkflowInstance adalah Aggregate Root untuk running workflow.
type WorkflowInstance struct {
	instanceID    string
	definitionID  string
	referenceID   string // ID objek yang di-approve (e.g., documentId)
	referenceType string // Tipe objek (e.g., "Document")
	status        WorkflowStatus
	currentStep   int
	steps         []*WorkflowStep
	requestedBy   string
	events        []types.DomainEvent
	createdAt     time.Time
	updatedAt     time.Time
	completedAt   *time.Time
}

// WorkflowInstanceData adalah data untuk membuat WorkflowInstance.
type WorkflowInstanceData struct {
	InstanceID    string
	DefinitionID  string
	ReferenceID   string
	ReferenceType string
	RequestedBy   string
}

// NewWorkflowInstance membuat WorkflowInstance baru dari definition.
func NewWorkflowInstance(data WorkflowInstanceData, definition *WorkflowDefinition, ctx *context.OperationalContext) (*WorkflowInstance, error) {
	if data.InstanceID == "" {
		data.InstanceID = uuid.New().String()
	}

	if definition == nil {
		return nil, errors.New("workflow definition tidak boleh nil")
	}

	// Create step instances from definition
	steps := make([]*WorkflowStep, 0, len(definition.Steps()))
	for _, stepDef := range definition.Steps() {
		step := NewWorkflowStep(stepDef, data.InstanceID)
		steps = append(steps, step)
	}

	instance := &WorkflowInstance{
		instanceID:    data.InstanceID,
		definitionID:  data.DefinitionID,
		referenceID:   data.ReferenceID,
		referenceType: data.ReferenceType,
		status:        WorkflowStatusInProgress,
		currentStep:   1,
		steps:         steps,
		requestedBy:   data.RequestedBy,
		events:        []types.DomainEvent{},
		createdAt:     time.Now().UTC(),
		updatedAt:     time.Now().UTC(),
	}

	instance.recordEvent(workflowevents.NewWorkflowStarted(workflowevents.WorkflowStartedPayload{
		InstanceID:    instance.instanceID,
		DefinitionID:  instance.definitionID,
		ReferenceID:   instance.referenceID,
		ReferenceType: instance.referenceType,
		RequestedBy:   instance.requestedBy,
		TotalSteps:    len(instance.steps),
	}, ctx))

	return instance, nil
}

// Getters

func (i *WorkflowInstance) InstanceID() string      { return i.instanceID }
func (i *WorkflowInstance) DefinitionID() string    { return i.definitionID }
func (i *WorkflowInstance) ReferenceID() string     { return i.referenceID }
func (i *WorkflowInstance) ReferenceType() string   { return i.referenceType }
func (i *WorkflowInstance) Status() WorkflowStatus  { return i.status }
func (i *WorkflowInstance) CurrentStep() int        { return i.currentStep }
func (i *WorkflowInstance) Steps() []*WorkflowStep  { return i.steps }
func (i *WorkflowInstance) RequestedBy() string     { return i.requestedBy }
func (i *WorkflowInstance) CreatedAt() time.Time    { return i.createdAt }
func (i *WorkflowInstance) UpdatedAt() time.Time    { return i.updatedAt }
func (i *WorkflowInstance) CompletedAt() *time.Time { return i.completedAt }

// GetCurrentStepInstance mengambil step instance saat ini.
func (i *WorkflowInstance) GetCurrentStepInstance() *WorkflowStep {
	for _, step := range i.steps {
		if step.Order() == i.currentStep {
			return step
		}
	}
	return nil
}

// CompleteStep menyelesaikan step saat ini.
func (i *WorkflowInstance) CompleteStep(action, comment, userID string, ctx *context.OperationalContext) error {
	if i.status != WorkflowStatusInProgress {
		return errors.New("workflow sudah selesai (status: " + string(i.status) + ")")
	}

	currentStep := i.GetCurrentStepInstance()
	if currentStep == nil {
		return errors.New("step saat ini tidak ditemukan")
	}

	// Complete step
	if err := currentStep.Complete(action, comment, userID); err != nil {
		return err
	}

	// Generate StepCompleted event
	// i.recordEvent(workflowevents.NewStepCompleted(i, currentStep, ctx))
	i.recordEvent(workflowevents.NewStepCompleted(workflowevents.StepCompletedPayload{
		InstanceID:  i.instanceID,
		StepID:      currentStep.StepID(),
		StepName:    currentStep.Name(),
		StepOrder:   currentStep.Order(),
		Action:      currentStep.Action(),
		Comment:     currentStep.Comment(),
		CompletedBy: currentStep.CompletedBy(),
		CompletedAt: time.Now().UTC().Format(time.RFC3339),
	}, ctx))
	// Check if workflow completed
	if i.currentStep >= len(i.steps) {
		i.status = WorkflowStatusApproved
		now := time.Now().UTC()
		i.completedAt = &now
		i.recordEvent(workflowevents.NewWorkflowApproved(workflowevents.WorkflowApprovedPayload{
			InstanceID:    i.instanceID,
			ReferenceID:   i.referenceID,
			ReferenceType: i.referenceType,
			ApprovedAt:    time.Now().UTC().Format(time.RFC3339),
		}, ctx))
	} else {
		i.currentStep++
	}

	return nil
}

// RejectStep menolak step saat ini.
func (i *WorkflowInstance) RejectStep(reason, userID string, ctx *context.OperationalContext) error {
	if i.status != WorkflowStatusInProgress {
		return errors.New("workflow sudah selesai")
	}

	currentStep := i.GetCurrentStepInstance()
	if currentStep == nil {
		return errors.New("step saat ini tidak ditemukan")
	}

	if err := currentStep.Reject(reason, userID); err != nil {
		return err
	}

	i.status = WorkflowStatusRejected
	now := time.Now().UTC()
	i.completedAt = &now
	i.updatedAt = now

	i.recordEvent(workflowevents.NewWorkflowRejected(workflowevents.WorkflowRejectedPayload{
		InstanceID:  i.instanceID,
		ReferenceID: i.referenceID,
		Reason:      reason,
		RejectedAt:  time.Now().UTC().Format(time.RFC3339),
	}, ctx))

	return nil
}

// Cancel membatalkan workflow.
func (i *WorkflowInstance) Cancel(reason, userID string, ctx *context.OperationalContext) error {
	if i.status != WorkflowStatusInProgress && i.status != WorkflowStatusPending {
		return errors.New("workflow sudah selesai")
	}

	i.status = WorkflowStatusCancelled
	now := time.Now().UTC()
	i.completedAt = &now
	i.updatedAt = now

	i.recordEvent(workflowevents.NewWorkflowCancelled(workflowevents.WorkflowCancelledPayload{
		InstanceID:  i.instanceID,
		ReferenceID: i.referenceID,
		Reason:      reason,
		CancelledAt: time.Now().UTC().Format(time.RFC3339),
	}, ctx))
	return nil
}

// IsCompleted mengecek apakah workflow sudah selesai.
func (i *WorkflowInstance) IsCompleted() bool {
	return i.status == WorkflowStatusApproved ||
		i.status == WorkflowStatusRejected ||
		i.status == WorkflowStatusCancelled ||
		i.status == WorkflowStatusTimeout
}

// Event management

func (i *WorkflowInstance) recordEvent(event types.DomainEvent) {
	i.events = append(i.events, event)
}

func (i *WorkflowInstance) UncommittedEvents() []types.DomainEvent {
	return i.events
}

func (i *WorkflowInstance) ClearEvents() {
	i.events = []types.DomainEvent{}
}
