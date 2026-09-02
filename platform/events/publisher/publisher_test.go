package publisher

import (
	"context"
	"sync/atomic"
	"testing"

	"sekolah-platform/platform/events/types"
)

func createTestEvent(name string) types.DomainEvent {
	return types.NewBaseEventBuilder(name, "TestAggregate", "agg-001").
		WithEventID("evt-001").
		Build()
}

func TestMemoryPublisherPublish(t *testing.T) {
	p := NewMemoryPublisher()
	var called int32

	p.OnPublish(func(event types.DomainEvent) {
		atomic.AddInt32(&called, 1)
	})

	event := createTestEvent("StudentCreated")
	if err := p.Publish(context.Background(), event); err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	if atomic.LoadInt32(&called) != 1 {
		t.Errorf("Expected handler called once, got %d", called)
	}
}

func TestMemoryPublisherPublishBatch(t *testing.T) {
	p := NewMemoryPublisher()
	var called int32

	p.OnPublish(func(event types.DomainEvent) {
		atomic.AddInt32(&called, 1)
	})

	events := []types.DomainEvent{
		createTestEvent("Event1"),
		createTestEvent("Event2"),
		createTestEvent("Event3"),
	}

	if err := p.PublishBatch(context.Background(), events); err != nil {
		t.Fatalf("PublishBatch failed: %v", err)
	}

	if atomic.LoadInt32(&called) != 3 {
		t.Errorf("Expected handler called 3 times, got %d", called)
	}
}

func TestMemoryPublisherClose(t *testing.T) {
	p := NewMemoryPublisher()
	_ = p.Close()

	event := createTestEvent("TestEvent")
	if err := p.Publish(context.Background(), event); err == nil {
		t.Error("Expected error after close")
	}
}

func TestMemoryPublisherContextCancellation(t *testing.T) {
	p := NewMemoryPublisher()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	event := createTestEvent("TestEvent")
	if err := p.Publish(ctx, event); err == nil {
		t.Error("Expected error for cancelled context")
	}
}
