// Package db provides a unified set of abstractions, interfaces, and utilities for database access,
// transaction management, migrations, and efficient bulk operations across multiple backends.
//
// Features:
//   - Backend-agnostic interfaces for SQL (Postgres, ClickHouse), MongoDB, and Redis.
//   - Context and transaction management with hooks for after-commit actions.
//   - Executor interfaces for read-only, read-write, and transactional operations.
//   - Bulk DML helpers: efficient bulk insert, update, and upsert with conflict handling.
//   - Migration helpers for Postgres and ClickHouse using golang-migrate.
//   - Type utilities for nullable JSON, array types, and custom scanning.
//   - Common error types and helpers for consistent error handling.
//   - Query helpers for generic, type-safe data access patterns.
//
// Main Interfaces:
//   - Handler, Querier, Tx, Result, Rows, Row: Abstract over database drivers.
//   - Context: Extends context.Context with database query and transaction hooks.
//   - Executor, ExecutorRO, ExecutorRW: Transactional execution patterns.
//   - Bulkable, BulkableRanger: Bulk DML abstractions.
//
// Backend Adapters:
//   - db_pgx.go: Postgres (pgx) support
//   - db_clickhouse.go: ClickHouse support
//   - db_mongo.go: MongoDB support
//   - db_redis.go: Redis support
//
// Utilities:
//   - utils_bulk_insert.go, utils_bulk_update.go, utils_bulk_upsert.go: Bulk DML
//   - utils_in_clause.go, utils_order_clause.go, utils_query.go: Query helpers
//   - types.go: NullJSON, NullJSONArray, and more
//   - errors.go: Common error values and helpers
//   - migrate.go, migrations_postgres.go, migrations_clickhouse.go: Migration helpers
//
// Example:
//
//	import (
//	    "github.com/sonirico/vago/db"
//	    "github.com/sonirico/vago/lol"
//	)
//
//	// Setup a logger and a database handler (e.g., pgx)
//	logger := lol.NewLogger()
//	handler, _ := db.OpenPgxConn(logger, "postgres://user:pass@localhost/db", false)
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
