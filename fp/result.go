// Package fp provides functional programming primitives for Go.
package fp

type (
	// Result represents a computation that may succeed with type T or fail with an error.
	// It is similar to Option but includes error information when the value is absent.
	Result[T any] struct {
		value T
		err   error
	}
)

var (
	// OkAny is a predefined Result with a nil error and a zero any value.
	// Useful as a placeholder when only the success/failure state matters.
	OkAny = Result[any]{}
)

// IsOk checks if the result represents a successful computation.
// Returns true if there is no error, false otherwise.
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// IsErr checks if the result represents a failed computation.
// Returns true if there is an error, false otherwise.
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// UnwrapUnsafe returns the contained value or panics if there is an error.
// Use this method when you are certain the Result is in the Ok state.
func (r Result[T]) UnwrapUnsafe() T {
	if r.err != nil {
		panic("result is error: " + r.err.Error())
	}

	return r.value
}

// Unwrap extracts the contained value and any error.
// Returns the value (which may be the zero value) and the error (which may be nil).
func (r Result[T]) Unwrap() (T, error) {
	return r.value, r.err
}

// UnwrapOr returns the contained value or a provided default.
// If the result is Ok, returns the contained value, otherwise returns the provided default.
func (r Result[T]) UnwrapOr(other T) T {
	if r.err == nil {
		return r.value
	}

	return other
}

// UnwrapOrElse returns the contained value or computes it from a closure.
// If the result is Ok, returns the contained value, otherwise returns the result of calling fn.
func (r Result[T]) UnwrapOrElse(fn func() T) T {
	if r.err == nil {
		return r.value
	}

	return fn()
}

// UnwrapOrDefault returns the contained value or the zero value.
// If the result is Ok, returns the contained value, otherwise returns the zero value for type T.
func (r Result[T]) UnwrapOrDefault() T {
	if r.err == nil {
		return r.value
	}

	var res T
	return res
}

// Or returns the result if it is Ok, otherwise returns the provided result.
// This is useful for providing a fallback result when the current one might be an error.
func (r Result[T]) Or(other Result[T]) Result[T] {
	if r.err == nil {
		return r
	}

	return other
}

// OrElse returns the result if it is Ok, otherwise calls the function and returns its result.
// This is a lazy version of Or, as the fallback result is only computed if needed.
func (r Result[T]) OrElse(fn func() Result[T]) Result[T] {
	if r.err == nil {
		return r
	}

	return fn()
}

// Match allows pattern matching on the Result.
// If the result is Ok, calls handleOk with the contained value,
// otherwise calls handleErr with the error.
func (r Result[T]) Match(
	handleOk func(T) Result[T],
	handleErr func(error) Result[T],
) Result[T] {
	if r.err == nil {
		return handleOk(r.value)
	}

	return handleErr(r.err)
}

// And returns other if the result is Ok, otherwise returns the error result.
// This is useful for chaining results where both must succeed.
func (r Result[T]) And(other Result[T]) Result[T] {
	if r.err == nil {
		return other
	}

	return r
}

// AndThen returns a new result with the value computed from fn if the result is Ok,
// otherwise returns the error result without calling fn.
// This is useful for chaining operations that might fail.
func (r Result[T]) AndThen(fn func() T) Result[T] {
	if r.err == nil {
		return Ok(fn())
	}

	return r
}

// Map transforms the contained value using the provided function if the result is Ok.
// If the result is an error, returns the error result without calling the function.
func (r Result[T]) Map(fn func(T) T) Result[T] {
	if r.err == nil {
		return Ok(fn(r.value))
	}

	return r
}

// MapOr transforms the contained value or returns a default.
// If the result is Ok, applies the function to the contained value,
// otherwise returns an Ok result with the provided default value.
func (r Result[T]) MapOr(value T, fn func(T) T) Result[T] {
	if r.err == nil {
		return Ok(fn(r.value))
	}

	return Ok(value)
}

// MapOrElse transforms the contained value or handles the error.
// If the result is Ok, applies handleOk to the contained value,
// otherwise calls handleErr with the error.
// In both cases, returns an Ok result with the computed value.
func (r Result[T]) MapOrElse(
	handleErr func(error) T,
	handleOk func(T) T,
) Result[T] {
	if r.err == nil {
		return Ok(handleOk(r.value))
	}

	return Ok(handleErr(r.err))
}

// Ok creates a new Result in the Ok state with the given value.
// This is a constructor function for creating a Result that represents success.
func Ok[T any](v T) Result[T] {
	return Result[T]{value: v, err: nil}
}

// OkZero creates a new Result in the Ok state with the zero value.
// This is a constructor function for creating a Result that represents success
// but doesn't carry a meaningful value.
func OkZero[T any]() Result[T] {
	return Result[T]{err: nil}
}

// Err creates a new Result in the error state with the given error.
// This is a constructor function for creating a Result that represents failure.
func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}
