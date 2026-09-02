package models

import "time"

// Enrollment adalah Entity dalam Student Aggregate.
// SYMMETRIC dengan Enrollment.js di frontend.
type Enrollment struct {
	enrollmentID   string
	schoolID       string
	periodID       string
	classID        string
	enrollmentDate time.Time
	status         string // active, completed, transferred
}

// NewEnrollment membuat Enrollment baru.
func NewEnrollment(enrollmentID, schoolID, periodID, classID string) *Enrollment {
	return &Enrollment{
		enrollmentID:   enrollmentID,
		schoolID:       schoolID,
		periodID:       periodID,
		classID:        classID,
		enrollmentDate: time.Now().UTC(),
		status:         "active",
	}
}

// RehydrateEnrollment merekonstruksi Enrollment dari data persisten.
func RehydrateEnrollment(enrollmentID, schoolID, periodID, classID, status string, enrollmentDate time.Time) *Enrollment {
	return &Enrollment{
		enrollmentID:   enrollmentID,
		schoolID:       schoolID,
		periodID:       periodID,
		classID:        classID,
		enrollmentDate: enrollmentDate,
		status:         status,
	}
}

// Getters

func (e *Enrollment) EnrollmentID() string      { return e.enrollmentID }
func (e *Enrollment) SchoolID() string          { return e.schoolID }
func (e *Enrollment) PeriodID() string          { return e.periodID }
func (e *Enrollment) ClassID() string           { return e.classID }
func (e *Enrollment) EnrollmentDate() time.Time { return e.enrollmentDate }
func (e *Enrollment) Status() string            { return e.status }

// IsActive mengecek apakah enrollment aktif.
func (e *Enrollment) IsActive() bool {
	return e.status == "active"
}

// IsForPeriod mengecek apakah enrollment untuk periode tertentu.
func (e *Enrollment) IsForPeriod(periodID string) bool {
	return e.periodID == periodID
}

// Complete menandai enrollment sebagai selesai.
func (e *Enrollment) Complete() {
	if e.status != "active" {
		return
	}
	e.status = "completed"
}
