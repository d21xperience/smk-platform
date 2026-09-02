// Package consumer menyediakan ResilientConsumer yang menggabungkan:
// - Retry dengan exponential backoff
// - Idempotency checking
// - Dead Letter Queue
//
// Consumer ini adalah decorator untuk Handler yang menambahkan resilience.
package consumer

import (
	"context"
	"fmt"
	"time"

	// "sekolah-platform/platform/shared/errors"
	"sekolah-platform/platform/events/dlq"
	"sekolah-platform/platform/events/idempotency"
	"sekolah-platform/platform/events/retry"
	"sekolah-platform/platform/events/subscriber"
	"sekolah-platform/platform/events/types"
)

// Config adalah konfigurasi untuk ResilientConsumer.
type Config struct {
	// Name adalah nama consumer (untuk logging dan DLQ).
	Name string

	// RetryPolicy adalah retry policy.
	RetryPolicy retry.Policy

	// IdempotencyTTL adalah TTL untuk idempotency key.
	IdempotencyTTL time.Duration

	// DLQ adalah dead letter queue.
	DLQ dlq.DLQ

	// IdempotencyChecker adalah idempotency checker.
	IdempotencyChecker idempotency.Checker

	// Logger untuk logging (optional).
	Logger Logger

	// OnRetry callback dipanggil setiap retry (optional).
	OnRetry func(attempt int, err error, delay time.Duration)

	// OnDLQ callback dipanggil saat event masuk DLQ (optional).
	OnDLQ func(deadLetter dlq.DeadLetter)
}

// Logger adalah interface untuk logging.
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type defaultLogger struct{}

func (l *defaultLogger) Info(msg string, args ...interface{})  {}
func (l *defaultLogger) Error(msg string, args ...interface{}) {}

// DefaultConfig adalah konfigurasi default.
var DefaultConfig = Config{
	Name:           "default-consumer",
	RetryPolicy:    retry.DefaultPolicy,
	IdempotencyTTL: 24 * time.Hour,
}

// ResilientConsumer adalah consumer dengan retry, idempotency, dan DLQ.
type ResilientConsumer struct {
	config   Config
	handler  subscriber.Handler
	executor *retry.Executor
}

// NewResilientConsumer membuat ResilientConsumer baru.
func NewResilientConsumer(handler subscriber.Handler, config Config) *ResilientConsumer {
	if config.Logger == nil {
		config.Logger = &defaultLogger{}
	}
	if config.IdempotencyTTL == 0 {
		config.IdempotencyTTL = DefaultConfig.IdempotencyTTL
	}
	if config.Name == "" {
		config.Name = DefaultConfig.Name
	}

	return &ResilientConsumer{
		config:   config,
		handler:  handler,
		executor: retry.NewExecutor(config.RetryPolicy),
	}
}

// Handle menangani event dengan retry, idempotency, dan DLQ.
func (c *ResilientConsumer) Handle(ctx context.Context, event types.DomainEvent) error {
	// 1. Generate idempotency key
	key := idempotency.GenerateKey(event.GetEventName(), event.GetEventID())

	// 2. Check idempotency
	if c.config.IdempotencyChecker != nil {
		processed, err := c.config.IdempotencyChecker.IsProcessed(key)
		if err != nil {
			c.config.Logger.Error("[%s] Failed to check idempotency: %v", c.config.Name, err)
			// Lanjutkan proses, jangan block karena idempotency check gagal
		} else if processed {
			// Sudah diproses, skip
			return nil
		}
	}

	// 3. Execute dengan retry
	result := c.executor.ExecuteWithCallback(ctx, func() error {
		return c.handler(ctx, event)
	}, c.config.OnRetry)

	if result.Success {
		// 4. Mark as processed
		if c.config.IdempotencyChecker != nil {
			if _, err := c.config.IdempotencyChecker.Mark(key, c.config.IdempotencyTTL); err != nil {
				c.config.Logger.Error("[%s] Failed to mark idempotency: %v", c.config.Name, err)
			}
		}
		return nil
	}

	// 5. Retry habis, kirim ke DLQ
	if c.config.DLQ != nil && result.LastError != nil {
		deadLetter := dlq.DeadLetter{
			OriginalEvent: event,
			ConsumerName:  c.config.Name,
			FailedAt:      time.Now().UTC(),
			FailureCount:  result.Attempts,
			LastError:     result.LastError.Error(),
		}

		if err := c.config.DLQ.Store(deadLetter); err != nil {
			c.config.Logger.Error("[%s] Failed to store to DLQ: %v", c.config.Name, err)
		} else {
			c.config.Logger.Error("[%s] Event sent to DLQ after %d attempts: %v",
				c.config.Name, result.Attempts, result.LastError)
			if c.config.OnDLQ != nil {
				c.config.OnDLQ(deadLetter)
			}
		}
	}

	return fmt.Errorf("[%s] consumer failed after %d attempts: %w",
		c.config.Name, result.Attempts, result.LastError)
}

// AsSubscriberHandler mengembalikan function yang cocok untuk subscriber.Handler.
func (c *ResilientConsumer) AsSubscriberHandler() subscriber.Handler {
	return c.Handle
}
