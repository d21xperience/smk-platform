package context

import "testing"

func TestExtract(t *testing.T) {
	provider := MapProvider{
		HeaderSchoolID:       "smk-001",
		HeaderSchoolName:     "SMK Pasundan",
		HeaderAcademicYear:   "2026/2027",
		HeaderSemester:       "1",
		HeaderUserID:         "user-123",
		HeaderUserRole:       "operator",
	}

	ctx, err := Extract(provider)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if ctx.SchoolID != "smk-001" {
		t.Errorf("Expected schoolId 'smk-001', got '%s'", ctx.SchoolID)
	}
	if ctx.AcademicPeriodID != "20261" {
		t.Errorf("Expected periodId '20261', got '%s'", ctx.AcademicPeriodID)
	}
	if ctx.UserID != "user-123" {
		t.Errorf("Expected userId 'user-123', got '%s'", ctx.UserID)
	}
}

func TestExtractFromPeriodID(t *testing.T) {
	provider := MapProvider{
		HeaderSchoolID:       "smk-001",
		HeaderAcademicPeriod: "20261",
		HeaderUserID:         "user-123",
		HeaderUserRole:       "operator",
	}

	ctx, err := Extract(provider)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if ctx.AcademicYear != "2026/2027" {
		t.Errorf("Expected academicYear '2026/2027', got '%s'", ctx.AcademicYear)
	}
	if ctx.Semester != 1 {
		t.Errorf("Expected semester 1, got %d", ctx.Semester)
	}
}

func TestParsePeriodID(t *testing.T) {
	tests := []struct {
		periodID     string
		expectedYear string
		expectedSem  int
	}{
		{"20261", "2026/2027", 1},
		{"20262", "2026/2027", 2},
		{"20251", "2025/2026", 1},
	}

	for _, tt := range tests {
		year, sem := ParsePeriodID(tt.periodID)
		if year != tt.expectedYear {
			t.Errorf("PeriodID %s: expected year %s, got %s", tt.periodID, tt.expectedYear, year)
		}
		if sem != tt.expectedSem {
			t.Errorf("PeriodID %s: expected semester %d, got %d", tt.periodID, tt.expectedSem, sem)
		}
	}
}

func TestToHeaders(t *testing.T) {
	ctx := &OperationalContext{
		SchoolID:         "smk-001",
		SchoolName:       "SMK Pasundan",
		AcademicYear:     "2026/2027",
		Semester:         1,
		AcademicPeriodID: "20261",
		UserID:           "user-123",
		Role:             "operator",
	}

	headers := ctx.ToHeaders()
	if headers[HeaderSchoolID] != "smk-001" {
		t.Error("Expected schoolId header")
	}
	if headers[HeaderAcademicPeriod] != "20261" {
		t.Error("Expected period header")
	}
}