package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sonirico/vago/lol"
)

type (
	executorFactory func(log lol.Logger, handler Handler) Executor

	dbExecutorRWO struct {
		ro     Handler
		rw     Handler
		logger lol.Logger
	}

	ExecutorRWO interface {
		RO() ExecutorRO
		RW() ExecutorRW

		DoRead(ctx context.Context, fn func(ctx Context) error) error
		DoWrite(ctx context.Context, fn func(ctx Context) error) error
		DoTx(ctx context.Context, fn func(ctx Context) error) error
	}
)

func (ds *dbExecutorRWO) RO() ExecutorRO {
	return newExecutor(ds.logger.WithField("mode", "ro"), ds.ro)
}

func (ds *dbExecutorRWO) RW() ExecutorRW {
	return newExecutor(ds.logger.WithField("mode", "rw"), ds.ro)
}

func (ds *dbExecutorRWO) DoRead(ctx context.Context, fn func(ctx Context) error) error {
	return newExecutor(ds.logger.WithField("mode", "ro"), ds.ro).Do(ctx, fn)
}

func (ds *dbExecutorRWO) DoWrite(ctx context.Context, fn func(ctx Context) error) error {
	return newExecutor(ds.logger.WithField("mode", "rw"), ds.rw).Do(ctx, fn)
}

func (ds *dbExecutorRWO) DoTx(ctx context.Context, fn func(ctx Context) error) error {
	return newExecutor(ds.logger.WithField("mode", "rw"), ds.rw).DoWithTx(ctx, fn)
}

func NewRWOExecutor(log lol.Logger, ro, rw Handler) ExecutorRWO {
	return &dbExecutorRWO{
		ro:     ro,
		rw:     rw,
		logger: log,
	}
}

func NewRWOExecutorPgx(log lol.Logger, ro, rw *pgxpool.Pool) ExecutorRWO {
	return &dbExecutorRWO{
		ro:     newPgxAdapter(ro),
		rw:     newPgxAdapter(rw),
		logger: log,
	}
}
