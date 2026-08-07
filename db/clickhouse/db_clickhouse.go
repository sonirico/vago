package clickhouse

import (
	"database/sql"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/sonirico/vago/lol"

	"github.com/sonirico/vago/db"
)

func OpenClickhouse(url string, log lol.Logger) (*sql.DB, error) {
	opts, err := clickhouse.ParseDSN(url)

	if err != nil {
		return nil, fmt.Errorf("unable to parse clickhouse DSN: %w", err)
	}

	conn := clickhouse.OpenDB(opts)

	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping clickhouse database: %w", err)
	} else {
		log.Debug("clickhouse connected")
	}

	return conn, nil
}

func OpenCH(url string, log lol.Logger) (db.Handler, error) {
	opts, err := clickhouse.ParseDSN(url)

	if err != nil {
		return nil, fmt.Errorf("unable to parse clickhouse DSN: %w", err)
	}

	conn := clickhouse.OpenDB(opts)

	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping clickhouse database: %w", err)
	} else {
		log.Debug("clickhouse connected")
	}

	return db.NewDatabaseSqlHandler(conn), nil
}
