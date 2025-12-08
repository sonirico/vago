package rp

import (
	"context"
	"fmt"

	"github.com/sonirico/vago/lol"
)

type (
	Config struct {
		Brokers []string

		producerPublishSync    bool
		producerOnPublishAsync func(Msg, error)

		ConsumerGroup                string
		ConsumerBlockRebalanceOnPoll bool

		Log lol.Logger

		WithInternalLogger bool
		WithLoggingHooks   bool
		InternalLogLevel   LogLevel
		WithMaxPollRecords int
	}

	Transport struct {
		cfg      Config
		producer Producer
		consumer Consumer
		apm      *APMConfig
	}

	Redpanda interface {
		Publish(ctx context.Context, msg Msg) error
		Flush(ctx context.Context) error
		Subscribe(
			ctx context.Context,
			topic string,
			handler ConsumerHandler,
		) error
		Close() error
	}
)

func New(producer Producer, cfg Config, opts ...Option) (Redpanda, error) {
	k := &Transport{
		cfg:      cfg,
		producer: producer,
	}

	if cfg.Log == nil {
		return nil, ErrLoggerRequired
	}

	for _, opt := range opts {
		opt.Apply(k)
	}

	if !k.cfg.producerPublishSync &&
		k.cfg.producerOnPublishAsync == nil {
		k.cfg.producerOnPublishAsync = NoCallback
	}

	return k, nil
}

func FromOpts(log lol.Logger, opts ...Option) (Redpanda, error) {
	k := &Transport{}

	k.cfg.Log = log

	for _, opt := range opts {
		opt.Apply(k)

	}

	return k, nil
}

func (k *Transport) Ping(ctx context.Context) error {
	if k.producer != nil {
		if err := k.producer.Ping(ctx); err != nil {
			return fmt.Errorf("producer Ping failed: %w", err)
		}
	}

	if k.consumer != nil {
		if err := k.consumer.Ping(ctx); err != nil {
			return fmt.Errorf("consumer Ping failed: %w", err)
		}
	}

	return nil
}

func (k *Transport) Publish(ctx context.Context, m Msg) error {
	if k.cfg.producerPublishSync {
		return k.producer.Publish(ctx, m)
	}

	return k.producer.PublishAsync(ctx, m, k.cfg.producerOnPublishAsync)
}

func (k *Transport) Flush(ctx context.Context) error {
	return k.producer.Flush(ctx)
}

func (k *Transport) Subscribe(
	ctx context.Context,
	topic string,
	handler ConsumerHandler,
) (err error) {
	if k.consumer != nil {
		// Nothing to do
		err = ErrConsumerAlreadyCreated
		return
	}

	cfg := ConsumerConfig{
		Brokers:                      k.cfg.Brokers,
		ConsumerGroup:                k.cfg.ConsumerGroup,
		ConsumerBlockRebalanceOnPoll: k.cfg.ConsumerBlockRebalanceOnPoll,
		WithLogger:                   k.cfg.WithInternalLogger,
		WithLogLevel:                 k.cfg.InternalLogLevel,
		MaxPollRecords:               k.cfg.WithMaxPollRecords,
		WithLoggingHooks:             k.cfg.WithLoggingHooks,
	}

	k.consumer, err = NewConsumer(
		k.cfg.Log.WithField("type", "consumer"),
		cfg,
		[]string{topic},
		k.apm,
	)

	if err != nil {
		return
	}

	return k.consumer.Subscribe(ctx, handler)
}

func (k *Transport) Close() error {
	if k.consumer != nil {
		k.consumer.Close()
	}

	if k.producer != nil {
		k.producer.Close()
	}

	return nil
}
