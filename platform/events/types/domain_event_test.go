package types

import (
	"testing"
	"time"

	platformctx "sekolah-platform/platform/context"
)

func TestBaseEventBuilder(t *testing.T) {
	ctx := &platformctx.OperationalContext{
		SchoolID:         "smk-001",
		AcademicYear:     "2026/2027",
		Semester:         1,
		AcademicPeriodID: "20261",
		UserID:           "user-123",
	}

	payload := map[string]interface{}{
		"studentId": "student-001",
		"nisn":      "1234567890",
	}

	event := NewBaseEventBuilder("StudentCreated", "Student", "student-001").
		WithEventID("evt-001").
		WithVersion(1).
		WithPayload(payload).
		WithContext(ctx).
		WithObservability("corr-001", "trace-001", "req-001").
		Build()

	if event.GetEventName() != "StudentCreated" {
		t.Errorf("Expected eventName 'StudentCreated', got '%s'", event.GetEventName())
	}
	if event.GetAggregateType() != "Student" {
		t.Errorf("Expected aggregateType 'Student', got '%s'", event.GetAggregateType())
	}
	if event.GetSchoolID() != "smk-001" {
		t.Errorf("Expected schoolId 'smk-001', got '%s'", event.GetSchoolID())
	}
	if event.GetAcademicPeriodID() != "20261" {
		t.Errorf("Expected periodId '20261', got '%s'", event.GetAcademicPeriodID())
	}
	if event.GetCorrelationID() != "corr-001" {
		t.Errorf("Expected correlationId 'corr-001', got '%s'", event.GetCorrelationID())
	}
	if event.GetVersion() != 1 {
		t.Errorf("Expected version 1, got %d", event.GetVersion())
	}
}

func TestBaseEventToMap(t *testing.T) {
	event := &BaseEvent{
		EventID:          "evt-001",
		EventName:        "StudentCreated",
		Version:          1,
		OccurredAt:       time.Date(2026, 7, 29, 10, 0, 0, 0, time.UTC),
		AggregateID:      "student-001",
		AggregateType:    "Student",
		SchoolID:         "smk-001",
		AcademicPeriodID: "20261",
		UserID:           "user-123",
		CorrelationID:    "corr-001",
		TraceID:          "trace-001",
		RequestID:        "req-001",
		Payload:          map[string]string{"name": "Budi"},
	}

	m := event.ToMap()

	if m["eventName"] != "StudentCreated" {
		t.Error("Expected eventName in map")
	}
	if m["schoolId"] != "smk-001" {
		t.Error("Expected schoolId in map")
	}
	if m["academicPeriodId"] != "20261" {
		t.Error("Expected periodId in map")
	}
	if m["payload"] == nil {
		t.Error("Expected payload in map")
	}
}

func TestBaseEventBuilderNilContext(t *testing.T) {
	event := NewBaseEventBuilder("TestEvent", "Test", "id-1").
		WithContext(nil).
		Build()

	if event.GetSchoolID() != "" {
		t.Error("Expected empty schoolId with nil context")
	}
}
