package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sonirico/vago/db"
)

// NewRWOExecutorPgx creates a new db.ExecutorRWO backed by separate
// read-only and read-write pgxpool.Pool databases.
func NewRWOExecutorPgx(log db.Logger, ro, rw *pgxpool.Pool) db.ExecutorRWO {
	return db.NewRWOExecutor(log, newPgxAdapter(ro), newPgxAdapter(rw))
}
