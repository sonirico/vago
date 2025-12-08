package testit

import (
	"fmt"
	"log"

	"github.com/sonirico/vago/lol"
)

type DockerResourcesPool struct {
	pool      *Pool
	resources []*Resource
	logger    lol.Logger
}

func NewDockerResourcesPool(
	logger lol.Logger,
	dockerhost string,
	res ...*Resource,
) *DockerResourcesPool {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := NewPool(dockerhost, logger)
	if err != nil {
		log.Panicf("Could not construct pool: %s", err)
	}

	return &DockerResourcesPool{
		pool:      pool,
		resources: res,
		logger:    logger,
	}
}

func (s *DockerResourcesPool) Down() {
	s.pool.Down()
}

func (s *DockerResourcesPool) Up() (err error) {
	if err = s.up(); err != nil {
		s.Down()
		return err
	}

	return nil
}

func (s *DockerResourcesPool) up() error {
	for _, res := range s.resources {
		if err := s.pool.RunResource(res); err != nil {
			return fmt.Errorf("error running resource: %w", err)
		}
	}

	return nil
}
