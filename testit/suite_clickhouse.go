package testit

import (
	"os"

	"github.com/ory/dockertest/v3"
	"github.com/sonirico/vago/db"

	"github.com/sonirico/vago/lol"
	"github.com/sonirico/vago/opts"
)

type ClickhouseHTTPTestSuite struct {
	DB     db.Handler
	DBExec db.Executor

	Log lol.Logger

	pool *DockerResourcesPool

	migrationsPath *string
	fixtures       [][]byte

	SetEnvFunc SetEnvFunc
}

func WithChLogger(log lol.Logger) opts.Configurator[ClickhouseHTTPTestSuite] {
	return opts.Fn[ClickhouseHTTPTestSuite](func(suite *ClickhouseHTTPTestSuite) {
		suite.Log = log
	})
}

func WithChMigrationsPath(path string) opts.Configurator[ClickhouseHTTPTestSuite] {
	return opts.Fn[ClickhouseHTTPTestSuite](func(suite *ClickhouseHTTPTestSuite) {
		suite.migrationsPath = &path
	})
}

func WithChFixtures(fixtures ...[]byte) opts.Configurator[ClickhouseHTTPTestSuite] {
	return opts.Fn[ClickhouseHTTPTestSuite](func(suite *ClickhouseHTTPTestSuite) {
		suite.fixtures = fixtures
	})
}

func (s *ClickhouseHTTPTestSuite) Setup(
	fns ...opts.Configurator[ClickhouseHTTPTestSuite],
) {
	opts.ApplyAll(s, fns...)

	defaultPath := "file://../../../migrations/clickhouse/"

	if s.migrationsPath != nil {
		defaultPath = *s.migrationsPath
	}

	if s.SetEnvFunc == nil {
		s.SetEnvFunc = func(dockerhost string, resource *dockertest.Resource) {

		}
	}

	logger := s.Log
	if logger == nil {
		logger = lol.ZeroTestLogger
	}

	s.pool = NewDockerResourcesPool(
		logger,
		os.Getenv("DOCKER_HOSTNAME"),
		NewClickhouseResource(defaultPath, logger, s.SetEnvFunc),
	)
	if err := s.pool.Up(); err != nil {
		logger.Panicln(err)
	}

	ch, err := db.OpenCH(
		os.Getenv("BROCK_CLICKHOUSE_HTTP_URL"),
		logger,
	) // TODO: parametrize
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

func (s *ClickhouseHTTPTestSuite) TearDown() {
	s.pool.Down()
}
