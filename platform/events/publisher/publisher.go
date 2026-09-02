// Package publisher menyediakan abstraction untuk publish domain events ke message broker.
// Implementasi saat ini: MemoryPublisher (untuk testing).
// Implementasi NATS akan ditambahkan di Langkah 3.
package publisher

import (
	"context"

	"sekolah-platform/platform/events/types"
)

// Publisher adalah interface untuk publish domain events.
type Publisher interface {
	// Publish mengirim event ke message broker.
	Publish(ctx context.Context, event types.DomainEvent) error

	// PublishBatch mengirim beberapa events sekaligus.
	PublishBatch(ctx context.Context, events []types.DomainEvent) error

	// Close menutup publisher dan membersihkan resource.
	Close() error
}
