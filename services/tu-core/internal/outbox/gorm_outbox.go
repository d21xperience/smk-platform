// Package outbox menyediakan GORM implementation untuk outbox pattern.
package outbox

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	dbmodels "sekolah-platform/services/tu-core/internal/adapter/database/models"
	"sekolah-platform/platform/events/outbox"
	"sekolah-platform/platform/events/types"
)

// GormOutbox adalah GORM implementation dari outbox.Outbox.
type GormOutbox struct {
	db *gorm.DB
}

// NewGormOutbox membuat GormOutbox baru.
func NewGormOutbox(db *gorm.DB) *GormOutbox {
	return &GormOutbox{db: db}
}

// Save mengimplementasikan outbox.Outbox.
func (o *GormOutbox) Save(ctx context.Context, entry outbox.OutboxEntry) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	payload, err := json.Marshal(entry.Event.ToMap())
	if err != nil {
		return err
	}

	model := &dbmodels.OutboxEventModel{
		ID:          entry.ID,
		EventType:   entry.Event.GetEventName(),
		AggregateID: entry.Event.GetAggregateID(),
		Payload:     string(payload),
		Attempts:    entry.Attempts,
		LastError:   entry.LastError,
		NextRetry:   entry.NextRetry,
		CreatedAt:   entry.CreatedAt,
	}

	return o.db.WithContext(ctx).Create(model).Error
}

// FetchBatch mengimplementasikan outbox.Outbox.
func (o *GormOutbox) FetchBatch(ctx context.Context, batchSize int) ([]outbox.OutboxEntry, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var models []dbmodels.OutboxEventModel
	err := o.db.WithContext(ctx).
		Where("next_retry <= ? OR next_retry IS NULL", time.Now().UTC()).
		Order("created_at ASC").
		Limit(batchSize).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	entries := make([]outbox.OutboxEntry, 0, len(models))
	for _, m := range models {
		event, err := unmarshalEvent(m.Payload)
		if err != nil {
			continue // Skip invalid entries
		}

		entries = append(entries, outbox.OutboxEntry{
			ID:        m.ID,
			Event:     event,
			CreatedAt: m.CreatedAt,
			Attempts:  m.Attempts,
			LastError: m.LastError,
			NextRetry: m.NextRetry,
		})
	}

	return entries, nil
}

// Remove mengimplementasikan outbox.Outbox.
func (o *GormOutbox) Remove(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return o.db.WithContext(ctx).
		Delete(&dbmodels.OutboxEventModel{}, "id = ?", id).Error
}

// UpdateAttempt mengimplementasikan outbox.Outbox.
func (o *GormOutbox) UpdateAttempt(ctx context.Context, id string, attempts int, lastError string, nextRetry time.Time) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return o.db.WithContext(ctx).
		Model(&dbmodels.OutboxEventModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"attempts":   attempts,
			"last_error": lastError,
			"next_retry": nextRetry,
		}).Error
}

// Count mengimplementasikan outbox.Outbox.
func (o *GormOutbox) Count(ctx context.Context) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	var count int64
	err := o.db.WithContext(ctx).
		Model(&dbmodels.OutboxEventModel{}).
		Count(&count).Error
	return int(count), err
}

// unmarshalEvent deserialize event dari JSON payload.
func unmarshalEvent(payload string) (types.DomainEvent, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &raw); err != nil {
		return nil, err
	}

	event := &types.BaseEvent{}

	if v, ok := raw["eventId"].(string); ok {
		event.EventID = v
	} else {
		event.EventID = uuid.New().String()
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
	if v, ok := raw["aggregateId"].(string); ok {
		event.AggregateID = v
	}
	if v, ok := raw["aggregateType"].(string); ok {
		event.AggregateType = v
	}
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
	if v, ok := raw["correlationId"].(string); ok {
		event.CorrelationID = v
	}
	if v, ok := raw["traceId"].(string); ok {
		event.TraceID = v
	}
	if v, ok := raw["requestId"].(string); ok {
		event.RequestID = v
	}

	event.Payload = raw["payload"]

	return event, nil
}
