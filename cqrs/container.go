package cqrs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"runtime/debug"
	"sync"
	"time"

	"go.elastic.co/apm/v2"

	"github.com/sonirico/vago/cond"
	"github.com/sonirico/vago/ent"
	"github.com/sonirico/vago/fp"
	optslib "github.com/sonirico/vago/opts"
	"github.com/sonirico/vago/zero"

	maps "github.com/sonirico/stadio/ds/map"
	"github.com/sonirico/vago/lol"
	"github.com/sonirico/vago/rp"
)

type Container struct {
	log lol.Logger

	producer rp.Producer

	commandBuses maps.Map[string, CommandBus]
	eventBuses   maps.Map[string, EventBus]

	restartOnError       bool
	warnUnprocessed      bool
	mustProcessOrFail    bool
	errorCaptureDisabled bool
	apmDisabled          bool

	closeC    chan error
	closeOnce sync.Once
}

func NewContainer(log lol.Logger, opts ...optslib.Configurator[Container]) *Container {
	container := &Container{
		log:          log,
		commandBuses: maps.NewConcurrent[string, CommandBus](maps.NewNative[string, CommandBus]()),
		eventBuses:   maps.NewConcurrent[string, EventBus](maps.NewNative[string, EventBus]()),
		closeC:       make(chan error),
	}

	optslib.ApplyAll(container, opts...)

	if !container.errorCaptureDisabled {
		var err error
		container.producer, err = rp.NewProducer(
			context.Background(),
			rp.ProducerConfig{
				FlushTimeout: ent.Duration(
					"VAGO_CQRS_ERROR_PRODUCER_Redpanda_FLUSH_TIMEOUT",
					time.Second*5,
				),
				Timeout: ent.Duration(
					"VAGO_CQRS_ERROR_PRODUCER_Redpanda_TIMEOUT",
					time.Second*5,
				),
				ProduceRequestTimeout: ent.Duration(
					"VAGO_CQRS_ERROR_PRODUCER_Redpanda_PRODUCE_REQUEST_TIMEOUT",
					time.Second*5,
				),
				ConnIdleTimeout: ent.Duration(
					"VAGO_CQRS_ERROR_PRODUCER_Redpanda_CONN_IDLE_TIMEOUT",
					time.Second*5,
				),
				RequestTimeoutOverhead: ent.Duration(
					"VAGO_CQRS_ERROR_PRODUCER_Redpanda_REQUEST_TIMEOUT_OVERHEAD",
					time.Second*5,
				),
				RecordDeliveryTimeout: ent.Duration(
					"VAGO_CQRS_ERROR_PRODUCER_Redpanda_RECORD_DELIVERY_TIMEOUT",
					time.Second*5,
				),
				SessionTimeout: ent.Duration(
					"VAGO_CQRS_ERROR_PRODUCER_Redpanda_SESSION_TIMEOUT",
					time.Second*5,
				),
				MaxBufferedRecords: ent.Int(
					"VAGO_CQRS_ERROR_PRODUCER_Redpanda_MAX_BUFFERED_RECORDS",
					10_000,
				),
				Brokers: ent.SliceStr(
					"VAGO_CQRS_ERROR_PRODUCER_Redpanda_BROKERS",
					[]string{"redpanda-cluster:9092"},
				),
				WithLogger:  false,
				WithNoAPM:   true,
				Compression: 0,
				App:         "vago/cqrs",
				Version:     "0.1.0",
			},
			log.WithField("producer", "errors"),
		)

		if err != nil {
			log.Errorf("error creating Redpanda producer for errors: %v", err)
		}
	}

	return container
}

// Done should be used if ContainerMustProcessOrFail config option is set, otherwise don't use it
func (c *Container) Done() error {
	err := <-c.closeC
	return err
}

func (c *Container) Command(ctx context.Context, busID string, cmd CommandPayload) error {
	bus, ok := c.commandBuses.Get(busID)
	if !ok {
		return fmt.Errorf("%w: command bus %s not found", ErrBusNotFound, busID)
	}

	return c.send(ctx, KindCommands, cmd.sendMsg, bus)
}

