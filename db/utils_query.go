package db

import (
	"context"
	"errors"
)

// Query executes a function within the context of an Executor and returns its result and error.
//
// If fn returns an error, that error is returned. Otherwise, the error (if any) reported by the
// executor itself is returned.
func Query[T any](ctx context.Context, executor Executor, fn func(Context) (T, error)) (T, error) {
	var (
		data     T
		errQuery error
	)
	errDo := executor.Do(ctx, func(tx Context) error {
		data, errQuery = fn(tx)
		return errQuery
	})

	if errQuery != nil {
		return data, errQuery
	}

	return data, errDo
}

// QueryRO executes a function within the context of an ExecutorRO (read-only) and returns its result and error.
//
// If fn returns an error, that error is returned. Otherwise, the error (if any) reported by the
// executor itself is returned.
func QueryRO[T any](
	ctx context.Context,
	executor ExecutorRO,
	fn func(Context) (T, error),
) (T, error) {
	var (
		data     T
		errQuery error
	)
	errDo := executor.Do(ctx, func(tx Context) error {
		data, errQuery = fn(tx)
		return errQuery
	})

	if errQuery != nil {
		return data, errQuery
	}

	return data, errDo
}

// QueryRW executes a function within the context of an ExecutorRW (read-write) and returns its result and error.
//
// If fn returns an error, that error is returned. Otherwise, the error (if any) reported by the
// executor itself is returned.
func QueryRW[T any](
	ctx context.Context,
	executor ExecutorRW,
	fn func(Context) (T, error),
) (T, error) {
	var (
		data     T
		errQuery error
	)
	errDo := executor.Do(ctx, func(tx Context) error {
		data, errQuery = fn(tx)
		return errQuery
	})

	if errQuery != nil {
		return data, errQuery
	}

	return data, errDo
}

// QueryTx executes a function within the context of an Executor using a transaction and returns its result and error.
//
// fn's error is propagated to the transaction machinery, so a failing fn always rolls back instead
// of committing. When fn fails and the rollback itself also fails, both errors are joined with
// errors.Join so neither is silently lost. When fn succeeds but the commit fails, the commit error
// is returned even though fn itself reported no error.
func QueryTx[T any](
	ctx context.Context,
	executor Executor,
	fn func(Context) (T, error),
) (T, error) {
	var (
		data     T
		errQuery error
	)
	errExec := executor.DoWithTx(ctx, func(tx Context) error {
		data, errQuery = fn(tx)
		return errQuery
	})

	if errQuery == nil {
		return data, errExec
	}

	if errExec != nil && !errors.Is(errExec, errQuery) {
		return data, errors.Join(errQuery, errExec)
	}

	return data, errQuery
}
