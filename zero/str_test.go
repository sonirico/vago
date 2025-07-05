package zero

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestS2B(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  string
		expect []byte
	}{
		{"ascii", "hello", []byte{'h', 'e', 'l', 'l', 'o'}},
		{"empty", "", []byte{}},
		{"utf8", "¡Hola!", []byte{0xc2, 0xa1, 'H', 'o', 'l', 'a', '!'}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := S2B(tt.input)
			if tt.input == "" {
				assert.Equal(t, 0, len(b))
			} else {
				assert.Equal(t, tt.expect, b)
				if len(b) > 0 {
					assert.Equal(t,
						uintptr(unsafe.Pointer(unsafe.StringData(tt.input))),
						uintptr(unsafe.Pointer(&b[0])),
					)
				}
			}
		})
	}
}

func TestB2S(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  []byte
		expect string
	}{
		{"ascii", []byte{'h', 'e', 'l', 'l', 'o'}, "hello"},
		{"empty", []byte{}, ""},
		{"utf8", []byte{0xc2, 0xa1, 'H', 'o', 'l', 'a', '!'}, "¡Hola!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := B2S(tt.input)
			assert.Equal(t, tt.expect, s)
			if len(tt.input) > 0 {
				assert.Equal(t,
					uintptr(unsafe.Pointer(&tt.input[0])),
					uintptr(unsafe.Pointer(unsafe.StringData(s))),
				)
			}
		})
	}
}

// ExampleS2B demonstrates converting a string to []byte without memory allocation
func ExampleS2B() {
	// Convert string to []byte using zero-allocation conversion
	s := "Hello, World!"
	b := S2B(s)

	fmt.Printf("Original string: %s\n", s)
	fmt.Printf("Converted to bytes: %v\n", b)
	fmt.Printf("Bytes as string: %s\n", string(b))
	fmt.Printf(
		"Same underlying data: %t\n",
		uintptr(unsafe.Pointer(unsafe.StringData(s))) == uintptr(unsafe.Pointer(&b[0])),
	)

	// Output:
	// Original string: Hello, World!
	// Converted to bytes: [72 101 108 108 111 44 32 87 111 114 108 100 33]
	// Bytes as string: Hello, World!
	// Same underlying data: true
}

// ExampleB2S demonstrates converting []byte to string without memory allocation
func ExampleB2S() {
	// Convert []byte to string using zero-allocation conversion
	b := []byte("Hello, Gophers!")
	s := B2S(b)

	fmt.Printf("Original bytes: %v\n", b)
	fmt.Printf("Converted to string: %s\n", s)
	fmt.Printf("Length: %d\n", len(s))
	fmt.Printf(
		"Same underlying data: %t\n",
		uintptr(unsafe.Pointer(&b[0])) == uintptr(unsafe.Pointer(unsafe.StringData(s))),
	)

	// Output:
	// Original bytes: [72 101 108 108 111 44 32 71 111 112 104 101 114 115 33]
	// Converted to string: Hello, Gophers!
	// Length: 15
	// Same underlying data: true
}
