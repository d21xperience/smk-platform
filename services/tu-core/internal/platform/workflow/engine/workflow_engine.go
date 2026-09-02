package engine

import (
	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/tu-core/internal/platform/workflow/models"
)

// Result adalah hasil dari engine operation.
type Result struct {
	Instance *models.WorkflowInstance
	Events   []types.DomainEvent
	Error    error
}

func Success(instance *models.WorkflowInstance) Result {
	return Result{
		Instance: instance,
		Events:   instance.UncommittedEvents(),
		Error:    nil,
	}
}

func Failure(err error) Result {
	return Result{Instance: nil, Events: nil, Error: err}
}

// WorkflowEngine adalah business engine untuk Workflow.
type WorkflowEngine struct{}

func NewWorkflowEngine() *WorkflowEngine {
	return &WorkflowEngine{}
}

// StartWorkflow memulai workflow baru.
func (e *WorkflowEngine) StartWorkflow(
	definition *models.WorkflowDefinition,
	referenceID, referenceType, requestedBy string,
	ctx *context.OperationalContext,
) Result {
	if definition == nil {
		return Failure(errDefinitionNotFound)
	}
	if !definition.IsActive() {
		return Failure(errDefinitionInactive)
	}

	data := models.WorkflowInstanceData{
		DefinitionID:  definition.DefinitionID(),
		ReferenceID:   referenceID,
		ReferenceType: referenceType,
		RequestedBy:   requestedBy,
	}

	instance, err := models.NewWorkflowInstance(data, definition, ctx)
	if err != nil {
		return Failure(err)
	}

	return Success(instance)
}

// CompleteStep menyelesaikan step saat ini.
func (e *WorkflowEngine) CompleteStep(
	instance *models.WorkflowInstance,
	action, comment, userID string,
	ctx *context.OperationalContext,
) Result {
	if instance == nil {
		return Failure(errInstanceNotFound)
	}

	if err := instance.CompleteStep(action, comment, userID, ctx); err != nil {
		return Failure(err)
	}

	return Success(instance)
}

// RejectStep menolak step saat ini.
func (e *WorkflowEngine) RejectStep(
	instance *models.WorkflowInstance,
	reason, userID string,
	ctx *context.OperationalContext,
) Result {
	if instance == nil {
		return Failure(errInstanceNotFound)
	}

	if err := instance.RejectStep(reason, userID, ctx); err != nil {
		return Failure(err)
	}

	return Success(instance)
}

// CancelWorkflow membatalkan workflow.
func (e *WorkflowEngine) CancelWorkflow(
	instance *models.WorkflowInstance,
	reason, userID string,
	ctx *context.OperationalContext,
) Result {
	if instance == nil {
		return Failure(errInstanceNotFound)
	}

	if err := instance.Cancel(reason, userID, ctx); err != nil {
		return Failure(err)
	}

	return Success(instance)
}

var (
	errDefinitionNotFound = errString("workflow definition tidak ditemukan")
	errDefinitionInactive = errString("workflow definition tidak aktif")
	errInstanceNotFound   = errString("workflow instance tidak ditemukan")
)

type errString string

func (e errString) Error() string { return string(e) }
