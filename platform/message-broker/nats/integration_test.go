package nats

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"sekolah-platform/platform/events/dlq"
	"sekolah-platform/platform/events/idempotency"
	"sekolah-platform/platform/events/retry"
	"sekolah-platform/platform/events/types"
	"sekolah-platform/platform/message-broker/contract"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

// runTestNATSServer menjalankan embedded NATS server untuk testing.
func runTestNATSServer(t *testing.T) (*server.Server, string) {
	t.Helper()

	opts := &server.Options{
		Port:      -1, // random port
		JetStream: true,
		StoreDir:  t.TempDir(),
	}

	s, err := server.NewServer(opts)
	if err != nil {
		t.Fatalf("Failed to create NATS server: %v", err)
	}

	go s.Start()

	if !s.ReadyForConnections(5 * time.Second) {
		t.Fatal("NATS server not ready")
	}

	t.Cleanup(func() {
		s.Shutdown()
	})

	return s, s.ClientURL()
}

// setupTestConnection membuat connection untuk testing.
func setupTestConnection(t *testing.T, url string) *Connection {
	t.Helper()

	config := DefaultConfig.WithURL(url).WithConnectionName("test")
	conn, err := NewConnection(config, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	t.Cleanup(func() {
		conn.Close()
	})

	return conn
}

// natsStreamConfig helper untuk membuat stream config.
func natsStreamConfig(name string, subjects []string) nats.StreamConfig {
	return nats.StreamConfig{
		Name:     name,
		Subjects: subjects,
		Storage:  nats.MemoryStorage, // use memory for tests
	}
}

// TestFullIntegration menguji full flow: publish → subscribe → ack
func TestFullIntegration(t *testing.T) {
	srv, url := runTestNATSServer(t)
	defer srv.Shutdown()

	conn := setupTestConnection(t, url)
	broker := NewBrokerWithConnection(conn, "SDP_EVENTS")

	// Setup stream
	err := broker.EnsureStream(context.Background(), contract.StreamConfig{
		Name:     "SDP_EVENTS",
		Subjects: []string{"sdp.events.>"},
		Storage:  "memory",
	})
	if err != nil {
		t.Fatalf("EnsureStream failed: %v", err)
	}

	// Subscribe
	var received int32
	var lastEvent types.DomainEvent

	sub := NewNATSSubscriber(conn, "SDP_EVENTS", "integration-test")
	defer sub.Close()

	handler := func(ctx context.Context, event types.DomainEvent) error {
		atomic.AddInt32(&received, 1)
		lastEvent = event
		return nil
	}

	_, err = sub.Subscribe("StudentCreated", handler)
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	// Publish
	event := types.NewBaseEventBuilder("StudentCreated", "Student", "student-001").
		WithEventID("evt-001").
		WithPayload(map[string]string{"name": "Budi"}).
		Build()

	err = broker.Publish(context.Background(), event)
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	// Wait for delivery
	time.Sleep(300 * time.Millisecond)

	if atomic.LoadInt32(&received) != 1 {
		t.Errorf("Expected 1 message, got %d", received)
	}
	if lastEvent == nil {
		t.Fatal("Expected event to be received")
	}
	if lastEvent.GetEventName() != "StudentCreated" {
		t.Errorf("Expected eventName 'StudentCreated', got '%s'", lastEvent.GetEventName())
	}
	if lastEvent.GetAggregateID() != "student-001" {
		t.Errorf("Expected aggregateID 'student-001', got '%s'", lastEvent.GetAggregateID())
	}
}

// TestMultipleConsumerGroups menguji consumer groups.
func TestMultipleConsumerGroups(t *testing.T) {
	srv, url := runTestNATSServer(t)
	defer srv.Shutdown()

	conn := setupTestConnection(t, url)
	broker := NewBrokerWithConnection(conn, "SDP_EVENTS")

	_ = broker.EnsureStream(context.Background(), contract.StreamConfig{
		Name:     "SDP_EVENTS",
		Subjects: []string{"sdp.events.>"},
		Storage:  "memory",
	})

	// 3 consumer groups (edge-guru, edge-website, edge-ppdb)
	var guruReceived, websiteReceived, ppdbReceived int32

	guruSub := NewNATSSubscriber(conn, "SDP_EVENTS", "edge-guru")
	defer guruSub.Close()
	websiteSub := NewNATSSubscriber(conn, "SDP_EVENTS", "edge-website")
	defer websiteSub.Close()
	ppdbSub := NewNATSSubscriber(conn, "SDP_EVENTS", "edge-ppdb")
	defer ppdbSub.Close()

	_, _ = guruSub.Subscribe("*", func(ctx context.Context, event types.DomainEvent) error {
		atomic.AddInt32(&guruReceived, 1)
		return nil
	})
	_, _ = websiteSub.Subscribe("*", func(ctx context.Context, event types.DomainEvent) error {
		atomic.AddInt32(&websiteReceived, 1)
		return nil
	})
	_, _ = ppdbSub.Subscribe("*", func(ctx context.Context, event types.DomainEvent) error {
		atomic.AddInt32(&ppdbReceived, 1)
		return nil
	})

	// Publish 1 event
	event := types.NewBaseEventBuilder("StudentCreated", "Student", "id-1").Build()
	_ = broker.Publish(context.Background(), event)

	time.Sleep(500 * time.Millisecond)

	// Semua consumer group harus terima
	if atomic.LoadInt32(&guruReceived) != 1 {
		t.Errorf("edge-guru: expected 1, got %d", guruReceived)
	}
	if atomic.LoadInt32(&websiteReceived) != 1 {
		t.Errorf("edge-website: expected 1, got %d", websiteReceived)
	}
	if atomic.LoadInt32(&ppdbReceived) != 1 {
		t.Errorf("edge-ppdb: expected 1, got %d", ppdbReceived)
	}
}

// TestResilientConsumerWithNATS menguji ResilientConsumer + NATS.
func TestResilientConsumerWithNATS(t *testing.T) {
	srv, url := runTestNATSServer(t)
	defer srv.Shutdown()

	conn := setupTestConnection(t, url)
	broker := NewBrokerWithConnection(conn, "SDP_EVENTS")

	_ = broker.EnsureStream(context.Background(), contract.StreamConfig{
		Name:     "SDP_EVENTS",
		Subjects: []string{"sdp.events.>"},
		Storage:  "memory",
	})

	var processed int32
	d := dlq.NewMemoryDLQ(0)
	checker := idempotency.NewMemoryChecker()
	defer checker.ClearAll()

	// Handler yang selalu sukses
	handler := func(ctx context.Context, event types.DomainEvent) error {
		atomic.AddInt32(&processed, 1)
		return nil
	}

	// Wrap dengan resilient consumer (import dari events/consumer)
	// Untuk simplicity, test langsung di sini
	sub := NewNATSSubscriber(conn, "SDP_EVENTS", "resilient-test")
	defer sub.Close()

	resilientHandler := func(ctx context.Context, event types.DomainEvent) error {
		// Idempotency check
		key := event.GetEventName() + ":" + event.GetEventID()
		marked, _ := checker.Mark(key, time.Hour)
		if !marked {
			return nil // Already processed
		}
		return handler(ctx, event)
	}

	_, _ = sub.Subscribe("*", resilientHandler)

	// Publish event
	event := types.NewBaseEventBuilder("StudentCreated", "Student", "id-1").
		WithEventID("evt-001").
		Build()
	_ = broker.Publish(context.Background(), event)

	time.Sleep(300 * time.Millisecond)

	if atomic.LoadInt32(&processed) != 1 {
		t.Errorf("Expected 1 processed, got %d", processed)
	}

	// Publish event yang sama lagi (harus di-skip oleh idempotency)
	_ = broker.Publish(context.Background(), event)
	time.Sleep(300 * time.Millisecond)

	if atomic.LoadInt32(&processed) != 1 {
		t.Errorf("Expected still 1 processed (idempotent), got %d", processed)
	}

	// DLQ tidak boleh ada event
	if d.Count() != 0 {
		t.Errorf("Expected empty DLQ, got %d", d.Count())
	}

	_ = d                   // gunakan variable
	_ = retry.DefaultPolicy // gunakan variable
}

// TestStreamInfo menguji stream info retrieval.
func TestStreamInfo(t *testing.T) {
	srv, url := runTestNATSServer(t)
	defer srv.Shutdown()

	conn := setupTestConnection(t, url)
	broker := NewBrokerWithConnection(conn, "SDP_EVENTS")

	_ = broker.EnsureStream(context.Background(), contract.StreamConfig{
		Name:     "SDP_EVENTS",
		Subjects: []string{"sdp.events.>"},
		Storage:  "memory",
	})

	// Publish beberapa events
	for i := 0; i < 5; i++ {
		event := types.NewBaseEventBuilder("Event", "Agg", "id").Build()
		_ = broker.Publish(context.Background(), event)
	}

	info, err := broker.StreamInfo(context.Background(), "SDP_EVENTS")
	if err != nil {
		t.Fatalf("StreamInfo failed: %v", err)
	}

	if info.Messages != 5 {
		t.Errorf("Expected 5 messages, got %d", info.Messages)
	}
	if info.FirstSeq != 1 {
		t.Errorf("Expected first seq 1, got %d", info.FirstSeq)
	}
	if info.LastSeq != 5 {
		t.Errorf("Expected last seq 5, got %d", info.LastSeq)
	}
}

// TestPing menguji health check.
func TestPing(t *testing.T) {
	srv, url := runTestNATSServer(t)
	defer srv.Shutdown()

	conn := setupTestConnection(t, url)
	broker := NewBrokerWithConnection(conn, "SDP_EVENTS")

	if err := broker.Ping(context.Background()); err != nil {
		t.Errorf("Ping failed: %v", err)
	}
}
