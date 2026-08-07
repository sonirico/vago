package db

import (
	"testing"

	"github.com/sonirico/vago/lol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLaunchPostgresql_UnknownAction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		action string
	}{
		{"empty", ""},
		{"wrong case", "Up"},
		{"typo", "aplly"},
		{"unrelated word", "reset"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange: an unreachable URL proves the function never opens a
			// connection before rejecting the action - if it did, this test
			// would hang or fail with a connection error instead of the
			// validation error asserted below.
			cfg := MigrationsConfig{Url: "postgres://unreachable.invalid:1/db"}

			// Act
			err := LaunchPostgresql(cfg, tt.action, lol.ZeroTestLogger)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), "unknown migration action")
			assert.Contains(t, err.Error(), tt.action)
		})
	}
}
