package models

import "time"

// ClassStats adalah aggregated statistics per kelas per periode.
type ClassStats struct {
	ID       string `gorm:"primaryKey"`
	SchoolID string `gorm:"index;not null"`
	PeriodID string `gorm:"index;not null"`
	ClassID  string `gorm:"index;not null"`

	// Student stats
	TotalStudents int

	// Attendance stats
	TotalAttendanceRecords int
	TotalPresent           int
	AttendanceRate         float64

	// Assessment stats
	TotalAssessments int
	AverageScore     float64
	PassRate         float64

	// Metadata
	LastUpdatedAt time.Time
	CreatedAt     time.Time
}

func (ClassStats) TableName() string { return "dashboard_class_stats" }

// UniqueKey mengembalikan unique key.
func (c *ClassStats) UniqueKey() string {
	return c.SchoolID + ":" + c.PeriodID + ":" + c.ClassID
}
