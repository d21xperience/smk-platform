# Platform Message Broker

NATS JetStream implementation untuk School Digital Platform.

## Overview

Module ini menyediakan Message Broker production-ready menggunakan NATS JetStream
untuk distribusi Domain Events antar service (tu-core → edge-bff).

## Features

- ✅ NATS JetStream persistent messaging
- ✅ Consumer groups (multiple edge services)
- ✅ Durable consumer (resume after restart)
- ✅ Replay support (untuk Projection Rebuild)
- ✅ Auto-reconnect
- ✅ Graceful shutdown
- ✅ Compatible dengan `Publisher` & `Subscriber` interfaces dari `platform/events`

## Usage

### Basic Publish

```go
import "sekolah-platform/platform/message-broker/nats"

config := nats.DefaultConfig.WithURL("nats://localhost:4222")
broker, err := nats.NewBroker(config, nil)
if err != nil {
    log.Fatal(err)
}
defer broker.Close()

// Ensure stream
broker.EnsureStream(ctx, contract.StreamConfig{
    Name:     "SDP_EVENTS",
    Subjects: []string{"sdp.events.>"},
    Storage:  "file",
})

// Publish event
event := types.NewBaseEventBuilder("StudentCreated", "Student", "student-001").
    WithEventID("evt-001").
    Build()

broker.Publish(ctx, event)
```
