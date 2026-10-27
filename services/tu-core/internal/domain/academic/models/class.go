package models

import (
	"errors"
	"time"

	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/tu-core/internal/domain/academic/events"

	"github.com/google/uuid"
)

// =============================================================================
// AGGREGATE ROOT
// =============================================================================

// Class adalah Aggregate Root untuk kelas/rombel.
type Class struct {
	id         string
	schoolID   string
	periodID   string // Academic Period ID (e.g. "20261")
	name       string // e.g. "X IPA 1", "XI RPL 2"
	grade      int    // 10, 11, 12
	level      string // "X", "XI", "XII"
	major      string // "IPA", "RPL", "TKJ", dll
	homeroomID string // Teacher ID sebagai wali kelas
	capacity   int
	status     string // "active" | "inactive" (tambahan)
	events     []types.DomainEvent
	createdAt  time.Time
	UpdatedAt  time.Time
}

// IsActive returns true jika status active
func (c *Class) IsActive() bool {
	return c.status == "active"
}

// CanEnroll memeriksa apakah kapasitas masih tersedia
func (c *Class) CanEnroll(currentCount int) bool {
	return c.IsActive() && currentCount < c.capacity
}

// =============================================================================
// COMMAND DATA (Input DTO)
// =============================================================================

// ClassData adalah data untuk membuat Class.
type ClassData struct {
	ID         string
	SchoolID   string
	PeriodID   string
	Name       string
	Grade      int
	Level      string
	Major      string
	HomeroomID string
	Capacity   int
}

// =============================================================================
// FACTORY & BEHAVIORS
// =============================================================================

// NewClass membuat Class baru.
func NewClass(data ClassData, ctx *context.OperationalContext) (*Class, error) {
	if data.ID == "" {
		data.ID = uuid.New().String()
	}

	if data.Name == "" {
		return nil, errors.New("nama kelas wajib diisi")
	}
	if data.PeriodID == "" {
		return nil, errors.New("periodID wajib diisi")
	}
	if data.Grade < 10 || data.Grade > 12 {
		return nil, errors.New("grade harus 10, 11, atau 12")
	}
	if data.Capacity <= 0 {
		return nil, errors.New("capacity harus lebih dari 0")
	}

	c := &Class{
		id:         data.ID,
		schoolID:   ctx.SchoolID,
		periodID:   data.PeriodID,
		name:       data.Name,
		grade:      data.Grade,
		level:      data.Level,
		major:      data.Major,
		homeroomID: data.HomeroomID,
		capacity:   data.Capacity,
		status:     "active", // default aktif
		events:     []types.DomainEvent{},
		createdAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	c.recordEvent(events.NewClassCreated(events.ClassCreatedPayload{
		ClassID:    c.id,
		PeriodID:   c.periodID,
		Name:       c.name,
		Grade:      c.grade,
		Level:      c.level,
		Major:      c.major,
		HomeroomID: c.homeroomID,
		Capacity:   c.capacity,
	}, ctx))
	return c, nil
}
func LoadClass(id, schoolID, periodID, name string, grade int, level string, major, homeroomID string, capacity int, status string, createdAt time.Time) *Class {
	return &Class{
		id: id, schoolID: schoolID, periodID: periodID, name: name,
		grade: grade, level: level, major: major, homeroomID: homeroomID,
		capacity: capacity, status: status, events: []types.DomainEvent{},
		createdAt: createdAt, UpdatedAt: time.Now().UTC(),
	}
}

// =============================================================================
// GETTERS (Existing)
// =============================================================================

func (c *Class) ID() string           { return c.id }
func (c *Class) SchoolID() string     { return c.schoolID }
func (c *Class) PeriodID() string     { return c.periodID }
func (c *Class) Name() string         { return c.name }
func (c *Class) Grade() int           { return c.grade }
func (c *Class) Level() string        { return c.level }
func (c *Class) Major() string        { return c.major }
func (c *Class) HomeroomID() string   { return c.homeroomID }
func (c *Class) Capacity() int        { return c.capacity }
func (c *Class) CreatedAt() time.Time { return c.createdAt }
func (c *Class) Status() string       { return c.status } // tambahan

// =============================================================================
// BEHAVIORS (Domain Logic) - Existing + Tambahan
// =============================================================================

// SetHomeroom mengatur wali kelas.
func (c *Class) SetHomeroom(teacherID string, ctx *context.OperationalContext) error {
	if teacherID == "" {
		return errors.New("teacherId wajib diisi")
	}
	c.homeroomID = teacherID
	c.UpdatedAt = time.Now().UTC()

	c.recordEvent(events.NewClassHomeroomChanged(events.ClassHomeroomChangedPayload{
		ClassID:    c.id,
		HomeroomID: teacherID,
		ChangedAt:  time.Now().UTC(),
	}, ctx))
	return nil
}

// UpdateClass memperbarui atribut class yang diizinkan.
// Ini adalah method baru untuk update massal (name, grade, level, major, homeroom, capacity, status).
func (c *Class) UpdateClass(
	name string,
	grade int,
	level string,
	major string,
	homeroomID string,
	capacity int,
	status string,
	ctx *context.OperationalContext,
) error {
	if name == "" {
		return errors.New("nama kelas wajib diisi")
	}
	if grade < 10 || grade > 12 {
		return errors.New("grade harus 10, 11, atau 12")
	}
	if capacity <= 0 {
		return errors.New("capacity harus lebih dari 0")
	}
	if status != "" && status != "active" && status != "inactive" {
		return errors.New("status harus 'active' atau 'inactive'")
	}

	c.name = name
	c.grade = grade
	c.level = level
	c.major = major
	c.homeroomID = homeroomID
	c.capacity = capacity
	if status != "" {
		c.status = status
	}
	c.UpdatedAt = time.Now().UTC()

	// Record event (bisa ditambahkan jika perlu)
	c.recordEvent(events.NewClassUpdated(events.ClassUpdatedPayload{
		ClassID:    c.id,
		Name:       c.name,
		Grade:      c.grade,
		Level:      c.level,
		Major:      c.major,
		HomeroomID: c.homeroomID,
		Capacity:   c.capacity,
		Status:     c.status,
	}, ctx))
	return nil
}

// SetStatus mengubah status class (active/inactive).
func (c *Class) SetStatus(status string, ctx *context.OperationalContext) error {
	if status != "active" && status != "inactive" {
		return errors.New("status harus 'active' atau 'inactive'")
	}
	c.status = status
	c.UpdatedAt = time.Now().UTC()
	// optional event
	return nil
}

// =============================================================================
// EVENT MANAGEMENT
// =============================================================================

func (c *Class) recordEvent(event types.DomainEvent) {
	c.events = append(c.events, event)
}

func (c *Class) UncommittedEvents() []types.DomainEvent {
	return c.events
}

func (c *Class) ClearEvents() {
	c.events = []types.DomainEvent{}
}
