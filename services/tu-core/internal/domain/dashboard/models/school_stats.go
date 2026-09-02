// Package models menyediakan models untuk Dashboard projection.
package models

import "time"

// SchoolStats adalah aggregated statistics per sekolah per periode.
type SchoolStats struct {
	ID       string `gorm:"primaryKey"`
	SchoolID string `gorm:"index;not null"`
	PeriodID string `gorm:"index;not null"`

	// Student stats
	TotalStudents       int
	ActiveStudents      int
	GraduatedStudents   int
	TransferredStudents int

	// Teacher stats
	TotalTeachers  int
	ActiveTeachers int

	// Attendance stats
	TotalAttendanceRecords int
	TotalPresent           int
	TotalAbsent            int
	TotalSick              int
	TotalLeave             int
	AttendanceRate         float64 // Percentage

	// Assessment stats
	TotalAssessments int
	TotalScores      int
	AverageScore     float64
	PassRate         float64 // Percentage >= 75

	// Class stats
	TotalClasses int

	// Metadata
	LastUpdatedAt time.Time
	CreatedAt     time.Time
}

func (SchoolStats) TableName() string { return "dashboard_school_stats" }

// UniqueKey mengembalikan unique key untuk stats.
func (s *SchoolStats) UniqueKey() string {
	return s.SchoolID + ":" + s.PeriodID
}
