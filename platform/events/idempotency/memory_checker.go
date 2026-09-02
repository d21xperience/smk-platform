package idempotency

import (
	"sync"
	"time"
)

type entry struct {
	expiresAt time.Time
}

// MemoryChecker adalah in-memory implementation dari Checker.
// Menggunakan map dengan TTL untuk auto-cleanup.
type MemoryChecker struct {
	mu      sync.Mutex
	entries map[string]entry
}

// NewMemoryChecker membuat MemoryChecker baru.
func NewMemoryChecker() *MemoryChecker {
	c := &MemoryChecker{
		entries: make(map[string]entry),
	}
	// Start cleanup goroutine
	go c.cleanupLoop()
	return c
}

// Mark mengimplementasikan Checker.
func (c *MemoryChecker) Mark(key string, ttl time.Duration) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check apakah sudah ada dan belum expired
	if existing, exists := c.entries[key]; exists {
		if time.Now().Before(existing.expiresAt) {
			return false, nil // Sudah diproses
		}
		// Expired, boleh di-mark ulang
	}

	c.entries[key] = entry{
		expiresAt: time.Now().Add(ttl),
	}
	return true, nil
}

// IsProcessed mengimplementasikan Checker.
func (c *MemoryChecker) IsProcessed(key string) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	existing, exists := c.entries[key]
	if !exists {
		return false, nil
	}
	return time.Now().Before(existing.expiresAt), nil
}

// Clear mengimplementasikan Checker.
func (c *MemoryChecker) Clear(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.entries, key)
	return nil
}

// ClearAll mengimplementasikan Checker.
func (c *MemoryChecker) ClearAll() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]entry)
	return nil
}

// cleanupLoop menghapus entry yang sudah expired secara periodik.
func (c *MemoryChecker) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanup()
	}
}

// cleanup menghapus entry yang sudah expired.
func (c *MemoryChecker) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.entries {
		if now.After(entry.expiresAt) {
			delete(c.entries, key)
		}
	}
}

// Size mengembalikan jumlah entry (untuk testing/monitoring).
func (c *MemoryChecker) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.entries)
}
