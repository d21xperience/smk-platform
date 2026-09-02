package events

import (
	"time"

	"github.com/google/uuid"
	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
)

const (
	EventSessionCreated      = "SessionCreated"
	EventSessionSubmitted    = "SessionSubmitted"
	EventSessionClosed       = "SessionClosed"
	EventAttendanceRecorded  = "AttendanceRecorded"
	EventAttendanceSubmitted = "AttendanceSubmitted"
	EventAttendanceCorrected = "AttendanceCorrected"
	EventAttendanceStatistics = "AttendanceStatistics"
)

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

// === SESSION EVENTS ===

type SessionCreatedPayload struct {
	SessionID string
	ClassID   string
	Date      time.Time
}

func NewSessionCreated(p SessionCreatedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventSessionCreated, "AttendanceSession", p.SessionID, ctx).
		WithPayload(map[string]interface{}{
			"sessionId": p.SessionID,
			"classId":   p.ClassID,
			"date":      p.Date.Format("2006-01-02"),
		}).Build()
}

type SessionSubmittedPayload struct {
	SessionID    string
	ClassID      string
	TotalRecords int
}

func NewSessionSubmitted(p SessionSubmittedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventSessionSubmitted, "AttendanceSession", p.SessionID, ctx).
		WithPayload(map[string]interface{}{
			"sessionId":    p.SessionID,
			"classId":      p.ClassID,
			"totalRecords": p.TotalRecords,
		}).Build()
}

type SessionClosedPayload struct {
	SessionID string
	ClassID   string
}

func NewSessionClosed(p SessionClosedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventSessionClosed, "AttendanceSession", p.SessionID, ctx).
		WithPayload(map[string]interface{}{
			"sessionId": p.SessionID,
			"classId":   p.ClassID,
		}).Build()
}

// === ATTENDANCE RECORD EVENTS ===

type AttendanceRecordPayload struct {
	SessionID string
	StudentID string
	Status    string
}

func NewAttendanceRecord(p AttendanceRecordPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventAttendanceRecorded, "AttendanceRecord", p.SessionID, ctx).
		WithPayload(map[string]interface{}{
			"sessionId": p.SessionID,
			"studentId": p.StudentID,
			"status":    p.Status,
		}).Build()
}

// === ATTENDANCE SUBMITTED (Bulk) ===

type AttendanceSubmittedPayload struct {
	SessionID    string
	ClassID      string
	TotalRecords int
	SubmittedBy  string
}

func NewAttendanceSubmitted(p AttendanceSubmittedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventAttendanceSubmitted, "AttendanceSession", p.SessionID, ctx).
		WithPayload(map[string]interface{}{
			"sessionId":    p.SessionID,
			"classId":      p.ClassID,
			"totalRecords": p.TotalRecords,
			"submittedBy":  p.SubmittedBy,
		}).Build()
}

// === ATTENDANCE CORRECTED ===

type AttendanceCorrectedPayload struct {
	SessionID  string
	StudentID  string
	OldStatus  string
	NewStatus  string
	CorrectedBy string
}

func NewAttendanceCorrected(p AttendanceCorrectedPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventAttendanceCorrected, "AttendanceRecord", p.SessionID, ctx).
		WithPayload(map[string]interface{}{
			"sessionId":   p.SessionID,
			"studentId":   p.StudentID,
			"oldStatus":   p.OldStatus,
			"newStatus":   p.NewStatus,
			"correctedBy": p.CorrectedBy,
		}).Build()
}

// === ATTENDANCE STATISTICS ===

type AttendanceStatisticsPayload struct {
	ClassID    string
	Month      int
	Year       int
	TotalDays  int
	Present    int
	Absent     int
	Late       int
	Sick       int
	Permission int
}

func NewAttendanceStatistics(p AttendanceStatisticsPayload, ctx *context.OperationalContext) *types.BaseEvent {
	return baseBuilder(EventAttendanceStatistics, "AttendanceStatistics", p.ClassID, ctx).
		WithPayload(map[string]interface{}{
			"classId":    p.ClassID,
			"month":      p.Month,
			"year":       p.Year,
			"totalDays":  p.TotalDays,
			"present":    p.Present,
			"absent":     p.Absent,
			"late":       p.Late,
			"sick":       p.Sick,
			"permission": p.Permission,
		}).Build()
}
