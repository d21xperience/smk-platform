package engine

import (
	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/tu-core/internal/domain/assessment/models"
)

type Result struct {
	Assessment *models.Assessment
	Events     []types.DomainEvent
	Error      error
}

func Success(assessment *models.Assessment) Result {
	return Result{
		Assessment: assessment,
		Events:     assessment.UncommittedEvents(),
		Error:      nil,
	}
}

func Failure(err error) Result {
	return Result{
		Assessment: nil,
		Events:     nil,
		Error:      err,
	}
}

type AssessmentEngine struct{}

func NewAssessmentEngine() *AssessmentEngine {
	return &AssessmentEngine{}
}

func (e *AssessmentEngine) CreateAssessment(data models.AssessmentData, ctx *context.OperationalContext) Result {
	a, err := models.NewAssessment(data, ctx)
	if err != nil {
		return Failure(err)
	}
	return Success(a)
}

func (e *AssessmentEngine) UpdateAssessment(a *models.Assessment, score, maxScore float64, notes string) Result {
	if a == nil {
		return Failure(errAssessmentNotFound)
	}
	if err := a.Update(score, maxScore, notes); err != nil {
		return Failure(err)
	}
	return Success(a)
}

var errAssessmentNotFound = errString("assessment tidak ditemukan")

type errString string

func (e errString) Error() string { return string(e) }
