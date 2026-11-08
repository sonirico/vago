package sqlite

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migratesqlite "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sonirico/vago/lol"

	"github.com/sonirico/vago/db"
)

// LaunchSqlite runs the migrations under cfg.MigrationsPath (a golang-migrate
// source URL such as "file://migrations") against the SQLite database at
// cfg.Url, which is created when missing. ActionUp applies every pending
// migration; ActionDown reverts the last applied one.
func LaunchSqlite(cfg db.MigrationsConfig, action string, logger lol.Logger) error {
	if err := db.ValidateMigrationAction(action); err != nil {
		return err
	}

	log := logger.WithField("database", "sqlite")

	write, err := OpenSqlite(log, cfg.Url)
	if err != nil {
		return fmt.Errorf("cannot connect to sqlite: %w", err)
	}

	defer func() { _ = write.Close() }()

	driver, err := migratesqlite.WithInstance(write, &migratesqlite.Config{})
	if err != nil {
		return fmt.Errorf("failed to create sqlite driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(cfg.MigrationsPath, "sqlite", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	switch action {
	case string(db.ActionUp):
		err = m.Up()
	case string(db.ActionDown):
		err = m.Steps(-1)
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
