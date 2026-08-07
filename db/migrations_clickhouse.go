package db

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sonirico/vago/lol"
)

func LaunchClickhouse(cfg MigrationsConfig, action string, logger lol.Logger) error {
	if err := validateMigrationAction(action); err != nil {
		return err
	}

	log := logger.WithField("database", "clickhouse")
	conn, err := OpenClickhouse(cfg.Url, log)

	if err != nil {
		return fmt.Errorf("cannot connect to clickhouse: %w", err)
	}

	driver, err := clickhouse.WithInstance(conn, &clickhouse.Config{})
	if err != nil {
		return fmt.Errorf("with instance: %w", err)
	}

	defer func() { _ = driver.Close() }()

	m, err := migrate.NewWithDatabaseInstance(cfg.MigrationsPath, "clickhouse", driver)
	if err != nil {
		return fmt.Errorf("new with database instance: %w", err)
	}

	switch action {
	case string(ActionUp):
		err = m.Up()
	case string(ActionDown):
		err = m.Steps(-1)
	default:
		return fmt.Errorf(
			"unknown migration action %q: expected %q or %q",
			action, ActionUp, ActionDown,
		)
	}

	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Warnln("No changes to apply")
		} else {
			return fmt.Errorf("migrate action failed: %w", err)
		}
	}

	return nil
}
