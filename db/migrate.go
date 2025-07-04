package db

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
