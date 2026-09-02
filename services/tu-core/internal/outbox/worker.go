package outbox

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
)

// OutboxWorker membaca event dari tabel outbox dan publish ke NATS.
type OutboxWorker struct {
	db        *pgxpool.Pool
	js        nats.JetStreamContext
	interval  time.Duration
	batchSize int
}

// NewOutboxWorker membuat instance baru dari OutboxWorker.
func NewOutboxWorker(db *pgxpool.Pool, js nats.JetStreamContext) *OutboxWorker {
	return &OutboxWorker{
		db:        db,
		js:        js,
		interval:  1 * time.Second, // Polling setiap 1 detik
		batchSize: 100,
	}
}

// Start memulai worker loop.
func (w *OutboxWorker) Start(ctx context.Context) {
	log.Println("[Outbox Worker] Dimulai.")
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("[Outbox Worker] Dihentikan.")
			return
		case <-ticker.C:
			w.processBatch(ctx)
		}
	}
}

func (w *OutboxWorker) processBatch(ctx context.Context) {
	// Ambil batch event yang belum dipublish
	// Query ini HARUS match dengan schema di migrations/003_create_outbox_events.sql
	query := `
		SELECT id, aggregate_type, aggregate_id, event_type, payload,
		       correlation_id, trace_id, created_at
		FROM outbox_events
		WHERE published_at IS NULL
		ORDER BY created_at ASC
		LIMIT $1
		FOR UPDATE SKIP LOCKED
	`

	rows, err := w.db.Query(ctx, query, w.batchSize)
	if err != nil {
		log.Printf("[Outbox Worker] Query error: %v", err)
		return
	}
	defer rows.Close()

	var events []outboxEvent
	for rows.Next() {
		var e outboxEvent
		if err := rows.Scan(
			&e.ID, &e.AggregateType, &e.AggregateID,
			&e.EventType, &e.Payload, &e.CorrelationID,
			&e.TraceID, &e.CreatedAt,
		); err != nil {
			log.Printf("[Outbox Worker] Scan error: %v", err)
			return
		}
		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[Outbox Worker] Rows iteration error: %v", err)
		return
	}

	if len(events) == 0 {
		return
	}

	log.Printf("[Outbox Worker] Memproses %d event(s)...", len(events))

	// Publish setiap event ke NATS
	for _, e := range events {
		subject := fmt.Sprintf("sdp.event.%s.%s", e.AggregateType, e.EventType)

		// Set headers untuk tracing (ADR-009)
		msg := nats.Msg{
			Subject: subject,
			Data:    e.Payload,
			Header:  nats.Header{},
		}
		msg.Header.Set("X-Correlation-ID", e.CorrelationID)
		msg.Header.Set("X-Trace-ID", e.TraceID)
		msg.Header.Set("X-Event-ID", e.ID)
		msg.Header.Set("X-Event-Type", e.EventType)

		if _, err := w.js.PublishMsg(&msg); err != nil {
			log.Printf("[Outbox Worker] Publish error untuk event %s: %v", e.ID, err)

			// Increment attempts dan set next_retry
			_, updateErr := w.db.Exec(ctx,
				`UPDATE outbox_events
				 SET attempts = attempts + 1,
				     last_error = $2,
				     next_retry = NOW() + INTERVAL '1 minute' * (attempts + 1)
				 WHERE id = $1`,
				e.ID, err.Error(),
			)
			if updateErr != nil {
				log.Printf("[Outbox Worker] Update retry error: %v", updateErr)
			}
			continue // Jangan update published_at, akan dicoba lagi di batch berikutnya
		}

		// Tandai sebagai published
		_, err := w.db.Exec(ctx,
			"UPDATE outbox_events SET published_at = NOW() WHERE id = $1",
			e.ID,
		)
		if err != nil {
			log.Printf("[Outbox Worker] Update published_at error: %v", err)
		} else {
			log.Printf("[Outbox Worker] ✅ Event %s published ke %s", e.ID, subject)
		}
	}
}

// outboxEvent merepresentasikan event yang akan dipublish.
type outboxEvent struct {
	ID            string
	AggregateType string
	AggregateID   string
	EventType     string
	Payload       []byte
	CorrelationID string
	TraceID       string
	CreatedAt     time.Time
}
