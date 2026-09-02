package templates

import (
	"bytes"
	"fmt"
	"text/template"
)

// TemplateEngine bertanggung jawab merender template dengan data dinamis.
type TemplateEngine struct {
	registry *Registry
}

// NewTemplateEngine membuat engine baru dengan registry.
func NewTemplateEngine(registry *Registry) *TemplateEngine {
	return &TemplateEngine{
		registry: registry,
	}
}

// Register menambahkan template ke registry (convenience method).
func (te *TemplateEngine) Register(tmpl Template) error {
	return te.registry.Register(tmpl)
}

// Render menghasilkan subject dan body final dari template ID dan data.
func (te *TemplateEngine) Render(templateID string, data map[string]interface{}) (subject, body string, err error) {
	// Ambil template dari registry
	tmpl, exists := te.registry.Get(templateID)
	if !exists {
		return "", "", fmt.Errorf("template dengan ID '%s' tidak ditemukan", templateID)
	}

	// Render Subject (jika ada)
	if tmpl.Subject != "" {
		subjectTmpl, err := template.New(templateID + "_subject").Parse(tmpl.Subject)
		if err != nil {
			return "", "", fmt.Errorf("gagal parse subject template %s: %w", templateID, err)
		}

		var subjectBuf bytes.Buffer
		if err := subjectTmpl.Execute(&subjectBuf, data); err != nil {
			return "", "", fmt.Errorf("gagal execute subject template %s: %w", templateID, err)
		}
		subject = subjectBuf.String()
	}

	// Render Body
	bodyTmpl, err := template.New(templateID + "_body").Parse(tmpl.Body)
	if err != nil {
		return "", "", fmt.Errorf("gagal parse body template %s: %w", templateID, err)
	}

	var bodyBuf bytes.Buffer
	if err := bodyTmpl.Execute(&bodyBuf, data); err != nil {
		return "", "", fmt.Errorf("gagal execute body template %s: %w", templateID, err)
	}
	body = bodyBuf.String()

	return subject, body, nil
}

// List mengembalikan semua template yang terdaftar.
func (te *TemplateEngine) List() []Template {
	return te.registry.List()
}

// Count mengembalikan jumlah template yang terdaftar.
func (te *TemplateEngine) Count() int {
	return te.registry.Count()
}
