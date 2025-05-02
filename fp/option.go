// Package fp provides functional programming primitives for Go.
// It implements monadic types like Option and Result for more expressive error handling.
package fp

type (
	// Option represents an optional value: either Some value or None.
	// It's a type that encapsulates an optional value, avoiding the need for nil checks.
	Option[T any] struct {
		value  T
		isSome bool
	}
)

// IsSome checks if the option is in the Some state.
// Returns true if the option contains a value, false otherwise.
func (o Option[T]) IsSome() bool {
	return o.isSome
}

// IsNone checks if the option is in the None state.
// Returns true if the option does not contain a value, false otherwise.
func (o Option[T]) IsNone() bool {
	return !o.isSome
}

// Unwrap extracts the contained value and a boolean indicating if it exists.
// Returns the value and true if the option is Some, a zero value and false if None.
func (o Option[T]) Unwrap() (T, bool) {
	return o.value, o.isSome
}

// UnwrapOr returns the contained value or a provided default.
// If the option is Some, returns the contained value, otherwise returns the provided default.
func (o Option[T]) UnwrapOr(value T) T {
	if o.isSome {
		return o.value
	}
	return value
}

// UnwrapOrElse returns the contained value or computes it from a closure.
// If the option is Some, returns the contained value, otherwise returns the result of calling fn.
func (o Option[T]) UnwrapOrElse(fn func() T) T {
	if o.isSome {
		return o.value
	}
	return fn()
}

// UnwrapOrDefault returns the contained value or the zero value.
// If the option is Some, returns the contained value, otherwise returns the zero value for type T.
func (o Option[T]) UnwrapOrDefault() T {
	if o.isSome {
		return o.value
	}
	var res T
	return res
}

// UnwrapUnsafe returns the contained value or panics.
// If the option is Some, returns the contained value, otherwise panics with a message.
// Use this method when you are certain the Option is in the Some state.
func (o Option[T]) UnwrapUnsafe() T {
	if !o.isSome {
		panic("option is none")
	}
	return o.value
}

// Or returns the option if it is Some, otherwise returns the provided option.
// This is useful for providing a fallback option when the current one might be None.
func (o Option[T]) Or(other Option[T]) Option[T] {
	if !o.isSome {
		return other
	}
	return o
}

// OrElse returns the option if it is Some, otherwise calls the function and returns its result.
// This is a lazy version of Or, as the fallback option is only computed if needed.
func (o Option[T]) OrElse(fn func() Option[T]) Option[T] {
	if !o.isSome {
		return fn()
	}
	return o
}

// Map transforms the contained value using the provided function if the option is Some.
// If the option is None, returns None without calling the function.
func (o Option[T]) Map(fn func(T) T) Option[T] {
	if o.isSome {
		return Some(fn(o.value))
	}
	return o
}

// MapOr transforms the contained value or returns a default.
// If the option is Some, applies the function to the contained value,
// otherwise returns the provided default value.
func (o Option[T]) MapOr(value T, fn func(T) T) T {
	if o.isSome {
		return fn(o.value)
	}
	return value
}

// MapOrElse transforms the contained value or computes a default.
// If the option is Some, applies handleSome to the contained value,
// otherwise calls handleNone to compute a default.
func (o Option[T]) MapOrElse(handleNone func() T, handleSome func(T) T) T {
	if o.isSome {
		return handleSome(o.value)
	}
	return handleNone()
}

// OkOr converts an Option to a Result.
// If the option is Some, returns Ok with the contained value,
// otherwise returns Err with the provided error.
func (o Option[T]) OkOr(err error) Result[T] {
	if o.isSome {
		return Ok(o.value)
	}
	return Err[T](err)
}

// OkOrElse converts an Option to a Result, with a computed error.
// If the option is Some, returns Ok with the contained value,
// otherwise calls the function to create an error and returns Err with that error.
func (o Option[T]) OkOrElse(fn func() error) Result[T] {
	if o.isSome {
		return Ok(o.value)
	}
	return Err[T](fn())
}

// Match allows pattern matching on the Option.
// If the option is Some, calls handleSome with the contained value,
// otherwise calls handleNone.
func (o Option[T]) Match(handleSome func(T) Option[T], handleNone func() Option[T]) Option[T] {
	if o.isSome {
		return handleSome(o.value)
	}

	return handleNone()
}

// Some creates a new Option in the Some state with the given value.
// This is a constructor function for creating an Option that contains a value.
func Some[T any](t T) Option[T] {
	return Option[T]{value: t, isSome: true}
}

// None creates a new Option in the None state.
// This is a constructor function for creating an Option that does not contain a value.
func None[T any]() Option[T] {
	return Option[T]{}
}

// OptionFromTuple creates an Option from a tuple-like return (value, ok).
// If ok is true, returns Some(x), otherwise returns None.
// This is useful for converting Go's common (value, ok) pattern to an Option.
func OptionFromTuple[T any](x T, ok bool) Option[T] {
	if ok {
		return Some(x)
	}
	return None[T]()
}

// OptionFromPtr creates an Option from a pointer.
// If the pointer is nil, returns None, otherwise returns Some with the dereferenced value.
// This is useful for converting nullable pointers to the Option type.
func OptionFromPtr[T any](x *T) Option[T] {
	if x == nil {
		return None[T]()
	}
	return Some(*x)
}

// OptionFromZero creates an Option from a value, treating zero values as None.
// If the value equals the zero value for its type, returns None, otherwise returns Some(x).
// This is useful when zero values are treated as invalid or unset.
func OptionFromZero[T comparable](x T) Option[T] {
	var zero T
	if x == zero {
		return None[T]()
	}
	return Some(x)
}
