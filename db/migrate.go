package db

import "fmt"

type Action string
type Database string

const (
	ActionUp   Action = "up"
	ActionDown Action = "down"

	DatabaseClickhouse     Database = "ch"
	DatabasePostgresBroker Database = "psql-broker"
	DataBasePostgresAuth   Database = "psql-auth"
)

type MigrationsConfig struct {
	// Url is the connection URL
	Url string
	// Path to folder containing migrations
	MigrationsPath string
}

// ValidateMigrationAction reports an error if action is neither ActionUp nor
// ActionDown. Backend migration runners call it before opening any database
// connection, so a bad action never reports success having applied nothing.
func ValidateMigrationAction(action string) error {
	switch action {
	case string(ActionUp), string(ActionDown):
		return nil
	default:
		return fmt.Errorf(
			"unknown migration action %q: expected %q or %q",
			action, ActionUp, ActionDown,
		)
	}
}
