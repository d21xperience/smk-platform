package templates

import (
	"testing"
)

func TestRender(t *testing.T) {
	registry := NewRegistry()
	engine := NewTemplateEngine(registry)

	// Register template
	err := engine.Register(Template{
		ID:      "test_template",
		Subject: "Hello {{.name}}",
		Body:    "Welcome {{.name}}, your invoice is Rp {{.amount}}",
	})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Render
	subject, body, err := engine.Render("test_template", map[string]interface{}{
		"name":   "Budi",
		"amount": "500.000",
	})
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if subject != "Hello Budi" {
		t.Errorf("Expected subject 'Hello Budi', got '%s'", subject)
	}

	expectedBody := "Welcome Budi, your invoice is Rp 500.000"
	if body != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, body)
	}
}

func TestRenderNotFound(t *testing.T) {
	registry := NewRegistry()
	engine := NewTemplateEngine(registry)

	_, _, err := engine.Render("nonexistent", nil)
	if err == nil {
		t.Error("Expected error for nonexistent template")
	}
}

func TestRenderEmptySubject(t *testing.T) {
	registry := NewRegistry()
	engine := NewTemplateEngine(registry)

	// Register template tanpa subject (untuk WhatsApp)
	err := engine.Register(Template{
		ID:      "wa_template",
		Subject: "",
		Body:    "Hello {{.name}}",
	})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	subject, body, err := engine.Render("wa_template", map[string]interface{}{
		"name": "Budi",
	})
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if subject != "" {
		t.Errorf("Expected empty subject, got '%s'", subject)
	}

	if body != "Hello Budi" {
		t.Errorf("Expected body 'Hello Budi', got '%s'", body)
	}
}

func TestDefaultTemplates(t *testing.T) {
	templates := DefaultTemplates()

	if len(templates) == 0 {
		t.Error("Expected default templates")
	}

	// Register all default templates
	registry := NewRegistry()
	engine := NewTemplateEngine(registry)
	for _, tmpl := range templates {
		if err := engine.Register(tmpl); err != nil {
			t.Errorf("Failed to register template %s: %v", tmpl.ID, err)
		}
	}

	// Verify all registered
	list := engine.List()
	if len(list) != len(templates) {
		t.Errorf("Expected %d templates, got %d", len(templates), len(list))
	}

	// Verify specific templates exist
	expectedIDs := []string{
		"invoice_created_wa",
		"invoice_created_email",
		"payment_received_wa",
		"payment_received_email",
	}

	for _, id := range expectedIDs {
		if _, exists := registry.Get(id); !exists {
			t.Errorf("Expected template %s to exist", id)
		}
	}
}

func TestRenderWithComplexData(t *testing.T) {
	registry := NewRegistry()
	engine := NewTemplateEngine(registry)

	err := engine.Register(Template{
		ID:      "complex_template",
		Subject: "Report for {{.schoolName}}",
		Body:    "School: {{.schoolName}}, Period: {{.period}}, Students: {{.studentCount}}",
	})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	data := map[string]interface{}{
		"schoolName":   "SMP Negeri 1",
		"period":       "2026/2027",
		"studentCount": 500,
	}

	subject, body, err := engine.Render("complex_template", data)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if subject != "Report for SMP Negeri 1" {
		t.Errorf("Expected subject 'Report for SMP Negeri 1', got '%s'", subject)
	}

	expectedBody := "School: SMP Negeri 1, Period: 2026/2027, Students: 500"
	if body != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, body)
	}
}
