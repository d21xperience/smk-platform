package models

import (
	"errors"
	"time"

	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/tu-core/internal/domain/assessment/events"
	"sekolah-platform/services/tu-core/internal/domain/assessment/valueobjects"

	"github.com/google/uuid"
)

type Assessment struct {
	ID               string
	SchoolID         string
	AcademicPeriodID string
	StudentID        string
	SubjectID        string
	AssessmentType   valueobjects.AssessmentType
	Score            float64
	MaxScore         float64
	AssessmentDate   time.Time
	Semester         int
	Notes            string
	CreatedBy        string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	domainEvents     []types.DomainEvent
}

type AssessmentData struct {
	ID             string
	StudentID      string
	SubjectID      string
	AssessmentType valueobjects.AssessmentType
	Score          float64
	MaxScore       float64
	AssessmentDate time.Time
	Semester       int
	Notes          string
}

func NewAssessment(data AssessmentData, ctx *context.OperationalContext) (*Assessment, error) {
	if data.ID == "" {
		data.ID = uuid.New().String()
	}

	// Validasi
	if data.StudentID == "" {
		return nil, errors.New("student_id wajib diisi")
	}
	if data.SubjectID == "" {
		return nil, errors.New("subject_id wajib diisi")
	}
	if data.Score < 0 {
		return nil, errors.New("score tidak boleh negatif")
	}
	if data.MaxScore <= 0 {
		return nil, errors.New("max_score harus lebih dari 0")
	}
	if data.Score > data.MaxScore {
		return nil, errors.New("score tidak boleh melebihi max_score")
	}
	if data.AssessmentDate.IsZero() {
		return nil, errors.New("assessment_date wajib diisi")
	}
	if data.Semester != 1 && data.Semester != 2 {
		return nil, errors.New("semester harus 1 atau 2")
	}

	a := &Assessment{
		ID:               data.ID,
		SchoolID:         ctx.SchoolID,
		AcademicPeriodID: ctx.AcademicPeriodID,
		StudentID:        data.StudentID,
		SubjectID:        data.SubjectID,
		AssessmentType:   data.AssessmentType,
		Score:            data.Score,
		MaxScore:         data.MaxScore,
		AssessmentDate:   data.AssessmentDate,
		Semester:         data.Semester,
		Notes:            data.Notes,
		CreatedBy:        ctx.UserID,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
		domainEvents:     []types.DomainEvent{},
	}

	a.recordEvent(events.NewAssessmentCreated(events.AssessmentCreatedPayload{
		AssessmentID:   a.ID,
		StudentID:      a.StudentID,
		SubjectID:      a.SubjectID,
		AssessmentType: a.AssessmentType.String(),
		Score:          a.Score,
		MaxScore:       a.MaxScore,
		Semester:       a.Semester,
	}, ctx))

	return a, nil
}

func (a *Assessment) Update(score float64, maxScore float64, notes string) error {
	if score < 0 {
		return errors.New("score tidak boleh negatif")
	}
	if maxScore <= 0 {
		return errors.New("max_score harus lebih dari 0")
	}
	if score > maxScore {
		return errors.New("score tidak boleh melebihi max_score")
	}

	oldScore := a.Score
	a.Score = score
	a.MaxScore = maxScore
	a.Notes = notes
	a.UpdatedAt = time.Now().UTC()

	a.recordEvent(events.NewAssessmentUpdated(events.AssessmentUpdatedPayload{
		AssessmentID: a.ID,
		StudentID:    a.StudentID,
		SubjectID:    a.SubjectID,
		OldScore:     oldScore,
		NewScore:     score,
	}))

	return nil
}

func (a *Assessment) Percentage() float64 {
	if a.MaxScore == 0 {
		return 0
	}
	return (a.Score / a.MaxScore) * 100
}

func (a *Assessment) Grade() string {
	return valueobjects.CalculateGrade(a.Percentage())
}

func (a *Assessment) recordEvent(event types.DomainEvent) {
	a.domainEvents = append(a.domainEvents, event)
}

func (a *Assessment) UncommittedEvents() []types.DomainEvent {
	return a.domainEvents
}

func (a *Assessment) ClearEvents() {
	a.domainEvents = []types.DomainEvent{}
}
