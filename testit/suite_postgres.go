package testit

import (
	"os"

	"github.com/ory/dockertest/v3"
	"github.com/sonirico/vago/db"

	"github.com/sonirico/vago/lol"
	"github.com/sonirico/vago/opts"
)

type PostgresTestSuite struct {
	DB     db.Handler
	DBExec db.Executor

	Log lol.Logger

	pool *DockerResourcesPool

	migrationsPath *string
	fixtures       [][]byte

	MigrateFunc MigrateFunc
	SetEnvFunc  SetEnvFunc
}

func WithPsqlLogger(log lol.Logger) opts.Configurator[PostgresTestSuite] {
	return opts.Fn[PostgresTestSuite](func(suite *PostgresTestSuite) {
		suite.Log = log
	})
}

func WithPsqlMigrationsPath(path string) opts.Configurator[PostgresTestSuite] {
	return opts.Fn[PostgresTestSuite](func(suite *PostgresTestSuite) {
		suite.migrationsPath = &path
	})
}

func WithPsqlFixtures(fixtures ...[]byte) opts.Configurator[PostgresTestSuite] {
	return opts.Fn[PostgresTestSuite](func(suite *PostgresTestSuite) {
		suite.fixtures = fixtures
	})
}

func (s *PostgresTestSuite) Setup(
	fns ...opts.Configurator[PostgresTestSuite],
) {
	opts.ApplyAll(s, fns...)

	psqlMigrationsPath := "file://../../../migrations/postgresql/"

	if s.migrationsPath != nil {
		psqlMigrationsPath = *s.migrationsPath
	}

	if s.MigrateFunc == nil {
		s.MigrateFunc = func(dockerhost string, resource *dockertest.Resource) error {
			return nil
		}
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
		NewPostgresResource(
			psqlMigrationsPath,
			s.Log,
			"BROCK_POSTGRES_URL",
			s.MigrateFunc,
			s.SetEnvFunc,
		),
	)
	if err := s.pool.Up(); err != nil {
		logger.Panicln(err)
	}

	psql, err := db.OpenPgx(s.Log, os.Getenv("BROCK_POSTGRES_URL"))
	if err != nil {
		logger.Panicln("cannot connect to postgresql db", err)
	}

	for _, val := range s.fixtures {
		fix := NewSQLFixture(psql, val)
		if err = fix.Load(); err != nil {
			logger.Panicln(err)
		}
	}

	s.Log = logger
	s.DB = psql
	s.DBExec = db.NewExecutor(logger, psql)
}

func (s *PostgresTestSuite) TearDown() {
	s.pool.Down()
}
