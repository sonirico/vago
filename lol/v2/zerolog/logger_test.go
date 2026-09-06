package zerolog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	lol "github.com/sonirico/vago/lol/v2"
	"github.com/sonirico/vago/lol/v2/zerolog"
)

type testFixture struct {
	buf *bytes.Buffer
	log lol.Logger
}

func newTestLogger(level lol.Level, opts ...lol.Opt) testFixture {
	buf := new(bytes.Buffer)
	cfg := lol.NewConfig(append([]lol.Opt{
		lol.WithLevel(level), lol.WithEnv(lol.EnvProd), lol.WithWriter(buf),
	}, opts...)...)
	return testFixture{buf: buf, log: zerolog.New(cfg)}
}

func (f testFixture) records(t *testing.T) []map[string]any {
	t.Helper()
	var out []map[string]any
	for _, line := range bytes.Split(bytes.TrimSpace(f.buf.Bytes()), []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		var rec map[string]any
		require.NoError(t, json.Unmarshal(line, &rec), "line %q", line)
		out = append(out, rec)
	}
	return out
}

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("every typed field reaches the record", func(t *testing.T) {
		t.Parallel()

		f := newTestLogger(lol.LevelInfo)
		at := time.Date(2026, 9, 6, 12, 0, 0, 0, time.UTC)

		f.log.Info().Str("s", "v").Int("i", 1).Int64("i64", 2).Uint64("u", 3).
			Float64("f", 4.5).Bool("b", true).Dur("d", time.Second).
			Time("t", at).Err(errors.New("boom")).Any("a", []int{1}).Msg("hello")

		recs := f.records(t)
		require.Len(t, recs, 1)
		rec := recs[0]
		require.Equal(t, "hello", rec["message"])
		require.Equal(t, "v", rec["s"])
		require.EqualValues(t, 1, rec["i"])
		require.EqualValues(t, 2, rec["i64"])
		require.EqualValues(t, 3, rec["u"])
		require.EqualValues(t, 4.5, rec["f"])
		require.Equal(t, true, rec["b"])
		require.NotNil(t, rec["d"])
		require.NotNil(t, rec["t"])
		require.Equal(t, "boom", rec["error"])
		require.Equal(t, []any{float64(1)}, rec["a"])
	})

	t.Run("Msgf formats", func(t *testing.T) {
		t.Parallel()

		f := newTestLogger(lol.LevelInfo)

		f.log.Warn().Msgf("%d of %d", 1, 2)

		recs := f.records(t)
		require.Len(t, recs, 1)
		require.Equal(t, "1 of 2", recs[0]["message"])
	})

	type levelCase struct {
		name    string
		min     lol.Level
		emit    func(lol.Logger) lol.Event
		written bool
	}
	levelCases := []levelCase{
		{"debug below info is dropped", lol.LevelInfo, lol.Logger.Debug, false},
		{"trace below debug is dropped", lol.LevelDebug, lol.Logger.Trace, false},
		{"info at info is written", lol.LevelInfo, lol.Logger.Info, true},
		{"error above warn is written", lol.LevelWarn, lol.Logger.Error, true},
		{"warn below error is dropped", lol.LevelError, lol.Logger.Warn, false},
		{"trace at trace is written", lol.LevelTrace, lol.Logger.Trace, true},
	}
	for _, tc := range levelCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := newTestLogger(tc.min)

			tc.emit(f.log).Msg("x")

			require.Equal(t, tc.written, f.buf.Len() > 0, f.buf.String())
		})
	}

	t.Run("config fields and With fields stamp every record", func(t *testing.T) {
		t.Parallel()

		f := newTestLogger(lol.LevelInfo, lol.WithFields(lol.Fields{"app": "x"}))

		f.log.With(lol.Fields{"req": "r1"}).Info().Msg("a")
		f.log.Info().Msg("b")

		recs := f.records(t)
		require.Len(t, recs, 2)
		require.Equal(t, "x", recs[0]["app"])
		require.Equal(t, "r1", recs[0]["req"])
		require.Equal(t, "x", recs[1]["app"])
		require.NotContains(t, recs[1], "req")
	})

	t.Run("WithTrace without a hook returns a working logger", func(t *testing.T) {
		t.Parallel()

		f := newTestLogger(lol.LevelInfo)

		f.log.WithTrace(context.Background()).Info().Msg("a")

		require.Len(t, f.records(t), 1)
	})

	t.Run("Panic writes then panics with the message", func(t *testing.T) {
		t.Parallel()

		f := newTestLogger(lol.LevelInfo)

		require.PanicsWithValue(t, "bye", func() { f.log.Panic().Msg("bye") })
		recs := f.records(t)
		require.Len(t, recs, 1)
		require.Equal(t, "bye", recs[0]["message"])
	})

	t.Run("EnvLocal writes for a terminal, not JSON", func(t *testing.T) {
		t.Parallel()

		buf := new(bytes.Buffer)
		log := zerolog.New(lol.NewConfig(lol.WithEnv(lol.EnvLocal), lol.WithWriter(buf)))

		log.Info().Str("k", "v").Msg("hi")

		require.Contains(t, buf.String(), "hi")
		require.False(t, json.Valid(bytes.TrimSpace(buf.Bytes())))
		require.NotContains(t, buf.String(), "\x1b[", "no colour into a writer that is not a terminal")
	})

	t.Run("TimeFormat shapes the timestamp", func(t *testing.T) {
		t.Parallel()

		f := newTestLogger(lol.LevelInfo, lol.WithTimeFormat("2006"))

		f.log.Info().Msg("a")

		recs := f.records(t)
		require.Len(t, recs, 1)
		require.Len(t, recs[0]["time"], 4, recs[0])
	})
}

// TestNew_Allocs stands alone and never runs in parallel: testing.AllocsPerRun
// refuses a parallel test, its parent included.
func TestNew_Allocs(t *testing.T) {
	f := newTestLogger(lol.LevelInfo)
	err := errors.New("boom")

	allocs := testing.AllocsPerRun(100, func() {
		f.log.Debug().Str("s", "v").Int("i", 1).Int64("i64", 2).Uint64("u", 3).
			Float64("f", 4.5).Bool("b", true).Dur("d", time.Second).
			Time("t", time.Time{}).Err(err).Any("a", nil).Msg("m")
	})

	require.Zero(t, allocs)
	require.Zero(t, f.buf.Len())
}
