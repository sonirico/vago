package db

import (
	"context"

	"github.com/sonirico/vago/lol"
)

type (
	NoopExecutor struct {
		logger lol.Logger
	}
)

func (ds NoopExecutor) Do(ctx context.Context, fn func(ctx Context) error) error {
	return ds.do(func(querier Querier) error {
		return fn(NewRepoContext(ctx, querier))
	})
}

func (ds NoopExecutor) RO() ExecutorRO {
	return ds
}

func (ds NoopExecutor) RW() ExecutorRW {
	return ds
}

func (ds NoopExecutor) DoRead(ctx context.Context, fn func(ctx Context) error) error {
	return ds.Do(ctx, fn)
}

func (ds NoopExecutor) DoWrite(ctx context.Context, fn func(ctx Context) error) error {
	return ds.Do(ctx, fn)
}

func (ds NoopExecutor) DoTx(ctx context.Context, fn func(ctx Context) error) error {
	return ds.DoWithTx(ctx, fn)
}

func (ds NoopExecutor) DoWithTx(
	ctx context.Context,
	fn func(ctx Context) error,
) error {
	var dbCtx *RepoContext

	if err := ds.doWithTx(func(querier Querier) error {
		dbCtx = NewRepoContext(ctx, querier)
		return fn(dbCtx)
	}); err != nil {
		return err
	}

	return ds.Do(ctx, func(repoCtx Context) error {
		dbCtx.ExecAfterCommit(repoCtx)
		return nil
	})
}

func NewNoopExecutor() NoopExecutor {
	return NoopExecutor{}
}

func (ds NoopExecutor) do(fn func(querier Querier) error) error {
	return fn(nil)
}

func (ds NoopExecutor) doWithTx(fn func(querier Querier) error) error {
	return ds.do(fn)
}
