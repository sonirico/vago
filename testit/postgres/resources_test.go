package postgres

import (
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/sonirico/vago/lol"
	"github.com/sonirico/vago/testit"
)

func TestPostgres(t *testing.T) {
	t.Skip()

	const psqlMigrationsPath = "file://../../../migrations/postgresql/"

	host := testit.NewDockerResourcesPool(
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
