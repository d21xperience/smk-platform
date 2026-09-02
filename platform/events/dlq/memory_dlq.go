package dlq

import (
	"errors"
	"sync"
	"time"
)

// MemoryDLQ adalah in-memory implementation dari DLQ.
type MemoryDLQ struct {
	mu        sync.RWMutex
	letters   map[string]*DeadLetter
	maxSize   int
	idCounter int
}

// NewMemoryDLQ membuat MemoryDLQ baru.
// maxSize adalah jumlah maksimum dead letters yang disimpan (0 = unlimited).
func NewMemoryDLQ(maxSize int) *MemoryDLQ {
	return &MemoryDLQ{
		letters: make(map[string]*DeadLetter),
		maxSize: maxSize,
	}
}

// Store mengimplementasikan DLQ.
func (d *MemoryDLQ) Store(dl DeadLetter) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Check max size
	if d.maxSize > 0 && len(d.letters) >= d.maxSize {
		return errDLQFull
	}

	// Generate ID jika belum ada
	if dl.ID == "" {
		d.idCounter++
		dl.ID = generateID(d.idCounter)
	}

	if dl.FailedAt.IsZero() {
		dl.FailedAt = time.Now().UTC()
	}

	// Copy untuk menghindari mutasi eksternal
	dlCopy := dl
	if dl.OriginalEvent != nil {
		// Event disimpan apa adanya (immutable)
	}
	d.letters[dl.ID] = &dlCopy
	return nil
}

// List mengimplementasikan DLQ.
func (d *MemoryDLQ) List() []DeadLetter {
	d.mu.RLock()
	defer d.mu.RUnlock()

	result := make([]DeadLetter, 0, len(d.letters))
	for _, dl := range d.letters {
		result = append(result, *dl)
	}
	return result
}

// ListByConsumer mengimplementasikan DLQ.
func (d *MemoryDLQ) ListByConsumer(consumerName string) []DeadLetter {
	d.mu.RLock()
	defer d.mu.RUnlock()

	result := []DeadLetter{}
	for _, dl := range d.letters {
		if dl.ConsumerName == consumerName {
			result = append(result, *dl)
		}
	}
	return result
}

// Get mengimplementasikan DLQ.
func (d *MemoryDLQ) Get(id string) (*DeadLetter, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	dl, exists := d.letters[id]
	if !exists {
		return nil, errNotFound
	}

	dlCopy := *dl
	return &dlCopy, nil
}

// Remove mengimplementasikan DLQ.
func (d *MemoryDLQ) Remove(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.letters[id]; !exists {
		return errNotFound
	}
	delete(d.letters, id)
	return nil
}

// Count mengimplementasikan DLQ.
func (d *MemoryDLQ) Count() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.letters)
}

// Clear mengimplementasikan DLQ.
func (d *MemoryDLQ) Clear() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.letters = make(map[string]*DeadLetter)
	return nil
}

// generateID membuat ID sederhana untuk dead letter.
func generateID(counter int) string {
	const prefix = "dlq-"
	s := prefix
	for counter > 0 {
		s = string(rune('0'+counter%10)) + s
		counter /= 10
	}
	// Simple approach: just use prefix + counter
	return prefix + itoaSimple(counter)
}

func itoaSimple(i int) string {
	if i == 0 {
		return "0"
	}
	s := ""
	for i > 0 {
		s = string(rune('0'+i%10)) + s
		i /= 10
	}
	return s
}

var (
	errDLQFull  = errors.New("DLQ penuh")
	errNotFound = errors.New("dead letter tidak ditemukan")
)
