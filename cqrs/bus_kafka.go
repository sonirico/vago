package cqrs

import (
	"context"
	"errors"
	"fmt"
	"time"

	maps "github.com/sonirico/stadio/ds/map"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/sonirico/vago/lol"
	"github.com/sonirico/vago/rp"
)

type (
	opts struct {
		startUpPing bool

		consumerConf    *rp.ConsumerConfig
		consumerAPMConf *rp.APMConfig
		producerConf    *rp.ProducerConfig
		flushTimeout    time.Duration
	}
)

// RedpandaBus
// Can produce to any topic
// Can only consume from 1 topic -> buses and topics follow a 1-1 relationship
type (
	RedpandaBus struct {
		idx string

		mtopic string
		mcodec Codec

		log lol.Logger

		p rp.Producer
		c rp.Consumer

		opts opts
	}

	EventRedpandaBus struct {
		*RedpandaBus

		eventHandlers maps.Map[string, EventHandler]
		sagaHandlers  maps.Map[string, SagaHandler]
	}

	CommandRedpandaBus struct {
		*RedpandaBus

		handlers maps.Map[string, CommandHandler]
	}
)

func (b *CommandRedpandaBus) CommandHandler(h CommandHandler) CommandBus {
	b.handlers.Set(hashKey(h), h)
	return b

}

func (b *CommandRedpandaBus) hasHandlers() bool {
	return len(b.handlers.Keys()) > 0
}

func (b *CommandRedpandaBus) commandHandler(h Handler) (CommandHandler, bool) {
	return b.handlers.Get(hashKey(h))
}

func (b *EventRedpandaBus) EventHandler(h EventHandler) EventBus {
	b.eventHandlers.Set(hashKey(h), h)
	return b
}

func (b *EventRedpandaBus) SagaHandler(h SagaHandler) EventBus {
	b.sagaHandlers.Set(hashKey(h), h)
	return b
}

func (b *EventRedpandaBus) eventHandler(h Handler) (EventHandler, bool) {
	return b.eventHandlers.Get(hashKey(h))
}
func (b *EventRedpandaBus) sagaHandler(h Handler) (SagaHandler, bool) {
	return b.sagaHandlers.Get(hashKey(h))
}

func (b *EventRedpandaBus) hasHandlers() bool {
	return len(b.sagaHandlers.Keys()) > 0 || len(b.eventHandlers.Keys()) > 0
}

type (
	EventBus interface {
		bus

		SagaHandler(sagaHandler SagaHandler) EventBus
		EventHandler(eventHandler EventHandler) EventBus

		sagaHandler(x Handler) (SagaHandler, bool)
		eventHandler(x Handler) (EventHandler, bool)
	}

	CommandBus interface {
		bus

		CommandHandler(cmdHandler CommandHandler) CommandBus

		commandHandler(x Handler) (CommandHandler, bool)
	}
)

func newBus(
	id string,
	topic string,
	log lol.Logger,
	opts ...Configurator[RedpandaBus],
) (*RedpandaBus, error) {
	bus := &RedpandaBus{
		idx:    id,
		mtopic: topic,
		log:    log.WithFields(lol.Fields{"bus_id": id, "topic": topic}),
	}

	for _, opt := range opts {
		opt.Apply(bus)
	}

	var err error

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if bus.opts.producerConf != nil {
		bus.p, err = rp.NewProducer(
			ctx,
			*bus.opts.producerConf,
			log,
		)

		if err != nil {
			return nil, err
		}
	}

	if bus.opts.consumerConf != nil {

		bus.c, err = rp.NewConsumer(
			log,
			*bus.opts.consumerConf,
			[]string{bus.topic()},
			bus.opts.consumerAPMConf,
		)

		if err != nil {
			return nil, err
		}
	}

	if bus.opts.startUpPing {
		if pingErr := bus.p.Ping(ctx); pingErr != nil {
			return nil, fmt.Errorf("%w: %v", rp.ErrPingFailed, pingErr)
		}
	}

	return bus, err
}

func NewCommandBus(
	id string,
	topic string,
	log lol.Logger,
	opts ...Configurator[RedpandaBus],
) (*CommandRedpandaBus, error) {
	bus, err := newBus(id, topic, log, opts...)
	if err != nil {
		return nil, err
	}

	return &CommandRedpandaBus{
		handlers: maps.NewConcurrent[string, CommandHandler](
			maps.NewNative[string, CommandHandler](),
		),
		RedpandaBus: bus,
	}, nil
}

