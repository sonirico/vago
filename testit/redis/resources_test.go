package redis

import (
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/sonirico/vago/lol"
	"github.com/sonirico/vago/testit"
)

func TestRedis(t *testing.T) {
	t.Skip()
	host := testit.NewDockerResourcesPool(
		lol.ZeroTestLogger,
		os.Getenv("DOCKER_HOSTNAME"),
		NewRedisResource(func(dockerhost string, resource *dockertest.Resource) {
		}),
	)

	defer host.Down()
	if err := host.Up(); err != nil {
		t.Fatal(err)
	}
}
