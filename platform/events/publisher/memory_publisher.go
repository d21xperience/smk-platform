package publisher

import (
	"context"
	"sync"

	"sekolah-platform/platform/events/types"
)

// MemoryPublisher adalah in-memory implementation dari Publisher.
// Digunakan untuk testing dan single-binary deployment.
type MemoryPublisher struct {
	mu       sync.RWMutex
	handlers []func(types.DomainEvent)
	closed   bool
}

// NewMemoryPublisher membuat MemoryPublisher baru.
func NewMemoryPublisher() *MemoryPublisher {
	return &MemoryPublisher{}
}

// Publish mengimplementasikan Publisher.
func (p *MemoryPublisher) Publish(ctx context.Context, event types.DomainEvent) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.closed {
		return errPublisherClosed
	}

	for _, handler := range p.handlers {
		handler(event)
	}
	return nil
}

// PublishBatch mengimplementasikan Publisher.
func (p *MemoryPublisher) PublishBatch(ctx context.Context, events []types.DomainEvent) error {
	for _, event := range events {
		if err := p.Publish(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// Close mengimplementasikan Publisher.
func (p *MemoryPublisher) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.closed = true
	p.handlers = nil
	return nil
}

// OnPublish mendaftarkan handler yang dipanggil setiap kali event di-publish.
// Digunakan untuk testing dan untuk menghubungkan ke MemorySubscriber.
func (p *MemoryPublisher) OnPublish(handler func(types.DomainEvent)) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.handlers = append(p.handlers, handler)
}

// PublishedCount mengembalikan jumlah handler yang terdaftar.
func (p *MemoryPublisher) PublishedCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.handlers)
}

var errPublisherClosed = errString("publisher sudah ditutup")

type errString string

func (e errString) Error() string { return string(e) }
