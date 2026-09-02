package consumer

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"sekolah-platform/platform/events/dlq"
	"sekolah-platform/platform/events/idempotency"
	"sekolah-platform/platform/events/retry"
	"sekolah-platform/platform/events/types"
)

func createTestEvent(name, id string) types.DomainEvent {
	return types.NewBaseEventBuilder(name, "TestAggregate", "agg-001").
		WithEventID(id).
		Build()
}

func TestResilientConsumerSuccess(t *testing.T) {
	var called int32
	handler := func(ctx context.Context, event types.DomainEvent) error {
		atomic.AddInt32(&called, 1)
		return nil
	}

	consumer := NewResilientConsumer(handler, Config{
		Name:        "test-consumer",
		RetryPolicy: retry.Policy{MaxRetries: 3, InitialDelay: time.Millisecond, MaxDelay: time.Millisecond, BackoffType: retry.BackoffConstant},
	})

	event := createTestEvent("StudentCreated", "evt-001")
	err := consumer.Handle(context.Background(), event)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if atomic.LoadInt32(&called) != 1 {
		t.Errorf("Expected handler called once, got %d", called)
	}
}

func TestResilientConsumerRetryThenSuccess(t *testing.T) {
	var attempts int32
	handler := func(ctx context.Context, event types.DomainEvent) error {
		count := atomic.AddInt32(&attempts, 1)
		if count < 3 {
			return errors.New("temporary error")
		}
		return nil
	}

	consumer := NewResilientConsumer(handler, Config{
		Name:        "test-consumer",
		RetryPolicy: retry.Policy{MaxRetries: 3, InitialDelay: time.Millisecond, MaxDelay: time.Millisecond, BackoffType: retry.BackoffConstant},
	})

	event := createTestEvent("StudentCreated", "evt-001")
	err := consumer.Handle(context.Background(), event)

	if err != nil {
		t.Fatalf("Expected success after retry, got: %v", err)
	}
	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestResilientConsumerSendToDLQ(t *testing.T) {
	handler := func(ctx context.Context, event types.DomainEvent) error {
		return errors.New("permanent error")
	}

	d := dlq.NewMemoryDLQ(0)
	consumer := NewResilientConsumer(handler, Config{
		Name:        "test-consumer",
		RetryPolicy: retry.Policy{MaxRetries: 2, InitialDelay: time.Millisecond, MaxDelay: time.Millisecond, BackoffType: retry.BackoffConstant},
		DLQ:         d,
	})

	event := createTestEvent("StudentCreated", "evt-001")
	err := consumer.Handle(context.Background(), event)

	if err == nil {
		t.Error("Expected error after DLQ")
	}
	if d.Count() != 1 {
		t.Errorf("Expected 1 dead letter, got %d", d.Count())
	}

	letters := d.List()
	if len(letters) != 1 {
		t.Fatalf("Expected 1 dead letter, got %d", len(letters))
	}
	if letters[0].ConsumerName != "test-consumer" {
		t.Errorf("Expected consumer name 'test-consumer', got '%s'", letters[0].ConsumerName)
	}
}

func TestResilientConsumerIdempotency(t *testing.T) {
	var called int32
	handler := func(ctx context.Context, event types.DomainEvent) error {
		atomic.AddInt32(&called, 1)
		return nil
	}

	checker := idempotency.NewMemoryChecker()
	defer checker.ClearAll()

	consumer := NewResilientConsumer(handler, Config{
		Name:               "test-consumer",
		RetryPolicy:        retry.Policy{MaxRetries: 0},
		IdempotencyChecker: checker,
		IdempotencyTTL:     time.Hour,
	})

	event := createTestEvent("StudentCreated", "evt-001")

	// Call pertama
	if err := consumer.Handle(context.Background(), event); err != nil {
		t.Fatalf("First call failed: %v", err)
	}

	// Call kedua dengan event yang sama (harus skip)
	if err := consumer.Handle(context.Background(), event); err != nil {
		t.Fatalf("Second call failed: %v", err)
	}

	if atomic.LoadInt32(&called) != 1 {
		t.Errorf("Expected handler called once (idempotent), got %d", called)
	}
}

func TestResilientConsumerIdempotencyDifferentEvents(t *testing.T) {
	var called int32
	handler := func(ctx context.Context, event types.DomainEvent) error {
		atomic.AddInt32(&called, 1)
		return nil
	}

	checker := idempotency.NewMemoryChecker()
	defer checker.ClearAll()

	consumer := NewResilientConsumer(handler, Config{
		Name:               "test-consumer",
		RetryPolicy:        retry.Policy{MaxRetries: 0},
		IdempotencyChecker: checker,
		IdempotencyTTL:     time.Hour,
	})

	event1 := createTestEvent("StudentCreated", "evt-001")
	event2 := createTestEvent("StudentCreated", "evt-002")

	_ = consumer.Handle(context.Background(), event1)
	_ = consumer.Handle(context.Background(), event2)

	if atomic.LoadInt32(&called) != 2 {
		t.Errorf("Expected handler called twice for different events, got %d", called)
	}
}

func TestResilientConsumerOnRetryCallback(t *testing.T) {
	var retryCalls int32
	handler := func(ctx context.Context, event types.DomainEvent) error {
		return errors.New("error")
	}

	consumer := NewResilientConsumer(handler, Config{
		Name:        "test-consumer",
		RetryPolicy: retry.Policy{MaxRetries: 2, InitialDelay: time.Millisecond, MaxDelay: time.Millisecond, BackoffType: retry.BackoffConstant},
		OnRetry: func(attempt int, err error, delay time.Duration) {
			atomic.AddInt32(&retryCalls, 1)
		},
	})

	event := createTestEvent("StudentCreated", "evt-001")
	_ = consumer.Handle(context.Background(), event)

	if atomic.LoadInt32(&retryCalls) != 2 {
		t.Errorf("Expected 2 retry callbacks, got %d", retryCalls)
	}
}

func TestResilientConsumerOnDLQCallback(t *testing.T) {
	var dlqCalls int32
	handler := func(ctx context.Context, event types.DomainEvent) error {
		return errors.New("error")
	}

	d := dlq.NewMemoryDLQ(0)
	consumer := NewResilientConsumer(handler, Config{
		Name:        "test-consumer",
		RetryPolicy: retry.Policy{MaxRetries: 1, InitialDelay: time.Millisecond, MaxDelay: time.Millisecond, BackoffType: retry.BackoffConstant},
		DLQ:         d,
		OnDLQ: func(deadLetter dlq.DeadLetter) {
			atomic.AddInt32(&dlqCalls, 1)
		},
	})

	event := createTestEvent("StudentCreated", "evt-001")
	_ = consumer.Handle(context.Background(), event)

	if atomic.LoadInt32(&dlqCalls) != 1 {
		t.Errorf("Expected 1 DLQ callback, got %d", dlqCalls)
	}
}

func TestAsSubscriberHandler(t *testing.T) {
	handler := func(ctx context.Context, event types.DomainEvent) error {
		return nil
	}

	consumer := NewResilientConsumer(handler, Config{Name: "test"})
	subHandler := consumer.AsSubscriberHandler()

	if subHandler == nil {
		t.Error("Expected non-nil subscriber handler")
	}

	event := createTestEvent("TestEvent", "evt-001")
	if err := subHandler(context.Background(), event); err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}
