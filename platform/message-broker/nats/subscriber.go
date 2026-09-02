package nats

import (
	"context"
	"fmt"
	"sync"
	"time"

	"sekolah-platform/platform/events/subscriber"

	"github.com/nats-io/nats.go"
)

// NATSSubscriber adalah NATS JetStream implementation dari Subscriber.
type NATSSubscriber struct {
	conn          *Connection
	streamName    string
	durableName   string
	mu            sync.Mutex
	subscriptions []*natsSubscription
	closed        bool
}

// natsSubscription wrapper untuk NATS subscription.
type natsSubscription struct {
	sub    *nats.Subscription
	cancel context.CancelFunc
}

// NewNATSSubscriber membuat NATSSubscriber baru.
func NewNATSSubscriber(conn *Connection, streamName, durableName string) *NATSSubscriber {
	return &NATSSubscriber{
		conn:        conn,
		streamName:  streamName,
		durableName: durableName,
	}
}

// Subscribe mengimplementasikan subscriber.Subscriber.
func (s *NATSSubscriber) Subscribe(eventName string, handler subscriber.Handler) (subscriber.Subscription, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil, errSubscribeFailed
	}

	if !s.conn.IsConnected() {
		return nil, errNotConnected
	}

	js := s.conn.JetStream()

	// Determine subject
	var subject string
	if eventName == "*" {
		subject = "sdp.events.>"
	} else {
		subject = EventToSubject(eventName)
	}

	// Create context for this subscription
	ctx, cancel := context.WithCancel(context.Background())

	// Subscribe dengan queue group (consumer group)
	// Queue group = durable name, sehingga semua instance dengan durable name yang sama
	// akan share load (load balancing)
	sub, err := js.Subscribe(
		subject,
		func(msg *nats.Msg) {
			s.handleMessage(ctx, msg, handler)
		},
		nats.Durable(s.durableName),
		nats.ManualAck(),
		nats.AckWait(30*time.Second),
		nats.MaxAckPending(1000),
	)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("%w: %v", errSubscribeFailed, err)
	}

	natsSub := &natsSubscription{
		sub:    sub,
		cancel: cancel,
	}
	s.subscriptions = append(s.subscriptions, natsSub)

	s.conn.logger.Info("Subscribed to %s with durable %s", subject, s.durableName)

	return &natsSubscriptionWrapper{
		subscriber: s,
		natsSub:    natsSub,
	}, nil
}

// handleMessage menangani message dari NATS.
func (s *NATSSubscriber) handleMessage(ctx context.Context, msg *nats.Msg, handler subscriber.Handler) {
	// Convert NATS message ke DomainEvent
	event, err := BytesToEvent(msg.Data)
	if err != nil {
		s.conn.logger.Error("Failed to parse event from message: %v", err)
		// Ack untuk remove dari queue (invalid message)
		_ = msg.Ack()
		return
	}

	// Call handler
	if err := handler(ctx, event); err != nil {
		s.conn.logger.Error("Handler failed for event %s: %v", event.GetEventName(), err)
		// Nak untuk redelivery
		_ = msg.Nak()
		return
	}

	// Ack untuk konfirmasi sukses
	if err := msg.Ack(); err != nil {
		s.conn.logger.Error("Failed to ack message: %v", err)
	}
}

// Close mengimplementasikan subscriber.Subscriber.
func (s *NATSSubscriber) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.closed = true

	for _, sub := range s.subscriptions {
		sub.cancel()
		if err := sub.sub.Unsubscribe(); err != nil {
			s.conn.logger.Error("Failed to unsubscribe: %v", err)
		}
	}
	s.subscriptions = nil

	return nil
}

// ReplayFrom memulai replay dari sequence tertentu.
func (s *NATSSubscriber) ReplayFrom(ctx context.Context, startSeq uint64, handler subscriber.Handler) error {
	if !s.conn.IsConnected() {
		return errNotConnected
	}

	js := s.conn.JetStream()

	subject := "sdp.events.>"
	ctxSub, cancel := context.WithCancel(ctx)
	defer cancel()

	sub, err := js.Subscribe(
		subject,
		func(msg *nats.Msg) {
			s.handleMessage(ctxSub, msg, handler)
		},
		nats.Durable(s.durableName),
		nats.ManualAck(),
		nats.StartSequence(startSeq),
		nats.AckWait(30*time.Second),
	)
	if err != nil {
		return fmt.Errorf("%w: %v", errSubscribeFailed, err)
	}

	// Tunggu sampai context di-cancel
	<-ctxSub.Done()
	_ = sub.Unsubscribe()
	return nil
}

// natsSubscriptionWrapper adalah wrapper untuk Subscription interface.
type natsSubscriptionWrapper struct {
	subscriber *NATSSubscriber
	natsSub    *natsSubscription
}

// Unsubscribe mengimplementasikan subscriber.Subscription.
func (w *natsSubscriptionWrapper) Unsubscribe() error {
	w.subscriber.mu.Lock()
	defer w.subscriber.mu.Unlock()

	w.natsSub.cancel()
	err := w.natsSub.sub.Unsubscribe()

	// Remove from subscriptions list
	for i, sub := range w.subscriber.subscriptions {
		if sub == w.natsSub {
			w.subscriber.subscriptions = append(
				w.subscriber.subscriptions[:i],
				w.subscriber.subscriptions[i+1:]...,
			)
			break
		}
	}

	return err
}
