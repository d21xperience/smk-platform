package outbox

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// MemoryOutbox adalah in-memory implementation dari Outbox.
type MemoryOutbox struct {
	mu        sync.RWMutex
	entries   map[string]*OutboxEntry
	idCounter int64
}

// NewMemoryOutbox membuat MemoryOutbox baru.
func NewMemoryOutbox() *MemoryOutbox {
	return &MemoryOutbox{
		entries: make(map[string]*OutboxEntry),
	}
}

// Save mengimplementasikan Outbox.
func (o *MemoryOutbox) Save(ctx context.Context, entry OutboxEntry) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	if entry.ID == "" {
		entry.ID = fmt.Sprintf("outbox-%d", atomic.AddInt64(&o.idCounter, 1))
	}

	entryCopy := entry
	o.entries[entry.ID] = &entryCopy
	return nil
}

// FetchBatch mengimplementasikan Outbox.
func (o *MemoryOutbox) FetchBatch(ctx context.Context, batchSize int) ([]OutboxEntry, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	o.mu.RLock()
	defer o.mu.RUnlock()

	now := timeNow()
	result := []OutboxEntry{}

	for _, entry := range o.entries {
		// Skip entry yang masih dalam retry delay
		if !entry.NextRetry.IsZero() && now.Before(entry.NextRetry) {
			continue
		}

		result = append(result, *entry)
		if len(result) >= batchSize {
			break
		}
	}

	return result, nil
}

// Remove mengimplementasikan Outbox.
func (o *MemoryOutbox) Remove(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	delete(o.entries, id)
	return nil
}

// UpdateAttempt mengimplementasikan Outbox.
func (o *MemoryOutbox) UpdateAttempt(ctx context.Context, id string, attempts int, lastError string, nextRetry time.Time) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	entry, exists := o.entries[id]
	if !exists {
		return nil // Entry sudah dihapus, skip
	}

	entry.Attempts = attempts
	entry.LastError = lastError
	entry.NextRetry = nextRetry
	return nil
}

// Count mengimplementasikan Outbox.
func (o *MemoryOutbox) Count(ctx context.Context) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	o.mu.RLock()
	defer o.mu.RUnlock()
	return len(o.entries), nil
}

// Clear menghapus semua entries (untuk testing).
func (o *MemoryOutbox) Clear() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.entries = make(map[string]*OutboxEntry)
}

// Entries mengembalikan semua entries (untuk testing/inspection).
func (o *MemoryOutbox) Entries() []OutboxEntry {
	o.mu.RLock()
	defer o.mu.RUnlock()

	result := make([]OutboxEntry, 0, len(o.entries))
	for _, entry := range o.entries {
		result = append(result, *entry)
	}
	return result
}
