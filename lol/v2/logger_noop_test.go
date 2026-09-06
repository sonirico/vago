package lol_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	lol "github.com/sonirico/vago/lol/v2"
)

func TestDiscard(t *testing.T) {
	t.Parallel()

	t.Run("Fatal and Panic neither exit nor panic", func(t *testing.T) {
		t.Parallel()

		require.NotPanics(t, func() {
			lol.Discard.Panic().Msg("no")
			lol.Discard.Fatal().Msg("no")
		})
	})
}

// TestDiscard_Allocs stands alone and never runs in parallel: testing.AllocsPerRun
// refuses a parallel test, its parent included.
func TestDiscard_Allocs(t *testing.T) {
	log := lol.Discard.With(lol.Fields{"k": "v"}).WithTrace(context.Background())
	err := errors.New("boom")

	allocs := testing.AllocsPerRun(100, func() {
		for _, ev := range []lol.Event{
			log.Trace(), log.Debug(), log.Info(), log.Warn(),
			log.Error(), log.Fatal(), log.Panic(),
		} {
			ev.Str("s", "v").Int("i", 1).Int64("i64", 2).Uint64("u", 3).
				Float64("f", 4.5).Bool("b", true).Dur("d", time.Second).
				Time("t", time.Time{}).Err(err).Any("a", nil).Msg("m")
		}
	})

	require.Zero(t, allocs)
}
