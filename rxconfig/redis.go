package rxconfig

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sonirico/vago/codec"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"github.com/sonirico/vago/lol"
)

type (
	ClusteredNameGenerator func(service, node string) string
)

func (fn ClusteredNameGenerator) Key(service, node string) string {
	return fn(service, node)
}

type RedisItem[T any] struct {
	codec codec.Codec

	logger  lol.Logger
	redis   redis.UniversalClient
	changes chan Change[T]

	name    string
	service string
	node    string

	channel string
	hset    string

	latestConfigMutex sync.RWMutex
	latestConfig      Config[T]
	version           atomic.Int64

	closeOnce sync.Once
}

func (r *RedisItem[T]) close() {
	r.closeOnce.Do(func() {
		close(r.changes)
	})
}

func (r *RedisItem[T]) Path() string {
	return r.hset
}

func (r *RedisItem[T]) Put(ctx context.Context, data T) (value Config[T], err error) {
	r.latestConfigMutex.Lock()
	defer r.latestConfigMutex.Unlock()

	value, err = r.put(ctx, data)
	if err != nil {
		err = errors.Wrap(err, "unable to put")
		return
	}
	r.latestConfig = value
	return
}

func (r *RedisItem[T]) Get(ctx context.Context) (Config[T], error) {
	r.latestConfigMutex.Lock()
	defer r.latestConfigMutex.Unlock()
	value, err := r.get(ctx)
	if err != nil {
		err = errors.Wrap(err, "unable to get")

		r.logger.Errorf("cannot retrieve data from path: %s. Error: %s", r.node, err)
		return nil, err
	}
	r.latestConfig = value
	return value, nil
}

func (r *RedisItem[T]) Watch(ctx context.Context) {
	defer r.close()

	for {
		r.logger.Info("starting redis config watch")

		err := r.watch(ctx)

		if err != nil {
			r.logger.Errorf("watch failed due to %v",
				err)

			if errors.Is(err, context.Canceled) {
				r.logger.Info("context cancelled, terminating watch")
				return
			}
		}

		time.Sleep(time.Second)
	}
}

func (r *RedisItem[T]) Changes() <-chan Change[T] {
	return r.changes
}

func (r *RedisItem[T]) watch(ctx context.Context) error {
	conf, err := r.Get(ctx)

	if err != nil {
		r.logger.Errorf("cannot send first config to watch channel. Error: %s", err)
	} else if conf == nil {
		r.logger.Errorf("no config on path %s.", r.Path())
	} else {
		r.logger.Info("got initial config from hset")

		r.version.Store(-1)

		r.changes <- Change[T]{
			Prev:    nil,
			Next:    conf,
			Version: -1,
		}
	}

	r.logger.Debugf("redis subscribe to %s", r.channel)

	subscriber := r.redis.Subscribe(ctx, r.channel)

	defer func(ps *redis.PubSub) {
		pserr := ps.Close()
		if pserr != nil {
			r.logger.Errorf("cannot close redis subscriber due to '%v'", pserr)
		}
	}(subscriber)

	ch := subscriber.Channel(
		redis.WithChannelSize(32),
		redis.WithChannelHealthCheckInterval(time.Second*10),
		redis.WithChannelSendTimeout(time.Minute),
	)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-ch:
			if !ok {
				return errPubSubClosed
			}

			r.logger.Infof("got config from subscribe, on %s of length %d",
				msg.Channel, len(msg.Payload))

			if _, err = r.handleSub(ctx, []byte(msg.Payload)); err != nil {
				r.logger.Errorf("unable to handle sub due to '%v', waiting for next payload", err)
			}

		}
	}
}

// handleSub when publish event is received, means that config has been updated
// so fresh config shall be retrieved again and sent as a Change
func (r *RedisItem[T]) handleSub(ctx context.Context, data []byte) (Config[T], error) {
	r.latestConfigMutex.Lock()
	defer r.latestConfigMutex.Unlock()

	msg := pubSubMsg{}
	if err := r.codec.Decode(data, &msg); err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal pubSubMsg")
	}

	var prev Config[T]

	if r.latestConfig != nil {
		// initial config may not exist
		prev = r.latestConfig.Clone()
	}

	next, err := r.get(ctx)

	if err != nil {
		return nil, err
	}

	r.latestConfig = next

	r.changes <- Change[T]{
		Prev:    prev,
		Next:    next,
		Version: r.version.Add(1),
	}

	return next, nil
}

func (r *RedisItem[T]) put(ctx context.Context, data T) (Config[T], error) {
	bts, err := r.codec.Encode(data)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse config")
	}

	_, err = r.redis.HSet(ctx, r.hset, r.node, bts).Result()
	if err != nil {
		return nil, errors.Wrap(err, "unable to hset data")
	}

	msg := pubSubMsg{
		Service: r.service,
		Node:    r.node,
		Channel: r.channel,
		Time:    time.Now().UTC(),
	}

	payload, err := r.codec.Encode(msg)

	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal publish payload")
	}

	receivedBy, err := r.redis.Publish(ctx, r.channel, payload).Result()
	if err != nil {
		return nil, errors.Wrap(err, "unable to publish")
	}

	r.logger.Infof("publish on channel %s received by %d",
		r.channel, receivedBy)

	return Wrap[T](data), nil
}

func (r *RedisItem[T]) get(ctx context.Context) (Config[T], error) {
	bts, err := r.redis.HGet(ctx, r.hset, r.node).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errors.Wrapf(ErrDoesNotExist, "hset=%s, key=%s",
				r.hset, r.node)
		}
		return nil, err
	}

	var value T
	if err = r.codec.Decode(bts, &value); err != nil {
		return nil, errors.Wrap(err, "config factory failed")
	}

	return Wrap(value), nil
}

func NewRedis[T any](
	logger lol.Logger,
	redis redis.UniversalClient,
	codec codec.Codec,
	name string,
	service string,
	node string,
) *RedisItem[T] {
	return &RedisItem[T]{
		name:              name,
		codec:             codec,
		logger:            logger.WithField("rxconfig", "redis"),
		redis:             redis,
		changes:           make(chan Change[T]),
		service:           service,
		node:              node,
		channel:           hashKey(service, node, name),
		hset:              hashKey(service, name),
		latestConfigMutex: sync.RWMutex{},
		latestConfig:      nil,
	}
}

func hashKey(items ...string) string {
	return strings.Join(items, ":")
}
