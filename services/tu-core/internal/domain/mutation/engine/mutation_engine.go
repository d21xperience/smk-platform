package engine

import (
	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/tu-core/internal/domain/mutation/models"
)

type Result struct {
	Mutation *models.Mutation
	Events   []types.DomainEvent
	Error    error
}

func Success(mutation *models.Mutation) Result {
	return Result{
		Mutation: mutation,
		Events:   mutation.UncommittedEvents(),
		Error:    nil,
	}
}

func Failure(err error) Result {
	return Result{
		Mutation: nil,
		Events:   nil,
		Error:    err,
	}
}

type MutationEngine struct{}

func NewMutationEngine() *MutationEngine {
	return &MutationEngine{}
}

func (e *MutationEngine) CreateMutation(data models.MutationData, ctx *context.OperationalContext) Result {
	m, err := models.NewMutation(data, ctx)
	if err != nil {
		return Failure(err)
	}
	return Success(m)
}

func (e *MutationEngine) ApproveMutation(m *models.Mutation, catatan string, ctx *context.OperationalContext) Result {
	if m == nil {
		return Failure(errMutationNotFound)
	}
	if err := m.Approve(catatan, ctx); err != nil {
		return Failure(err)
	}
	return Success(m)
}

func (e *MutationEngine) RejectMutation(m *models.Mutation, catatan string, ctx *context.OperationalContext) Result {
	if m == nil {
		return Failure(errMutationNotFound)
	}
	if err := m.Reject(catatan, ctx); err != nil {
		return Failure(err)
	}
	return Success(m)
}

var errMutationNotFound = errString("mutasi tidak ditemukan")

type errString string

func (e errString) Error() string { return string(e) }
