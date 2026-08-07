package db

import (
	"context"
	"database/sql"

	"errors"
)

type (
	dbExecutor struct {
		db     Handler
		logger Logger
	}

	ExecutorRO interface {
		Do(ctx context.Context, fn func(ctx Context) error) error
	}

	ExecutorRW interface {
		ExecutorRO
		DoWithTx(ctx context.Context, fn func(ctx Context) error) error
	}

	Executor interface {
		ExecutorRW
	}
)

// NewExecutor creates a new Executor using the provided logger and Handler.
func NewExecutor(log Logger, db Handler) Executor {
	return newExecutor(log, db)
}

// NewDatabaseSqlExecutor creates a new Executor for a *sql.DB database.
func NewDatabaseSqlExecutor(log Logger, db *sql.DB) Executor {
	return newDatabaseSqlExecutor(log, db)
}

func newDatabaseSqlExecutor(log Logger, db *sql.DB) Executor {
	return newExecutor(log, NewDatabaseSqlHandler(db))
}

func newExecutor(log Logger, db Handler) Executor {
	return &dbExecutor{
		db:     db,
		logger: log,
	}
}

func (ds dbExecutor) do(fn func(querier Querier) error) error {
	return fn(ds.db)
}

func (ds dbExecutor) doWithTx(
	ctx context.Context,
	fn func(querier Querier) error,
) error {
	tx, err := ds.db.Begin(ctx)

	if err != nil {
		return err
	}

	if err = fn(tx); err != nil {
		if errors.Is(err, context.Canceled) {
			ctx = context.Background()
		}

		if rerr := tx.Rollback(ctx); rerr != nil {
			ds.logger.Errorf("unable to rollback tx due to '%s', previous error was %s",
				rerr, err)
			return rerr
		}
		return err
	}

	return tx.Commit(ctx)
}

// Do executes a db query
func (ds dbExecutor) Do(ctx context.Context, fn func(ctx Context) error) error {
	return ds.do(func(querier Querier) error {
		repoCtx := NewRepoContext(ctx, querier)
		if err := fn(repoCtx); err != nil {
			return err
		}

		repoCtx.ExecAfterCommit(repoCtx)

		return nil
	})
}

// DoWithTx Executes a db query using a transaction
func (ds dbExecutor) DoWithTx(ctx context.Context, fn func(ctx Context) error) error {
	var repoCtx *RepoContext

	if err := ds.doWithTx(ctx, func(querier Querier) error {
		repoCtx = NewRepoContext(ctx, querier)

		return fn(repoCtx)
	}); err != nil {
		return err
	}

	return ds.do(func(querier Querier) error {
		repoCtx.ExecAfterCommit(NewRepoContext(ctx, querier))
		return nil
	})
}
