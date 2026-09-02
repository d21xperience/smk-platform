package validator

import (
	"testing"
)

func TestRequired(t *testing.T) {
	v := New().Required("name", "")
	if v.IsValid() {
		t.Error("Expected validation to fail for empty value")
	}

	v = New().Required("name", "Budi")
	if !v.IsValid() {
		t.Error("Expected validation to pass")
	}
}

func TestMinLength(t *testing.T) {
	v := New().MinLength("name", "ab", 3)
	if v.IsValid() {
		t.Error("Expected validation to fail")
	}

	v = New().MinLength("name", "abc", 3)
	if !v.IsValid() {
		t.Error("Expected validation to pass")
	}
}

func TestNISN(t *testing.T) {
	tests := []struct {
		value string
		valid bool
	}{
		{"1234567890", true},
		{"12345", false},
		{"123456789a", false},
		{"", false}, // NISN tidak boleh kosong jika divalidasi
	}

	for _, tt := range tests {
		v := New().Required("nisn", tt.value).NISN("nisn", tt.value)
		if v.IsValid() != tt.valid {
			t.Errorf("NISN '%s': expected valid=%v, got %v", tt.value, tt.valid, v.IsValid())
		}
	}
}

func TestEmail(t *testing.T) {
	tests := []struct {
		value string
		valid bool
	}{
		{"a@b.com", true},
		{"invalid", false},
		{"", true}, // Email optional
	}

	for _, tt := range tests {
		v := New().Email("email", tt.value)
		if v.IsValid() != tt.valid {
			t.Errorf("Email '%s': expected valid=%v, got %v", tt.value, tt.valid, v.IsValid())
		}
	}
}

func TestOneOf(t *testing.T) {
	v := New().OneOf("gender", "MALE", []string{"MALE", "FEMALE"})
	if !v.IsValid() {
		t.Error("Expected validation to pass")
	}

	v = New().OneOf("gender", "UNKNOWN", []string{"MALE", "FEMALE"})
	if v.IsValid() {
		t.Error("Expected validation to fail")
	}
}

func TestChainedValidation(t *testing.T) {
	v := New().
		Required("name", "Budi").
		MinLength("name", "Budi", 2).
		Required("nisn", "1234567890").
		NISN("nisn", "1234567890")

	if !v.IsValid() {
		t.Errorf("Expected validation to pass, got errors: %v", v.Errors())
	}
}

// func TestToAppError(t *testing.T) {
// 	v := New().Required("name", "")
// 	err := v.Validate()
// 	if err == nil {
// 		t.Fatal("Expected error")
// 	}

// 	appErr, ok := err.(ValidationErrors).ToAppError()
// 	if !ok {
// 		t.Fatal("Expected AppError")
// 	}
// 	if appErr.Code != errors.CodeValidation {
// 		t.Errorf("Expected code %s, got %s", errors.CodeValidation, appErr.Code)
// 	}
// }
