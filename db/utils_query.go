package db

import "context"

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
