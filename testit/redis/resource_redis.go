package redis

import (
	"context"
	"fmt"

	goredis "github.com/go-redis/redis/v8"
	"github.com/ory/dockertest/v3"
	"github.com/sonirico/vago/testit"
)

func NewRedisResource(envFunc testit.SetEnvFunc) *testit.Resource {
	return &testit.Resource{
		RunOptions: &dockertest.RunOptions{
			Repository:   "redis",
			Tag:          "6.2",
			Hostname:     "redis",
			ExposedPorts: []string{"6379/tcp"},
		},
		RetryFunc: func(dockerhost string, resource *dockertest.Resource) testit.RetryOp {
			redisAddr := fmt.Sprintf("%s:%s", dockerhost, resource.GetPort("6379/tcp"))
			return func() error {
				db := goredis.NewClient(&goredis.Options{
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
