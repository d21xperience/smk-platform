// Package dlq menyediakan Dead Letter Queue untuk menyimpan event yang gagal diproses.
// Event masuk ke DLQ setelah retry policy habis.
// Admin bisa inspect DLQ dan melakukan manual retry atau discard.
package dlq

import (
	"time"

	"sekolah-platform/platform/events/types"
)

// DeadLetter adalah event yang gagal diproses dan masuk ke DLQ.
type DeadLetter struct {
	// Event asli yang gagal
	OriginalEvent types.DomainEvent `json:"originalEvent"`

	// Metadata kegagalan
	ConsumerName string    `json:"consumerName"`
	FailedAt     time.Time `json:"failedAt"`
	FailureCount int       `json:"failureCount"`
	LastError    string    `json:"lastError"`

	// Identitas
	ID string `json:"id"`
}

// DLQ adalah interface untuk Dead Letter Queue.
type DLQ interface {
	// Store menyimpan event yang gagal ke DLQ.
	Store(deadLetter DeadLetter) error

	// List mengembalikan semua dead letters.
	List() []DeadLetter

	// ListByConsumer mengembalikan dead letters berdasarkan consumer.
	ListByConsumer(consumerName string) []DeadLetter

	// Get mengambil dead letter berdasarkan ID.
	Get(id string) (*DeadLetter, error)

	// Remove menghapus dead letter dari DLQ (setelah manual retry atau discard).
	Remove(id string) error

	// Count mengembalikan jumlah dead letters.
	Count() int

	// Clear menghapus semua dead letters (untuk testing).
	Clear() error
}
