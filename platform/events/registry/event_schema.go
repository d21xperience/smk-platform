// Package registry menyediakan Event Schema Registry untuk mengelola metadata event.
// Registry ini menyimpan informasi tentang setiap event type:
// - Versi dan skema payload
// - Producer dan consumer
// - Retention policy
// - Retry policy
//
// Registry digunakan untuk:
// - Validasi event sebelum publish
// - Lookup retry policy untuk consumer
// - Monitoring dan audit
// - Dokumentasi otomatis
package registry

import (
	"time"

	"sekolah-platform/platform/events/retry"
)

// EventSchema adalah metadata lengkap untuk satu event type.
type EventSchema struct {
	// Identitas
	EventName string `json:"eventName"`
	Version   int    `json:"version"`

	// Payload
	PayloadSchema string `json:"payloadSchema"` // JSON Schema (optional)

	// Ownership
	Producer  string   `json:"producer"`  // Service yang produce (e.g. "tu-core")
	Consumers []string `json:"consumers"` // Services yang consume (e.g. ["edge-guru", "edge-website"])

	// Lifecycle
	Retention time.Duration `json:"retention"` // Berapa lama event disimpan di broker

	// Resilience
	RetryPolicy retry.Policy `json:"retryPolicy"`
	DLQEnabled  bool         `json:"dlqEnabled"`

	// Metadata
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Validate memvalidasi schema.
func (s *EventSchema) Validate() error {
	if s.EventName == "" {
		return errEmptyEventName
	}
	if s.Version < 1 {
		return errInvalidVersion
	}
	if s.Producer == "" {
		return errEmptyProducer
	}
	return nil
}

// Key mengembalikan unique key untuk schema (eventName:version).
func (s *EventSchema) Key() string {
	return s.EventName + ":v" + itoa(s.Version)
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	s := ""
	for i > 0 {
		s = string(rune('0'+i%10)) + s
		i /= 10
	}
	return s
}
