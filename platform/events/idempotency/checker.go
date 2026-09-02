// Package idempotency menyediakan idempotency checker untuk consumer.
// Idempotency memastikan event yang sama tidak diproses ganda,
// meskipun event di-deliver berkali-kali (at-least-once delivery).
package idempotency

import "time"

// Checker adalah interface untuk idempotency checker.
type Checker interface {
	// Mark menandai event sebagai sudah diproses.
	// Return true jika berhasil di-mark (belum pernah diproses sebelumnya).
	// Return false jika sudah pernah diproses.
	Mark(key string, ttl time.Duration) (bool, error)

	// IsProcessed mengecek apakah event sudah diproses.
	IsProcessed(key string) (bool, error)

	// Clear menghapus标记 (untuk testing).
	Clear(key string) error

	// ClearAll menghapus semua标记 (untuk testing).
	ClearAll() error
}

// GenerateKey membuat idempotency key dari event.
// Format: "eventName:eventID"
func GenerateKey(eventName, eventID string) string {
	return eventName + ":" + eventID
}
