package zerolog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	rs "github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	lol "github.com/sonirico/vago/lol/v2"
	"github.com/sonirico/vago/lol/v2/zerolog"
)

type ctxKey struct{}

type stampHook struct {
	id string
}

func (h stampHook) Run(e *rs.Event, _ rs.Level, _ string) {
	e.Str("trace", h.id)
}

func TestOptions(t *testing.T) {
	t.Parallel()

	t.Run("WithTraceHook stamps what the hook reads from the context", func(t *testing.T) {
		t.Parallel()

		buf := new(bytes.Buffer)
		log := zerolog.New(
			lol.NewConfig(lol.WithEnv(lol.EnvProd), lol.WithWriter(buf)),
			zerolog.WithTraceHook(func(ctx context.Context) rs.Hook {
				return stampHook{id: ctx.Value(ctxKey{}).(string)}
			}),
		)
		ctx := context.WithValue(context.Background(), ctxKey{}, "t-1")

		log.WithTrace(ctx).Info().Msg("a")
		log.Info().Msg("b")

		lines := bytes.Split(bytes.TrimSpace(buf.Bytes()), []byte("\n"))
		require.Len(t, lines, 2)
		var traced, plain map[string]any
		require.NoError(t, json.Unmarshal(lines[0], &traced))
		require.NoError(t, json.Unmarshal(lines[1], &plain))
		require.Equal(t, "t-1", traced["trace"])
		require.NotContains(t, plain, "trace")
	})

	t.Run("WithExtraWriter tees every record", func(t *testing.T) {
		t.Parallel()

		primary, extra := new(bytes.Buffer), new(bytes.Buffer)
		log := zerolog.New(
			lol.NewConfig(lol.WithEnv(lol.EnvProd), lol.WithWriter(primary)),
			zerolog.WithExtraWriter(extra),
		)

		log.Info().Msg("a")

		require.Equal(t, primary.String(), extra.String())
		require.Contains(t, primary.String(), `"message":"a"`)
	})
}
