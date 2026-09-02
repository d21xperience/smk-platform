package nats

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"sekolah-platform/platform/events/types"
)

func TestNATSSubscriberSubscribe(t *testing.T) {
	srv, url := runTestNATSServer(t)
	defer srv.Shutdown()

	conn := setupTestConnection(t, url)
	defer conn.Close()

	js := conn.JetStream()
	cfg := natsStreamConfig("SDP_EVENTS", []string{"sdp.events.>"})
	_, _ = js.AddStream(&cfg)

	sub := NewNATSSubscriber(conn, "SDP_EVENTS", "test-consumer")
	defer sub.Close()

	var received int32
	handler := func(ctx context.Context, event types.DomainEvent) error {
		atomic.AddInt32(&received, 1)
		return nil
	}

	_, err := sub.Subscribe("StudentCreated", handler)
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	// Publish event
	pub := NewNATSPublisher(conn, "SDP_EVENTS")
	event := types.NewBaseEventBuilder("StudentCreated", "Student", "student-001").
		WithEventID("evt-001").
		Build()
	_ = pub.Publish(context.Background(), event)

	// Wait for delivery
	time.Sleep(200 * time.Millisecond)

	if atomic.LoadInt32(&received) != 1 {
		t.Errorf("Expected 1 message received, got %d", received)
	}
}

func TestNATSSubscriberGlobalSubscription(t *testing.T) {
	srv, url := runTestNATSServer(t)
	defer srv.Shutdown()

	conn := setupTestConnection(t, url)
	defer conn.Close()

	js := conn.JetStream()
	// _, _ = js.AddStream(&natsStreamConfig("SDP_EVENTS", []string{"sdp.events.>"}))
	cfg := natsStreamConfig("SDP_EVENTS", []string{"sdp.events.>"})
	_, _ = js.AddStream(&cfg)
	sub := NewNATSSubscriber(conn, "SDP_EVENTS", "global-consumer")
	defer sub.Close()

	var received int32
	handler := func(ctx context.Context, event types.DomainEvent) error {
		atomic.AddInt32(&received, 1)
		return nil
	}

	_, _ = sub.Subscribe("*", handler)

	// Publish multiple events
	pub := NewNATSPublisher(conn, "SDP_EVENTS")
	_ = pub.Publish(context.Background(), types.NewBaseEventBuilder("Event1", "Agg", "id-1").Build())
	_ = pub.Publish(context.Background(), types.NewBaseEventBuilder("Event2", "Agg", "id-2").Build())
	_ = pub.Publish(context.Background(), types.NewBaseEventBuilder("Event3", "Agg", "id-3").Build())

	time.Sleep(300 * time.Millisecond)

	if atomic.LoadInt32(&received) != 3 {
		t.Errorf("Expected 3 messages, got %d", received)
	}
}

func TestNATSSubscriberConsumerGroups(t *testing.T) {
	srv, url := runTestNATSServer(t)
	defer srv.Shutdown()

	conn := setupTestConnection(t, url)
	defer conn.Close()

	js := conn.JetStream()
	// _, _ = js.AddStream(&natsStreamConfig("SDP_EVENTS", []string{"sdp.events.>"}))
	cfg := natsStreamConfig("SDP_EVENTS", []string{"sdp.events.>"})
	_, _ = js.AddStream(&cfg)

	// Dua subscriber dengan durable name berbeda = dua consumer groups
	sub1 := NewNATSSubscriber(conn, "SDP_EVENTS", "consumer-group-1")
	defer sub1.Close()
	sub2 := NewNATSSubscriber(conn, "SDP_EVENTS", "consumer-group-2")
	defer sub2.Close()

	var received1, received2 int32

	handler1 := func(ctx context.Context, event types.DomainEvent) error {
		atomic.AddInt32(&received1, 1)
		return nil
	}
	handler2 := func(ctx context.Context, event types.DomainEvent) error {
		atomic.AddInt32(&received2, 1)
		return nil
	}

	_, _ = sub1.Subscribe("StudentCreated", handler1)
	_, _ = sub2.Subscribe("StudentCreated", handler2)

	// Publish 1 event
	pub := NewNATSPublisher(conn, "SDP_EVENTS")
	_ = pub.Publish(context.Background(), types.NewBaseEventBuilder("StudentCreated", "Student", "id-1").Build())

	time.Sleep(300 * time.Millisecond)

	// Kedua consumer group harus menerima event
	if atomic.LoadInt32(&received1) != 1 {
		t.Errorf("Consumer group 1: expected 1, got %d", received1)
	}
	if atomic.LoadInt32(&received2) != 1 {
		t.Errorf("Consumer group 2: expected 1, got %d", received2)
	}
}
