package events

import (
	"time"

	"github.com/google/uuid"
	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
)

// Event names
const (
	EventStudentCreated        = "StudentCreated"
	EventStudentEnrolled       = "StudentEnrolled"
	EventStudentProfileUpdated = "StudentProfileUpdated"
	EventStudentGraduated      = "StudentGraduated"
	EventStudentTransferred    = "StudentTransferred"
)

// === PAYLOAD STRUCTS (Pemisah dari models) ===

type StudentCreatedPayload struct {
	StudentID string
	NISN      string
	NIS       string
	FullName  string
	Gender    string
	BirthDate time.Time
}

type StudentEnrolledPayload struct {
	StudentID    string
	EnrollmentID string
	PeriodID     string
	ClassID      string
}

type StudentProfileUpdatedPayload struct {
	StudentID     string
	UpdatedFields []string
	UpdatedAt     time.Time
}

type StudentGraduatedPayload struct {
	StudentID      string
	GraduationDate time.Time
}

type StudentTransferredPayload struct {
	StudentID    string
	TargetSchool string
	Reason       string
	TransferDate time.Time
}

// === HELPER ===

func baseBuilder(eventName, aggregateType, aggregateID string, ctx *context.OperationalContext) *types.BaseEventBuilder {
	builder := types.NewBaseEventBuilder(eventName, aggregateType, aggregateID).
		WithEventID(uuid.New().String()).
		WithVersion(1)

	if ctx != nil {
		builder.WithContext(ctx)
	}

	return builder
}

// === EVENT FACTORIES ===

func NewStudentCreated(p StudentCreatedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventStudentCreated, "Student", p.StudentID, ctx).
		WithPayload(map[string]interface{}{
			"studentId": p.StudentID,
			"nisn":      p.NISN,
			"nis":       p.NIS,
			"fullName":  p.FullName,
			"gender":    p.Gender,
			"birthDate": p.BirthDate.Format("2006-01-02"),
		}).
		Build()
}

func NewStudentEnrolled(p StudentEnrolledPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventStudentEnrolled, "Student", p.StudentID, ctx).
		WithPayload(map[string]interface{}{
			"studentId":    p.StudentID,
			"enrollmentId": p.EnrollmentID,
			"periodId":     p.PeriodID,
			"classId":      p.ClassID,
		}).
		Build()
}

func NewStudentProfileUpdated(p StudentProfileUpdatedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventStudentProfileUpdated, "Student", p.StudentID, ctx).
		WithPayload(map[string]interface{}{
			"studentId":     p.StudentID,
			"updatedFields": p.UpdatedFields,
			"updatedAt":     p.UpdatedAt.Format(time.RFC3339),
		}).
		Build()
}

func NewStudentGraduated(p StudentGraduatedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventStudentGraduated, "Student", p.StudentID, ctx).
		WithPayload(map[string]interface{}{
			"studentId":      p.StudentID,
			"graduationDate": p.GraduationDate.Format(time.RFC3339),
		}).
		Build()
}

func NewStudentTransferred(p StudentTransferredPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventStudentTransferred, "Student", p.StudentID, ctx).
		WithPayload(map[string]interface{}{
			"studentId":    p.StudentID,
			"targetSchool": p.TargetSchool,
			"reason":       p.Reason,
			"transferDate": p.TransferDate.Format(time.RFC3339),
		}).
		Build()
}
