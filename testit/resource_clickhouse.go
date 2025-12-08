package testit

import (
	"fmt"
	"os"
	"strings"

	"github.com/ory/dockertest/v3"
	"github.com/sonirico/vago/db"

	ch "github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sonirico/vago/lol"
)

func NewClickhouseResource(migrationsPath string, logger lol.Logger, envFunc SetEnvFunc) *Resource {
	pwd, _ := os.Getwd()

	volume := fmt.Sprintf(
		"%sbrock/it/internal/resources/clickhouse/config.d:/etc/clickhouse-server/config.d",
		pwd[0:strings.Index(pwd, "brock")])

	return &Resource{
		RunOptions: &dockertest.RunOptions{
			Mounts:       []string{volume},
			ExposedPorts: []string{"9000", "9009", "8123"},
			Repository:   "clickhouse/clickhouse-server",
			Tag:          "24.12-alpine",
			Hostname:     "clickhouse01", //Do not change, referenced in clickhouse config.xml file
		},

		HostConfig: &docker.HostConfig{
			Mounts: []docker.HostMount{
				{
					Type:   "tmpfs",
					Target: "/var/lib/clickhouse/",
				},
				{
					Type:   "tmpfs",
					Target: "/var/log/clickhouse-server/",
				},
			}},

		RetryFunc: func(dockerhost string, resource *dockertest.Resource) retryFunc {
			databaseUrl := fmt.Sprintf("tcp://%s:%s", dockerhost, resource.GetPort("9000/tcp"))

			return func() error {
				opts, err := ch.ParseDSN(databaseUrl)

				if err != nil {
					return err
				}

				db := ch.OpenDB(opts)

				if err != nil {
					return err
				}
				defer db.Close()
				return db.Ping()
			}
		},

		MigrateFunc: func(dockerhost string, resource *dockertest.Resource) error {
			log := logger.WithField("database", "clickhouse")

			url := fmt.Sprintf("tcp://%s:%s", dockerhost, resource.GetPort("9000/tcp"))

			cfg := db.MigrationsConfig{
				Url:            url,
				MigrationsPath: migrationsPath,
			}

			return db.LaunchClickhouse(cfg, "up", log)
		},

		SetEnvFunc: envFunc,
	}
}
