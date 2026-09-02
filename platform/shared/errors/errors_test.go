package errors

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	err := New(CodeNotFound, "Student not found", nil)

	if err.Code != CodeNotFound {
		t.Errorf("Expected code %s, got %s", CodeNotFound, err.Code)
	}
	if err.Message != "Student not found" {
		t.Errorf("Expected message 'Student not found', got '%s'", err.Message)
	}
}

func TestWrap(t *testing.T) {
	cause := errors.New("database connection failed")
	err := Wrap(CodeInternal, "Failed to save student", cause)

	if err.Cause != cause {
		t.Error("Expected cause to be set")
	}
	if !errors.Is(err, cause) {
		t.Error("Expected errors.Is to find cause")
	}
}

func TestIsHelpers(t *testing.T) {
	notFound := NotFound("Student", "123")
	if !IsNotFound(notFound) {
		t.Error("Expected IsNotFound to return true")
	}

	validation := Validation("Invalid NISN", nil)
	if !IsValidation(validation) {
		t.Error("Expected IsValidation to return true")
	}

	duplicate := Duplicate("Student", "NISN", "1234567890")
	if !IsDuplicate(duplicate) {
		t.Error("Expected IsDuplicate to return true")
	}
}

func TestErrorString(t *testing.T) {
	err := New(CodeNotFound, "Student not found", nil)
	expected := "[NOT_FOUND] Student not found"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}
}

func TestConvenienceConstructors(t *testing.T) {
	tests := []struct {
		name     string
		err      *AppError
		code     string
		contains string
	}{
		{"NotFound", NotFound("Student", "123"), CodeNotFound, "123"},
		{"Validation", Validation("Invalid", nil), CodeValidation, "Invalid"},
		{"Duplicate", Duplicate("Student", "NISN", "123"), CodeDuplicate, "NISN"},
		{"InvalidContext", InvalidContext("Missing school"), CodeInvalidContext, "Missing"},
		{"BusinessRule", BusinessRule("Cannot enroll"), CodeBusinessRule, "Cannot"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.code {
				t.Errorf("Expected code %s, got %s", tt.code, tt.err.Code)
			}
		})
	}
}