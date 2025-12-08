package testit

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/ory/dockertest/v3"
)

func NewRedisResource(envFunc SetEnvFunc) *Resource {
	return &Resource{
		RunOptions: &dockertest.RunOptions{
			Repository:   "redis",
			Tag:          "6.2",
			Hostname:     "redis",
			ExposedPorts: []string{"6379/tcp"},
		},
		RetryFunc: func(dockerhost string, resource *dockertest.Resource) retryFunc {
			redisAddr := fmt.Sprintf("%s:%s", dockerhost, resource.GetPort("6379/tcp"))
			return func() error {
				db := redis.NewClient(&redis.Options{
					Addr: redisAddr,
				})

				return db.Ping(context.Background()).Err()
			}
		},
		MigrateFunc: func(dockerhost string, resource *dockertest.Resource) error {
			return nil
		},
		SetEnvFunc: envFunc,
	}
}
