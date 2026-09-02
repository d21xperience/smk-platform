package nats

import (
	"encoding/json"
	"fmt"
	"time"

	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
)

// Subject untuk event.
const SubjectPrefix = "sdp.events."

// EventToSubject mengkonversi event name ke NATS subject.
// Contoh: "StudentCreated" → "sdp.events.student.created"
func EventToSubject(eventName string) string {
	// Convert CamelCase ke dot-separated lowercase
	result := make([]byte, 0, len(eventName)+10)
	for i, r := range eventName {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				result = append(result, '.')
			}
			result = append(result, byte(r-'A'+'a'))
		} else {
			result = append(result, byte(r))
		}
	}
	return SubjectPrefix + string(result)
}

// SubjectToEventName mengkonversi NATS subject ke event name.
// Contoh: "sdp.events.student.created" → "StudentCreated"
func SubjectToEventName(subject string) string {
	if len(subject) <= len(SubjectPrefix) {
		return subject
	}
	rest := subject[len(SubjectPrefix):]
	// Convert dot-separated ke CamelCase
	result := make([]byte, 0, len(rest))
	capitalizeNext := true
	for _, r := range rest {
		if r == '.' {
			capitalizeNext = true
			continue
		}
		if capitalizeNext {
			if r >= 'a' && r <= 'z' {
				result = append(result, byte(r-'a'+'A'))
			} else {
				result = append(result, byte(r))
			}
			capitalizeNext = false
		} else {
			result = append(result, byte(r))
		}
	}
	return string(result)
}

// EventToBytes mengkonversi DomainEvent ke bytes untuk NATS message.
func EventToBytes(event types.DomainEvent) ([]byte, error) {
	data := event.ToMap()
	return json.Marshal(data)
}

// BytesToEvent mengkonversi NATS message bytes ke DomainEvent.
func BytesToEvent(data []byte) (*types.BaseEvent, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("gagal unmarshal event: %w", err)
	}

	event := &types.BaseEvent{}

	// Extract metadata
	if v, ok := raw["eventId"].(string); ok {
		event.EventID = v
	}
	if v, ok := raw["eventName"].(string); ok {
		event.EventName = v
	}
	if v, ok := raw["version"].(float64); ok {
		event.Version = int(v)
	}
	if v, ok := raw["occurredAt"].(string); ok {
		if t, err := time.Parse(time.RFC3339Nano, v); err == nil {
			event.OccurredAt = t
		}
	}

	// Extract aggregate
	if v, ok := raw["aggregateId"].(string); ok {
		event.AggregateID = v
	}
	if v, ok := raw["aggregateType"].(string); ok {
		event.AggregateType = v
	}

	// Extract operational context
	if v, ok := raw["schoolId"].(string); ok {
		event.SchoolID = v
	}
	if v, ok := raw["academicYear"].(string); ok {
		event.AcademicYear = v
	}
	if v, ok := raw["semester"].(float64); ok {
		event.Semester = int(v)
	}
	if v, ok := raw["academicPeriodId"].(string); ok {
		event.AcademicPeriodID = v
	}
	if v, ok := raw["tenantId"].(string); ok {
		event.TenantID = v
	}
	if v, ok := raw["userId"].(string); ok {
		event.UserID = v
	}

	// Extract observability
	if v, ok := raw["correlationId"].(string); ok {
		event.CorrelationID = v
	}
	if v, ok := raw["traceId"].(string); ok {
		event.TraceID = v
	}
	if v, ok := raw["requestId"].(string); ok {
		event.RequestID = v
	}

	// Extract payload
	event.Payload = raw["payload"]

	return event, nil
}

// EventToContext mengekstrak OperationalContext dari event.
func EventToContext(event types.DomainEvent) *context.OperationalContext {
	return &context.OperationalContext{
		SchoolID:         event.GetSchoolID(),
		AcademicPeriodID: event.GetAcademicPeriodID(),
		UserID:           event.GetUserID(),
	}
}
