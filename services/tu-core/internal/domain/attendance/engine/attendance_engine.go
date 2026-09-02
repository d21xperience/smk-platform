package engine

import (
	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/tu-core/internal/domain/attendance/models"
)

type Result struct {
	Session *models.Session
	Events  []types.DomainEvent
	Error   error
}

func Success(session *models.Session) Result {
	return Result{
		Session: session,
		Events:  session.UncommittedEvents(),
		Error:   nil,
	}
}

func Failure(err error) Result {
	return Result{
		Session: nil,
		Events:  nil,
		Error:   err,
	}
}

type AttendanceEngine struct{}

func NewAttendanceEngine() *AttendanceEngine {
	return &AttendanceEngine{}
}

func (e *AttendanceEngine) CreateSession(data models.SessionData, ctx *context.OperationalContext) Result {
	s, err := models.NewSession(data, ctx)
	if err != nil {
		return Failure(err)
	}
	return Success(s)
}

func (e *AttendanceEngine) SubmitSession(s *models.Session, ctx *context.OperationalContext) Result {
	if s == nil {
		return Failure(errSessionNotFound)
	}
	if err := s.Submit(ctx); err != nil {
		return Failure(err)
	}
	return Success(s)
}

func (e *AttendanceEngine) CloseSession(s *models.Session, ctx *context.OperationalContext) Result {
	if s == nil {
		return Failure(errSessionNotFound)
	}
	if err := s.Close(ctx); err != nil {
		return Failure(err)
	}
	return Success(s)
}

var errSessionNotFound = errString("session tidak ditemukan")

type errString string

func (e errString) Error() string { return string(e) }
