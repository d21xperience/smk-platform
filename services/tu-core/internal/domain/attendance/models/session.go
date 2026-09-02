package models

import (
	"errors"
	"time"

	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/tu-core/internal/domain/attendance/events"
	"sekolah-platform/services/tu-core/internal/domain/attendance/valueobjects"

	"github.com/google/uuid"
)

type Session struct {
	ID               string
	SchoolID         string
	AcademicPeriodID string
	ClassID          string
	SubjectID        string
	TeacherID        string
	SessionDate      time.Time
	StartTime        time.Time
	EndTime          time.Time
	Status           valueobjects.SessionStatus
	Notes            string
	Records          []*Record
	CreatedAt        time.Time
	UpdatedAt        time.Time
	domainEvents     []types.DomainEvent
}

type SessionData struct {
	ID          string
	ClassID     string
	SubjectID   string
	TeacherID   string
	SessionDate time.Time
	StartTime   time.Time
	EndTime     time.Time
	Notes       string
}

func NewSession(data SessionData, ctx *context.OperationalContext) (*Session, error) {
	if data.ID == "" {
		data.ID = uuid.New().String()
	}

	if data.ClassID == "" {
		return nil, errors.New("class_id wajib diisi")
	}

	if data.SessionDate.IsZero() {
		return nil, errors.New("session_date wajib diisi")
	}

	s := &Session{
		ID:               data.ID,
		SchoolID:         ctx.SchoolID,
		AcademicPeriodID: ctx.AcademicPeriodID,
		ClassID:          data.ClassID,
		SubjectID:        data.SubjectID,
		TeacherID:        data.TeacherID,
		SessionDate:      data.SessionDate,
		StartTime:        data.StartTime,
		EndTime:          data.EndTime,
		Status:           valueobjects.SessionDraft,
		Notes:            data.Notes,
		Records:          []*Record{},
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
		domainEvents:     []types.DomainEvent{},
	}

	s.recordEvent(events.NewSessionCreated(events.SessionCreatedPayload{
		SessionID: s.ID,
		ClassID:   s.ClassID,
		Date:      s.SessionDate,
	}, ctx))

	return s, nil
}

func (s *Session) AddRecord(record *Record) error {
	if !s.Status.IsDraft() {
		return errors.New("tidak dapat menambah record pada session yang sudah submitted/closed")
	}

	// Cek duplikasi student
	for _, r := range s.Records {
		if r.StudentID == record.StudentID {
			return errors.New("student sudah memiliki record dalam session ini")
		}
	}

	s.Records = append(s.Records, record)
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *Session) Submit(ctx *context.OperationalContext) error {
	if !s.Status.CanSubmit() {
		return errors.New("hanya session dengan status DRAFT yang dapat disubmit")
	}

	if len(s.Records) == 0 {
		return errors.New("session harus memiliki minimal 1 record sebelum disubmit")
	}

	s.Status = valueobjects.SessionSubmitted
	s.UpdatedAt = time.Now().UTC()

	s.recordEvent(events.NewSessionSubmitted(events.SessionSubmittedPayload{
		SessionID:    s.ID,
		ClassID:      s.ClassID,
		TotalRecords: len(s.Records),
	}, ctx))

	return nil
}

func (s *Session) Close(ctx *context.OperationalContext) error {
	if !s.Status.CanClose() {
		return errors.New("hanya session dengan status SUBMITTED yang dapat ditutup")
	}

	s.Status = valueobjects.SessionClosed
	s.UpdatedAt = time.Now().UTC()

	s.recordEvent(events.NewSessionClosed(events.SessionClosedPayload{
		SessionID: s.ID,
		ClassID:   s.ClassID,
	}, ctx))

	return nil
}

func (s *Session) recordEvent(event types.DomainEvent) {
	s.domainEvents = append(s.domainEvents, event)
}

func (s *Session) UncommittedEvents() []types.DomainEvent {
	return s.domainEvents
}

func (s *Session) ClearEvents() {
	s.domainEvents = []types.DomainEvent{}
}
