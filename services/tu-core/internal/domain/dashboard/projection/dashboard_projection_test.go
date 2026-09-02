package projection

import (
	"context"
	platformctx "sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/tu-core/internal/domain/dashboard/models"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test DB: %v", err)
	}

	if err := db.AutoMigrate(
		&models.SchoolStats{},
		&models.ClassStats{},
		&models.TeacherStats{},
	); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func createTestEvent(eventName string, payload map[string]interface{}) types.DomainEvent {
	return types.NewBaseEventBuilder(eventName, "Test", "agg-001").
		WithEventID("evt-001").
		WithContext(&platformctx.OperationalContext{
			SchoolID:         "smk-001",
			AcademicPeriodID: "20261",
			UserID:           "user-123",
		}).
		WithPayload(payload).
		Build()
}

func TestHandleAttendanceSubmitted(t *testing.T) {
	db := setupTestDB(t)
	dp := NewDashboardProjection(db)

	payload := map[string]interface{}{
		"classId":   "class-x-1",
		"teacherId": "teacher-001",
		"statistics": map[string]interface{}{
			"total":          36.0,
			"present":        30.0,
			"absent":         3.0,
			"sick":           2.0,
			"leave":          1.0,
			"attendanceRate": 83.33,
		},
	}

	event := createTestEvent("AttendanceSubmitted", payload)

	err := dp.Handle(context.Background(), event)
	if err != nil {
		t.Fatalf("Handle failed: %v", err)
	}

	// Check school stats
	schoolStats, err := dp.GetSchoolStats(context.Background(), "smk-001", "20261")
	if err != nil {
		t.Fatalf("GetSchoolStats failed: %v", err)
	}
	if schoolStats == nil {
		t.Fatal("Expected school stats to be created")
	}
	if schoolStats.TotalAttendanceRecords != 36 {
		t.Errorf("Expected total records 36, got %d", schoolStats.TotalAttendanceRecords)
	}
	if schoolStats.TotalPresent != 30 {
		t.Errorf("Expected present 30, got %d", schoolStats.TotalPresent)
	}

	// Check class stats
	classStats, err := dp.GetClassStats(context.Background(), "smk-001", "20261", "class-x-1")
	if err != nil {
		t.Fatalf("GetClassStats failed: %v", err)
	}
	if classStats == nil {
		t.Fatal("Expected class stats to be created")
	}
	if classStats.TotalAttendanceRecords != 36 {
		t.Errorf("Expected class total records 36, got %d", classStats.TotalAttendanceRecords)
	}

	// Check teacher stats
	teacherStats, err := dp.GetTeacherStats(context.Background(), "smk-001", "20261", "teacher-001")
	if err != nil {
		t.Fatalf("GetTeacherStats failed: %v", err)
	}
	if teacherStats == nil {
		t.Fatal("Expected teacher stats to be created")
	}
	if teacherStats.AttendanceSubmissions != 1 {
		t.Errorf("Expected submissions 1, got %d", teacherStats.AttendanceSubmissions)
	}
}

func TestHandleStudentCreated(t *testing.T) {
	db := setupTestDB(t)
	dp := NewDashboardProjection(db)

	payload := map[string]interface{}{
		"studentId": "student-001",
		"nisn":      "1234567890",
	}

	event := createTestEvent("StudentCreated", payload)

	err := dp.Handle(context.Background(), event)
	if err != nil {
		t.Fatalf("Handle failed: %v", err)
	}

	schoolStats, err := dp.GetSchoolStats(context.Background(), "smk-001", "20261")
	if err != nil {
		t.Fatalf("GetSchoolStats failed: %v", err)
	}
	if schoolStats.TotalStudents != 1 {
		t.Errorf("Expected total students 1, got %d", schoolStats.TotalStudents)
	}
	if schoolStats.ActiveStudents != 1 {
		t.Errorf("Expected active students 1, got %d", schoolStats.ActiveStudents)
	}
}

func TestHandleAssessmentSubmitted(t *testing.T) {
	db := setupTestDB(t)
	dp := NewDashboardProjection(db)

	payload := map[string]interface{}{
		"classId":   "class-x-1",
		"teacherId": "teacher-001",
		"statistics": map[string]interface{}{
			"total":    36.0,
			"passing":  30.0,
			"average":  78.5,
			"passRate": 83.33,
			"min":      60.0,
			"max":      95.0,
		},
	}

	event := createTestEvent("AssessmentSubmitted", payload)

	err := dp.Handle(context.Background(), event)
	if err != nil {
		t.Fatalf("Handle failed: %v", err)
	}

	schoolStats, err := dp.GetSchoolStats(context.Background(), "smk-001", "20261")
	if err != nil {
		t.Fatalf("GetSchoolStats failed: %v", err)
	}
	if schoolStats.TotalAssessments != 1 {
		t.Errorf("Expected total assessments 1, got %d", schoolStats.TotalAssessments)
	}
	if schoolStats.TotalScores != 36 {
		t.Errorf("Expected total scores 36, got %d", schoolStats.TotalScores)
	}
}

func TestMultipleEventsAggregation(t *testing.T) {
	db := setupTestDB(t)
	dp := NewDashboardProjection(db)

	// Submit 3 attendance sessions
	for i := 0; i < 3; i++ {
		payload := map[string]interface{}{
			"classId":   "class-x-1",
			"teacherId": "teacher-001",
			"statistics": map[string]interface{}{
				"total":   36.0,
				"present": 30.0,
				"absent":  6.0,
			},
		}
		event := createTestEvent("AttendanceSubmitted", payload)
		_ = dp.Handle(context.Background(), event)
	}

	teacherStats, err := dp.GetTeacherStats(context.Background(), "smk-001", "20261", "teacher-001")
	if err != nil {
		t.Fatalf("GetTeacherStats failed: %v", err)
	}
	if teacherStats.AttendanceSubmissions != 3 {
		t.Errorf("Expected 3 submissions, got %d", teacherStats.AttendanceSubmissions)
	}
	if teacherStats.TotalStudentsAttended != 108 { // 36 * 3
		t.Errorf("Expected 108 students attended, got %d", teacherStats.TotalStudentsAttended)
	}
}

func TestUnknownEventIgnored(t *testing.T) {
	db := setupTestDB(t)
	dp := NewDashboardProjection(db)

	event := createTestEvent("UnknownEvent", map[string]interface{}{})

	// Should not error
	err := dp.Handle(context.Background(), event)
	if err != nil {
		t.Errorf("Expected no error for unknown event, got: %v", err)
	}
}

func TestHandlerCount(t *testing.T) {
	db := setupTestDB(t)
	dp := NewDashboardProjection(db)

	count := dp.HandlerCount()
	if count == 0 {
		t.Error("Expected handlers to be registered")
	}

	// Expected handlers:
	// - AttendanceSubmitted, AttendanceCorrected
	// - AssessmentSubmitted, ScoreUpdated
	// - StudentCreated, StudentGraduated, StudentTransferred
	// - TeacherCreated, TeacherDeactivated
	expectedMin := 9
	if count < expectedMin {
		t.Errorf("Expected at least %d handlers, got %d", expectedMin, count)
	}
}

// Placeholder for types.OperationalContext
type OperationalContext = platformctx.OperationalContext

// Placeholder for time
var _ = time.Now
