package engine

import (
	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/tu-core/internal/domain/registration/models"
)

type Result struct {
	Registration *models.Registration
	Events       []types.DomainEvent
	Error        error
}

func Success(registration *models.Registration) Result {
	return Result{
		Registration: registration,
		Events:       registration.UncommittedEvents(),
		Error:        nil,
	}
}

func Failure(err error) Result {
	return Result{
		Registration: nil,
		Events:       nil,
		Error:        err,
	}
}

type RegistrationEngine struct{}

func NewRegistrationEngine() *RegistrationEngine {
	return &RegistrationEngine{}
}

func (e *RegistrationEngine) CreateRegistration(data models.RegistrationData, ctx *context.OperationalContext) Result {
	r, err := models.NewRegistration(data, ctx)
	if err != nil {
		return Failure(err)
	}
	return Success(r)
}

func (e *RegistrationEngine) SubmitRegistration(r *models.Registration, ctx *context.OperationalContext) Result {
	if r == nil {
		return Failure(errRegistrationNotFound)
	}
	if err := r.Submit(ctx); err != nil {
		return Failure(err)
	}
	return Success(r)
}

func (e *RegistrationEngine) VerifyRegistration(r *models.Registration, catatan, verifiedBy string, ctx *context.OperationalContext) Result {
	if r == nil {
		return Failure(errRegistrationNotFound)
	}
	if err := r.Verify(catatan, verifiedBy, ctx); err != nil {
		return Failure(err)
	}
	return Success(r)
}

func (e *RegistrationEngine) ApproveRegistration(r *models.Registration, catatan, approvedBy string, ctx *context.OperationalContext) Result {
	if r == nil {
		return Failure(errRegistrationNotFound)
	}
	if err := r.Approve(catatan, approvedBy, ctx); err != nil {
		return Failure(err)
	}
	return Success(r)
}

func (e *RegistrationEngine) RejectRegistration(r *models.Registration, alasan, rejectedBy string, ctx *context.OperationalContext) Result {
	if r == nil {
		return Failure(errRegistrationNotFound)
	}
	if err := r.Reject(alasan, rejectedBy, ctx); err != nil {
		return Failure(err)
	}
	return Success(r)
}

var errRegistrationNotFound = errString("pendaftaran tidak ditemukan")

type errString string

func (e errString) Error() string { return string(e) }
