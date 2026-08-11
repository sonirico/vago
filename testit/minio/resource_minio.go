package minio

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sonirico/vago/opts"
	"github.com/sonirico/vago/testit"
)

const (
	defaultRepository = "minio/minio"
	defaultTag        = "RELEASE.2025-09-07T16-13-09Z"

	defaultRootUser     = "minioadmin"
	defaultRootPassword = "minioadmin"
)

// InternalEndpoint is the S3 endpoint as seen by containers sharing the
// testit docker network.
const InternalEndpoint = "minio:9000"

// Config holds configuration for the MinIO test resource.
type Config struct {
	Repository   string
	Tag          string
	Logger       testit.Logger
	SetEnvFunc   testit.SetEnvFunc
	RootUser     string
	RootPassword string
	Buckets      []string
}

// Opt configures a MinIO resource.
type Opt = opts.Configurator[Config]

// WithImage sets the MinIO image repository and tag.
func WithImage(repository, tag string) Opt {
	return opts.Fn[Config](func(c *Config) {
		c.Repository = repository
		c.Tag = tag
	})
}

// WithLogger sets the logger.
func WithLogger(log testit.Logger) Opt {
	return opts.Fn[Config](func(c *Config) {
		c.Logger = log
	})
}

// WithSetEnvFunc sets the environment setup function.
func WithSetEnvFunc(fn testit.SetEnvFunc) Opt {
	return opts.Fn[Config](func(c *Config) {
		c.SetEnvFunc = fn
	})
}

// WithCredentials sets the root user and password.
func WithCredentials(user, password string) Opt {
	return opts.Fn[Config](func(c *Config) {
		c.RootUser = user
		c.RootPassword = password
	})
}

// WithBuckets sets bucket names created once the server answers, before the
// resource is handed to the caller.
func WithBuckets(buckets ...string) Opt {
	return opts.Fn[Config](func(c *Config) {
		c.Buckets = buckets
	})
}

// NewResource builds a MinIO container resource. The advertised external
// listener must carry a host port known before the container starts, so the
// free port is picked here; that reservation can fail, hence the error.
func NewResource(dockerhost string, options ...Opt) (*testit.Resource, error) {
	cfg := Config{
		Repository:   defaultRepository,
		Tag:          defaultTag,
		Logger:       testit.NewNoopLogger(),
		SetEnvFunc:   func(dockerhost string, resource *dockertest.Resource) {},
		RootUser:     defaultRootUser,
		RootPassword: defaultRootPassword,
	}
	opts.ApplyAll(&cfg, options...)

	if dockerhost == "" {
		dockerhost = "localhost"
	}

	hostPort, err := freePort()
	if err != nil {
		return nil, fmt.Errorf("minio: reserve host port: %w", err)
	}

	log := cfg.Logger

	return &testit.Resource{
		RunOptions: &dockertest.RunOptions{
			Repository: cfg.Repository,
			Tag:        cfg.Tag,
			Hostname:   "minio",
			Env: []string{
				"MINIO_ROOT_USER=" + cfg.RootUser,
				"MINIO_ROOT_PASSWORD=" + cfg.RootPassword,
			},
			ExposedPorts: []string{
				"9000/tcp",
			},
			Cmd: []string{
				"server",
				"/data",
				"--address",
				":9000",
			},
			PortBindings: map[docker.Port][]docker.PortBinding{
				"9000/tcp": {{HostPort: fmt.Sprintf("%d/tcp", hostPort)}},
			},
		},
		RetryFunc: func(dockerhost string, resource *dockertest.Resource) testit.RetryOp {
			url := fmt.Sprintf("http://%s/minio/health/live", Endpoint(dockerhost, resource))
			return func() error {
				log.Infof("minio: trying to connect to %s", url)

				resp, err := http.Get(url)
				if err != nil {
					return err
				}
				defer func() {
					if closeErr := resp.Body.Close(); closeErr != nil {
						log.Errorf("minio: close health response body: %s", closeErr)
					}
				}()

				if resp.StatusCode != http.StatusOK {
					return fmt.Errorf("minio: health check returned status %d", resp.StatusCode)
				}

				return nil
			}
		},
		MigrateFunc: func(dockerhost string, resource *dockertest.Resource) error {
			if len(cfg.Buckets) == 0 {
				return nil
			}
			log.Infof("minio: creating %d buckets", len(cfg.Buckets))
			return createBuckets(dockerhost, resource, cfg)
		},
		SetEnvFunc: cfg.SetEnvFunc,
	}, nil
}

// Endpoint returns the host-visible S3 endpoint (host:port) for a running
// MinIO resource.
func Endpoint(dockerhost string, resource *dockertest.Resource) string {
	if dockerhost == "" {
		dockerhost = "localhost"
	}
	return fmt.Sprintf("%s:%s", dockerhost, resource.GetPort("9000/tcp"))
}

func createBuckets(dockerhost string, resource *dockertest.Resource, cfg Config) error {
	client, err := minio.New(Endpoint(dockerhost, resource), &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.RootUser, cfg.RootPassword, ""),
		Secure: false,
	})
	if err != nil {
		return fmt.Errorf("minio: build client: %w", err)
	}

	ctx := context.Background()
	for _, bucket := range cfg.Buckets {
		exists, err := client.BucketExists(ctx, bucket)
		if err != nil {
			return fmt.Errorf("minio: check bucket %s: %w", bucket, err)
		}
		if exists {
			continue
		}
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("minio: make bucket %s: %w", bucket, err)
		}
	}

	return nil
}

func freePort() (int, error) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	port := ln.Addr().(*net.TCPAddr).Port
	if err := ln.Close(); err != nil {
		return 0, err
	}
	return port, nil
}
