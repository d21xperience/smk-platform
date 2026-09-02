// Package subscriber menyediakan NATS event subscriber untuk update projection.
package subscriber

import (
	"context"
	"log"
	"sync"
	"time"

	"sekolah-platform/platform/events/types"
	"sekolah-platform/platform/message-broker/nats"
	"sekolah-platform/services/edge-bff/internal/projection"
)

// Logger interface.
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type defaultLogger struct{}

func (l *defaultLogger) Info(msg string, args ...interface{}) {
	log.Printf("[SUBSCRIBER][INFO] "+msg, args...)
}
func (l *defaultLogger) Error(msg string, args ...interface{}) {
	log.Printf("[SUBSCRIBER][ERROR] "+msg, args...)
}

// EventSubscriber adalah subscriber untuk event dari tu-core.
type EventSubscriber struct {
	subscriber  *nats.NATSSubscriber
	projections *projection.ProjectionManager
	logger      Logger
	stopCh      chan struct{}
	running     bool
	mu          sync.Mutex
}

// NewEventSubscriber membuat EventSubscriber baru.
func NewEventSubscriber(sub *nats.NATSSubscriber, projections *projection.ProjectionManager, logger Logger) *EventSubscriber {
	if logger == nil {
		logger = &defaultLogger{}
	}
	return &EventSubscriber{
		subscriber:  sub,
		projections: projections,
		logger:      logger,
		stopCh:      make(chan struct{}),
	}
}

// Start memulai subscriber.
func (s *EventSubscriber) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return nil
	}
	s.running = true
	s.mu.Unlock()

	// Subscribe semua events
	_, err := s.subscriber.Subscribe("*", s.handleEvent)
	if err != nil {
		return err
	}

	s.logger.Info("Event subscriber started")
	return nil
}

// Stop menghentikan subscriber.
func (s *EventSubscriber) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}
	s.running = false
	close(s.stopCh)
	s.logger.Info("Event subscriber stopped")
}

// handleEvent menangani event dari NATS.
func (s *EventSubscriber) handleEvent(ctx context.Context, event types.DomainEvent) error {
	eventName := event.GetEventName()
	s.logger.Info("Received event: %s (id: %s)", eventName, event.GetEventID())

	// Delegate ke projection manager
	if err := s.projections.Handle(ctx, event); err != nil {
		s.logger.Error("Failed to handle event %s: %v", eventName, err)
		return err
	}

	return nil
}

// Wait menunggu sampai subscriber berhenti.
func (s *EventSubscriber) Wait() {
	<-s.stopCh
}

// IsRunning mengecek status subscriber.
func (s *EventSubscriber) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

// Placeholder untuk time.Now (bisa di-mock di test)
var now = time.Now
