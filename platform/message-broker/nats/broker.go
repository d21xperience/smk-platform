package nats

import (
	"context"
	"time"

	"sekolah-platform/platform/events/subscriber"
	"sekolah-platform/platform/events/types"

	"github.com/nats-io/nats.go"
)

// BrokerImpl adalah implementasi contract.Broker menggunakan NATS JetStream.
// BrokerImpl mengimplementasikan publisher.Publisher dan subscriber.Subscriber
// sehingga bisa dipakai langsung dengan OutboxProcessor dan ResilientConsumer.
type BrokerImpl struct {
	conn       *Connection
	streamName string
}

// NewBroker membuat Broker NATS baru.
func NewBroker(config Config, logger Logger) (*BrokerImpl, error) {
	conn, err := NewConnection(config, logger)
	if err != nil {
		return nil, err
	}

	return &BrokerImpl{
		conn:       conn,
		streamName: "SDP_EVENTS",
	}, nil
}

// NewBrokerWithConnection membuat Broker dari Connection yang sudah ada.
func NewBrokerWithConnection(conn *Connection, streamName string) *BrokerImpl {
	return &BrokerImpl{
		conn:       conn,
		streamName: streamName,
	}
}

// === publisher.Publisher Implementation ===

// Publish mengimplementasikan publisher.Publisher.
func (b *BrokerImpl) Publish(ctx context.Context, event types.DomainEvent) error {
	pub := NewNATSPublisher(b.conn, b.streamName)
	return pub.Publish(ctx, event)
}

// PublishBatch mengimplementasikan publisher.Publisher.
func (b *BrokerImpl) PublishBatch(ctx context.Context, events []types.DomainEvent) error {
	pub := NewNATSPublisher(b.conn, b.streamName)
	return pub.PublishBatch(ctx, events)
}

// === subscriber.Subscriber Implementation ===

// Subscribe mengimplementasikan subscriber.Subscriber.
// durableName bisa di-set via WithDurableName sebelum Subscribe.
func (b *BrokerImpl) Subscribe(eventName string, handler subscriber.Handler) (subscriber.Subscription, error) {
	sub := NewNATSSubscriber(b.conn, b.streamName, "default-consumer")
	return sub.Subscribe(eventName, handler)
}

// SubscribeWithDurable membuat subscription dengan durable name spesifik.
func (b *BrokerImpl) SubscribeWithDurable(eventName, durableName string, handler subscriber.Handler) (subscriber.Subscription, error) {
	sub := NewNATSSubscriber(b.conn, b.streamName, durableName)
	return sub.Subscribe(eventName, handler)
}

// === Stream Management ===

// EnsureStream mengimplementasikan contract.Broker.
// func (b *BrokerImpl) EnsureStream(ctx context.Context, config contract.StreamConfig) error {
// 	if err := ctx.Err(); err != nil {
// 		return err
// 	}

// 	js := b.conn.JetStream()

// 	retention := nats.LimitsPolicy
// 	switch config.Retention {
// 	case "interest":
// 		retention = nats.InterestPolicy
// 	case "workqueue":
// 		retention = nats.WorkQueuePolicy
// 	}

// 	storage := nats.FileStorage
// 	if config.Storage == "memory" {
// 		storage = nats.MemoryStorage
// 	}

// 	streamConfig := &nats.StreamConfig{
// 		Name:      config.Name,
// 		Subjects:  config.Subjects,
// 		Retention: retention,
// 		Storage:   storage,
// 		MaxAge:    config.MaxAge,
// 		MaxBytes:  config.MaxBytes,
// 		MaxMsgs:   config.MaxMessages,
// 	}

// 	_, err := js.StreamInfo(config.Name)
// 	if err == nats.ErrStreamNotFound {
// 		_, err = js.AddStream(streamConfig)
// 		if err != nil {
// 			return err
// 		}
// 		b.conn.logger.Info("Created stream: %s", config.Name)
// 		return nil
// 	}
// 	if err != nil {
// 		return err
// 	}

// 	_, err = js.UpdateStream(streamConfig)
// 	if err != nil {
// 		return err
// 	}
// 	b.conn.logger.Info("Updated stream: %s", config.Name)
// 	return nil
// }

// DeleteStream mengimplementasikan contract.Broker.
// func (b *BrokerImpl) DeleteStream(ctx context.Context, name string) error {
// 	if err := ctx.Err(); err != nil {
// 		return err
// 	}
// 	return b.conn.JetStream().DeleteStream(name)
// }

