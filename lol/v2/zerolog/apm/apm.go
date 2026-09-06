// Package apm wires Elastic APM into the lol/v2/zerolog backend: error
// records are shipped to APM, error stacks are marshalled its way, and
// Logger.WithTrace stamps the transaction and span ids found in a context.
package apm

import (
	rs "github.com/rs/zerolog"
	"go.elastic.co/apm/module/apmzerolog/v2"

	"github.com/sonirico/vago/lol/v2/zerolog"
)

// Options returns what zerolog.New needs to report to APM. It also sets
// zerolog's global error stack marshaller, as apmzerolog requires.
func Options() []zerolog.Option {
	rs.ErrorStackMarshaler = apmzerolog.MarshalErrorStack
	return []zerolog.Option{
		zerolog.WithExtraWriter(new(apmzerolog.Writer)),
		zerolog.WithTraceHook(apmzerolog.TraceContextHook),
	}
}
