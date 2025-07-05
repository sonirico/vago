package db

import "context"

// Query executes a function within the context of an Executor and returns its result and error.
func Query[T any](ctx context.Context, executor Executor, fn func(Context) (T, error)) (T, error) {
	var (
		data     T
		errQuery error
	)
	_ = executor.Do(ctx, func(tx Context) error {
		data, errQuery = fn(tx)
		return nil
	})

	return data, errQuery
}

// QueryRO executes a function within the context of an ExecutorRO (read-only) and returns its result and error.
func QueryRO[T any](
	ctx context.Context,
	executor ExecutorRO,
	fn func(Context) (T, error),
) (T, error) {
	var (
		data     T
		errQuery error
	)
	_ = executor.Do(ctx, func(tx Context) error {
		data, errQuery = fn(tx)
		return nil
	})

	return data, errQuery
}

// QueryRW executes a function within the context of an ExecutorRW (read-write) and returns its result and error.
func QueryRW[T any](
	ctx context.Context,
	executor ExecutorRW,
	fn func(Context) (T, error),
) (T, error) {
	var (
		data     T
		errQuery error
	)
	_ = executor.Do(ctx, func(tx Context) error {
		data, errQuery = fn(tx)
		return nil
	})

	return data, errQuery
}

// QueryTx executes a function within the context of an Executor using a transaction and returns its result and error.
func QueryTx[T any](
	ctx context.Context,
	executor Executor,
	fn func(Context) (T, error),
) (T, error) {
	var (
		data     T
		errQuery error
	)
	_ = executor.DoWithTx(ctx, func(tx Context) error {
		data, errQuery = fn(tx)
		return nil
	})

	return data, errQuery
}
