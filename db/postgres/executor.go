package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sonirico/vago/db"
)

// NewExecutorPgx creates a new db.Executor for a pgxpool.Pool database.
func NewExecutorPgx(log db.Logger, pool *pgxpool.Pool) db.Executor {
	return db.NewExecutor(log, newPgxAdapter(pool))
}
