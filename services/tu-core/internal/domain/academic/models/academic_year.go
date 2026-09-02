// Package models menyediakan domain models untuk Academic aggregate.
// Academic adalah reference data yang digunakan oleh semua domain lain.
package models

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/tu-core/internal/domain/academic/events"

	"github.com/google/uuid"
)

// AcademicYear adalah Aggregate Root untuk tahun ajaran.
type AcademicYear struct {
	id        string
	schoolID  string
	startYear int
	endYear   int
	startDate time.Time
	endDate   time.Time
	isActive  bool
	semesters []*Semester
	events    []types.DomainEvent
	createdAt time.Time
}

var academicYearPattern = regexp.MustCompile(`^\d{4}/\d{4}$`)

// AcademicYearData adalah data untuk membuat AcademicYear.
type AcademicYearData struct {
	ID        string
	SchoolID  string
	StartYear int
	EndYear   int
	StartDate time.Time
	EndDate   time.Time
	IsActive  bool
}

// NewAcademicYear membuat AcademicYear baru.
func NewAcademicYear(data AcademicYearData, ctx *context.OperationalContext) (*AcademicYear, error) {
	if data.ID == "" {
		data.ID = uuid.New().String()
	}

	if data.StartYear <= 0 || data.EndYear <= 0 {
		return nil, errors.New("tahun ajaran tidak valid")
	}
	if data.EndYear != data.StartYear+1 {
		return nil, fmt.Errorf("tahun akhir harus tahun awal + 1, got: %d/%d", data.StartYear, data.EndYear)
	}
	if data.StartDate.IsZero() || data.EndDate.IsZero() {
		return nil, errors.New("tanggal mulai dan selesai wajib diisi")
	}
	if data.EndDate.Before(data.StartDate) {
		return nil, errors.New("tanggal selesai harus setelah tanggal mulai")
	}

	ay := &AcademicYear{
		id:        data.ID,
		schoolID:  ctx.SchoolID,
		startYear: data.StartYear,
		endYear:   data.EndYear,
		startDate: data.StartDate,
		endDate:   data.EndDate,
		isActive:  data.IsActive,
		semesters: []*Semester{},
		events:    []types.DomainEvent{},
		createdAt: time.Now().UTC(),
	}

	ay.recordEvent(events.NewAcademicYearCreated(events.AcademicYearCreatedPayload{
		AcademicYearID: ay.id,
		Label:          ay.Label(),
		StartYear:      ay.startYear,
		EndYear:        ay.endYear,
		IsActive:       ay.isActive,
	}, ctx))
	return ay, nil
}

// Getters

func (ay *AcademicYear) ID() string             { return ay.id }
func (ay *AcademicYear) SchoolID() string       { return ay.schoolID }
func (ay *AcademicYear) StartYear() int         { return ay.startYear }
func (ay *AcademicYear) EndYear() int           { return ay.endYear }
func (ay *AcademicYear) StartDate() time.Time   { return ay.startDate }
func (ay *AcademicYear) EndDate() time.Time     { return ay.endDate }
func (ay *AcademicYear) IsActive() bool         { return ay.isActive }
func (ay *AcademicYear) Semesters() []*Semester { return ay.semesters }

// Label mengembalikan label tahun ajaran (e.g. "2026/2027").
func (ay *AcademicYear) Label() string {
	return fmt.Sprintf("%d/%d", ay.startYear, ay.endYear)
}

// AddSemester menambahkan semester ke tahun ajaran.
func (ay *AcademicYear) AddSemester(number int, startDate, endDate time.Time, ctx *context.OperationalContext) (*Semester, error) {
	if number != 1 && number != 2 {
		return nil, errors.New("semester harus 1 atau 2")
	}

	// Cek duplikasi
	for _, s := range ay.semesters {
		if s.Number() == number {
			return nil, fmt.Errorf("semester %d sudah ada", number)
		}
	}

	semester := NewSemester(uuid.New().String(), ay.id, ay.schoolID, number, startDate, endDate)
	ay.semesters = append(ay.semesters, semester)
	ay.recordEvent(events.NewSemesterAdded(events.SemesterAddedPayload{
		AcademicYearID: ay.id,
		SemesterID:     semester.ID(),
		Number:         semester.Number(),
		PeriodID:       semester.PeriodID(ay.startYear),
		Label:          semester.Label(),
	}, ctx))
	return semester, nil
}

// Event management

func (ay *AcademicYear) recordEvent(event types.DomainEvent) {
	ay.events = append(ay.events, event)
}

func (ay *AcademicYear) UncommittedEvents() []types.DomainEvent {
	return ay.events
}

func (ay *AcademicYear) ClearEvents() {
	ay.events = []types.DomainEvent{}
}
