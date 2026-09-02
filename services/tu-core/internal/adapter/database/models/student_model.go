// Package models menyediakan GORM models untuk database schema.
package models

import (
	"time"

	"gorm.io/gorm"
)

// StudentModel adalah GORM model untuk tabel students.
type StudentModel struct {
	ID              string         `gorm:"primaryKey;size:36"`
	SchoolID        string         `gorm:"size:36;index;not null"`
	NISN            string         `gorm:"size:10;uniqueIndex:idx_school_nisn;not null"`
	NIS             string         `gorm:"size:20;uniqueIndex:idx_school_nis;not null"`
	FirstName       string         `gorm:"size:100;not null"`
	MiddleName      string         `gorm:"size:100"`
	LastName        string         `gorm:"size:100"`
	BirthDate       time.Time      `gorm:"not null"`
	Gender          string         `gorm:"size:10;not null"` // MALE, FEMALE
	AddressStreet   string         `gorm:"size:255"`
	AddressRtRw     string         `gorm:"size:20"`
	AddressVillage  string         `gorm:"size:100"`
	AddressDistrict string         `gorm:"size:100"`
	AddressCity     string         `gorm:"size:100"`
	AddressPostal   string         `gorm:"size:10"`
	ContactPhone    string         `gorm:"size:20"`
	ContactEmail    string         `gorm:"size:100"`
	GuardianName    string         `gorm:"size:100"`
	GuardianRel     string         `gorm:"size:20"` // father, mother, guardian
	GuardianPhone   string         `gorm:"size:20"`
	GuardianJob     string         `gorm:"size:100"`
	Status          string         `gorm:"size:20;not null;default:ACTIVE"`
	CreatedAt       time.Time      `gorm:"not null"`
	UpdatedAt       time.Time      `gorm:"not null"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`

	// Relationships
	Enrollments []EnrollmentModel `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE"`
}

// TableName mengembalikan nama tabel.
func (StudentModel) TableName() string {
	return "students"
}

// EnrollmentModel adalah GORM model untuk tabel enrollments.
type EnrollmentModel struct {
	ID             string    `gorm:"primaryKey;size:36"`
	StudentID      string    `gorm:"size:36;index;not null"`
	SchoolID       string    `gorm:"size:36;index;not null"`
	PeriodID       string    `gorm:"size:10;not null"`
	ClassID        string    `gorm:"size:36;not null"`
	EnrollmentDate time.Time `gorm:"not null"`
	Status         string    `gorm:"size:20;not null;default:active"` // active, completed, transferred
	CreatedAt      time.Time `gorm:"not null"`
	UpdatedAt      time.Time `gorm:"not null"`

	// Relationships
	Student StudentModel `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE"`
}

// TableName mengembalikan nama tabel.
func (EnrollmentModel) TableName() string {
	return "enrollments"
}

// OutboxEventModel adalah GORM model untuk tabel outbox_events.
type OutboxEventModel struct {
	ID         string    `gorm:"primaryKey;size:36"`
	EventType  string    `gorm:"size:100;not null;index"`
	AggregateID string   `gorm:"size:36;not null;index"`
	Payload    string    `gorm:"type:jsonb;not null"` // JSON serialized event
	Attempts   int       `gorm:"not null;default:0"`
	LastError  string    `gorm:"type:text"`
	NextRetry  time.Time `gorm:"index"`
	CreatedAt  time.Time `gorm:"not null;index"`
}

// TableName mengembalikan nama tabel.
func (OutboxEventModel) TableName() string {
	return "outbox_events"
}
