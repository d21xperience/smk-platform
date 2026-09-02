package registry

import (
	"errors"
	"sync"
	"time"

	"sekolah-platform/platform/events/retry"
)

var (
	errEmptyEventName = errors.New("eventName tidak boleh kosong")
	errInvalidVersion = errors.New("version harus >= 1")
	errEmptyProducer  = errors.New("producer tidak boleh kosong")
	errSchemaNotFound = errors.New("schema tidak ditemukan")
	errSchemaExists   = errors.New("schema sudah terdaftar")
)

// Registry adalah interface untuk Event Schema Registry.
type Registry interface {
	// Register mendaftarkan event schema baru.
	Register(schema EventSchema) error

	// Get mengambil schema berdasarkan eventName dan version.
	Get(eventName string, version int) (*EventSchema, error)

	// GetLatest mengambil schema versi terbaru untuk eventName.
	GetLatest(eventName string) (*EventSchema, error)

	// List mengembalikan semua schema yang terdaftar.
	List() []EventSchema

	// ListByProducer mengembalikan schema berdasarkan producer.
	ListByProducer(producer string) []EventSchema

	// Exists mengecek apakah schema sudah terdaftar.
	Exists(eventName string, version int) bool

	// Unregister menghapus schema (untuk testing).
	Unregister(eventName string, version int) error
}

// MemoryRegistry adalah in-memory implementation dari Registry.
type MemoryRegistry struct {
	mu      sync.RWMutex
	schemas map[string]*EventSchema // key: "eventName:vN"
}

// NewMemoryRegistry membuat MemoryRegistry baru.
func NewMemoryRegistry() *MemoryRegistry {
	return &MemoryRegistry{
		schemas: make(map[string]*EventSchema),
	}
}

// Register mengimplementasikan Registry.
func (r *MemoryRegistry) Register(schema EventSchema) error {
	if err := schema.Validate(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	key := schema.Key()
	if _, exists := r.schemas[key]; exists {
		return errSchemaExists
	}

	now := time.Now().UTC()
	schema.CreatedAt = now
	schema.UpdatedAt = now

	// Copy untuk menghindari mutasi eksternal
	schemaCopy := schema
	r.schemas[key] = &schemaCopy
	return nil
}

// Get mengimplementasikan Registry.
func (r *MemoryRegistry) Get(eventName string, version int) (*EventSchema, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := eventName + ":v" + itoa(version)
	schema, exists := r.schemas[key]
	if !exists {
		return nil, errSchemaNotFound
	}

	// Return copy
	schemaCopy := *schema
	return &schemaCopy, nil
}

// GetLatest mengimplementasikan Registry.
func (r *MemoryRegistry) GetLatest(eventName string) (*EventSchema, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var latest *EventSchema
	for _, schema := range r.schemas {
		if schema.EventName == eventName {
			if latest == nil || schema.Version > latest.Version {
				latest = schema
			}
		}
	}

	if latest == nil {
		return nil, errSchemaNotFound
	}

	schemaCopy := *latest
	return &schemaCopy, nil
}

// List mengimplementasikan Registry.
func (r *MemoryRegistry) List() []EventSchema {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]EventSchema, 0, len(r.schemas))
	for _, schema := range r.schemas {
		result = append(result, *schema)
	}
	return result
}

// ListByProducer mengimplementasikan Registry.
func (r *MemoryRegistry) ListByProducer(producer string) []EventSchema {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := []EventSchema{}
	for _, schema := range r.schemas {
		if schema.Producer == producer {
			result = append(result, *schema)
		}
	}
	return result
}

// Exists mengimplementasikan Registry.
func (r *MemoryRegistry) Exists(eventName string, version int) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := eventName + ":v" + itoa(version)
	_, exists := r.schemas[key]
	return exists
}

// Unregister mengimplementasikan Registry.
func (r *MemoryRegistry) Unregister(eventName string, version int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := eventName + ":v" + itoa(version)
	if _, exists := r.schemas[key]; !exists {
		return errSchemaNotFound
	}
	delete(r.schemas, key)
	return nil
}

// DefaultRetryPolicy adalah retry policy default untuk event.
var DefaultRetryPolicy = retry.Policy{
	MaxRetries:   3,
	InitialDelay: 1 * time.Second,
	MaxDelay:     30 * time.Second,
	BackoffType:  retry.BackoffExponential,
}
