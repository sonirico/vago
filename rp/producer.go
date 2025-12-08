package rp

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sonirico/vago/lol"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/apm/v2"
)

const (
	CompressionSnappy int = iota
	CompressionGzip
	CompressionLz4
	CompressionZstd
	CompressionNone
)

var (
	NoCallback = func(Msg, error) {}
)

type (
	ProducerConfig struct {
		FlushTimeout           time.Duration
		Timeout                time.Duration
		ProduceRequestTimeout  time.Duration
		ConnIdleTimeout        time.Duration
		RequestTimeoutOverhead time.Duration
		RecordDeliveryTimeout  time.Duration
		Linger                 time.Duration
		SessionTimeout         time.Duration
		MaxBufferedRecords     int
		MaxBytes               int32
		Brokers                []string
		WithLogger             bool
		WithNoAPM              bool
		Compression            int
		ProduceSync            bool
		App                    string
		Version                string
	}

	BasicProducer struct {
		cli *kgo.Client
		cfg ProducerConfig
		log lol.Logger
	}
)

func NewProducer(
	ctx context.Context,
	config ProducerConfig,
	log lol.Logger,
) (Producer, error) {
	var opts []kgo.Opt

	if len(config.Brokers) < 1 {
		return nil, fmt.Errorf("brokers: %w", ErrConfig)
	}

	opts = append(opts, kgo.SeedBrokers(config.Brokers...))

	if config.ProduceRequestTimeout > 0 {
		opts = append(opts, kgo.ProduceRequestTimeout(config.ProduceRequestTimeout))
	}

	if config.ConnIdleTimeout > 0 {
		opts = append(opts, kgo.ConnIdleTimeout(config.ConnIdleTimeout))
	}

	if config.RequestTimeoutOverhead > 0 {
		opts = append(opts, kgo.RequestTimeoutOverhead(config.RequestTimeoutOverhead))
	}

	if config.RecordDeliveryTimeout > 0 {
		opts = append(opts, kgo.RecordDeliveryTimeout(config.RecordDeliveryTimeout))
	}

	if config.SessionTimeout > 0 {
		opts = append(opts, kgo.SessionTimeout(config.SessionTimeout))
	}

	if config.App != "" && config.Version != "" {
		opts = append(opts, kgo.SoftwareNameAndVersion(config.App, config.Version))
	}

	if config.Linger > 0 {
		opts = append(opts, kgo.ProducerLinger(config.Linger))
	}

	if config.MaxBytes > 0 {
		opts = append(opts, kgo.ProducerBatchMaxBytes(config.MaxBytes))
	}

	if config.WithLogger {
		opts = append(
			opts,
			kgo.WithLogger(kgo.BasicLogger(os.Stderr, kgo.LogLevelInfo, func() string {
				return "redpanda[BasicProducer]"
			})),
		)
	}

	if config.MaxBufferedRecords > 0 {
		opts = append(opts, kgo.MaxBufferedRecords(config.MaxBufferedRecords))
	}

	switch config.Compression {
	case CompressionLz4:
		opts = append(opts, kgo.ProducerBatchCompression(kgo.Lz4Compression()))
	case CompressionGzip:
		opts = append(opts, kgo.ProducerBatchCompression(kgo.GzipCompression()))
	case CompressionZstd:
		opts = append(opts, kgo.ProducerBatchCompression(kgo.ZstdCompression()))
	case CompressionSnappy:
		opts = append(opts, kgo.ProducerBatchCompression(kgo.SnappyCompression()))
	case CompressionNone:
	}

	cli, err := kgo.NewClient(opts...)

	if err != nil {
		return nil, err
	}

	if err = cli.Ping(ctx); err != nil {
		return nil, err
	}

	producer := &BasicProducer{cli: cli, log: log.WithField("type", "BasicProducer"), cfg: config}

	return producer, nil
}

func (p *BasicProducer) Ping(ctx context.Context) (err error) {
	err = p.cli.Ping(ctx)
	return
}

func (p *BasicProducer) PublishAsync(
	ctx context.Context,
	msg Msg,
	onPublish func(Msg, error),
) (err error) {
	return p.publish(ctx, msg, onPublish)
}

func (p *BasicProducer) Publish(ctx context.Context, msg Msg) (err error) {
	return p.publish(ctx, msg, nil)
}

func (p *BasicProducer) publish(
	ctx context.Context,
	msg Msg,
	onPublished func(Msg, error),
) (err error) {
	var headers []kgo.RecordHeader

	if !p.cfg.WithNoAPM {
		// If APM is not disabled...
		if tx := apm.TransactionFromContext(ctx); tx != nil {
			headers = make([]kgo.RecordHeader, 0, 2)
			headers = append(
				headers,
				kgo.RecordHeader{
					Key:   apmhttp.ElasticTraceparentHeader,
					Value: []byte(apmhttp.FormatTraceparentHeader(tx.TraceContext())),
				},
			)
			if span, spanCtx := apm.StartSpan(ctx, "PUBLISH "+msg.Topic, apmTxType); span != nil {
				ctx = spanCtx
				span.Context.SetLabel("topic", msg.Topic)
				span.Context.SetDatabase(apm.DatabaseSpanContext{
					Instance:  "Transport-BasicProducer",
					Statement: string(msg.Value),
					Type:      apmTxType,
				})
				defer span.End()
			}
		}
	}

	if onPublished != nil {
		// Async Publish
		rec := &kgo.Record{
			Topic:   msg.Topic,
			Key:     msg.Key,
			Value:   msg.Value,
			Headers: headers,
		}

		p.cli.Produce(ctx, rec, func(record *kgo.Record, err error) {
			if err != nil {
				p.log.WithFields(
					lol.Fields{
						"topic": msg.Topic,
						"key":   string(msg.Key),
						"value": string(msg.Value),
					},
				).WithTrace(ctx).Errorf("publish async error: '%v'", err)
			}
			onPublished(Msg{
				Topic:     record.Topic,
				Key:       record.Key,
				Value:     record.Value,
				Ts:        record.Timestamp,
				Partition: record.Partition,
			}, err)
		})
	} else {
		if err = p.cli.ProduceSync(ctx, &kgo.Record{
			Topic:   msg.Topic,
			Key:     msg.Key,
			Value:   msg.Value,
			Headers: headers,
		}).FirstErr(); err != nil {
			p.log.WithFields(
				lol.Fields{
					"topic": msg.Topic,
					"key":   string(msg.Key),
					"value": string(msg.Value),
				},
			).WithTrace(ctx).Errorf("publish sync error: '%v'", err)
		}
	}

	return err
}

func (p *BasicProducer) Flush(ctx context.Context) error {
	return p.cli.Flush(ctx)
}

func (p *BasicProducer) Close() {
	p.cli.Close()
}

func (c ProducerConfig) GetFlushTimeout() time.Duration {
	if c.FlushTimeout != 0 {
		return c.FlushTimeout
	}

	if c.Timeout != 0 {
		return c.Timeout
	}

	return time.Minute
}
