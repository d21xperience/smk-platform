package registry

import (
	"testing"
	"time"

	"sekolah-platform/platform/events/retry"
)

func createTestSchema(name string, version int) EventSchema {
	return EventSchema{
		EventName: name,
		Version:   version,
		Producer:  "tu-core",
		Consumers: []string{"edge-guru", "edge-website"},
		Retention: 7 * 24 * time.Hour,
		RetryPolicy: retry.Policy{
			MaxRetries:   3,
			InitialDelay: time.Second,
			MaxDelay:     30 * time.Second,
			BackoffType:  retry.BackoffExponential,
		},
		DLQEnabled: true,
	}
}

func TestRegister(t *testing.T) {
	r := NewMemoryRegistry()
	schema := createTestSchema("StudentCreated", 1)

	if err := r.Register(schema); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if !r.Exists("StudentCreated", 1) {
		t.Error("Expected schema to exist")
	}
}

func TestRegisterInvalid(t *testing.T) {
	r := NewMemoryRegistry()

	tests := []struct {
		name   string
		schema EventSchema
	}{
		{"empty name", EventSchema{Version: 1, Producer: "p"}},
		{"invalid version", EventSchema{EventName: "e", Version: 0, Producer: "p"}},
		{"empty producer", EventSchema{EventName: "e", Version: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.Register(tt.schema); err == nil {
				t.Error("Expected error for invalid schema")
			}
		})
	}
}

func TestRegisterDuplicate(t *testing.T) {
	r := NewMemoryRegistry()
	schema := createTestSchema("StudentCreated", 1)

	if err := r.Register(schema); err != nil {
		t.Fatalf("First register failed: %v", err)
	}
	if err := r.Register(schema); err == nil {
		t.Error("Expected error for duplicate register")
	}
}

func TestGet(t *testing.T) {
	r := NewMemoryRegistry()
	schema := createTestSchema("StudentCreated", 1)
	_ = r.Register(schema)

	got, err := r.Get("StudentCreated", 1)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.EventName != "StudentCreated" {
		t.Errorf("Expected eventName 'StudentCreated', got '%s'", got.EventName)
	}
	if got.Version != 1 {
		t.Errorf("Expected version 1, got %d", got.Version)
	}
}

func TestGetNotFound(t *testing.T) {
	r := NewMemoryRegistry()

	_, err := r.Get("Unknown", 1)
	if err == nil {
		t.Error("Expected error for not found")
	}
}

func TestGetLatest(t *testing.T) {
	r := NewMemoryRegistry()
	_ = r.Register(createTestSchema("StudentCreated", 1))
	_ = r.Register(createTestSchema("StudentCreated", 2))
	_ = r.Register(createTestSchema("StudentCreated", 3))

	got, err := r.GetLatest("StudentCreated")
	if err != nil {
		t.Fatalf("GetLatest failed: %v", err)
	}
	if got.Version != 3 {
		t.Errorf("Expected latest version 3, got %d", got.Version)
	}
}

func TestList(t *testing.T) {
	r := NewMemoryRegistry()
	_ = r.Register(createTestSchema("Event1", 1))
	_ = r.Register(createTestSchema("Event2", 1))
	_ = r.Register(createTestSchema("Event1", 2))

	list := r.List()
	if len(list) != 3 {
		t.Errorf("Expected 3 schemas, got %d", len(list))
	}
}

func TestListByProducer(t *testing.T) {
	r := NewMemoryRegistry()

	s1 := createTestSchema("Event1", 1)
	s1.Producer = "tu-core"
	_ = r.Register(s1)

	s2 := createTestSchema("Event2", 1)
	s2.Producer = "edge-guru"
	_ = r.Register(s2)

	list := r.ListByProducer("tu-core")
	if len(list) != 1 {
		t.Errorf("Expected 1 schema for tu-core, got %d", len(list))
	}
}

func TestUnregister(t *testing.T) {
	r := NewMemoryRegistry()
	_ = r.Register(createTestSchema("Event1", 1))

	if err := r.Unregister("Event1", 1); err != nil {
		t.Fatalf("Unregister failed: %v", err)
	}
	if r.Exists("Event1", 1) {
		t.Error("Expected schema to be removed")
	}
}

func TestSchemaKey(t *testing.T) {
	schema := createTestSchema("StudentCreated", 2)
	expected := "StudentCreated:v2"
	if got := schema.Key(); got != expected {
		t.Errorf("Expected key '%s', got '%s'", expected, got)
	}
}
