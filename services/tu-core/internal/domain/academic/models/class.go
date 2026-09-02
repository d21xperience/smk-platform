package models

import (
	"errors"
	"time"

	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"

	"github.com/google/uuid"

	// ✅ Import ini aman sekarang karena events TIDAK mengimpor balik ke models
	"sekolah-platform/services/tu-core/internal/domain/academic/events"
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
	events     []types.DomainEvent
	createdAt  time.Time
}

// =============================================================================
// COMMAND DATA (Input DTO)
// =============================================================================

// ClassData adalah data untuk membuat Class.
// Ini merepresentasikan INTENT (niat) dari command, bukan fakta historis.
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
		events:     []types.DomainEvent{},
		createdAt:  time.Now().UTC(),
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

// =============================================================================
// GETTERS
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

// =============================================================================
// BEHAVIORS (Domain Logic)
// =============================================================================

// SetHomeroom mengatur wali kelas.
func (c *Class) SetHomeroom(teacherID string, ctx *context.OperationalContext) error {
	if teacherID == "" {
		return errors.New("teacherId wajib diisi")
	}
	c.homeroomID = teacherID

	c.recordEvent(events.NewClassHomeroomChanged(events.ClassHomeroomChangedPayload{
		ClassID:    c.id,
		HomeroomID: teacherID,
		ChangedAt:  time.Now().UTC(),
	}, ctx))
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
