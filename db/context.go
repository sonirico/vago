package db

import (
	"context"
	"database/sql"
)

type (
	Context interface {
		context.Context
		Querier() Querier
		AfterCommit(fn AfterCommitCallbackFn)
		AfterCommitDo(func(Context))
	}

	AfterCommitCallbackFn func()
)

type (
	RepoContext struct {
		context.Context
		querier          Querier
		afterCommitCbs   []AfterCommitCallbackFn
		afterCommitDoCbs []func(Context)
	}
)

func (r RepoContext) Querier() Querier {
	return r.querier
}

// AfterCommit adds a callback function to the queue that will be executed
// after a database transaction has been executed successfully.
//
// IMPORTANT: Only valid for DoWithTx().
func (r *RepoContext) AfterCommit(fn AfterCommitCallbackFn) {
	r.afterCommitCbs = append(r.afterCommitCbs, fn)
}

// AfterCommitDo adds a callback function to the queue that will be executed
// after a database transaction has been executed successfully.
//
// IMPORTANT: Only valid for DoWithTx().
func (r *RepoContext) AfterCommitDo(fn func(Context)) {
	r.afterCommitDoCbs = append(r.afterCommitDoCbs, fn)
}

func (r *RepoContext) ExecAfterCommit(ctx Context) {
	for i := range r.afterCommitCbs {
		r.afterCommitCbs[i]()
	}

	for _, fn := range r.afterCommitDoCbs {
		fn(ctx)
	}
}

func NewRepoContext(ctx context.Context, querier Querier) *RepoContext {
	return &RepoContext{Context: ctx, querier: querier}
}

func NewNoopRepoContext(ctx context.Context) *RepoContext {
	return &RepoContext{Context: ctx, querier: newSqlAdapter(&sql.DB{})}
}

func NewNoopRepoContextTx(ctx context.Context) *RepoContext {
	return &RepoContext{Context: ctx, querier: newSqlTxAdapter(&sql.Tx{})}
}
