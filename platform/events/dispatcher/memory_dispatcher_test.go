package dispatcher

import (
	"sync/atomic"
	"testing"
	"time"

	"sekolah-platform/platform/events/types"
)

func createTestEvent(name string) types.DomainEvent {
	return types.NewBaseEventBuilder(name, "TestAggregate", "agg-001").
		WithEventID("evt-001").
		WithPayload(map[string]string{"data": "test"}).
		Build()
}

func TestMemoryDispatcherSubscribe(t *testing.T) {
	d := NewMemoryDispatcher()
	var called int32

	listener := NewListenerFunc("StudentCreated", func(event types.DomainEvent) error {
		atomic.AddInt32(&called, 1)
		return nil
	})

	d.Subscribe("StudentCreated", listener)

	event := createTestEvent("StudentCreated")
	if err := d.Dispatch(event); err != nil {
		t.Fatalf("Dispatch failed: %v", err)
	}

	if atomic.LoadInt32(&called) != 1 {
		t.Errorf("Expected listener to be called once, got %d", called)
	}
}

func TestMemoryDispatcherMultipleListeners(t *testing.T) {
	d := NewMemoryDispatcher()
	var called1, called2 int32

	listener1 := NewListenerFunc("StudentCreated", func(event types.DomainEvent) error {
		atomic.AddInt32(&called1, 1)
		return nil
	})
	listener2 := NewListenerFunc("StudentCreated", func(event types.DomainEvent) error {
		atomic.AddInt32(&called2, 1)
		return nil
	})

	d.Subscribe("StudentCreated", listener1)
	d.Subscribe("StudentCreated", listener2)

	event := createTestEvent("StudentCreated")
	if err := d.Dispatch(event); err != nil {
		t.Fatalf("Dispatch failed: %v", err)
	}

	if atomic.LoadInt32(&called1) != 1 || atomic.LoadInt32(&called2) != 1 {
		t.Error("Expected both listeners to be called")
	}
}

func TestMemoryDispatcherGlobalListener(t *testing.T) {
	d := NewMemoryDispatcher()
	var called int32

	globalListener := NewListenerFunc("*", func(event types.DomainEvent) error {
		atomic.AddInt32(&called, 1)
		return nil
	})

	d.Subscribe("*", globalListener)

	// Dispatch beberapa event berbeda
	d.Dispatch(createTestEvent("StudentCreated"))
	d.Dispatch(createTestEvent("StudentUpdated"))
	d.Dispatch(createTestEvent("AttendanceSubmitted"))

	if atomic.LoadInt32(&called) != 3 {
		t.Errorf("Expected global listener to be called 3 times, got %d", called)
	}
}

func TestMemoryDispatcherNoListener(t *testing.T) {
	d := NewMemoryDispatcher()

	// Dispatch tanpa listener tidak boleh error
	event := createTestEvent("UnknownEvent")
	if err := d.Dispatch(event); err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestMemoryDispatcherUnsubscribe(t *testing.T) {
	d := NewMemoryDispatcher()
	var called int32

	listener := NewListenerFunc("StudentCreated", func(event types.DomainEvent) error {
		atomic.AddInt32(&called, 1)
		return nil
	})

	d.Subscribe("StudentCreated", listener)
	d.Dispatch(createTestEvent("StudentCreated"))

	d.Unsubscribe("StudentCreated", listener)
	d.Dispatch(createTestEvent("StudentCreated"))

	if atomic.LoadInt32(&called) != 1 {
		t.Errorf("Expected listener to be called once, got %d", called)
	}
}

func TestMemoryDispatcherDispatchAll(t *testing.T) {
	d := NewMemoryDispatcher()
	var called int32

	listener := NewListenerFunc("*", func(event types.DomainEvent) error {
		atomic.AddInt32(&called, 1)
		return nil
	})
	d.Subscribe("*", listener)

	events := []types.DomainEvent{
		createTestEvent("Event1"),
		createTestEvent("Event2"),
		createTestEvent("Event3"),
	}

	if err := d.DispatchAll(events); err != nil {
		t.Fatalf("DispatchAll failed: %v", err)
	}

	if atomic.LoadInt32(&called) != 3 {
		t.Errorf("Expected 3 calls, got %d", called)
	}
}

func TestMemoryDispatcherAsync(t *testing.T) {
	d := NewMemoryDispatcher()
	var called int32

	listener := NewListenerFunc("TestEvent", func(event types.DomainEvent) error {
		atomic.AddInt32(&called, 1)
		return nil
	})
	d.Subscribe("TestEvent", listener)

	event := createTestEvent("TestEvent")
	if err := d.DispatchAsync(event); err != nil {
		t.Fatalf("DispatchAsync failed: %v", err)
	}

	// Tunggu goroutine selesai
	time.Sleep(50 * time.Millisecond)

	if atomic.LoadInt32(&called) != 1 {
		t.Errorf("Expected async listener to be called, got %d", called)
	}
}

func TestMemoryDispatcherListenerCount(t *testing.T) {
	d := NewMemoryDispatcher()

	listener1 := NewListenerFunc("Event1", func(event types.DomainEvent) error { return nil })
	listener2 := NewListenerFunc("Event1", func(event types.DomainEvent) error { return nil })
	globalListener := NewListenerFunc("*", func(event types.DomainEvent) error { return nil })

	d.Subscribe("Event1", listener1)
	d.Subscribe("Event1", listener2)
	d.Subscribe("*", globalListener)

	if count := d.ListenerCount("Event1"); count != 3 {
		t.Errorf("Expected 3 listeners, got %d", count)
	}

	if count := d.ListenerCount("Unknown"); count != 1 {
		t.Errorf("Expected 1 global listener, got %d", count)
	}
}

func TestMemoryDispatcherReset(t *testing.T) {
	d := NewMemoryDispatcher()

	listener := NewListenerFunc("Event1", func(event types.DomainEvent) error { return nil })
	d.Subscribe("Event1", listener)

	d.Reset()

	if count := d.ListenerCount("Event1"); count != 0 {
		t.Errorf("Expected 0 listeners after reset, got %d", count)
	}
}
