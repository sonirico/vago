package zero

// Package zero provides utility functions for zero memory allocation conversions

import (
	"unsafe"
)

// S2B converts string to []byte without memory allocation.
// This is an unsafe operation and should be used with caution.
func S2B(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// B2S converts []byte to string without memory allocation.
// This is an unsafe operation and should be used with caution.
func B2S(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
