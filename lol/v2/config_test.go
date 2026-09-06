package lol_test

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	lol "github.com/sonirico/vago/lol/v2"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	t.Run("defaults are info, local, discard", func(t *testing.T) {
		t.Parallel()

		cfg := lol.NewConfig()

		require.Equal(
			t,
			lol.Config{Level: lol.LevelInfo, Env: lol.EnvLocal, Writer: io.Discard},
			cfg,
		)
	})

	t.Run("every option lands on its field", func(t *testing.T) {
		t.Parallel()

		fields := lol.Fields{"app": "x"}
		cfg := lol.NewConfig(
			lol.WithLevel(lol.LevelTrace),
			lol.WithEnv(lol.EnvProd),
			lol.WithWriter(os.Stdout),
			lol.WithFields(fields),
			lol.WithTimeFormat("15:04"),
		)

		require.Equal(t, lol.Config{
			Level: lol.LevelTrace, Env: lol.EnvProd, Writer: os.Stdout,
			Fields: fields, TimeFormat: "15:04",
		}, cfg)
	})
}
