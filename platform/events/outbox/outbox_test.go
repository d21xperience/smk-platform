package outbox

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"sekolah-platform/platform/events/dlq"
	"sekolah-platform/platform/events/publisher"
	"sekolah-platform/platform/events/retry"
	"sekolah-platform/platform/events/types"
)

func createTestEvent(name string) types.DomainEvent {
	return types.NewBaseEventBuilder(name, "TestAggregate", "agg-001").
		WithEventID("evt-001").
		Build()
}

// === MemoryOutbox Tests ===

func TestMemoryOutboxSave(t *testing.T) {
	o := NewMemoryOutbox()
	ctx := context.Background()

	entry := OutboxEntry{
		Event:     createTestEvent("StudentCreated"),
		CreatedAt: time.Now(),
	}

	if err := o.Save(ctx, entry); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	count, _ := o.Count(ctx)
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}
}

func TestMemoryOutboxFetchBatch(t *testing.T) {
	o := NewMemoryOutbox()
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		_ = o.Save(ctx, OutboxEntry{Event: createTestEvent("Event")})
	}

	entries, err := o.FetchBatch(ctx, 3)
	if err != nil {
		t.Fatalf("FetchBatch failed: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(entries))
	}
}

func TestMemoryOutboxRemove(t *testing.T) {
	o := NewMemoryOutbox()
	ctx := context.Background()

	entry := OutboxEntry{ID: "test-id", Event: createTestEvent("Event")}
	_ = o.Save(ctx, entry)

	if err := o.Remove(ctx, "test-id"); err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	count, _ := o.Count(ctx)
	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}

func TestMemoryOutboxUpdateAttempt(t *testing.T) {
	o := NewMemoryOutbox()
	ctx := context.Background()

	entry := OutboxEntry{ID: "test-id", Event: createTestEvent("Event")}
	_ = o.Save(ctx, entry)

	nextRetry := time.Now().Add(time.Minute)
	err := o.UpdateAttempt(ctx, "test-id", 2, "connection error", nextRetry)
	if err != nil {
		t.Fatalf("UpdateAttempt failed: %v", err)
	}

	entries := o.Entries()
	if len(entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(entries))
	}
	if entries[0].Attempts != 2 {
		t.Errorf("Expected attempts 2, got %d", entries[0].Attempts)
	}
	if entries[0].LastError != "connection error" {
		t.Errorf("Expected lastError 'connection error', got '%s'", entries[0].LastError)
	}
}

// === Writer Tests ===

func TestWriterWriteEvent(t *testing.T) {
	o := NewMemoryOutbox()
	w := NewWriter(o)
	ctx := context.Background()

	event := createTestEvent("StudentCreated")
	if err := w.WriteEvent(ctx, event); err != nil {
		t.Fatalf("WriteEvent failed: %v", err)
	}

	count, _ := o.Count(ctx)
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}
}

func TestWriterWriteInTransaction(t *testing.T) {
	o := NewMemoryOutbox()
	w := NewWriter(o)
	ctx := context.Background()

	op := func(ctx context.Context) (types.DomainEvent, error) {
		// Simulasi operasi bisnis
		return createTestEvent("StudentCreated"), nil
	}

	if err := w.WriteInTransaction(ctx, op); err != nil {
		t.Fatalf("WriteInTransaction failed: %v", err)
	}

	count, _ := o.Count(ctx)
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}
}

func TestWriterWriteInTransactionError(t *testing.T) {
	o := NewMemoryOutbox()
	w := NewWriter(o)
	ctx := context.Background()

	op := func(ctx context.Context) (types.DomainEvent, error) {
		return nil, errors.New("business error")
	}

	err := w.WriteInTransaction(ctx, op)
	if err == nil {
		t.Error("Expected error")
	}

	count, _ := o.Count(ctx)
	if count != 0 {
		t.Errorf("Expected count 0 after error, got %d", count)
	}
}

func TestWriterWriteNilEvent(t *testing.T) {
	o := NewMemoryOutbox()
	w := NewWriter(o)
	ctx := context.Background()

	if err := w.WriteEvent(ctx, nil); err != nil {
		t.Fatalf("WriteEvent with nil should not error: %v", err)
	}

	count, _ := o.Count(ctx)
	if count != 0 {
		t.Errorf("Expected count 0 for nil event, got %d", count)
	}
}

// === Processor Tests ===

type failingPublisher struct {
	failCount int32
	maxFails  int
}

func (p *failingPublisher) Publish(ctx context.Context, event types.DomainEvent) error {
	if int(atomic.LoadInt32(&p.failCount)) < p.maxFails {
		atomic.AddInt32(&p.failCount, 1)
		return errors.New("publish error")
	}
	return nil
}

