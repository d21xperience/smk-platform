// Package outbox menyediakan Transactional Outbox pattern.
// Outbox memastikan event tidak hilang saat publish ke message broker.
//
// Pattern:
// 1. Service menulis business data + event ke outbox dalam SATU transaksi
// 2. Background worker (Processor) membaca outbox dan publish ke broker
// 3. Jika publish sukses, event dihapus dari outbox
// 4. Jika publish gagal, retry sesuai policy, lalu masuk DLQ
package outbox

import (
	"context"
	"time"

	"sekolah-platform/platform/events/types"
)

// OutboxEntry adalah event yang menunggu untuk di-publish.
type OutboxEntry struct {
	// ID unik entry
	ID string `json:"id"`

	// Event yang akan di-publish
	Event types.DomainEvent `json:"event"`

	// Metadata
	CreatedAt  time.Time `json:"createdAt"`
	Attempts   int       `json:"attempts"`
	LastError  string    `json:"lastError,omitempty"`
	NextRetry  time.Time `json:"nextRetry,omitempty"`
	LockedUntil time.Time `json:"lockedUntil,omitempty"` // Untuk distributed processing
}

// Outbox adalah interface untuk outbox storage.
type Outbox interface {
	// Save menyimpan entry ke outbox (dalam transaksi yang sama dengan business data).
	Save(ctx context.Context, entry OutboxEntry) error

	// FetchBatch mengambil batch entry yang siap di-process.
	FetchBatch(ctx context.Context, batchSize int) ([]OutboxEntry, error)

	// Remove menghapus entry setelah sukses publish.
	Remove(ctx context.Context, id string) error

	// UpdateAttempt mengupdate attempt count dan last error.
	UpdateAttempt(ctx context.Context, id string, attempts int, lastError string, nextRetry time.Time) error

	// Count mengembalikan jumlah entry di outbox (untuk monitoring).
	Count(ctx context.Context) (int, error)
}
