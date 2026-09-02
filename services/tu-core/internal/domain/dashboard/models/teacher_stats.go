package models

import "time"

// TeacherStats adalah aggregated statistics per guru per periode.
type TeacherStats struct {
	ID        string `gorm:"primaryKey"`
	SchoolID  string `gorm:"index;not null"`
	PeriodID  string `gorm:"index;not null"`
	TeacherID string `gorm:"index;not null"`

	// Attendance submissions
	AttendanceSubmissions int
	TotalStudentsAttended int

	// Assessment submissions
	AssessmentsCompleted int
	TotalScoresSubmitted int
	AverageStudentScore  float64

	// Metadata
	LastUpdatedAt time.Time
	CreatedAt     time.Time
}

func (TeacherStats) TableName() string { return "dashboard_teacher_stats" }

// UniqueKey mengembalikan unique key.
func (t *TeacherStats) UniqueKey() string {
	return t.SchoolID + ":" + t.PeriodID + ":" + t.TeacherID
}
