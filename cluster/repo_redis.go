package cluster

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisRepo struct {
	cli redis.UniversalClient
}

func (r *redisRepo) Set(
	ctx context.Context,
	key string,
	value []byte,
	ttl time.Duration,
) (err error) {
	_, err = r.cli.Set(ctx, key, value, ttl).Result()
	return
}

func (r *redisRepo) Scan(ctx context.Context, pattern string) (map[string][]byte, error) {
	// Retrieve all keys matching the pattern.
	keys, err := r.cli.Keys(ctx, pattern+"*").Result()
	if err != nil {
		return nil, err
	}

	if len(keys) < 1 {
		return nil, nil
	}

	// Fetch values for all keys in one go.
	values, err := r.cli.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	res := make(map[string][]byte, len(values))
	for i, key := range keys {
		res[key] = []byte(values[i].(string))
	}

	return res, nil
}
