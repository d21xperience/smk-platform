package models

import (
	"fmt"
	"time"
)

// Semester adalah Entity dalam AcademicYear Aggregate.
type Semester struct {
	id        string
	yearID    string
	schoolID  string
	number    int // 1 atau 2
	startDate time.Time
	endDate   time.Time
}

// NewSemester membuat Semester baru.
func NewSemester(id, yearID, schoolID string, number int, startDate, endDate time.Time) *Semester {
	return &Semester{
		id:        id,
		yearID:    yearID,
		schoolID:  schoolID,
		number:    number,
		startDate: startDate,
		endDate:   endDate,
	}
}

// Getters

func (s *Semester) ID() string           { return s.id }
func (s *Semester) YearID() string       { return s.yearID }
func (s *Semester) SchoolID() string     { return s.schoolID }
func (s *Semester) Number() int          { return s.number }
func (s *Semester) StartDate() time.Time { return s.startDate }
func (s *Semester) EndDate() time.Time   { return s.endDate }

// PeriodID menurunkan periodID dari tahun ajaran + semester.
// Contoh: 2026 + semester 1 → "20261"
func (s *Semester) PeriodID(startYear int) string {
	return fmt.Sprintf("%d%d", startYear, s.number)
}

// Label mengembalikan label semester (e.g. "Semester 1").
func (s *Semester) Label() string {
	return fmt.Sprintf("Semester %d", s.number)
}

// IsActiveOn mengecek apakah tanggal tertentu berada dalam rentang semester.
func (s *Semester) IsActiveOn(date time.Time) bool {
	return !date.Before(s.startDate) && !date.After(s.endDate)
}
