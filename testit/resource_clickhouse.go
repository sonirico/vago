package testit

import (
	"fmt"

	"github.com/ory/dockertest/v3"
	"github.com/sonirico/vago/db"
	"github.com/sonirico/vago/opts"

	ch "github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sonirico/vago/lol"
)

// ClickhouseResourceConfig holds configuration for ClickHouse test resource.
type ClickhouseResourceConfig struct {
	MigrationsPath string
	ConfigVolume   string // Optional: path to config.d to mount
	Tag            string
	Logger         lol.Logger
	SetEnvFunc     SetEnvFunc
}

// ClickhouseResourceOpt configures a ClickHouse resource.
type ClickhouseResourceOpt = opts.Configurator[ClickhouseResourceConfig]

// WithChConfigVolume sets a custom config volume to mount.
func WithChConfigVolume(path string) ClickhouseResourceOpt {
	return opts.Fn[ClickhouseResourceConfig](func(c *ClickhouseResourceConfig) {
		c.ConfigVolume = path
	})
}

// WithChTag sets the ClickHouse image tag.
func WithChTag(tag string) ClickhouseResourceOpt {
	return opts.Fn[ClickhouseResourceConfig](func(c *ClickhouseResourceConfig) {
		c.Tag = tag
	})
}

// WithChMigrations sets the migrations path.
func WithChMigrations(path string) ClickhouseResourceOpt {
	return opts.Fn[ClickhouseResourceConfig](func(c *ClickhouseResourceConfig) {
		c.MigrationsPath = path
	})
}

// WithChResourceLogger sets the logger.
func WithChResourceLogger(log lol.Logger) ClickhouseResourceOpt {
	return opts.Fn[ClickhouseResourceConfig](func(c *ClickhouseResourceConfig) {
		c.Logger = log
	})
}

// WithChSetEnvFunc sets the environment setup function.
func WithChSetEnvFunc(fn SetEnvFunc) ClickhouseResourceOpt {
	return opts.Fn[ClickhouseResourceConfig](func(c *ClickhouseResourceConfig) {
		c.SetEnvFunc = fn
	})
}

// NewClickhouseResourceWithOpts creates a ClickHouse resource with options.
func NewClickhouseResourceWithOpts(options ...ClickhouseResourceOpt) *Resource {
	cfg := ClickhouseResourceConfig{
		Tag:    "24.12-alpine",
		Logger: lol.ZeroDiscardLogger,
		SetEnvFunc: func(dockerhost string, resource *dockertest.Resource) {
		},
	}
	opts.ApplyAll(&cfg, options...)

	runOpts := &dockertest.RunOptions{
		ExposedPorts: []string{"9000", "9009", "8123"},
		Repository:   "clickhouse/clickhouse-server",
		Tag:          cfg.Tag,
		Hostname:     "clickhouse01",
	}

	// Only mount config volume if specified
	if cfg.ConfigVolume != "" {
		runOpts.Mounts = []string{cfg.ConfigVolume + ":/etc/clickhouse-server/config.d"}
	}

	return &Resource{
		RunOptions: runOpts,

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
				chOpts, err := ch.ParseDSN(databaseUrl)
				if err != nil {
					return err
				}

				conn := ch.OpenDB(chOpts)
				defer conn.Close()
				return conn.Ping()
			}
		},

		MigrateFunc: func(dockerhost string, resource *dockertest.Resource) error {
			if cfg.MigrationsPath == "" {
				return nil // Skip migrations if not configured
			}

			log := cfg.Logger.WithField("database", "clickhouse")
			url := fmt.Sprintf("tcp://%s:%s", dockerhost, resource.GetPort("9000/tcp"))

			migCfg := db.MigrationsConfig{
				Url:            url,
				MigrationsPath: cfg.MigrationsPath,
			}

			return db.LaunchClickhouse(migCfg, "up", log)
		},

		SetEnvFunc: cfg.SetEnvFunc,
	}
}

// NewClickhouseResource creates a ClickHouse resource (legacy API, kept for compatibility).
func NewClickhouseResource(migrationsPath string, logger lol.Logger, envFunc SetEnvFunc) *Resource {
	return NewClickhouseResourceWithOpts(
		WithChMigrations(migrationsPath),
		WithChResourceLogger(logger),
		WithChSetEnvFunc(envFunc),
	)
}
