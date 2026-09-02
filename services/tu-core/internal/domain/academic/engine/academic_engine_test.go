package engine

import (
	"testing"
	"time"

	"sekolah-platform/platform/context"
	"sekolah-platform/services/tu-core/internal/domain/academic/models"
)

func TestCreateAcademicYear(t *testing.T) {
	engine := NewAcademicEngine()
	ctx, _ := context.New("smk-001", "2026/2027", 1, "user-123", "operator")

	startDate := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2027, 6, 30, 0, 0, 0, 0, time.UTC)

	result := engine.CreateAcademicYear(2026, 2027, startDate, endDate, true, ctx)

	if result.Error != nil {
		t.Fatalf("Expected success, got error: %v", result.Error)
	}
	if result.AcademicYear == nil {
		t.Fatal("Expected academic year to be created")
	}
	if result.AcademicYear.Label() != "2026/2027" {
		t.Errorf("Expected label '2026/2027', got '%s'", result.AcademicYear.Label())
	}
}

func TestCreateAcademicYearInvalid(t *testing.T) {
	engine := NewAcademicEngine()
	ctx, _ := context.New("smk-001", "2026/2027", 1, "user-123", "operator")

	// End year != start year + 1
	result := engine.CreateAcademicYear(2026, 2028, time.Now(), time.Now().AddDate(1, 0, 0), true, ctx)

	if result.Error == nil {
		t.Error("Expected error for invalid year range")
	}
}

func TestAddSemester(t *testing.T) {
	engine := NewAcademicEngine()
	ctx, _ := context.New("smk-001", "2026/2027", 1, "user-123", "operator")

	startDate := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2027, 6, 30, 0, 0, 0, 0, time.UTC)

	yearResult := engine.CreateAcademicYear(2026, 2027, startDate, endDate, true, ctx)
	if yearResult.Error != nil {
		t.Fatalf("Create year failed: %v", yearResult.Error)
	}

	sem1Start := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	sem1End := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)

	semResult := engine.AddSemester(yearResult.AcademicYear, 1, sem1Start, sem1End, ctx)

	if semResult.Error != nil {
		t.Fatalf("Add semester failed: %v", semResult.Error)
	}
	if semResult.Semester == nil {
		t.Fatal("Expected semester to be created")
	}
	if semResult.Semester.PeriodID(2026) != "20261" {
		t.Errorf("Expected periodID '20261', got '%s'", semResult.Semester.PeriodID(2026))
	}
}

func TestCreateClass(t *testing.T) {
	engine := NewAcademicEngine()
	ctx, _ := context.New("smk-001", "2026/2027", 1, "user-123", "operator")

	data := models.ClassData{
		PeriodID: "20261",
		Name:     "X IPA 1",
		Grade:    10,
		Level:    "X",
		Major:    "IPA",
		Capacity: 36,
	}

	result := engine.CreateClass(data, ctx)

	if result.Error != nil {
		t.Fatalf("Expected success, got error: %v", result.Error)
	}
	if result.Class == nil {
		t.Fatal("Expected class to be created")
	}
	if result.Class.Name() != "X IPA 1" {
		t.Errorf("Expected name 'X IPA 1', got '%s'", result.Class.Name())
	}
}
