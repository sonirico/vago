package db

import (
	"database/sql"
	"fmt"

	"github.com/sonirico/vago/lol"
	"go.elastic.co/apm/module/apmsql/v2"
	_ "go.elastic.co/apm/module/apmsql/v2/pq"
)

func OpenPostgres(uri string, log lol.Logger) (*sql.DB, error) {
	db, err := apmsql.Open("postgres", uri)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	} else {
		log.Debug("postgresql connected")
	}

	return db, nil
}

// OpenPg opens a PostgreSQL connection using apmsql and tunes the connection pool via query parameters.
// Example URI:
//
//	postgres://user:pass@localhost/dbname?pool_max_conns=10&pool_min_conns=2&pool_max_conn_lifetime=5m&pool_max_conn_idle_time=1m
func OpenPg(log lol.Logger, uri string) (*sql.DB, error) {
	// Extract pooling parameters and obtain the cleaned URI.
	params, err := parsePoolingParams(uri)
	if err != nil {
		return nil, err
	}

	return OpenPgOpts(log, params.CleanURI, params)
}

func OpenPgOpts(log lol.Logger, uri string, params PgPoolParams) (*sql.DB, error) {
	db, err := apmsql.Open("postgres", uri)
	if err != nil {
		return nil, err
	}

	// Configure the connection pool according to the parsed parameters.
	if params.MaxOpenConns > 0 {
		db.SetMaxOpenConns(params.MaxOpenConns)
	}
	if params.MinConns > 0 {
		db.SetMaxIdleConns(params.MinConns)
	}
	if params.MaxConnLifetime > 0 {
		db.SetConnMaxLifetime(params.MaxConnLifetime)
	}
	if params.MaxConnIdleTime > 0 {
		db.SetConnMaxIdleTime(params.MaxConnIdleTime)
	}

	// Check connectivity.
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	} else {
		log.Debug("postgres connected")
	}

	return db, nil
}
