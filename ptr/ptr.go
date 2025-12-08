// Package ptr provides pointer utility functions for Go.
package ptr

// Ptr returns a pointer to the given value.
// This is useful for getting pointers to literals or values that need to be passed as pointers.
func Ptr[T any](x T) *T {
	return &x
}
