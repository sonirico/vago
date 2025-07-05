package zero

import (
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
