package rp

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/sonirico/vago/lol"
	"github.com/sonirico/vago/slices"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/apm/v2"
)

const (
	defaultPollRecords = 100
)

type (
	ConsumerHandler func(ctx context.Context, m Msg) error

	ConsumerConfig struct {
		Brokers []string

		ConsumerGroup                string
		ConsumerBlockRebalanceOnPoll bool
		WithLogger                   bool
		WithLogLevel                 LogLevel
		MaxPollRecords               int
		WithLoggingHooks             bool
		WithAppName                  string
		WithVersion                  string

		APMConf *APMConfig
	}

	BasicConsumer struct {
		log lol.Logger

		cfg ConsumerConfig

		topics []string
		client *kgo.Client

		apm *APMConfig

		closed bool

		closeOnce sync.Once
	}

	ConsumerFactory func(
		log lol.Logger,
		cfg ConsumerConfig,
		topics []string,
		apmConf *APMConfig,
	) (*BasicConsumer, error)

	consumerHandler func(ctx context.Context, m *kgo.Record) error
)

func (c *BasicConsumer) Ping(ctx context.Context) (err error) {
	err = c.client.Ping(ctx)
	return
}

func (c *BasicConsumer) Close() {
	c.safeClose()
}

func (c *BasicConsumer) Subscribe(ctx context.Context, handler ConsumerHandler) error {
	if c.closed {
		return ErrConsumerClosed
	}

	return c.start(ctx, func(ctx context.Context, rec *kgo.Record) error {
		if c.apm != nil {
			if header, found := slices.Find(rec.Headers, func(item kgo.RecordHeader) bool {
				return item.Key == apmhttp.ElasticTraceparentHeader
			}); found {
				if traceCtx, err := apmhttp.ParseTraceparentHeader(string(header.Value)); err == nil {
					txOpts := apm.TransactionOptions{TraceContext: traceCtx}

					txName := c.apm.TxName
					if !isset(txName) {
						txName = rec.Topic
					}
					txType := c.apm.TxType
					if !isset(txType) {
						txType = "subscribe"
					}

					tx := apm.DefaultTracer().
						StartTransactionOptions(txName, txType, txOpts)
					ctx = apm.ContextWithTransaction(ctx, tx)
					defer tx.End()
				}
			}
		}

		return handler(
			ctx,
			Msg{
				Topic:     rec.Topic,
				Key:       rec.Key,
				Value:     rec.Value,
				Partition: rec.Partition,
				Ts:        rec.Timestamp,
			},
		)
	})
}

func NewConsumer(
	log lol.Logger,
	cfg ConsumerConfig,
	topics []string,
	apmConf *APMConfig,
) (*BasicConsumer, error) {
	if len(topics) < 1 {
		return nil, ErrTopicsRequired
	}

	c := &BasicConsumer{log: log, cfg: cfg, topics: topics, apm: apmConf}

	if cfg.APMConf != nil {
		c.apm = cfg.APMConf
	}

	if c.cfg.MaxPollRecords == 0 {
		c.cfg.MaxPollRecords = defaultPollRecords
	}

	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ConsumerGroup(cfg.ConsumerGroup),
		kgo.ConsumeTopics(topics...),
		kgo.DisableAutoCommit(),
	}

	if isset(cfg.WithAppName) && isset(cfg.WithVersion) {
		opts = append(opts, kgo.SoftwareNameAndVersion(cfg.WithAppName, cfg.WithVersion))
	}

	if cfg.ConsumerBlockRebalanceOnPoll {
		opts = append(opts, kgo.BlockRebalanceOnPoll())
	}

	if cfg.WithLoggingHooks {
		opts = append(
			opts,
			kgo.OnPartitionsAssigned(
				func(ctx context.Context, cl *kgo.Client, assigned map[string][]int32) {
					log.Debugf("PARTITIONS ASSIGNED => %v\n", assigned)
				},
			),
		)
		opts = append(
			opts,
			kgo.OnPartitionsRevoked(
				func(ctx context.Context, cl *kgo.Client, revoked map[string][]int32) {
					log.Debugf("PARTITIONS REVOKED => %v\n", revoked)
				},
			),
		)
		opts = append(
			opts,
			kgo.OnPartitionsLost(
				func(ctx context.Context, cl *kgo.Client, lost map[string][]int32) {
					log.Debugf("PARTITIONS LOST => %v\n", lost)
				},
			),
		)
	}

	if cfg.WithLogger {
		level := kgo.LogLevel(cfg.WithLogLevel)
		if cfg.WithLogLevel == LogLevelNone {
			level = kgo.LogLevelInfo
		}

		opts = append(
			opts,
			kgo.WithLogger(kgo.BasicLogger(os.Stderr, level, func() string {
				return "redpanda[consumer][" + time.Now().Format(time.RFC3339) + "]"
			})),
		)
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = client.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	c.client = client

	return c, nil
}

func (c *BasicConsumer) poll(ctx context.Context, handler consumerHandler) error {
	l := c.log.WithTrace(ctx)

	// PollRecords is strongly recommended when using
	// BlockRebalanceOnPoll. You can tune how many records to
	// process at once (upper bound -- could all be on one
	// partition), ensuring that your processor loops complete fast
	// enough to not block a rebalance too long.
	l.Debugf("polling records: %s -> %d", c.topics, c.cfg.MaxPollRecords)

	fetches := c.client.PollRecords(ctx, c.cfg.MaxPollRecords)
	if fetches.IsClientClosed() {
		l.Errorf("client is closed")
		return ErrConsumerClosed
	}

	fetches.EachError(func(topic string, partition int32, err error) {
		if errors.Is(err, context.Canceled) {
			return
		}

		l.Errorf(
			"failed to fetch records (topic: %s, partition: %d): %v",
			topic,
			partition,
			err,
		)
	})

	var (
		err  error
		recs []*kgo.Record
	)

	fetches.EachPartition(func(p kgo.FetchTopicPartition) {
		if err != nil {
			return
		}

		for i := range p.Records {
			rec := p.Records[i]

			if err2 := handler(ctx, rec); err2 != nil {
				err = fmt.Errorf(
					"topic: %s, partition: %d, offset: %d: %w",
					rec.Topic,
					rec.Partition,
					rec.Offset,
					err2,
				)

				break
			} else {
				recs = append(recs, rec)
			}
		}
	})

	l.Infof("committing %d records", len(recs))
	if err2 := c.client.CommitRecords(ctx, recs...); err2 != nil {
		return fmt.Errorf("failed to commit offsets: %w", err2)
	}

	l.Infof("committed %d records", len(recs))
	c.client.AllowRebalance()

	return err
}

func (c *BasicConsumer) start(ctx context.Context, handler consumerHandler) (err error) {
	defer func() {
		c.log.Errorf("consumer stopping after poll returned error: %v", err)
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err = c.poll(ctx, handler); err != nil {
				return err
			}
		}
	}
}

func (c *BasicConsumer) safeClose() {
	c.closeOnce.Do(func() {
		c.client.Close()
		c.closed = true
	})
}

func NewConsumerFactory() ConsumerFactory {
	return func(
		log lol.Logger,
		cfg ConsumerConfig,
		topics []string,
		apmConf *APMConfig,
	) (*BasicConsumer, error) {
		return NewConsumer(log, cfg, topics, apmConf)
	}
}
