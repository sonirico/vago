package testit

import (
	"context"
	"log"
	"os"

	goredis "github.com/go-redis/redis/v8"
	"github.com/sonirico/vago/db/redis"

	"github.com/sonirico/vago/lol"
)

type RedisTestSuite struct {
	DB         *goredis.Client
	Log        lol.Logger
	pool       *DockerResourcesPool
	SetEnvFunc SetEnvFunc
	Config     redis.RedisConfig
}

func (s *RedisTestSuite) Setup() {
	s.pool = NewDockerResourcesPool(
		lol.ZeroTestLogger,
		os.Getenv("DOCKER_HOSTNAME"),
		NewRedisResource(s.SetEnvFunc),
	)
	if err := s.pool.Up(); err != nil {
		log.Panicln("cannot run docker", err)
	}

	ctx := context.Background()

	client, err := redis.OpenRedis(ctx, s.Config)
	if err != nil {
		log.Panicln("cannot connect to redis", err)
	}

	s.DB = client
}

func (s *RedisTestSuite) TearDown() {
	s.pool.Down()
}
