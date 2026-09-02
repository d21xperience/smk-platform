// Package contract menyediakan high-level abstraction untuk Message Broker.
// Interface ini adalah facade yang menggabungkan Publisher dan Subscriber
// dengan kemampuan tambahan seperti stream management dan replay.
package contract

import (
	"context"
	"time"

	"sekolah-platform/platform/events/publisher"
	"sekolah-platform/platform/events/subscriber"
)

// StreamConfig adalah konfigurasi untuk stream.
type StreamConfig struct {
	// Name adalah nama stream (e.g. "SDP_EVENTS").
	Name string

	// Subjects adalah list subject patterns (e.g. ["sdp.events.>"]).
	Subjects []string

	// Retention adalah retention policy (limits, interest, workqueue).
	Retention string

	// MaxAge adalah maksimum umur message sebelum dihapus.
	MaxAge time.Duration

	// MaxBytes adalah maksimum ukuran stream.
	MaxBytes int64

	// MaxMessages adalah maksimum jumlah message dalam stream.
	MaxMessages int64

	// Storage adalah tipe storage (file, memory).
	Storage string
}

// ConsumerConfig adalah konfigurasi untuk consumer.
type ConsumerConfig struct {
	// DurableName adalah nama durable consumer (wajib).
	DurableName string

	// FilterSubject memfilter message berdasarkan subject (optional).
	FilterSubject string

	// AckWait adalah durasi menunggu ack sebelum redeliver.
	AckWait time.Duration

	// MaxDeliver adalah maksimum delivery attempts (0 = unlimited).
	MaxDeliver int

	// DeliverPolicy adalah policy delivery (all, last, new, byStartSequence, byTime).
	DeliverPolicy string

	// OptStartSeq adalah sequence awal untuk replay (jika DeliverPolicy=byStartSequence).
	OptStartSeq uint64

	// OptStartTime adalah waktu awal untuk replay (jika DeliverPolicy=byTime).
	OptStartTime time.Time
}

// Broker adalah high-level interface untuk Message Broker.
// Menggabungkan Publisher, Subscriber, dan stream management.
type Broker interface {
	publisher.Publisher
	subscriber.Subscriber

	// Stream Management
	EnsureStream(ctx context.Context, config StreamConfig) error
	DeleteStream(ctx context.Context, name string) error
	StreamInfo(ctx context.Context, name string) (*StreamInfo, error)

	// Consumer Management
	EnsureConsumer(ctx context.Context, streamName string, config ConsumerConfig) error
	DeleteConsumer(ctx context.Context, streamName, consumerName string) error

	// Replay Support (untuk Projection Rebuild - ADR-009)
	ReplayFrom(ctx context.Context, streamName, consumerName string, startSeq uint64) error
	ReplayFromTime(ctx context.Context, streamName, consumerName string, startTime time.Time) error

	// Health
	Ping(ctx context.Context) error
	Close() error
}

// StreamInfo adalah informasi tentang stream.
type StreamInfo struct {
	Name      string    `json:"name"`
	Messages  uint64    `json:"messages"`
	Bytes     uint64    `json:"bytes"`
	FirstSeq  uint64    `json:"firstSeq"`
	LastSeq   uint64    `json:"lastSeq"`
	CreatedAt time.Time `json:"createdAt"`
}

// DefaultStreamConfig adalah konfigurasi stream default.
var DefaultStreamConfig = StreamConfig{
	Name:        "SDP_EVENTS",
	Subjects:    []string{"sdp.events.>"},
	Retention:   "limits",
	MaxAge:      7 * 24 * time.Hour,     // 7 days
	MaxBytes:    1 * 1024 * 1024 * 1024, // 1 GB
	MaxMessages: 1_000_000,
	Storage:     "file",
}
