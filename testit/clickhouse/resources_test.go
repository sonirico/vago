package clickhouse

import (
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/sonirico/vago/lol"
	"github.com/sonirico/vago/testit"
)

func TestClickhouse(t *testing.T) {
	t.Skip()
	const chMigrationsPath = "file://../../../migrations/clickhouse/"

	host := testit.NewDockerResourcesPool(
		lol.ZeroTestLogger,
		os.Getenv("DOCKER_HOSTNAME"),
		NewClickhouseResource(
			chMigrationsPath,
			lol.ZeroTestLogger,
			func(dockerhost string, resource *dockertest.Resource) {}, // TODO
		),
	)
	defer host.Down()
	if err := host.Up(); err != nil {
		t.Fatal(err)
	}
}
