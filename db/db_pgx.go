package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sonirico/vago/lol"
	apmpgx "go.elastic.co/apm/module/apmpgxv5/v2"
)

// PgxConfig holds all configuration options for the database connection.
type PgxConfig struct {
	ApmEnabled        bool
	MaxOpenConns      int           // Maximum open connections (pool max size)
	MinPoolSize       int           // Minimum pool size (pre-opened connections)
	ConnMaxIdleTime   time.Duration // Connection idle timeout
	ConnMaxLifetime   time.Duration // Maximum lifetime for a connection
	HealthCheckPeriod time.Duration // Frequency of health checks for idle connections
	AcquireTimeout    time.Duration // Timeout for acquiring a connection (not directly supported in pgxpool.Config)
}

// OpenPgxConn opens a pgx pool with optional APM instrumentation and applies all configuration options.
// Conn string params:
//   - pool_max_conns: integer greater than 0
//   - pool_min_conns: integer 0 or greater
//   - pool_max_conn_lifetime: duration string
//   - pool_max_conn_idle_time: duration string
//   - pool_health_check_period: duration string
//   - pool_max_conn_lifetime_jitter: duration string
func OpenPgxConn(log lol.Logger, uri string, apm bool) (*pgx.Conn, error) {
	ctx := context.Background()

	// Parse the URI into a pgx pool configuration.
	conn, err := pgx.Connect(ctx, uri)
	if err != nil {
		log.Error("Failed to parse pgx config", err)
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	if apm {
		apmpgx.Instrument(conn.Config())
		log.Debug("APM instrumentation enabled for pgx connections")
	}

	if err = conn.Ping(ctx); err != nil {
		log.Error("Failed to ping PostgreSQL", err)
		return nil, fmt.Errorf("failed to ping pgx: %w", err)
	}

	log.Info("Successfully connected to PostgreSQL")

	return conn, nil
}

// OpenPgxPool opens a pgx pool with optional APM instrumentation and applies all configuration options.
// Conn string params:
//   - pool_max_conns: integer greater than 0
//   - pool_min_conns: integer 0 or greater
//   - pool_max_conn_lifetime: duration string
//   - pool_max_conn_idle_time: duration string
//   - pool_health_check_period: duration string
//   - pool_max_conn_lifetime_jitter: duration string
func OpenPgxPool(log lol.Logger, uri string, apm bool) (*pgxpool.Pool, error) {
	ctx := context.Background()

	// Parse the URI into a pgx pool configuration.
	poolConfig, err := pgxpool.ParseConfig(uri)
	if err != nil {
		log.Error("Failed to parse pgx pool config", err)
		return nil, fmt.Errorf("failed to parse pgx pool config: %w", err)
	}

	if apm {
		apmpgx.Instrument(poolConfig.ConnConfig)
		log.Debug("APM instrumentation enabled for pgx connections")
	}

	beforeConnect := poolConfig.BeforeConnect
	poolConfig.BeforeConnect = func(ctx context.Context, config *pgx.ConnConfig) error {
		log.Debugf("before.connect called for %s@%s:%d/%s",
			config.Config.User, config.Config.Host, config.Config.Port, config.Config.Database)

		if beforeConnect != nil {
			return beforeConnect(ctx, config)
		}

		return nil
	}

	afterConnect := poolConfig.AfterConnect
	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		log.Debugf("after.connect called, pid=%d", conn.PgConn().PID())

		if afterConnect != nil {
			return afterConnect(ctx, conn)
		}

		return nil
	}

	beforeAcquire := poolConfig.BeforeAcquire
	poolConfig.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		log.Debugf("before.acquire called, pid=%d", conn.PgConn().PID())

		if beforeAcquire != nil {
			return beforeAcquire(ctx, conn)
		}

		return true
	}

	afterRelease := poolConfig.AfterRelease
	poolConfig.AfterRelease = func(conn *pgx.Conn) bool {
		log.Debugf("after.release called, pid=%d", conn.PgConn().PID())

		if afterRelease != nil {
			return afterRelease(conn)
		}

		return true
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Error("Failed to create pgx pool", err)
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		log.Error("failed to ping postgresql", err)
		return nil, fmt.Errorf("failed to ping pgx pool: %w", err)
	}

	log.Info("successfully connected to postgresql using pgx driver")

	return pool, nil
}

// OpenPgx is an alias for OpenPgxPool with apm instrumented
func OpenPgx(log lol.Logger, uri string) (Handler, error) {
	pool, err := OpenPgxPool(log, uri, true)
	if err != nil {
		return nil, fmt.Errorf("failed to open pgx pool: %w", err)
	}

	return newPgxAdapter(pool), nil
}
