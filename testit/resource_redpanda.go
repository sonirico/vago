package testit

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/ory/dockertest/v3"

	"github.com/ory/dockertest/v3/docker"
	"github.com/sonirico/vago/lol"
	"github.com/twmb/franz-go/pkg/kgo"
)

func portAvailable(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return false
	}

	defer ln.Close()
	return true
}

func randomAvailablePort(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	port := rand.Intn(max-min) + min
	if !portAvailable(port) {
		return randomAvailablePort(min, max)
	}
	return port
}

func NewRedpandaResource(
	dockerhost string,
	logger lol.Logger,
	migrateFunc MigrateFunc,
	envFunc SetEnvFunc,
) *Resource {
	log := logger.WithFields(lol.Fields{"resource": "redpanda"})
	max := 10000
	min := 9000
	hostPort := randomAvailablePort(min, max)
	if dockerhost == "" {
		dockerhost = "localhost"
	}
	return &Resource{
		RunOptions: &dockertest.RunOptions{
			Repository: "docker.redpanda.com/redpandadata/redpanda",
			Tag:        "v24.2.10",
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
				"--Redpanda-addr PLAINTEXT://0.0.0.0:29092,OUTSIDE://0.0.0.0:9092",
				fmt.Sprintf(
					"--advertise-Redpanda-addr PLAINTEXT://redpanda:29092,OUTSIDE://%s:%d",
					dockerhost,
					hostPort,
				),
				"--pandaproxy-addr 0.0.0.0:8082",
				"--advertise-pandaproxy-addr localhost:8082",
			},
			PortBindings: map[docker.Port][]docker.PortBinding{
				"9092/tcp": {{HostIP: "localhost", HostPort: fmt.Sprintf("%d/tcp", hostPort)}},
			},
		},
		RetryFunc: func(dockerhost string, resource *dockertest.Resource) retryFunc {
			url := fmt.Sprintf("%s:%s", dockerhost, resource.GetPort("9092/tcp"))
			return func() error {
				log.Printf("Connecting to redpanda url: '%s'", url)
				seeds := []string{url}
				log.Info("trying to connect")

				cl, err := kgo.NewClient(
					kgo.SeedBrokers(seeds...),
				)
				if err != nil {
					return err
				}
				defer cl.Close()

				ctx := context.Background()

				return cl.Ping(ctx)
			}
		},
		MigrateFunc: migrateFunc,
		SetEnvFunc:  envFunc,
	}
}
