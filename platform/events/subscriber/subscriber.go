// Package subscriber menyediakan abstraction untuk subscribe domain events dari message broker.
package subscriber

import (
	"context"

	"sekolah-platform/platform/events/types"
)

// Handler adalah function untuk menangani event yang di-subscribe.
type Handler func(ctx context.Context, event types.DomainEvent) error

// Subscription merepresentasikan subscription aktif.
type Subscription interface {
	// Unsubscribe berhenti berlangganan.
	Unsubscribe() error
}

// Subscriber adalah interface untuk subscribe domain events.
type Subscriber interface {
	// Subscribe mendaftarkan handler untuk event tertentu.
	// Gunakan "*" untuk eventName untuk subscribe semua event.
	Subscribe(eventName string, handler Handler) (Subscription, error)

	// Close menutup subscriber dan semua subscription.
	Close() error
}
