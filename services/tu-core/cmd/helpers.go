package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/nats-io/nats.go"
	"sekolah-platform/services/tu-core/internal/service"
)

type natsEventPublisher struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

func (p *natsEventPublisher) Publish(ctx context.Context, event service.DomainEvent) error {
	msg := &nats.Msg{
		Subject: fmt.Sprintf("sdp.event.%s.%s", event.AggregateType, event.EventType),
		Data:    []byte(fmt.Sprintf("%+v", event.Payload)),
		Header:  make(nats.Header),
	}
	msg.Header.Set("X-Correlation-ID", event.CorrelationID)
	msg.Header.Set("X-Event-Type", event.EventType)
	_, err := p.js.PublishMsg(msg)
	return err
}

func (p *natsEventPublisher) PublishBatch(ctx context.Context, events []service.DomainEvent) error {
	for _, e := range events {
		if err := p.Publish(ctx, e); err != nil {
			return err
		}
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v, exists := os.LookupEnv(key); exists {
		return v
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if v, exists := os.LookupEnv(key); exists {
		var i int
		if _, err := fmt.Sscanf(v, "%d", &i); err == nil {
			return i
		}
	}
	return fallback
}

type dummyLogger struct{}

func (l *dummyLogger) Info(msg string, kv ...interface{}) {
	log.Printf("[INFO] %s %v", msg, kv)
}

func (l *dummyLogger) Error(msg string, kv ...interface{}) {
	log.Printf("[ERROR] %s %v", msg, kv)
}

func (l *dummyLogger) Debug(msg string, kv ...interface{}) {
	log.Printf("[DEBUG] %s %v", msg, kv)
}
