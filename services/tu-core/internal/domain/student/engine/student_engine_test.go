package engine

import (
	"testing"
	"time"

	"sekolah-platform/platform/context"
	"sekolah-platform/services/tu-core/internal/domain/student/engine/commands"
)

func createTestContext() *context.OperationalContext {
	ctx, _ := context.New("smk-001", "2026/2027", 1, "user-123", "operator")
	return ctx
}

func createTestCommand() commands.RegisterStudentCommand {
	return commands.RegisterStudentCommand{
		StudentID:        "student-001",
		NISN:             "1234567890",
		NIS:              "2026001",
		FirstName:        "Budi",
		MiddleName:       "",
		LastName:         "Santoso",
		BirthDate:        time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC),
		Gender:           "MALE",
		AddressStreet:    "Jl. Test No. 1",
		AddressVillage:   "Cibeusi",
		AddressCity:      "Sumedang",
		ContactPhone:     "081234567890",
		GuardianName:     "Ayah Budi",
		GuardianRelation: "father",
		GuardianPhone:    "081234567891",
	}
}

func TestRegisterStudent(t *testing.T) {
	engine := NewStudentEngine()
	ctx := createTestContext()
	cmd := createTestCommand()

	result := engine.RegisterStudent(cmd, ctx)

	if result.Error != nil {
		t.Fatalf("Expected success, got error: %v", result.Error)
	}
	if result.Student == nil {
		t.Fatal("Expected student to be created")
	}
	if result.Student.StudentID() != "student-001" {
		t.Errorf("Expected studentID 'student-001', got '%s'", result.Student.StudentID())
	}
	if result.Student.SchoolID() != "smk-001" {
		t.Errorf("Expected schoolID 'smk-001', got '%s'", result.Student.SchoolID())
	}
	if len(result.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(result.Events))
	}
	if result.Events[0].GetEventName() != "StudentCreated" {
		t.Errorf("Expected event 'StudentCreated', got '%s'", result.Events[0].GetEventName())
	}
}

func TestRegisterStudentInvalidCommand(t *testing.T) {
	engine := NewStudentEngine()
	ctx := createTestContext()
	cmd := commands.RegisterStudentCommand{
		NISN: "12345", // invalid
	}

	result := engine.RegisterStudent(cmd, ctx)

	if result.Error == nil {
		t.Error("Expected error for invalid command")
	}
}

func TestEnrollStudent(t *testing.T) {
	engine := NewStudentEngine()
	ctx := createTestContext()
	cmd := createTestCommand()

	// Register first
	regResult := engine.RegisterStudent(cmd, ctx)
	if regResult.Error != nil {
		t.Fatalf("Register failed: %v", regResult.Error)
	}

	// Enroll
	enrollCmd := commands.EnrollStudentCommand{
		StudentID:    "student-001",
		EnrollmentID: "enr-001",
		ClassID:      "class-x-1",
	}

	enrollResult := engine.EnrollStudent(enrollCmd, ctx, regResult.Student)

	if enrollResult.Error != nil {
		t.Fatalf("Enroll failed: %v", enrollResult.Error)
	}
	if len(enrollResult.Student.Enrollments()) != 1 {
		t.Errorf("Expected 1 enrollment, got %d", len(enrollResult.Student.Enrollments()))
	}
}

func TestEnrollStudentDoubleEnrollment(t *testing.T) {
	engine := NewStudentEngine()
	ctx := createTestContext()
	cmd := createTestCommand()

	regResult := engine.RegisterStudent(cmd, ctx)
	enrollCmd := commands.EnrollStudentCommand{
		StudentID:    "student-001",
		EnrollmentID: "enr-001",
		ClassID:      "class-x-1",
	}

	// First enroll
	engine.EnrollStudent(enrollCmd, ctx, regResult.Student)

	// Second enroll (same period)
	enrollCmd.EnrollmentID = "enr-002"
	result := engine.EnrollStudent(enrollCmd, ctx, regResult.Student)

	if result.Error == nil {
		t.Error("Expected error for double enrollment")
	}
}

func TestGraduateStudent(t *testing.T) {
	engine := NewStudentEngine()
	ctx := createTestContext()
	cmd := createTestCommand()

	regResult := engine.RegisterStudent(cmd, ctx)
	result := engine.GraduateStudent(ctx, regResult.Student)

	if result.Error != nil {
		t.Fatalf("Graduate failed: %v", result.Error)
	}
	if result.Student.Status().Code() != "GRADUATED" {
		t.Errorf("Expected status GRADUATED, got %s", result.Student.Status().Code())
	}
}

func TestTransferStudent(t *testing.T) {
	engine := NewStudentEngine()
	ctx := createTestContext()
	cmd := createTestCommand()

	regResult := engine.RegisterStudent(cmd, ctx)
	result := engine.TransferStudent("SMK Negeri Bandung", "Pindah domisili", ctx, regResult.Student)

	if result.Error != nil {
		t.Fatalf("Transfer failed: %v", result.Error)
	}
	if result.Student.Status().Code() != "TRANSFERRED" {
		t.Errorf("Expected status TRANSFERRED, got %s", result.Student.Status().Code())
	}
}
