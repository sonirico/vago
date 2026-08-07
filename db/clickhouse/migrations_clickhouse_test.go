package clickhouse

import (
	"testing"

	"github.com/sonirico/vago/lol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sonirico/vago/db"
)

func TestLaunchClickhouse_UnknownAction(t *testing.T) {
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
			// connection before rejecting the action.
			cfg := db.MigrationsConfig{Url: "tcp://unreachable.invalid:1"}

			// Act
			err := LaunchClickhouse(cfg, tt.action, lol.ZeroTestLogger)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), "unknown migration action")
			assert.Contains(t, err.Error(), tt.action)
		})
	}
}