func (p *failingPublisher) PublishBatch(ctx context.Context, events []types.DomainEvent) error {
	for _, e := range events {
		if err := p.Publish(ctx, e); err != nil {
			return err
		}
	}
	return nil
}

func (p *failingPublisher) Close() error { return nil }

func TestProcessorSuccess(t *testing.T) {
	o := NewMemoryOutbox()
	pub := publisher.NewMemoryPublisher()
	d := dlq.NewMemoryDLQ(0)
	ctx := context.Background()

	config := ProcessorConfig{
		PollInterval: 10 * time.Millisecond,
		BatchSize:    10,
		RetryPolicy:  retry.Policy{MaxRetries: 3, InitialDelay: time.Millisecond, MaxDelay: time.Millisecond, BackoffType: retry.BackoffConstant},
		DLQ:          d,
		Logger:       &defaultLogger{},
	}

	proc := NewProcessor(o, pub, config)

	// Save event ke outbox
	_ = o.Save(ctx, OutboxEntry{Event: createTestEvent("StudentCreated")})

	proc.Start(ctx)
	time.Sleep(50 * time.Millisecond)
	proc.Stop()

	count, _ := o.Count(ctx)
	if count != 0 {
		t.Errorf("Expected outbox empty after successful publish, got %d", count)
	}
}

func TestProcessorRetryAndSuccess(t *testing.T) {
	o := NewMemoryOutbox()
	pub := &failingPublisher{maxFails: 2} // Gagal 2x, lalu sukses
	d := dlq.NewMemoryDLQ(0)
	ctx := context.Background()

	config := ProcessorConfig{
		PollInterval: 10 * time.Millisecond,
		BatchSize:    10,
		RetryPolicy:  retry.Policy{MaxRetries: 3, InitialDelay: time.Millisecond, MaxDelay: time.Millisecond, BackoffType: retry.BackoffConstant},
		DLQ:          d,
		Logger:       &defaultLogger{},
	}

	proc := NewProcessor(o, pub, config)

	_ = o.Save(ctx, OutboxEntry{Event: createTestEvent("StudentCreated")})

	proc.Start(ctx)
	time.Sleep(100 * time.Millisecond)
	proc.Stop()

	count, _ := o.Count(ctx)
	if count != 0 {
		t.Errorf("Expected outbox empty after retry success, got %d", count)
	}
	if d.Count() != 0 {
		t.Errorf("Expected DLQ empty, got %d", d.Count())
	}
}

func TestProcessorRetryExhaustedToDLQ(t *testing.T) {
	o := NewMemoryOutbox()
	pub := &failingPublisher{maxFails: 100} // Selalu gagal
	d := dlq.NewMemoryDLQ(0)
	ctx := context.Background()

	config := ProcessorConfig{
		PollInterval: 10 * time.Millisecond,
		BatchSize:    10,
		RetryPolicy:  retry.Policy{MaxRetries: 2, InitialDelay: time.Millisecond, MaxDelay: time.Millisecond, BackoffType: retry.BackoffConstant},
		DLQ:          d,
		Logger:       &defaultLogger{},
	}

	proc := NewProcessor(o, pub, config)

	_ = o.Save(ctx, OutboxEntry{Event: createTestEvent("StudentCreated")})

	proc.Start(ctx)
	time.Sleep(150 * time.Millisecond)
	proc.Stop()

	count, _ := o.Count(ctx)
	if count != 0 {
		t.Errorf("Expected outbox empty after DLQ, got %d", count)
	}
	if d.Count() != 1 {
		t.Errorf("Expected 1 dead letter, got %d", d.Count())
	}

	letters := d.List()
	if len(letters) != 1 {
		t.Fatalf("Expected 1 dead letter, got %d", len(letters))
	}
	if letters[0].FailureCount != 3 { // 1 initial + 2 retries
		t.Errorf("Expected failure count 3, got %d", letters[0].FailureCount)
	}
}

func TestProcessorContextCancellation(t *testing.T) {
	o := NewMemoryOutbox()
	pub := publisher.NewMemoryPublisher()
	ctx, cancel := context.WithCancel(context.Background())

	config := ProcessorConfig{
		PollInterval: 10 * time.Millisecond,
		BatchSize:    10,
		RetryPolicy:  retry.DefaultPolicy,
		Logger:       &defaultLogger{},
	}

	proc := NewProcessor(o, pub, config)
	proc.Start(ctx)

	cancel()
	time.Sleep(50 * time.Millisecond)

	if proc.IsRunning() {
		t.Error("Expected processor to stop after context cancellation")
	}
}