func (c *Container) Event(ctx context.Context, busID string, event EventPayload) error {
	bus, ok := c.eventBuses.Get(busID)
	if !ok {
		return fmt.Errorf("%w: event bus %s not found", ErrBusNotFound, busID)
	}

	return c.send(ctx, KindEvents, event.sendMsg, bus)
}

func (c *Container) EventBus(b EventBus) *Container {
	c.eventBuses.Set(b.id(), b)
	return c
}

func (c *Container) CommandBus(b CommandBus) *Container {
	c.commandBuses.Set(b.id(), b)
	return c
}

func (c *Container) Start(ctx context.Context) error {
	c.commandBuses.Range(func(busID string, bus CommandBus, _ int) bool {
		l := c.log.WithFields(lol.Fields{"bus": busID, "op": "command"})

		if bus.hasHandlers() {
			go c.wrapConsume(func() error {
				return c.commandBusSubscribe(ctx, l, bus)
			})
		}

		l.Info("[ok] bus set up")

		return true
	})

	c.eventBuses.Range(func(busID string, bus EventBus, _ int) bool {
		l := c.log.WithFields(lol.Fields{"bus": busID, "op": "event"})
		l.Info("setting up bus")

		if bus.hasHandlers() {
			go c.wrapConsume(func() error {
				return c.eventBusSubscribe(ctx, l, bus)
			})
		}

		l.Info("[ok] bus set up")

		return true
	})

	return nil
}

func (c *Container) Close() {
	c.eventBuses.Range(func(s string, bus EventBus, i int) bool {
		bus.close()
		return true
	})

	c.commandBuses.Range(func(s string, bus CommandBus, i int) bool {
		bus.close()
		return true
	})

	// Done method is the only one listening to close and it's not called
	// unless mustProcessOrFail is true. Otherwise, this would block forever
	if c.mustProcessOrFail {
		c.close(nil)
	}
}

func (c *Container) close(err error) {
	c.closeOnce.Do(func() {
		c.closeC <- err
		close(c.closeC)
	})
}

func (c *Container) commandBusSubscribe(ctx context.Context, l lol.Logger, bus CommandBus) error {
	return bus.subscribe(ctx, func(ctx context.Context, m rp.Msg) error {
		recv := recvMsg{
			recordKey:       m.Key,
			recordPartition: m.Partition,
			recordTs:        m.Ts,
		}
		if err := bus.codec().Decode(m.Value, &recv); err != nil {
			//stop consuming
			return fmt.Errorf("%w: unable to decode msg: %v", ErrSubscribeNonRecoverable, err)
		}

		k := hashKey(recv)
		l.Debugf("[c][<] %v", k)

		// Find command Handler
		cmdHandler, ok := bus.commandHandler(recv)
		if !ok {
			if c.warnUnprocessed {
				l.Warnf("no command handler found for " + k)
			}
			return nil
		}

		if !c.apmDisabled {
			if tx := apm.TransactionFromContext(ctx); tx != nil {
				spanName := fmt.Sprintf("%s %s.%s.%s",
					"COMMAND HANDLER",
					recv.Version(),
					recv.Resource(),
					recv.Action(),
				)

				tx.Context.SetLabel("cqrs.kind", KindCommands)
				tx.Context.SetLabel("cqrs.id", recv.ID())
				if userID, ok := recv.User().Unwrap(); ok {
					tx.Context.SetLabel("cqrs.user", userID)
				}
				// Add Redpanda message details as labels
				tx.Context.SetCustom("Redpanda_topic", m.Topic)
				tx.Context.SetCustom("Redpanda_partition", m.Partition)
				tx.Context.SetCustom("Redpanda_key", string(m.Key))
				tx.Context.SetCustom("Redpanda_timestamp", m.Ts.String())

				if span, spanCtx := apm.StartSpan(ctx, spanName, cqrs); span != nil {
					ctx = spanCtx
					defer span.End()
				}
			}
		}

		w := &containerWrapper{Container: c}

		if err := cmdHandler.Handle(ctx, recv.Command(), w); err != nil {
			err = fmt.Errorf("%w: %v", ErrHandleCommand, err)

			c.logError(ctx, KindCommands, err, m, fp.Some[Handler](cmdHandler))

			if c.mustProcessOrFail {
				c.close(err)
				return err
			}

			l.Errorln(err)
		}

		if w.err != nil {
			if c.mustProcessOrFail {
				// Could not emit events
				c.close(w.err)
				return w.err
			}

			l.Errorln(w.err)
		}

		l.Debugf("[c][ok] handled by cmd handler %s", k)

		return nil
	})
}

