package projection

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// ProjectionManager mengelola semua projection handlers.
type ProjectionManager struct {
	mu          sync.RWMutex
	projections map[string]ProjectionMetadata // key: event type
}

// NewProjectionManager membuat instance baru dari ProjectionManager.
func NewProjectionManager() *ProjectionManager {
	return &ProjectionManager{
		projections: make(map[string]ProjectionMetadata),
	}
}

// Register mendaftarkan projection handler ke manager.
func (m *ProjectionManager) Register(handler ProjectionHandler) error {
	if handler == nil {
		return fmt.Errorf("projection handler tidak boleh nil")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	events := handler.SupportedEvents()
	if len(events) == 0 {
		return fmt.Errorf("projection %s tidak mendukung event apapun", handler.Name())
	}

	for _, eventType := range events {
		if _, exists := m.projections[eventType]; exists {
			return fmt.Errorf("%w: event type '%s' sudah didaftarkan", ErrProjectionAlreadyRegistered, eventType)
		}

		m.projections[eventType] = ProjectionMetadata{
			Name:            handler.Name(),
			SupportedEvents: events,
			Handler:         handler,
		}
	}

	log.Printf("[ProjectionManager] Registered projection '%s' untuk events: %v", handler.Name(), events)
	return nil
}

// GetHandler mengambil handler berdasarkan event type.
func (m *ProjectionManager) GetHandler(eventType string) (ProjectionHandler, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	meta, exists := m.projections[eventType]
	if !exists {
		return nil, fmt.Errorf("%w: event type '%s'", ErrProjectionNotFound, eventType)
	}

	return meta.Handler, nil
}

// HandleEvent mendelegasikan event ke handler yang sesuai.
func (m *ProjectionManager) HandleEvent(ctx context.Context, event DomainEvent) error {
	handler, err := m.GetHandler(event.EventType)
	if err != nil {
		return err
	}

	if err := handler.HandleEvent(ctx, event); err != nil {
		return fmt.Errorf("%w: projection '%s' gagal memproses event '%s': %v",
			ErrProjectionFailed, handler.Name(), event.EventType, err)
	}

	return nil
}

// ListProjections mengembalikan daftar semua projection yang terdaftar.
func (m *ProjectionManager) ListProjections() []ProjectionMetadata {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Deduplicate berdasarkan handler name
	seen := make(map[string]bool)
	var result []ProjectionMetadata

	for _, meta := range m.projections {
		if !seen[meta.Name] {
			seen[meta.Name] = true
			result = append(result, meta)
		}
	}

	return result
}

// GetSupportedEvents mengembalikan semua event type yang didukung.
func (m *ProjectionManager) GetSupportedEvents() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	events := make([]string, 0, len(m.projections))
	for eventType := range m.projections {
		events = append(events, eventType)
	}
	return events
}
