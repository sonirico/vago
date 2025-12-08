package testit

import (
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/sonirico/vago/lol"
)

func TestPostgres(t *testing.T) {
	t.Skip()

	const psqlMigrationsPath = "file://../../../migrations/postgresql/"

	host := NewDockerResourcesPool(
		lol.ZeroTestLogger,
		os.Getenv("DOCKER_HOSTNAME"),
		NewPostgresResource(
			psqlMigrationsPath,
			lol.ZeroTestLogger,
			"BROCK_POSTGRES_URL",
			func(dockerhost string, resource *dockertest.Resource) error { return nil }, // TODO
			func(dockerhost string, resource *dockertest.Resource) {},                   // TODO
		),
	)

	defer host.Down()
	if err := host.Up(); err != nil {
		t.Fatal(err)
	}

}

func TestClickhouse(t *testing.T) {
	t.Skip()
	const chMigrationsPath = "file://../../../migrations/clickhouse/"

	host := NewDockerResourcesPool(
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

func TestRedis(t *testing.T) {
	t.Skip()
	host := NewDockerResourcesPool(
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

func TestRedpanda(t *testing.T) {
	t.Skip()
	host := NewDockerResourcesPool(
		lol.ZeroTestLogger,
		os.Getenv("DOCKER_HOSTNAME"),
		NewRedpandaResource(
			os.Getenv("DOCKER_HOSTNAME"),
			lol.ZeroTestLogger,
			func(dockerhost string, resource *dockertest.Resource) error { return nil }, // TODO
			func(dockerhost string, resource *dockertest.Resource) {},
		),
	)
	defer host.Down()

	if err := host.Up(); err != nil {
		t.Fatal(err)
	}
}
