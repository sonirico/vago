package db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-redis/redis/v8"
)

type (
	RedisTLSConfig struct {
		Enabled  bool
		CertFile string
		KeyFile  string
		CAFile   string
	}

	RedisConfig struct {
		URL string
		TLS RedisTLSConfig
	}
)

func newRedisTLSConfig(cfg RedisTLSConfig) (*tls.Config, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("cannot load TLS key pair from certFile=%q and keyFile=%q: %w",
			cfg.CertFile, cfg.KeyFile, err,
		)
	}

	serverCACert, err := os.ReadFile(filepath.Clean(cfg.CAFile))
	if err != nil {
		return nil, fmt.Errorf("cannot load TLS caFile=%q: %w", cfg.CAFile, err)
	}

	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(serverCACert)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS13,
		RootCAs:    clientCertPool,
		Certificates: []tls.Certificate{
			cert,
		},
	}

	return tlsConfig, nil
}

func OpenRedis(ctx context.Context, cfg RedisConfig) (*redis.Client, error) {
	tlsConfig, err := newRedisTLSConfig(cfg.TLS)
	if err != nil {
		return nil, err
	}

	opts, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, err
	}
	opts.Username = "" // default user
	opts.TLSConfig = tlsConfig

	// TODO: Consider RedisCluster
	conn := redis.NewClient(opts)
	conn.AddHook(NewApmRedisHook())

	if err = conn.Ping(ctx).Err(); err != nil {
		conn.Close()

		return nil, err
	}

	return conn, nil
}