func (c *Container) eventBusSubscribe(ctx context.Context, l lol.Logger, bus EventBus) error {
	return bus.subscribe(ctx, func(ctx context.Context, m rp.Msg) error {
		msg := recvMsg{
			recordKey:       m.Key,
			recordPartition: m.Partition,
			recordTs:        m.Ts,
		}

		if err := bus.codec().Decode(m.Value, &msg); err != nil {
			//stop consuming
			return fmt.Errorf("%w: unable to decode msg: %v", ErrSubscribeNonRecoverable, err)
		}

		l.Debugf("[e][<] %v", msg)
		// Process event handlers
		var (
			k         = hashKey(msg)
			errEvent  error
			errSaga   error
			processed bool
		)

		eventHandler, eventHandlerExists := bus.eventHandler(msg)
		sagaHandler, sagaHandlerExists := bus.sagaHandler(msg)

		if !c.apmDisabled && (eventHandlerExists || sagaHandlerExists) {
			if tx := apm.TransactionFromContext(ctx); tx != nil {
				spanName := fmt.Sprintf("%s %s.%s.%s",
					cond.If(eventHandlerExists, "EVENT HANDLER", "SAGA HANDLER"),
					msg.Version(),
					msg.Resource(),
					msg.Action(),
				)

				tx.Context.SetLabel("cqrs.kind", KindEvents)
				tx.Context.SetLabel("cqrs.id", msg.ID())
				if userID, ok := msg.User().Unwrap(); ok {
					tx.Context.SetLabel("cqrs.user", userID)
				}
				// Add Redpanda message details as labels
				tx.Context.SetCustom("Redpanda_topic", m.Topic)
				tx.Context.SetCustom("Redpanda_partition", m.Partition)
				tx.Context.SetCustom("Redpanda_key", string(m.Key))
				tx.Context.SetCustom("Redpanda_timestamp", m.Ts.String())

				if span, spanCtx := apm.StartSpan(ctx, spanName, cqrs); span != nil {
					ctx = spanCtx
					defer span.End()
				}
			}
		}

		// THOUGHT: Consider using a sync.WaitGroup to process in parallel
		if eventHandlerExists {
			errEvent = eventHandler.Handle(ctx, msg.Event())
			processed = true
			l.Debugf("[c][ok] handled by event handler %s", k)

			if errEvent != nil && !c.errorCaptureDisabled {
				c.logError(ctx, KindEvents, errEvent, m, fp.Some[Handler](eventHandler))
			}
		}

		w := &containerWrapper{Container: c}

		if sagaHandlerExists {
			errSaga = sagaHandler.Handle(ctx, msg.Event(), w)
			processed = true
			l.Debugf("[c][ok] handled by saga handler %s", k)

			if errSaga != nil && !c.errorCaptureDisabled {
				c.logError(ctx, KindEvents, errSaga, m, fp.Some[Handler](sagaHandler))
			}
		}
		// End THOUGHT

		if errEvent != nil || errSaga != nil {
			eventError := "<nil>"
			if errEvent != nil {
				eventError = errEvent.Error()
			}

			sagaError := "<nil>"
			if errSaga != nil {
				sagaError = errSaga.Error()
			}

			err := fmt.Errorf("%w: event err = '%s', saga err = '%s'",
				ErrHandleEvent, eventError, sagaError)

			if c.mustProcessOrFail {
				c.close(err)
				return err
			}

			l.Errorln(err)
		}

		if w.err != nil {
			// Command called from saga failed
			if c.mustProcessOrFail {
				c.close(w.err)
				return w.err
			}

			l.Errorln(w.err)
		}

		if !processed && c.warnUnprocessed {
			l.Warnf("no event handler found for " + k)
		}

		return nil
	})

}

func (c *Container) wrapConsume(fn func() error) {
	err := fn()

	if err == nil {
		c.log.Warningf("subscribe returned without and error")
		return
	}

	// when subscribe returns, consuming has stopped. Log the error and consider
	// in the future to implement consumer restart.
	c.log.Errorf("subscribe returned error %v", err)
}

