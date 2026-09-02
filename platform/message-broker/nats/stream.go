package nats

import (
	"context"
	"time"

	"sekolah-platform/platform/message-broker/contract"

	"github.com/nats-io/nats.go"
)

// EnsureStream membuat stream jika belum ada, atau update jika sudah ada.
func (b *BrokerImpl) EnsureStream(ctx context.Context, config contract.StreamConfig) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	js := b.conn.JetStream()

	// Convert retention string ke nats.RetentionPolicy
	retention := nats.LimitsPolicy
	switch config.Retention {
	case "interest":
		retention = nats.InterestPolicy
	case "workqueue":
		retention = nats.WorkQueuePolicy
	}

	// Convert storage string ke nats.StorageType
	storage := nats.FileStorage
	if config.Storage == "memory" {
		storage = nats.MemoryStorage
	}

	streamConfig := &nats.StreamConfig{
		Name:      config.Name,
		Subjects:  config.Subjects,
		Retention: retention,
		Storage:   storage,
		MaxAge:    config.MaxAge,
		MaxBytes:  config.MaxBytes,
		MaxMsgs:   config.MaxMessages,
	}

	// Try to get existing stream
	_, err := js.StreamInfo(config.Name)
	if err == nats.ErrStreamNotFound {
		// Create new stream
		_, err = js.AddStream(streamConfig)
		if err != nil {
			return err
		}
		b.conn.logger.Info("Created stream: %s", config.Name)
		return nil
	}
	if err != nil {
		return err
	}

	// Update existing stream
	_, err = js.UpdateStream(streamConfig)
	if err != nil {
		return err
	}
	b.conn.logger.Info("Updated stream: %s", config.Name)
	return nil
}

// DeleteStream menghapus stream.
func (b *BrokerImpl) DeleteStream(ctx context.Context, name string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	js := b.conn.JetStream()
	return js.DeleteStream(name)
}

// StreamInfo mengambil informasi stream.
func (b *BrokerImpl) StreamInfo(ctx context.Context, name string) (*contract.StreamInfo, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	js := b.conn.JetStream()
	info, err := js.StreamInfo(name)
	if err != nil {
		return nil, err
	}

	return &contract.StreamInfo{
		Name:      info.Config.Name,
		Messages:  info.State.Msgs,
		Bytes:     info.State.Bytes,
		FirstSeq:  info.State.FirstSeq,
		LastSeq:   info.State.LastSeq,
		CreatedAt: info.Created,
	}, nil
}

// EnsureConsumer membuat consumer jika belum ada.
func (b *BrokerImpl) EnsureConsumer(ctx context.Context, streamName string, config contract.ConsumerConfig) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	js := b.conn.JetStream()

	// Convert deliver policy
	deliverPolicy := nats.DeliverAllPolicy
	switch config.DeliverPolicy {
	case "last":
		deliverPolicy = nats.DeliverLastPolicy
	case "new":
		deliverPolicy = nats.DeliverNewPolicy
	case "byStartSequence":
		deliverPolicy = nats.DeliverByStartSequencePolicy
	case "byTime":
		deliverPolicy = nats.DeliverByStartTimePolicy
	}

	consumerConfig := &nats.ConsumerConfig{
		Durable:       config.DurableName,
		FilterSubject: config.FilterSubject,
		AckPolicy:     nats.AckExplicitPolicy,
		AckWait:       config.AckWait,
		MaxDeliver:    config.MaxDeliver,
		DeliverPolicy: deliverPolicy,
	}

	if config.DeliverPolicy == "byStartSequence" {
		consumerConfig.OptStartSeq = config.OptStartSeq
	}
	if config.DeliverPolicy == "byTime" && !config.OptStartTime.IsZero() {
		consumerConfig.OptStartTime = &config.OptStartTime
	}

	// Default AckWait
	if consumerConfig.AckWait == 0 {
		consumerConfig.AckWait = 30 * time.Second
	}

	_, err := js.AddConsumer(streamName, consumerConfig)
	if err != nil && err != nats.ErrConsumerNameAlreadyInUse {
		return err
	}

	b.conn.logger.Info("Ensured consumer: %s on stream %s", config.DurableName, streamName)
	return nil
}

// DeleteConsumer menghapus consumer.
func (b *BrokerImpl) DeleteConsumer(ctx context.Context, streamName, consumerName string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	js := b.conn.JetStream()
	return js.DeleteConsumer(streamName, consumerName)
}
