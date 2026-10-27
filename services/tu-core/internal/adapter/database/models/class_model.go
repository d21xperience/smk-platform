package models

import (
	"time"
)

// ClassModel untuk mapping tabel classes
type ClassModel struct {
	ID         string    `gorm:"column:id;primaryKey"`
	PeriodID   string    `gorm:"column:period_id"`
	Name       string    `gorm:"column:name"`
	Grade      int       `gorm:"column:grade"`
	Level      int       `gorm:"column:level"`
	Major      string    `gorm:"column:major"`
	HomeroomID string    `gorm:"column:homeroom_id"`
	Capacity   int       `gorm:"column:capacity"`
	Status     string    `gorm:"column:status;default:active"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (ClassModel) TableName() string {
	return "classes"
}

// ToDomain mengubah GORM model ke domain model
func (m *ClassModel) ToDomain() *ClassModel {
	return &ClassModel{
		ID:         m.ID,
		PeriodID:   m.PeriodID,
		Name:       m.Name,
		Grade:      m.Grade,
		Level:      m.Level,
		Major:      m.Major,
		HomeroomID: m.HomeroomID,
		Capacity:   m.Capacity,
		Status:     m.Status,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

// FromDomain mengisi GORM model dari domain
func (m *ClassModel) FromDomain(c *ClassModel) {
	m.ID = c.ID
	m.PeriodID = c.PeriodID
	m.Name = c.Name
	m.Grade = c.Grade
	m.Level = c.Level
	m.Major = c.Major
	m.HomeroomID = c.HomeroomID
	m.Capacity = c.Capacity
	m.Status = c.Status
	m.CreatedAt = c.CreatedAt
	m.UpdatedAt = c.UpdatedAt
}
