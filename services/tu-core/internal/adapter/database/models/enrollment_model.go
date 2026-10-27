package models

import (
	"time"
)

// EnrollmentModel untuk mapping tabel enrollments
type EnrollmentModel struct {
	EnrollmentID   string    `gorm:"column:id;primaryKey"`
	StudentID      string    `gorm:"column:student_id"`
	SchoolID       string    `gorm:"size:36;index;not null"`
	PeriodID       string    `gorm:"column:period_id"`
	ClassID        string    `gorm:"column:class_id"`
	EnrollmentDate time.Time `gorm:"not null"`
	Status         string    `gorm:"column:status"` // active, completed, transferred
	Reason         string    `gorm:"column:reason"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
	// Relationships
	Student StudentModel `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE"`
}

func (EnrollmentModel) TableName() string {
	return "enrollments"
}

// ToDomain mengubah GORM model ke domain model
func (m *EnrollmentModel) ToDomain() *EnrollmentModel {
	return &EnrollmentModel{
		EnrollmentID: m.EnrollmentID,
		StudentID:    m.StudentID,
		PeriodID:     m.PeriodID,
		ClassID:      m.ClassID,
		Status:       m.Status,
		Reason:       m.Reason,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

// FromDomain mengisi GORM model dari domain
func (m *EnrollmentModel) FromDomain(e *EnrollmentModel) {
	m.EnrollmentID = e.EnrollmentID
	m.StudentID = e.StudentID
	m.PeriodID = e.PeriodID
	m.ClassID = e.ClassID
	m.Status = e.Status
	m.Reason = e.Reason
	m.CreatedAt = e.CreatedAt
	m.UpdatedAt = e.UpdatedAt
}
