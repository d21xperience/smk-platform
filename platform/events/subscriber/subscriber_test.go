package subscriber

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"sekolah-platform/platform/events/types"
)

func createTestEvent(name string) types.DomainEvent {
	return types.NewBaseEventBuilder(name, "TestAggregate", "agg-001").
		WithEventID("evt-001").
		Build()
}

func TestMemorySubscriberSubscribe(t *testing.T) {
	s := NewMemorySubscriber()
	var called int32

	handler := func(ctx context.Context, event types.DomainEvent) error {
		atomic.AddInt32(&called, 1)
		return nil
	}

	_, err := s.Subscribe("StudentCreated", handler)
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	event := createTestEvent("StudentCreated")
	if err := s.Dispatch(context.Background(), event); err != nil {
		t.Fatalf("Dispatch failed: %v", err)
	}

	if atomic.LoadInt32(&called) != 1 {
		t.Errorf("Expected handler called once, got %d", called)
	}
}

func TestMemorySubscriberGlobalHandler(t *testing.T) {
	s := NewMemorySubscriber()
	var called int32

	handler := func(ctx context.Context, event types.DomainEvent) error {
		atomic.AddInt32(&called, 1)
		return nil
	}

	_, _ = s.Subscribe("*", handler)

	s.Dispatch(context.Background(), createTestEvent("Event1"))
	s.Dispatch(context.Background(), createTestEvent("Event2"))
	s.Dispatch(context.Background(), createTestEvent("Event3"))

	if atomic.LoadInt32(&called) != 3 {
		t.Errorf("Expected global handler called 3 times, got %d", called)
	}
}

func TestMemorySubscriberUnsubscribe(t *testing.T) {
	s := NewMemorySubscriber()
	var called int32

	handler := func(ctx context.Context, event types.DomainEvent) error {
		atomic.AddInt32(&called, 1)
		return nil
	}

	sub, _ := s.Subscribe("StudentCreated", handler)

	s.Dispatch(context.Background(), createTestEvent("StudentCreated"))
	_ = sub.Unsubscribe()
	s.Dispatch(context.Background(), createTestEvent("StudentCreated"))

	if atomic.LoadInt32(&called) != 1 {
		t.Errorf("Expected handler called once, got %d", called)
	}
}

func TestMemorySubscriberHandlerError(t *testing.T) {
	s := NewMemorySubscriber()

	handler := func(ctx context.Context, event types.DomainEvent) error {
		return errors.New("handler error")
	}

	_, _ = s.Subscribe("StudentCreated", handler)

	err := s.Dispatch(context.Background(), createTestEvent("StudentCreated"))
	if err == nil {
		t.Error("Expected error from handler")
	}
}

func TestMemorySubscriberClose(t *testing.T) {
	s := NewMemorySubscriber()

	handler := func(ctx context.Context, event types.DomainEvent) error {
		return nil
	}

	_, err := s.Subscribe("StudentCreated", handler)
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	_ = s.Close()

	_, err = s.Subscribe("AnotherEvent", handler)
	if err == nil {
		t.Error("Expected error after close")
	}
}
