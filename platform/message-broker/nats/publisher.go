package nats

import (
	"context"
	"fmt"

	"sekolah-platform/platform/events/types"

	"github.com/nats-io/nats.go"
)

// NATSPublisher adalah NATS JetStream implementation dari Publisher.
type NATSPublisher struct {
	conn   *Connection
	stream string // stream name untuk publish
}

// NewNATSPublisher membuat NATSPublisher baru.
func NewNATSPublisher(conn *Connection, streamName string) *NATSPublisher {
	return &NATSPublisher{
		conn:   conn,
		stream: streamName,
	}
}

// Publish mengimplementasikan publisher.Publisher.
func (p *NATSPublisher) Publish(ctx context.Context, event types.DomainEvent) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if !p.conn.IsConnected() {
		return errNotConnected
	}

	// Convert event ke subject dan bytes
	subject := EventToSubject(event.GetEventName())
	data, err := EventToBytes(event)
	if err != nil {
		return fmt.Errorf("%w: %v", errPublishFailed, err)
	}

	// Publish ke JetStream
	js := p.conn.JetStream()
	ack, err := js.Publish(subject, data)
	if err != nil {
		return fmt.Errorf("%w: %v", errPublishFailed, err)
	}

	p.conn.logger.Debug("Published event %s (seq=%d) to %s",
		event.GetEventName(), ack.Sequence, subject)
	return nil
}

// PublishBatch mengimplementasikan publisher.Publisher.
func (p *NATSPublisher) PublishBatch(ctx context.Context, events []types.DomainEvent) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	js := p.conn.JetStream()

	// Gunakan JetStream batch publish untuk efisiensi
	msgs := make([]*nats.Msg, 0, len(events))
	for _, event := range events {
		subject := EventToSubject(event.GetEventName())
		data, err := EventToBytes(event)
		if err != nil {
			return fmt.Errorf("%w: %v", errPublishFailed, err)
		}
		msgs = append(msgs, &nats.Msg{
			Subject: subject,
			Data:    data,
		})
	}

	// Publish semua messages
	for _, msg := range msgs {
		if _, err := js.PublishMsg(msg); err != nil {
			return fmt.Errorf("%w: %v", errPublishFailed, err)
		}
	}

	p.conn.logger.Debug("Published batch of %d events", len(events))
	return nil
}

// Close mengimplementasikan publisher.Publisher.
func (p *NATSPublisher) Close() error {
	// Connection di-close oleh Broker, bukan Publisher
	return nil
}
