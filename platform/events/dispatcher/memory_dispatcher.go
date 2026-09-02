package dispatcher

import (
	"fmt"
	"sync"

	"sekolah-platform/platform/events/types"
)

// MemoryDispatcher adalah in-memory implementation dari Dispatcher.
// Digunakan untuk single-binary deployment dan testing.
// Untuk distribusi ke Edge BFF, gunakan NATS publisher (Langkah 3).
type MemoryDispatcher struct {
	mu        sync.RWMutex
	listeners map[string][]Listener // eventName → listeners
	global    []Listener            // listeners untuk semua event ("*")
	async     bool                  // apakah dispatch async
}

// NewMemoryDispatcher membuat MemoryDispatcher baru.
func NewMemoryDispatcher() *MemoryDispatcher {
	return &MemoryDispatcher{
		listeners: make(map[string][]Listener),
	}
}

// Subscribe mendaftarkan listener untuk event tertentu.
func (d *MemoryDispatcher) Subscribe(eventName string, listener Listener) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if eventName == "*" {
		d.global = append(d.global, listener)
		return
	}

	d.listeners[eventName] = append(d.listeners[eventName], listener)
}

// Unsubscribe menghapus listener.
func (d *MemoryDispatcher) Unsubscribe(eventName string, listener Listener) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if eventName == "*" {
		d.global = filterListeners(d.global, listener)
		return
	}

	d.listeners[eventName] = filterListeners(d.listeners[eventName], listener)
}

// Dispatch mengimplementasikan Dispatcher.
func (d *MemoryDispatcher) Dispatch(event types.DomainEvent) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	eventName := event.GetEventName()
	var errs []error

	// Dispatch ke specific listeners
	if listeners, ok := d.listeners[eventName]; ok {
		for _, listener := range listeners {
			if err := listener.HandleEvent(event); err != nil {
				errs = append(errs, fmt.Errorf("listener error for %s: %w", eventName, err))
			}
		}
	}

	// Dispatch ke global listeners
	for _, listener := range d.global {
		if err := listener.HandleEvent(event); err != nil {
			errs = append(errs, fmt.Errorf("global listener error: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("dispatch errors: %v", errs)
	}
	return nil
}

// DispatchAsync mengimplementasikan Dispatcher.
func (d *MemoryDispatcher) DispatchAsync(event types.DomainEvent) error {
	go func() {
		_ = d.Dispatch(event)
	}()
	return nil
}

// DispatchAll mengimplementasikan Dispatcher.
func (d *MemoryDispatcher) DispatchAll(events []types.DomainEvent) error {
	var errs []error
	for _, event := range events {
		if err := d.Dispatch(event); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("dispatch all errors: %v", errs)
	}
	return nil
}

// ListenerCount mengembalikan jumlah listener untuk event tertentu.
func (d *MemoryDispatcher) ListenerCount(eventName string) int {
	d.mu.RLock()
	defer d.mu.RUnlock()

	count := len(d.global)
	if listeners, ok := d.listeners[eventName]; ok {
		count += len(listeners)
	}
	return count
}

// Reset menghapus semua listener (untuk testing).
func (d *MemoryDispatcher) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.listeners = make(map[string][]Listener)
	d.global = nil
}

func filterListeners(listeners []Listener, target Listener) []Listener {
	result := make([]Listener, 0, len(listeners))
	for _, l := range listeners {
		if l != target {
			result = append(result, l)
		}
	}
	return result
}