func NewEventBus(
	id string,
	topic string,
	log lol.Logger,
	opts ...Configurator[RedpandaBus],
) (*EventRedpandaBus, error) {
	bus, err := newBus(id, topic, log, opts...)
	if err != nil {
		return nil, err
	}

	return &EventRedpandaBus{
		eventHandlers: maps.NewConcurrent[string, EventHandler](
			maps.NewNative[string, EventHandler](),
		),
		sagaHandlers: maps.NewConcurrent[string, SagaHandler](
			maps.NewNative[string, SagaHandler](),
		),
		RedpandaBus: bus,
	}, nil
}

func (b RedpandaBus) id() string { return b.idx }

func (b RedpandaBus) topic() string { return b.mtopic }

func (b RedpandaBus) codec() Codec { return b.mcodec }

func (b RedpandaBus) publish(ctx context.Context, m rp.Msg) error {
	if b.opts.producerConf != nil && b.opts.producerConf.ProduceSync {
		return b.parsePublishError(b.p.Publish(ctx, m))
	}

	return b.parsePublishError(b.p.PublishAsync(ctx, m, rp.NoCallback))
}

func (b RedpandaBus) subscribe(ctx context.Context, handler rp.ConsumerHandler) error {
	return b.parseSubscribeError(b.c.Subscribe(ctx, handler))
}

func (b RedpandaBus) close() {
	if b.p != nil {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			b.opts.producerConf.GetFlushTimeout(),
		)
		defer cancel()

		err := b.p.Flush(ctx)
		if err != nil {
			b.log.Errorf("error when flushing records: %v", err)
		}

		b.p.Close()
		b.log.Info("closed producer topic")
	}

	if b.c != nil {
		b.c.Close()
		b.log.Infof("closed consumer topic")
	}
}

func (b RedpandaBus) parseSubscribeError(err error) error {
	if err == nil {
		return nil
	}

	var (
		errFirstReadEOF = &kgo.ErrFirstReadEOF{}
		errGroupSession = &kgo.ErrGroupSession{}
	)

	switch {
	case errors.Is(kgo.ErrClientClosed, err):
	case errors.Is(kgo.ErrAborting, err):
	case errors.As(err, &errFirstReadEOF):
		return fmt.Errorf("%w: %v", ErrSubscribeNonRecoverable, err)
	case errors.As(err, &errGroupSession):
		return fmt.Errorf("%w: %v", ErrSubscribeNonRecoverable, err)
	default:
		return fmt.Errorf("unknown error: %w", err)
	}

	return nil
}

func (b RedpandaBus) parsePublishError(err error) error {
	if err == nil {
		return nil
	}

	var (
		errFirstReadEOF = kgo.ErrFirstReadEOF{}
	)

	switch {
	case errors.Is(kgo.ErrMaxBuffered, err):
		return fmt.Errorf("%w: %v", ErrSubscribeNonRecoverable, err)
	case errors.Is(kgo.ErrClientClosed, err):
	case errors.Is(kgo.ErrAborting, err):
	case errors.As(err, &errFirstReadEOF):
		return fmt.Errorf("%w: %v", ErrSubscribeNonRecoverable, err)
	default:
		return fmt.Errorf("unknown error: %w", err)
	}

	return nil
}

func BusWithConsumerConfig(c *rp.ConsumerConfig) Configurator[RedpandaBus] {
	return configureFn[RedpandaBus](func(bus *RedpandaBus) {
		bus.opts.consumerConf = c
	})
}

func BusWithProducerConfig(c *rp.ProducerConfig) Configurator[RedpandaBus] {
	return configureFn[RedpandaBus](func(bus *RedpandaBus) {
		bus.opts.producerConf = c
	})
}

func BusWithTopic(topic string) Configurator[RedpandaBus] {
	return configureFn[RedpandaBus](func(bus *RedpandaBus) {
		bus.mtopic = topic
	})
}

func BusWithCodec(codec Codec) Configurator[RedpandaBus] {
	return configureFn[RedpandaBus](func(bus *RedpandaBus) {
		bus.mcodec = codec
	})
}

func BusWithJsonCodec() Configurator[RedpandaBus] {
	return configureFn[RedpandaBus](func(bus *RedpandaBus) {
		bus.mcodec = NewJson()
	})
}

func BusWithDisabledStartupPing() Configurator[RedpandaBus] {
	return configureFn[RedpandaBus](func(bus *RedpandaBus) {
		bus.opts.startUpPing = false
	})
}

func BusWithProducerFlushTimeout(d time.Duration) Configurator[RedpandaBus] {
	return configureFn[RedpandaBus](func(bus *RedpandaBus) {
		bus.opts.producerConf.FlushTimeout = d
	})
}
