package db

import (
	"database/sql"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/sonirico/vago/lol"
)

func OpenClickhouse(url string, log lol.Logger) (*sql.DB, error) {
	opts, err := clickhouse.ParseDSN(url)

	if err != nil {
		return nil, fmt.Errorf("unable to parse clickhouse DSN: %w", err)
	}

	db := clickhouse.OpenDB(opts)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping clickhouse database: %w", err)
	} else {
		log.Debug("clickhouse connected")
	}

	return db, nil
}

func OpenCH(url string, log lol.Logger) (Handler, error) {
	opts, err := clickhouse.ParseDSN(url)

	if err != nil {
		return nil, fmt.Errorf("unable to parse clickhouse DSN: %w", err)
	}

	db := clickhouse.OpenDB(opts)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping clickhouse database: %w", err)
	} else {
		log.Debug("clickhouse connected")
	}

	return newSqlAdapter(db), nil
}
