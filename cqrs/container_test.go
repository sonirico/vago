package cqrs

import (
	"context"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/sonirico/vago/lol"
	"github.com/sonirico/vago/ptr"
	"github.com/sonirico/vago/rp"
	"github.com/stretchr/testify/assert"
	"go.elastic.co/apm/v2"
)

type (
	orderCreatedEvent struct {
		Status string  `json:"status"`
		Qty    float64 `json:"qty"`
	}

	lockBalanceCommand struct {
		Qty float64 `json:"qty"`
	}

	balanceLockedEvent struct {
		Locked float64 `json:"locked"`
		Free   float64 `json:"free"`
	}
)

var (
	producerConf = &rp.ProducerConfig{
		Timeout:                time.Minute,
		ConnIdleTimeout:        time.Minute,
		RecordDeliveryTimeout:  time.Minute,
		ProduceRequestTimeout:  time.Minute,
		RequestTimeoutOverhead: time.Minute,
		SessionTimeout:         time.Minute,
		Brokers:                []string{"localhost:30070"},
	}

	consumerConf = &rp.ConsumerConfig{
		Brokers:                      []string{"localhost:30070"},
		ConsumerGroup:                "cqrs-test",
		ConsumerBlockRebalanceOnPoll: true,
		WithLogger:                   false,
		MaxPollRecords:               100,
		APMConf: &rp.APMConfig{
			TxType: "cqrs-Redpanda",
			TxName: "cqrs-Redpanda-test",
		},
	}
)

func Test_Container(t *testing.T) {
	t.Skip()

	ctx := context.Background()
	log := lol.NewZerolog(lol.WithWriter(io.Discard))
	container := NewContainer(
		log,
		ContainerDisableErrorCapture(),
	)

	closeC := make(chan struct{})

	var (
		brokerEventBusID = "atani-events-broker-test"
		kycCommandBusID  = "atani-commands-kyc-test"
		kycEventBusID    = "atani-events-kyc-test"
	)

	brokerEventBus, err := NewEventBus(
		brokerEventBusID,
		"test.atani.events.broker",
		log,
		BusWithJsonCodec(),
		BusWithDisabledStartupPing(),
		BusWithProducerConfig(producerConf),
		BusWithConsumerConfig(consumerConf),
	)

	if err != nil {
		log.Fatal(err)
	}

	kycEventBus, err := NewEventBus(
		kycEventBusID,
		"test.atani.events.kyc",
		log,
		BusWithJsonCodec(),
		BusWithDisabledStartupPing(),
		BusWithProducerConfig(producerConf),
		BusWithConsumerConfig(consumerConf),
	)

	if err != nil {
		log.Fatal(err)
	}

	kycCommandsBus, err := NewCommandBus(
		kycCommandBusID,
		"test.atani.commands.kyc",
		log,
		BusWithJsonCodec(),
		BusWithDisabledStartupPing(),
		BusWithProducerConfig(producerConf),
		BusWithConsumerConfig(consumerConf),
	)

	if err != nil {
		log.Fatal(err)
	}

	var (
		actualOrderCreated       orderCreatedEvent
		actualLockBalanceCommand lockBalanceCommand
		actualBalanceLockedEvent balanceLockedEvent
	)

	orderCreatedSaga := NewSagaHandler(
		"0",
		"order",
		"created",
		func(ctx context.Context, event Event, commander Commander) error {
			err := json.Unmarshal(event.Payload(), &actualOrderCreated)
			if err != nil {
				return err
			}

			commander.Command(
				ctx,
				kycCommandBusID,
				NewSimpleCommand(
					"0",
					"applicant",
					"lock_balance",
					lockBalanceCommand{Qty: actualOrderCreated.Qty},
					nil,
				),
			)

			return nil
		},
		nil,
	)

	lockBalanceCommandHandler := NewCommandHandler(
		"0",
		"applicant",
		"lock_balance",
		func(ctx context.Context, cmd Command, eventer Eventer) error {
			err = json.Unmarshal(cmd.Payload(), &actualLockBalanceCommand)
			if err != nil {
				return err
			}

			eventer.Event(
				ctx,
				kycEventBusID,
				NewSimpleEvent("0", "applicant", "balance_locked", balanceLockedEvent{
					Locked: actualLockBalanceCommand.Qty,
					Free:   987,
				}, nil),
			)

			return nil
		},
	)

	balanceLockedEventHandler := NewEventHandler(
		"0",
		"applicant",
		"balance_locked",
		func(ctx context.Context, event Event) error {
			err = json.Unmarshal(event.Payload(), &actualBalanceLockedEvent)
			if err != nil {
				return err
			}

			close(closeC)

			return nil
		},
	)

	err = container.
		EventBus(kycEventBus.
			EventHandler(balanceLockedEventHandler)).
		EventBus(brokerEventBus.
			SagaHandler(orderCreatedSaga)).
		CommandBus(kycCommandsBus.
			CommandHandler(lockBalanceCommandHandler)).
		Start(ctx)

	if err != nil {
		close(closeC)
	}

	// Send initial event
	if c := consumerConf.APMConf; c != nil {
		tx := apm.DefaultTracer().
			StartTransactionOptions(c.TxName, c.TxType, apm.TransactionOptions{})
		defer tx.End()
		ctx = apm.ContextWithTransaction(ctx, tx)
	}

	err = container.Event(
		ctx,
		brokerEventBusID,
		NewSimpleEvent("0", "order", "created", orderCreatedEvent{
			Status: "created",
			Qty:    123.456,
		}, nil),
	)

	if err != nil {
		log.Fatal(err)
	}

	<-closeC

	assert.Equal(t, orderCreatedEvent{Status: "created", Qty: 123.456}, actualOrderCreated)
	assert.Equal(t, lockBalanceCommand{Qty: 123.456}, actualLockBalanceCommand)
	assert.Equal(t, balanceLockedEvent{Locked: 123.456, Free: 987}, actualBalanceLockedEvent)
}

