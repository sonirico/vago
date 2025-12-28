package testit

import (
	"fmt"
	"os"

	"github.com/ory/dockertest/v3"
	"github.com/sonirico/vago/db"

	"github.com/sonirico/vago/lol"
	"github.com/sonirico/vago/opts"
)

// ClickhouseHTTPTestSuite provides a test suite with a real ClickHouse instance.
type ClickhouseHTTPTestSuite struct {
	DB     db.Handler
	DBExec db.Executor

	Log lol.Logger

	pool *DockerResourcesPool

	migrationsPath *string
	configVolume   *string
	fixtures       [][]byte
	dsnEnvVar      string // Environment variable name for DSN (default: CLICKHOUSE_DSN)

	SetEnvFunc SetEnvFunc
}

// WithChLogger sets the logger for the suite.
func WithChLogger(log lol.Logger) opts.Configurator[ClickhouseHTTPTestSuite] {
	return opts.Fn[ClickhouseHTTPTestSuite](func(suite *ClickhouseHTTPTestSuite) {
		suite.Log = log
	})
}

// WithChMigrationsPath sets the migrations path.
func WithChMigrationsPath(path string) opts.Configurator[ClickhouseHTTPTestSuite] {
	return opts.Fn[ClickhouseHTTPTestSuite](func(suite *ClickhouseHTTPTestSuite) {
		suite.migrationsPath = &path
	})
}

// WithChSuiteConfigVolume sets a custom config volume to mount.
func WithChSuiteConfigVolume(path string) opts.Configurator[ClickhouseHTTPTestSuite] {
	return opts.Fn[ClickhouseHTTPTestSuite](func(suite *ClickhouseHTTPTestSuite) {
		suite.configVolume = &path
	})
}

// WithChFixtures sets fixtures to load after setup.
func WithChFixtures(fixtures ...[]byte) opts.Configurator[ClickhouseHTTPTestSuite] {
	return opts.Fn[ClickhouseHTTPTestSuite](func(suite *ClickhouseHTTPTestSuite) {
		suite.fixtures = fixtures
	})
}

// WithChDSNEnvVar sets the environment variable name for the DSN.
func WithChDSNEnvVar(envVar string) opts.Configurator[ClickhouseHTTPTestSuite] {
	return opts.Fn[ClickhouseHTTPTestSuite](func(suite *ClickhouseHTTPTestSuite) {
		suite.dsnEnvVar = envVar
	})
}

// Setup initializes the test suite with a real ClickHouse container.
func (s *ClickhouseHTTPTestSuite) Setup(
	fns ...opts.Configurator[ClickhouseHTTPTestSuite],
) {
	opts.ApplyAll(s, fns...)

	if s.SetEnvFunc == nil {
		s.SetEnvFunc = func(dockerhost string, resource *dockertest.Resource) {}
	}

	if s.dsnEnvVar == "" {
		s.dsnEnvVar = "CLICKHOUSE_DSN"
	}

	logger := s.Log
	if logger == nil {
		logger = lol.ZeroTestLogger
	}

	// Build resource options
	resourceOpts := []ClickhouseResourceOpt{
		WithChResourceLogger(logger),
		WithChSetEnvFunc(func(dockerhost string, resource *dockertest.Resource) {
			dsn := fmt.Sprintf("clickhouse://default:@%s:%s/default",
				dockerhost, resource.GetPort("9000/tcp"))
			os.Setenv(s.dsnEnvVar, dsn)
			s.SetEnvFunc(dockerhost, resource)
		}),
	}

	if s.migrationsPath != nil {
		resourceOpts = append(resourceOpts, WithChMigrations(*s.migrationsPath))
	}

	if s.configVolume != nil {
		resourceOpts = append(resourceOpts, WithChConfigVolume(*s.configVolume))
	}

	s.pool = NewDockerResourcesPool(
		logger,
		os.Getenv("DOCKER_HOSTNAME"),
		NewClickhouseResourceWithOpts(resourceOpts...),
	)
	if err := s.pool.Up(); err != nil {
		logger.Panicln(err)
	}

	ch, err := db.OpenCH(os.Getenv(s.dsnEnvVar), logger)
	if err != nil {
		logger.Panicln("cannot connect to clickhouse db", err)
	}

	for _, fixture := range s.fixtures {
		fix := NewCHFixture(ch, fixture)
		if err = fix.Load(); err != nil {
			logger.Panicln(err)
		}
	}

	s.Log = logger
	s.DB = ch
	s.DBExec = db.NewExecutor(s.Log, s.DB)
}

// TearDown cleans up the test suite.
func (s *ClickhouseHTTPTestSuite) TearDown() {
	if s.DB != nil {
		s.DB.Close()
	}
	if s.pool != nil {
		s.pool.Down()
	}
}
