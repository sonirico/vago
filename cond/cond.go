package cond

// If returns the first value if the condition is true, otherwise returns the second value.
// This is a generic ternary operator.
func If[T any](condition bool, ifTrue, ifFalse T) T {
	if condition {
		return ifTrue
	}
	return ifFalse
}

// IfElse returns the first value if the condition is true, otherwise returns the second value.
// Alias for If for better readability in some contexts.
func IfElse[T any](condition bool, ifTrue, ifFalse T) T {
	return If(condition, ifTrue, ifFalse)
}

// IfFunc returns the result of calling the first function if the condition is true,
// otherwise calls the second function. This is lazy evaluation.
func IfFunc[T any](condition bool, ifTrue, ifFalse func() T) T {
	if condition {
		return ifTrue()
	}
	return ifFalse()
}

// IfPtr returns the first pointer if the condition is true, otherwise returns the second pointer.
// If the selected pointer is nil, returns nil.
func IfPtr[T any](condition bool, ifTrue, ifFalse *T) *T {
	if condition {
		return ifTrue
	}
	return ifFalse
}

// OrDefault returns the value if the condition is true, otherwise returns the zero value of T.
func OrDefault[T any](condition bool, value T) T {
	if condition {
		return value
	}
	var zero T
	return zero
}
