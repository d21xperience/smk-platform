package context

import "testing"

func TestNew(t *testing.T) {
	ctx, err := New("smk-001", "2026/2027", 1, "user-123", "operator")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if ctx.SchoolID != "smk-001" {
		t.Errorf("Expected schoolId 'smk-001', got '%s'", ctx.SchoolID)
	}
	if ctx.AcademicPeriodID != "20261" {
		t.Errorf("Expected periodId '20261', got '%s'", ctx.AcademicPeriodID)
	}
}

func TestNewInvalid(t *testing.T) {
	tests := []struct {
		name       string
		schoolID   string
		year       string
		semester   int
		userID     string
		role       string
	}{
		{"empty school", "", "2026/2027", 1, "user", "role"},
		{"invalid year", "smk", "2026", 1, "user", "role"},
		{"invalid semester", "smk", "2026/2027", 3, "user", "role"},
		{"empty user", "smk", "2026/2027", 1, "", "role"},
		{"empty role", "smk", "2026/2027", 1, "user", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.schoolID, tt.year, tt.semester, tt.userID, tt.role)
			if err == nil {
				t.Error("Expected error")
			}
		})
	}
}

func TestDerivePeriodID(t *testing.T) {
	tests := []struct {
		year     string
		semester int
		expected string
	}{
		{"2026/2027", 1, "20261"},
		{"2026/2027", 2, "20262"},
		{"2025/2026", 1, "20251"},
	}

	for _, tt := range tests {
		ctx := &OperationalContext{AcademicYear: tt.year, Semester: tt.semester}
		if got := ctx.DerivePeriodID(); got != tt.expected {
			t.Errorf("Expected %s, got %s", tt.expected, got)
		}
	}
}

func TestEquals(t *testing.T) {
	ctx1 := &OperationalContext{SchoolID: "smk-001", AcademicPeriodID: "20261"}
	ctx2 := &OperationalContext{SchoolID: "smk-001", AcademicPeriodID: "20261"}
	ctx3 := &OperationalContext{SchoolID: "smk-002", AcademicPeriodID: "20261"}

	if !ctx1.Equals(ctx2) {
		t.Error("Expected equal")
	}
	if ctx1.Equals(ctx3) {
		t.Error("Expected not equal")
	}
	if ctx1.Equals(nil) {
		t.Error("Expected not equal to nil")
	}
}

func TestDisplayLabel(t *testing.T) {
	ctx := &OperationalContext{
		SchoolID:     "smk-001",
		SchoolName:   "SMK Pasundan",
		AcademicYear: "2026/2027",
		Semester:     1,
	}
	expected := "SMK Pasundan | 2026/2027 - Semester 1"
	if got := ctx.DisplayLabel(); got != expected {
		t.Errorf("Expected '%s', got '%s'", expected, got)
	}
}

func TestToMap(t *testing.T) {
	ctx := &OperationalContext{
		SchoolID:         "smk-001",
		AcademicYear:     "2026/2027",
		Semester:         1,
		AcademicPeriodID: "20261",
		UserID:           "user-123",
		Role:             "operator",
	}
	m := ctx.ToMap()
	if m["schoolId"] != "smk-001" {
		t.Error("Expected schoolId in map")
	}
	if m["academicPeriodId"] != "20261" {
		t.Error("Expected periodId in map")
	}
}