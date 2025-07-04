package db

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// PgPoolParams holds the extracted connection pool parameters.
type PgPoolParams struct {
	MaxOpenConns    int
	MinConns        int
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
	CleanURI        string
}

// parsePoolingParams extracts pooling parameters from the provided URI and returns the cleaned URI (with pooling parameters removed)
// along with the parsed values.
func parsePoolingParams(uri string) (PgPoolParams, error) {
	var params PgPoolParams

	u, err := url.Parse(uri)
	if err != nil {
		return params, fmt.Errorf("invalid URI: %w", err)
	}

	q := u.Query()

	// pool_max_conns: maximum number of open connections.
	if v := q.Get("pool_max_conns"); v != "" {
		params.MaxOpenConns, err = strconv.Atoi(v)
		if err != nil {
			return params, fmt.Errorf("invalid pool_max_conns: %w", err)
		}
		q.Del("pool_max_conns")
	}

	// pool_min_conns: maximum number of idle connections.
	if v := q.Get("pool_min_conns"); v != "" {
		params.MinConns, err = strconv.Atoi(v)
		if err != nil {
			return params, fmt.Errorf("invalid pool_min_conns: %w", err)
		}
		q.Del("pool_min_conns")
	}

	// pool_max_conn_lifetime: maximum lifetime of a connection.
	if v := q.Get("pool_max_conn_lifetime"); v != "" {
		params.MaxConnLifetime, err = time.ParseDuration(v)
		if err != nil {
			return params, fmt.Errorf("invalid pool_max_conn_lifetime: %w", err)
		}
		q.Del("pool_max_conn_lifetime")
	}

	// pool_max_conn_idle_time: maximum idle time for a connection.
	if v := q.Get("pool_max_conn_idle_time"); v != "" {
		params.MaxConnIdleTime, err = time.ParseDuration(v)
		if err != nil {
			return params, fmt.Errorf("invalid pool_max_conn_idle_time: %w", err)
		}
		q.Del("pool_max_conn_idle_time")
	}

	// Remove any other unsupported pooling parameters.
	q.Del("pool_health_check_period")
	q.Del("pool_max_conn_lifetime_jitter")

	// Rebuild the URI without the pooling parameters.
	u.RawQuery = q.Encode()
	params.CleanURI = u.String()

	return params, nil
}
