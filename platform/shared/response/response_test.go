package response

import (
	"encoding/json"
	"net/http"
	"testing"

	"sekolah-platform/platform/shared/errors"
)

func TestSuccess(t *testing.T) {
	data := map[string]string{"name": "Budi"}
	resp := Success(data)

	if !resp.Success {
		t.Error("Expected success=true")
	}
	if resp.Data == nil {
		t.Error("Expected data to be set")
	}
	if resp.Error != nil {
		t.Error("Expected error to be nil")
	}
}

func TestError(t *testing.T) {
	appErr := errors.New(errors.CodeNotFound, "Student not found", nil)
	resp := Error(appErr)

	if resp.Success {
		t.Error("Expected success=false")
	}
	if resp.Data != nil {
		t.Error("Expected data to be nil")
	}
	if resp.Error == nil {
		t.Error("Expected error to be set")
	}
	if resp.Error.Code != errors.CodeNotFound {
		t.Errorf("Expected error code %s, got %s", errors.CodeNotFound, resp.Error.Code)
	}
}

func TestResponseJSON(t *testing.T) {
	resp := Success(map[string]string{"id": "123"})
	bytes, err := resp.ToJSON()
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(bytes, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if parsed["success"] != true {
		t.Error("Expected success=true in JSON")
	}
	if parsed["data"] == nil {
		t.Error("Expected data in JSON")
	}
}

func TestHTTPStatus(t *testing.T) {
	tests := []struct {
		name     string
		resp     *Response
		expected int
	}{
		{"success", Success(nil), http.StatusOK},
		{"validation error", Error(errors.New(errors.CodeValidation, "invalid", nil)), http.StatusBadRequest},
		{"not found", Error(errors.New(errors.CodeNotFound, "not found", nil)), http.StatusNotFound},
		{"duplicate", Error(errors.New(errors.CodeDuplicate, "duplicate", nil)), http.StatusConflict},
		{"unauthorized", Error(errors.New(errors.CodeUnauthorized, "unauthorized", nil)), http.StatusUnauthorized},
		{"forbidden", Error(errors.New(errors.CodeForbidden, "forbidden", nil)), http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.resp.HTTPStatus(); got != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, got)
			}
		})
	}
}

func TestEmpty(t *testing.T) {
	resp := Empty()
	if !resp.Success {
		t.Error("Expected success=true")
	}
	if resp.Data != nil {
		t.Error("Expected data=nil")
	}
	if resp.Error != nil {
		t.Error("Expected error=nil")
	}
}
