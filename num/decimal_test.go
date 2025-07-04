package num

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestMustDecFromAny(t *testing.T) {
	testCases := []struct {
		name        string
		input       any
		expect      Dec
		shouldPanic bool
	}{
		{
			name:        "Input as int",
			input:       123,
			expect:      NewDecFromInt(123),
			shouldPanic: false,
		},
		{
			name:        "Input as int64",
			input:       int64(12345678901234),
			expect:      NewDecFromInt(int64(12345678901234)),
			shouldPanic: false,
		},
		{
			name:        "Input as float64",
			input:       float64(123.456),
			expect:      NewDecFromFloat(float64(123.456)),
			shouldPanic: false,
		},
		{
			name:        "Input as string",
			input:       "123.456",
			expect:      MustDecFromString("123.456"),
			shouldPanic: false,
		},
		{
			name:        "Input as invalid type",
			input:       []int{1, 2, 3},
			expect:      Zero,
			shouldPanic: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var actual Dec
			shouldPanic := false

			// Check if function panics
			func() {
				defer func() {
					if r := recover(); r != nil {
						shouldPanic = true
					}
				}()
				actual = MustDecFromAny(tc.input)
			}()

			assert.Equal(t, tc.shouldPanic, shouldPanic)

			if !shouldPanic {
				assert.Equal(t, tc.expect, actual)
			}
		})
	}
}

func TestDec_NumberOfDecimals(t *testing.T) {
	type fields struct {
		dec   decimal.Decimal
		isset bool
	}
	tests := []struct {
		name string
		dec  Dec
		want int32
	}{
		{
			name: "ok",
			dec:  MustDecFromString("123.456"),
			want: 3,
		},
		{
			name: "with trailing zeros",
			dec:  MustDecFromString("123.456000"),
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.dec.NumberOfDecimals(); got != tt.want {
				t.Errorf("Dec.NumberOfDecimals() = %v, want %v", got, tt.want)
			}
		})
	}
}
