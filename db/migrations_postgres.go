package db

import (
	"fmt"
	"net/url"
	"strings"

	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sonirico/vago/lol"
)

func LaunchPostgresql(cfg MigrationsConfig, action string, logger lol.Logger) error {
	log := logger.WithField("database", "postgresql")

	managementUrl, err := url.Parse(cfg.Url)
	if err != nil {
		return fmt.Errorf("failed to parse management url: %w", err)
	}

	databaseName := strings.TrimPrefix(managementUrl.Path, "/")
	managementUrl.Path = "/postgres"

	// create the database
	writeWithoutDb, err := OpenPostgres(managementUrl.String(), log)

	if err != nil {
		return fmt.Errorf("cannot connect to postgres management: %w", err)
	}

	defer func() { _ = writeWithoutDb.Close() }()

	log.Infoln("Connected to postgres")

	if res, err := writeWithoutDb.Query(fmt.Sprintf(`SELECT FROM pg_database WHERE datname = '%s'`, databaseName)); err != nil {
		return fmt.Errorf(
			"error while checking the existence of database %s: %w",
			databaseName,
			err,
		)
	} else if !res.Next() {
		if _, err = writeWithoutDb.Exec(fmt.Sprintf(`CREATE DATABASE %s`, databaseName)); err != nil {
			return fmt.Errorf("error while creating database %s: %w", databaseName, err)
		}
	}

	// connect to the database
	write, err := OpenPostgres(cfg.Url, log)
	if err != nil {
		return fmt.Errorf("cannot connect to Postgres write: %w", err)
	}

	defer func() { _ = write.Close() }()

	log.Infoln("Connected to database: ", databaseName)

	// run migrations
	driver, err := postgres.WithInstance(write, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(cfg.MigrationsPath, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	switch action {
	case string(ActionUp):
		err = m.Up()
	case string(ActionDown):
		err = m.Steps(-1)
	}

	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Warnln("No changes to apply")
		} else {
			return fmt.Errorf("migrate action failed: %w", err)
		}
	}

	postgresUrl, err := url.Parse(cfg.Url)
	if err != nil {
		return fmt.Errorf("failed to parse postgres url: %w", err)
	}

	// if the management user is different from the write user then we need to assign
	// all the tables to the write user because they only have the bare
	// minimum permissions and can't do anything on someone else's tables
	if managementUrl.User.Username() != postgresUrl.User.Username() &&
		postgresUrl.User.Username() != "" {
		changeOwnerQuery := fmt.Sprintf(`REASSIGN OWNED BY "%s" TO "%s"`,
			managementUrl.User.Username(),
			postgresUrl.User.Username(),
		)
		if _, err = write.Exec(changeOwnerQuery); err != nil {
			return fmt.Errorf("failed to change owner of tables: %w", err)
		}
	}

	return nil
}
