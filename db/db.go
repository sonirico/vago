// Package db provides the backend-agnostic contract, transaction management,
// and efficient bulk operations shared by every database backend this
// module supports. It has no database driver dependency: connecting to an
// actual backend is the job of the sibling backend modules below.
//
// Features:
//   - Backend-agnostic interfaces (Handler, Querier, Tx, Result, Rows, Row).
//   - Context and transaction management with hooks for after-commit actions.
//   - Executor interfaces for read-only, read-write, and transactional operations.
//   - Bulk DML helpers: efficient bulk insert, update, and upsert with conflict handling.
//   - Migration primitives (Action, MigrationsConfig, ValidateMigrationAction) shared by
//     the backend-specific migration runners.
//   - Type utilities for nullable JSON, array types, and custom scanning.
//   - Common error types and helpers for consistent error handling.
//   - Query helpers for generic, type-safe data access patterns.
//
// Main Interfaces:
//   - Handler, Querier, Tx, Result, Rows, Row: Abstract over database drivers.
//   - Context: Extends context.Context with database query and transaction hooks.
//   - Executor, ExecutorRO, ExecutorRW: Transactional execution patterns.
//   - Bulkable, BulkableRanger: Bulk DML abstractions.
//   - Logger: the minimal logging contract this package needs from a caller-supplied logger.
//
// Backend Modules:
//   - github.com/sonirico/vago/db/postgres: Postgres (pgx and database/sql) support and migrations
//   - github.com/sonirico/vago/db/clickhouse: ClickHouse support, including its HTTP client and migrations
//   - github.com/sonirico/vago/db/mongo: MongoDB support
//   - github.com/sonirico/vago/db/redis: Redis support
//
// Utilities:
//   - utils_bulk_insert.go, utils_bulk_update.go, utils_bulk_upsert.go: Bulk DML
//   - utils_in_clause.go, utils_order_clause.go, utils_query.go: Query helpers
//   - types.go: NullJSON, NullJSONArray, and more
//   - errors.go: Common error values and helpers
//   - migrate.go: Migration primitives shared by the backend-specific migration runners
//
// Example:
//
//	import (
//	    "github.com/sonirico/vago/db"
//	    "github.com/sonirico/vago/db/postgres"
//	    "github.com/sonirico/vago/lol"
//	)
//
//	// Setup a logger and a database handler (e.g., pgx)
//	logger := lol.NewLogger()
//	handler, _ := postgres.OpenPgx(logger, "postgres://user:pass@localhost/db")
//	executor := db.NewExecutor(logger, handler)
//
//	// Run a transactional operation: either all operations succeed, or none are applied
//	err := executor.DoWithTx(ctx, func(ctx db.Context) error {
//	    // Multiple DB operations in a transaction
//	    if _, err := ctx.Querier().ExecContext(ctx, "INSERT INTO users (name) VALUES ($1)", "alice"); err != nil {
//	        return err
//	    }
//	    if _, err := ctx.Querier().ExecContext(ctx, "INSERT INTO accounts (user) VALUES ($1)", "alice"); err != nil {
//	        return err
//	    }
//	    // If any error is returned, all changes are rolled back
//	    return nil
//	})
package db
