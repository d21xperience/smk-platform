package dispatcher

import "sekolah-platform/platform/events/types"

// Listener adalah interface untuk listener domain events.
type Listener interface {
	// HandleEvent menangani domain event.
	HandleEvent(event types.DomainEvent) error

	// EventName mengembalikan nama event yang didengarkan.
	// Return "*" untuk mendengarkan semua event.
	EventName() string
}

// ListenerFunc adalah adapter untuk menggunakan function sebagai Listener.
type ListenerFunc struct {
	name    string
	handler func(types.DomainEvent) error
}

// NewListenerFunc membuat Listener dari function.
func NewListenerFunc(eventName string, handler func(types.DomainEvent) error) *ListenerFunc {
	return &ListenerFunc{
		name:    eventName,
		handler: handler,
	}
}

// HandleEvent mengimplementasikan Listener.
func (l *ListenerFunc) HandleEvent(event types.DomainEvent) error {
	return l.handler(event)
}

// EventName mengimplementasikan Listener.
func (l *ListenerFunc) EventName() string {
	return l.name
}
