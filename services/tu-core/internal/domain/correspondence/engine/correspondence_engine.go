package engine

import (
	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/tu-core/internal/domain/correspondence/models"
)

type Result struct {
	Correspondence *models.Correspondence
	Events         []types.DomainEvent
	Error          error
}

func Success(c *models.Correspondence) Result {
	return Result{
		Correspondence: c,
		Events:         c.UncommittedEvents(),
		Error:          nil,
	}
}

func Failure(err error) Result {
	return Result{
		Correspondence: nil,
		Events:         nil,
		Error:          err,
	}
}

type CorrespondenceEngine struct{}

func NewCorrespondenceEngine() *CorrespondenceEngine {
	return &CorrespondenceEngine{}
}

func (e *CorrespondenceEngine) CreateCorrespondence(data models.CorrespondenceData, ctx *context.OperationalContext) Result {
	c, err := models.NewCorrespondence(data, ctx)
	if err != nil {
		return Failure(err)
	}
	return Success(c)
}

func (e *CorrespondenceEngine) ProcessCorrespondence(c *models.Correspondence, ctx *context.OperationalContext) Result {
	if c == nil {
		return Failure(errCorrespondenceNotFound)
	}
	if err := c.Process(ctx); err != nil {
		return Failure(err)
	}
	return Success(c)
}

func (e *CorrespondenceEngine) ArchiveCorrespondence(c *models.Correspondence, ctx *context.OperationalContext) Result {
	if c == nil {
		return Failure(errCorrespondenceNotFound)
	}
	if err := c.Archive(ctx); err != nil {
		return Failure(err)
	}
	return Success(c)
}

var errCorrespondenceNotFound = errString("surat tidak ditemukan")

type errString string

func (e errString) Error() string { return string(e) }