func Test_Container_PartitionAffinity(t *testing.T) {
	t.Skip()

	ctx := context.Background()
	log := lol.NewZerolog(lol.WithWriter(io.Discard))
	container := NewContainer(log)
	closeC := make(chan struct{})

	type want struct {
		Key       string
		Partition int32
	}

	expected := []want{
		{
			Key:       "applicant_balance/user-123",
			Partition: 0,
		},
		{
			Key:       "applicant_balance/user-555",
			Partition: 1,
		},
	}
	actual := []want{}

	var (
		kycCommandBusID = "atani-commands-kyc-test"
	)

	kycCommandsBus, err := NewCommandBus(
		kycCommandBusID,
		"test.atani.commands.kyc",
		log,
		BusWithJsonCodec(),
		BusWithDisabledStartupPing(),
		BusWithProducerConfig(producerConf),
		BusWithConsumerConfig(consumerConf),
	)

	if err != nil {
		log.Fatal(err)
	}

	lockBalanceCommandHandler := NewCommandHandler(
		"0",
		"applicant_balance",
		"lock",
		func(ctx context.Context, cmd Command, _ Eventer) error {
			actual = append(actual, want{
				Partition: cmd.recordPartition,
				Key:       string(cmd.recordKey),
			})

			if len(actual) == len(expected) {
				close(closeC)
			}

			return nil
		},
	)

	err = container.
		CommandBus(kycCommandsBus.
			CommandHandler(lockBalanceCommandHandler)).
		Start(ctx)

	if err != nil {
		close(closeC)
	}

	command := lockBalanceCommand{Qty: 456.789}

	err = container.Command(
		ctx,
		kycCommandBusID,
		NewSimpleCommand(
			"0",
			"applicant_balance",
			"lock",
			command,
			ptr.Ptr("user-123"),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	err = container.Command(
		ctx,
		kycCommandBusID,
		NewSimpleCommand(
			"0",
			"applicant_balance",
			"lock",
			command,
			ptr.Ptr("user-555"),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	<-closeC

	assert.Equal(t, len(expected), len(actual))
	assert.Equal(t, expected, actual)
}
