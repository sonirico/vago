package testit

import (
	"fmt"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/sonirico/vago/slices"

	"github.com/ory/dockertest/v3/docker"
	"github.com/sonirico/vago/lol"
)

type Pool struct {
	dockerhost string
	exp        time.Duration
	pool       *dockertest.Pool
	resource   []*dockertest.Resource
	network    *docker.Network
	log        lol.Logger
}

type (
	RetryFunc   func(dockerhost string, resource *dockertest.Resource) retryFunc
	MigrateFunc func(dockerhost string, resource *dockertest.Resource) error
	SetEnvFunc  func(dockerhost string, resource *dockertest.Resource)
)

type Resource struct {
	RunOptions  *dockertest.RunOptions
	HostConfig  *docker.HostConfig
	RetryFunc   RetryFunc
	MigrateFunc MigrateFunc
	SetEnvFunc  SetEnvFunc
}

func NewPool(dockerhost string, log lol.Logger) (res *Pool, err error) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("could not construct pool: %w", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		return nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	network, err := pool.Client.CreateNetwork(docker.CreateNetworkOptions{
		Labels: map[string]string{"it_test": "true"},
		Name:   "integration_test_network",
	})
	if err != nil {
		return nil, fmt.Errorf("could not create a docker network: %w", err)
	}

	if dockerhost == "" {
		dockerhost = "localhost"
	}

	res = &Pool{
		dockerhost: dockerhost,
		pool:       pool,
		network:    network,
		exp:        300 * time.Second,
		log:        log,
	}

	return res, nil
}

func (s *Pool) Down() {
	s.log.Info("resource down")

	for _, resource := range s.resource {
		if resource == nil {
			continue
		}

		s.log.Infof("stopping resource: %s", resource.Container.Name)

		// You can't defer this because os.Exit doesn't care for defer
		if err := s.pool.Purge(resource); err != nil {
			s.log.Errorf("could not purge resource: %s", err)
		}
	}

	networks, _ := s.pool.Client.ListNetworks()
	toDelete := slices.Filter(networks, func(network docker.Network) bool {
		return network.Labels["it_test"] == "true"
	})

	s.log.Infof("networks to delete: %v", toDelete)

	for _, network := range toDelete {
		if err := s.pool.Client.RemoveNetwork(network.ID); err != nil {
			s.log.Errorf("could not remove %s network: %s", network.Name, err)
		}
	}
}

type retryFunc func() error

func (s *Pool) Retry(op retryFunc) error {
	s.pool.MaxWait = 120 * time.Second
	return s.pool.Retry(op)
}

func (s *Pool) run(
	options *dockertest.RunOptions,
	host *docker.HostConfig,
) (resource *dockertest.Resource, err error) {
	// pulls an image, creates a container based on it and runs it
	resource, err = s.pool.RunWithOptions(
		options,
		func(config *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
			if host == nil {
				return
			}
			config.Mounts = append(config.Mounts, host.Mounts...)
			config.PublishAllPorts = host.PublishAllPorts
		})
	if err != nil {
		return
	}

	_ = resource.Expire(uint(s.exp.Seconds()))
	s.resource = append(s.resource, resource)
	return
}

func (p *Pool) RunResource(res *Resource) error {
	options := res.RunOptions
	options.NetworkID = p.network.ID
	dtResource, err := p.run(options, res.HostConfig)
	if err != nil {
		return fmt.Errorf("%w: error running resource: %s", err, options.Repository)
	}

	if err := p.Retry(res.RetryFunc(p.dockerhost, dtResource)); err != nil {
		return fmt.Errorf("%w: error connecting to resource: %s", err, options.Repository)
	}

	res.SetEnvFunc(p.dockerhost, dtResource)

	if err := res.MigrateFunc(p.dockerhost, dtResource); err != nil {
		p.log.Infof("error migrating resource: %s", options.Repository)
		p.Down()

		return fmt.Errorf("%w: error initializing resource: %s", err, options.Repository)
	}

	return nil
}
