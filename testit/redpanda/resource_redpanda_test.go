package redpanda_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/sonirico/vago/testit"
	"github.com/sonirico/vago/testit/redpanda"
)

func TestRedpandaRoundTrip(t *testing.T) {
	probe, err := dockertest.NewPool("")
	if err != nil || probe.Client.Ping() != nil {
		t.Skipf("docker unavailable: %v", err)
	}

	var brokerAddr string
	res, err := redpanda.NewResource(
		os.Getenv("DOCKER_HOSTNAME"),
		redpanda.WithClusterConfig(map[string]string{"log_compaction_interval_ms": "500"}),
		redpanda.WithSetEnvFunc(func(dockerhost string, resource *dockertest.Resource) {
			brokerAddr = redpanda.BrokerAddr(dockerhost, resource)
		}),
	)
	require.NoError(t, err)

	pool := testit.NewDockerResourcesPool(
		testit.NewNoopLogger(),
		os.Getenv("DOCKER_HOSTNAME"),
		res,
	)
	require.NoError(t, pool.Up())
	t.Cleanup(pool.Down)
	require.NotEmpty(t, brokerAddr)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	const topic = "testit-smoke"
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokerAddr),
		kgo.AllowAutoTopicCreation(),
		kgo.ConsumeTopics(topic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	require.NoError(t, err)
	t.Cleanup(client.Close)

	produced := &kgo.Record{Topic: topic, Key: []byte("k1"), Value: []byte("v1")}
	require.NoError(t, client.ProduceSync(ctx, produced).FirstErr())

	var got *kgo.Record
	for got == nil {
		fetches := client.PollFetches(ctx)
		require.NoError(t, ctx.Err())
		require.Empty(t, fetches.Errors())
		if records := fetches.Records(); len(records) > 0 {
			got = records[0]
		}
	}

	require.Equal(t, produced.Key, got.Key)
	require.Equal(t, produced.Value, got.Value)
}