func (c *Container) send(ctx context.Context, kind string, msg sendMsg, bus bus) error {
	l := c.log.WithTrace(ctx).
		WithFields(lol.Fields{
			"kind":     kind,
			"resource": msg.Resource(),
			"action":   msg.Action(),
			"version":  msg.Version(),
			"key":      msg.Key(),
		})

	value, err := bus.codec().Encode(msg)

	if err != nil {
		l.Errorf("unable to encode: %v", err)
	}

	kmsg := rp.Msg{
		Topic: bus.topic(),
		Key:   partitionKey(msg),
		Value: value,
	}

	if !c.apmDisabled {
		if tx := apm.TransactionFromContext(ctx); tx != nil {
			spanName := fmt.Sprintf("%s %s.%s.%s",
				cond.If(kind == KindCommands, "COMMAND", "EVENT"),
				msg.Version(),
				msg.Resource(),
				msg.Action(),
			)

			tx.Context.SetLabel("cqrs.kind", kind)
			tx.Context.SetLabel("cqrs.id", msg.ID())
			if userID, ok := msg.User().Unwrap(); ok {
				tx.Context.SetLabel("cqrs.user", userID)
			}
			if key := msg.Key(); key != nil {
				tx.Context.SetCustom("cqrs.key", *key)
			}

			if span, spanCtx := apm.StartSpan(ctx, spanName, cqrs); span != nil {
				ctx = spanCtx
				defer span.End()
			}
		}
	}

	err = bus.publish(ctx, kmsg)

	if err != nil {
		l.Errorf("unable to publish to %s due to %v", bus.topic(), err)

		return fmt.Errorf("%w: unable to publish to %s: %v",
			ErrPublish, bus.topic(), err)
	}

	l.Debug("published")

	return nil
}

func (c *Container) logError(
	ctx context.Context,
	kind string,
	err error,
	msg rp.Msg,
	maybeHandler fp.Option[Handler],
) {
	if c.errorCaptureDisabled {
		return
	}

	now := time.Now().UTC()

	type RedpandaMsg struct {
		Topic     string          `json:"topic"`
		Key       string          `json:"key"`
		Value     json.RawMessage `json:"value"`
		Ts        time.Time       `json:"ts"`
		Partition int32           `json:"partition"`
	}

	type handlerInfo struct {
		Resource string `json:"resource"`
		Action   string `json:"action"`
		Version  string `json:"version"`
		Name     string `json:"name"`
	}

	type Error struct {
		Time     time.Time    `json:"time"`
		Kind     string       `json:"kind"`
		Side     any          `json:"side"`
		Payload  RedpandaMsg     `json:"payload"`
		Handler  *handlerInfo `json:"handler,omitempty"`
		Err      error        `json:"error"`
		ErrStack string       `json:"error_stack"`
		Hostname string       `json:"hostname"`
	}

	e := Error{
		Time: now,
		Kind: kind,
		Side: maybeHandler.MatchAny("consumer", "producer"),
		Payload: RedpandaMsg{
			Topic:     msg.Topic,
			Key:       string(msg.Key),
			Value:     msg.Value,
			Ts:        now,
			Partition: msg.Partition,
		},
		Err:      err,
		ErrStack: zero.B2S(debug.Stack()),
		Hostname: fp.OptionFromTupleErr(os.Hostname()).UnwrapOrDefault(),
	}

	if h, ok := maybeHandler.Unwrap(); ok {
		e.Handler = &handlerInfo{
			Resource: h.Resource(),
			Action:   h.Action(),
			Version:  h.Version(),
			Name:     reflect.TypeOf(h).Name(),
		}
	}

	value, errJson := json.Marshal(e)
	if errJson != nil {
		c.log.Errorf("unable to marshal err: %v", errJson)
		return
	}

	errProduce := c.producer.Publish(ctx, rp.Msg{
		Topic: TopicErrors,
		Value: value,
		Ts:    now,
	})

	if errProduce != nil {
		c.log.Errorf("unable to publish error log '%v' with payload '%v' due to '%v'",
			err, value, errProduce)
	}
}
