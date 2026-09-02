package engine

import (
	"testing"
	"time"

	"sekolah-platform/platform/context"
	"sekolah-platform/services/tu-core/internal/domain/teacher/engine/commands"
)

func createTestContext() *context.OperationalContext {
	ctx, _ := context.New("smk-001", "2026/2027", 1, "user-123", "operator")
	return ctx
}

func TestRegisterTeacher(t *testing.T) {
	engine := NewTeacherEngine()
	ctx := createTestContext()

	cmd := commands.RegisterTeacherCommand{
		NIP:       "123456789012345678",
		FirstName: "Budi",
		LastName:  "Santoso",
		BirthDate: time.Date(1980, 5, 15, 0, 0, 0, 0, time.UTC),
		Gender:    "MALE",
		Subject:   "Matematika",
		JoinDate:  time.Now(),
	}

	result := engine.RegisterTeacher(cmd, ctx)

	if result.Error != nil {
		t.Fatalf("Expected success, got error: %v", result.Error)
	}
	if result.Teacher == nil {
		t.Fatal("Expected teacher to be created")
	}
	if result.Teacher.SchoolID() != "smk-001" {
		t.Errorf("Expected schoolID 'smk-001', got '%s'", result.Teacher.SchoolID())
	}
	if len(result.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(result.Events))
	}
	if result.Events[0].GetEventName() != "TeacherCreated" {
		t.Errorf("Expected event 'TeacherCreated', got '%s'", result.Events[0].GetEventName())
	}
}

func TestRegisterTeacherInvalid(t *testing.T) {
	engine := NewTeacherEngine()
	ctx := createTestContext()

	cmd := commands.RegisterTeacherCommand{
		NIP: "12345", // Invalid
	}

	result := engine.RegisterTeacher(cmd, ctx)

	if result.Error == nil {
		t.Error("Expected error for invalid NIP")
	}
}
