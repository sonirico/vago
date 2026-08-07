package db

import (
	"context"
)

type (
	executorFactory func(log Logger, handler Handler) Executor

	dbExecutorRWO struct {
		ro     Handler
		rw     Handler
		logger Logger
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
	return newExecutor(ds.logger, ds.ro)
}

func (ds *dbExecutorRWO) RW() ExecutorRW {
	return newExecutor(ds.logger, ds.rw)
}

func (ds *dbExecutorRWO) DoRead(ctx context.Context, fn func(ctx Context) error) error {
	return newExecutor(ds.logger, ds.ro).Do(ctx, fn)
}

func (ds *dbExecutorRWO) DoWrite(ctx context.Context, fn func(ctx Context) error) error {
	return newExecutor(ds.logger, ds.rw).Do(ctx, fn)
}

func (ds *dbExecutorRWO) DoTx(ctx context.Context, fn func(ctx Context) error) error {
	return newExecutor(ds.logger, ds.rw).DoWithTx(ctx, fn)
}

func NewRWOExecutor(log Logger, ro, rw Handler) ExecutorRWO {
	return &dbExecutorRWO{
		ro:     ro,
		rw:     rw,
		logger: log,
	}
}
