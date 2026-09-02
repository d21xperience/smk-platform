package subscriber

import (
	"context"
	"sync"

	"sekolah-platform/platform/events/types"
)

// MemorySubscriber adalah in-memory implementation dari Subscriber.
type MemorySubscriber struct {
	mu            sync.RWMutex
	handlers      map[string][]Handler // eventName → handlers
	globalHandler []Handler            // handlers untuk "*"
	closed        bool
}

// NewMemorySubscriber membuat MemorySubscriber baru.
func NewMemorySubscriber() *MemorySubscriber {
	return &MemorySubscriber{
		handlers: make(map[string][]Handler),
	}
}

// Subscribe mengimplementasikan Subscriber.
func (s *MemorySubscriber) Subscribe(eventName string, handler Handler) (Subscription, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil, errSubscriberClosed
	}

	if eventName == "*" {
		s.globalHandler = append(s.globalHandler, handler)
	} else {
		s.handlers[eventName] = append(s.handlers[eventName], handler)
	}

	return &memorySubscription{
		subscriber: s,
		eventName:  eventName,
		handler:    handler,
	}, nil
}

// Dispatch mengirim event ke semua handler yang terdaftar.
// Dipanggil oleh MemoryPublisher untuk menghubungkan publisher dan subscriber.
func (s *MemorySubscriber) Dispatch(ctx context.Context, event types.DomainEvent) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil
	}

	eventName := event.GetEventName()

	// Dispatch ke specific handlers
	if handlers, ok := s.handlers[eventName]; ok {
		for _, handler := range handlers {
			if err := handler(ctx, event); err != nil {
				return err
			}
		}
	}

	// Dispatch ke global handlers
	for _, handler := range s.globalHandler {
		if err := handler(ctx, event); err != nil {
			return err
		}
	}

	return nil
}

// Close mengimplementasikan Subscriber.
func (s *MemorySubscriber) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.closed = true
	s.handlers = make(map[string][]Handler)
	s.globalHandler = nil
	return nil
}

// memorySubscription adalah subscription untuk MemorySubscriber.
type memorySubscription struct {
	subscriber *MemorySubscriber
	eventName  string
	handler    Handler
}

// Unsubscribe mengimplementasikan Subscription.
func (sub *memorySubscription) Unsubscribe() error {
	sub.subscriber.mu.Lock()
	defer sub.subscriber.mu.Unlock()

	if sub.eventName == "*" {
		newHandlers := []Handler{}
		for _, h := range sub.subscriber.globalHandler {
			if !sameHandler(h, sub.handler) {
				newHandlers = append(newHandlers, h)
			}
		}
		sub.subscriber.globalHandler = newHandlers
	} else {
		handlers := sub.subscriber.handlers[sub.eventName]
		newHandlers := []Handler{}
		for _, h := range handlers {
			if !sameHandler(h, sub.handler) {
				newHandlers = append(newHandlers, h)
			}
		}
		sub.subscriber.handlers[sub.eventName] = newHandlers
	}
	return nil
}

// sameHandler mengecek apakah dua handler sama (berdasarkan pointer).
func sameHandler(h1, h2 Handler) bool {
	// Simple pointer comparison via interface
	return &h1 == &h2 || (h1 == nil && h2 == nil)
}

var errSubscriberClosed = errString("subscriber sudah ditutup")

type errString string

func (e errString) Error() string { return string(e) }
