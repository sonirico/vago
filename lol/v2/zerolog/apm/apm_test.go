package apm_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	apmv2 "go.elastic.co/apm/v2"
	"go.elastic.co/apm/v2/apmtest"

	lol "github.com/sonirico/vago/lol/v2"
	"github.com/sonirico/vago/lol/v2/zerolog"
	"github.com/sonirico/vago/lol/v2/zerolog/apm"
)

func TestOptions(t *testing.T) {
	buf := new(bytes.Buffer)
	log := zerolog.New(
		lol.NewConfig(lol.WithEnv(lol.EnvProd), lol.WithWriter(buf)),
		apm.Options()...)
	tracer := apmtest.NewRecordingTracer()
	t.Cleanup(tracer.Close)
	tx := tracer.StartTransaction("tx", "test")
	ctx := apmv2.ContextWithTransaction(context.Background(), tx)

	log.WithTrace(ctx).Info().Msg("a")
	tx.End()

	var rec map[string]any
	require.NoError(t, json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &rec))
	require.Equal(t, tx.TraceContext().Trace.String(), rec["trace.id"])
	require.Equal(t, tx.TraceContext().Span.String(), rec["transaction.id"])
}
