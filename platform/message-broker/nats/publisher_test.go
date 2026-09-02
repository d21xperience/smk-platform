package nats

import (
	"context"
	"testing"

	"sekolah-platform/platform/events/types"
)

func TestNATSPublisherPublish(t *testing.T) {
	srv, url := runTestNATSServer(t)
	defer srv.Shutdown()

	conn := setupTestConnection(t, url)
	defer conn.Close()

	// Setup stream
	js := conn.JetStream()
	// _, err := js.AddStream(&natsStreamConfig("SDP_EVENTS", []string{"sdp.events.>"}))
	cfg := natsStreamConfig("SDP_EVENTS", []string{"sdp.events.>"})
	_, err := js.AddStream(&cfg) // Berikan pointer dari variabel cfg
	if err != nil {
		t.Fatalf("Failed to create stream: %v", err)
	}

	pub := NewNATSPublisher(conn, "SDP_EVENTS")

	event := types.NewBaseEventBuilder("StudentCreated", "Student", "student-001").
		WithEventID("evt-001").
		WithPayload(map[string]string{"name": "Budi"}).
		Build()

	err = pub.Publish(context.Background(), event)
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	// Verify message in stream
	info, err := js.StreamInfo("SDP_EVENTS")
	if err != nil {
		t.Fatalf("Failed to get stream info: %v", err)
	}
	if info.State.Msgs != 1 {
		t.Errorf("Expected 1 message in stream, got %d", info.State.Msgs)
	}
}

func TestNATSPublisherPublishBatch(t *testing.T) {
	srv, url := runTestNATSServer(t)
	defer srv.Shutdown()

	conn := setupTestConnection(t, url)
	defer conn.Close()

	js := conn.JetStream()
	// Tampung ke variabel terlebih dahulu
	// _, _ = js.AddStream(&natsStreamConfig("SDP_EVENTS", []string{"sdp.events.>"}))
	cfg := natsStreamConfig("SDP_EVENTS", []string{"sdp.events.>"})
	_, _ = js.AddStream(&cfg) // Berikan pointer dari variabel cfg

	pub := NewNATSPublisher(conn, "SDP_EVENTS")

	events := []types.DomainEvent{
		types.NewBaseEventBuilder("Event1", "Agg", "id-1").WithEventID("evt-1").Build(),
		types.NewBaseEventBuilder("Event2", "Agg", "id-2").WithEventID("evt-2").Build(),
		types.NewBaseEventBuilder("Event3", "Agg", "id-3").WithEventID("evt-3").Build(),
	}

	err := pub.PublishBatch(context.Background(), events)
	if err != nil {
		t.Fatalf("PublishBatch failed: %v", err)
	}

	info, _ := js.StreamInfo("SDP_EVENTS")
	if info.State.Msgs != 3 {
		t.Errorf("Expected 3 messages, got %d", info.State.Msgs)
	}
}

func TestNATSPublisherNotConnected(t *testing.T) {
	conn := &Connection{} // not connected
	pub := NewNATSPublisher(conn, "SDP_EVENTS")

	event := types.NewBaseEventBuilder("Event", "Agg", "id").Build()
	err := pub.Publish(context.Background(), event)
	if err == nil {
		t.Error("Expected error when not connected")
	}
}