// StreamInfo mengimplementasikan contract.Broker.
// func (b *BrokerImpl) StreamInfo(ctx context.Context, name string) (*contract.StreamInfo, error) {
// 	if err := ctx.Err(); err != nil {
// 		return nil, err
// 	}

// 	info, err := b.conn.JetStream().StreamInfo(name)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &contract.StreamInfo{
// 		Name:      info.Config.Name,
// 		Messages:  info.State.Msgs,
// 		Bytes:     info.State.Bytes,
// 		FirstSeq:  info.State.FirstSeq,
// 		LastSeq:   info.State.LastSeq,
// 		CreatedAt: info.Created,
// 	}, nil
// }

// === Consumer Management ===

// EnsureConsumer mengimplementasikan contract.Broker.
// func (b *BrokerImpl) EnsureConsumer(ctx context.Context, streamName string, config contract.ConsumerConfig) error {
// 	if err := ctx.Err(); err != nil {
// 		return err
// 	}

// 	js := b.conn.JetStream()

// 	deliverPolicy := nats.DeliverAllPolicy
// 	switch config.DeliverPolicy {
// 	case "last":
// 		deliverPolicy = nats.DeliverLastPolicy
// 	case "new":
// 		deliverPolicy = nats.DeliverNewPolicy
// 	case "byStartSequence":
// 		deliverPolicy = nats.DeliverByStartSequencePolicy
// 	case "byTime":
// 		deliverPolicy = nats.DeliverByStartTimePolicy
// 	}

// 	consumerConfig := &nats.ConsumerConfig{
// 		Durable:       config.DurableName,
// 		FilterSubject: config.FilterSubject,
// 		AckPolicy:     nats.AckExplicitPolicy,
// 		AckWait:       config.AckWait,
// 		MaxDeliver:    config.MaxDeliver,
// 		DeliverPolicy: deliverPolicy,
// 	}

// 	if config.AckWait == 0 {
// 		consumerConfig.AckWait = 30 * time.Second
// 	}
// 	if config.DeliverPolicy == "byStartSequence" {
// 		consumerConfig.OptStartSeq = config.OptStartSeq
// 	}
// 	if config.DeliverPolicy == "byTime" && !config.OptStartTime.IsZero() {
// 		consumerConfig.OptStartTime = &config.OptStartTime
// 	}

// 	_, err := js.AddConsumer(streamName, consumerConfig)
// 	if err != nil && err != nats.ErrConsumerNameAlreadyInUse {
// 		return err
// 	}

// 	b.conn.logger.Info("Ensured consumer: %s on stream %s", config.DurableName, streamName)
// 	return nil
// }

// // DeleteConsumer mengimplementasikan contract.Broker.
// func (b *BrokerImpl) DeleteConsumer(ctx context.Context, streamName, consumerName string) error {
// 	if err := ctx.Err(); err != nil {
// 		return err
// 	}
// 	return b.conn.JetStream().DeleteConsumer(streamName, consumerName)
// }

// === Replay Support ===

// ReplayFrom mengimplementasikan contract.Broker.
func (b *BrokerImpl) ReplayFrom(ctx context.Context, streamName, consumerName string, startSeq uint64) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	_, err := b.conn.JetStream().UpdateConsumer(streamName, &nats.ConsumerConfig{
		Durable:       consumerName,
		DeliverPolicy: nats.DeliverByStartSequencePolicy,
		OptStartSeq:   startSeq,
		AckPolicy:     nats.AckExplicitPolicy,
	})
	return err
}

// ReplayFromTime mengimplementasikan contract.Broker.
func (b *BrokerImpl) ReplayFromTime(ctx context.Context, streamName, consumerName string, startTime time.Time) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	_, err := b.conn.JetStream().UpdateConsumer(streamName, &nats.ConsumerConfig{
		Durable:       consumerName,
		DeliverPolicy: nats.DeliverByStartTimePolicy,
		OptStartTime:  &startTime,
		AckPolicy:     nats.AckExplicitPolicy,
	})
	return err
}

// === Health ===

// Ping mengecek koneksi ke NATS server.
func (b *BrokerImpl) Ping(ctx context.Context) error {
	if !b.conn.IsConnected() {
		return errNotConnected
	}
	return b.conn.Conn().Flush()
}

// Close mengimplementasikan contract.Broker.
func (b *BrokerImpl) Close() error {
	return b.conn.Close()
}

// Connection mengembalikan underlying Connection (untuk advanced usage).
func (b *BrokerImpl) Connection() *Connection {
	return b.conn
}
