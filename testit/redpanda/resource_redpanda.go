package redpanda

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"sort"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sonirico/vago/opts"
	"github.com/sonirico/vago/testit"
	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	defaultRepository = "docker.redpanda.com/redpandadata/redpanda"
	defaultTag        = "v26.1.15"
)

// Config holds configuration for the Redpanda test resource.
type Config struct {
	Repository    string
	Tag           string
	Logger        testit.Logger
	SetEnvFunc    testit.SetEnvFunc
	ClusterConfig map[string]string
}

// Opt configures a Redpanda resource.
type Opt = opts.Configurator[Config]

// WithImage sets the Redpanda image repository and tag.
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

// WithClusterConfig sets cluster config keys applied with rpk once the
// broker answers, before the resource is handed to the caller.
func WithClusterConfig(config map[string]string) Opt {
	return opts.Fn[Config](func(c *Config) {
		c.ClusterConfig = config
	})
}

// NewResource builds a Redpanda container resource. The advertised
// external listener must carry a host port known before the container
// starts, so the free port is picked here; that reservation can fail,
// hence the error.
func NewResource(dockerhost string, options ...Opt) (*testit.Resource, error) {
	cfg := Config{
		Repository: defaultRepository,
		Tag:        defaultTag,
		Logger:     testit.NewNoopLogger(),
		SetEnvFunc: func(dockerhost string, resource *dockertest.Resource) {},
	}
	opts.ApplyAll(&cfg, options...)

	if dockerhost == "" {
		dockerhost = "localhost"
	}

	hostPort, err := freePort()
	if err != nil {
		return nil, fmt.Errorf("redpanda: reserve host port: %w", err)
	}

	log := cfg.Logger

	return &testit.Resource{
		RunOptions: &dockertest.RunOptions{
			Repository: cfg.Repository,
			Tag:        cfg.Tag,
			Hostname:   "redpanda",
			ExposedPorts: []string{
				"8081/tcp",
				"8082/tcp",
				"9092/tcp",
				"9644/tcp",
				"29092/tcp",
			},
			Cmd: []string{
				"redpanda start",
				"--smp 1",
				"--overprovisioned",
				"--kafka-addr PLAINTEXT://0.0.0.0:29092,OUTSIDE://0.0.0.0:9092",
				fmt.Sprintf(
					"--advertise-kafka-addr PLAINTEXT://redpanda:29092,OUTSIDE://%s:%d",
					dockerhost,
					hostPort,
				),
				"--pandaproxy-addr 0.0.0.0:8082",
				"--advertise-pandaproxy-addr localhost:8082",
			},
			PortBindings: map[docker.Port][]docker.PortBinding{
				"9092/tcp": {{HostPort: fmt.Sprintf("%d/tcp", hostPort)}},
			},
		},
		RetryFunc: func(dockerhost string, resource *dockertest.Resource) testit.RetryOp {
			url := BrokerAddr(dockerhost, resource)
			return func() error {
				log.Infof("redpanda: trying to connect to %s", url)

				cl, err := kgo.NewClient(kgo.SeedBrokers(url))
				if err != nil {
					return err
				}
				defer cl.Close()

				return cl.Ping(context.Background())
			}
		},
		MigrateFunc: func(dockerhost string, resource *dockertest.Resource) error {
			if len(cfg.ClusterConfig) == 0 {
				return nil
			}
			log.Infof("redpanda: applying %d cluster config keys", len(cfg.ClusterConfig))
			return SetClusterConfig(resource, cfg.ClusterConfig)
		},
		SetEnvFunc: cfg.SetEnvFunc,
	}, nil
}

// BrokerAddr returns the host-visible bootstrap address for a running
// Redpanda resource, as advertised on the external listener.
func BrokerAddr(dockerhost string, resource *dockertest.Resource) string {
	if dockerhost == "" {
		dockerhost = "localhost"
	}
	return fmt.Sprintf("%s:%s", dockerhost, resource.GetPort("9092/tcp"))
}

// SetClusterConfig applies each key with rpk cluster config set inside
// the container, in sorted key order for determinism.
func SetClusterConfig(resource *dockertest.Resource, config map[string]string) error {
	keys := make([]string, 0, len(config))
	for k := range config {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		var out bytes.Buffer
		exitCode, err := resource.Exec(
			[]string{"rpk", "cluster", "config", "set", k, config[k]},
			dockertest.ExecOptions{StdOut: &out, StdErr: &out},
		)
		if err != nil {
			return fmt.Errorf("rpk cluster config set %s: %w", k, err)
		}
		if exitCode != 0 {
			return fmt.Errorf("rpk cluster config set %s: exit %d: %s", k, exitCode, out.String())
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
