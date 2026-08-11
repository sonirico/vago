package minio_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"

	"github.com/sonirico/vago/testit"
	testitminio "github.com/sonirico/vago/testit/minio"
)

func TestMinioRoundTrip(t *testing.T) {
	probe, err := dockertest.NewPool("")
	if err != nil || probe.Client.Ping() != nil {
		t.Skipf("docker unavailable: %v", err)
	}

	var endpoint string
	res, err := testitminio.NewResource(
		os.Getenv("DOCKER_HOSTNAME"),
		testitminio.WithBuckets("it-bucket"),
		testitminio.WithSetEnvFunc(func(dockerhost string, resource *dockertest.Resource) {
			endpoint = testitminio.Endpoint(dockerhost, resource)
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
	require.NotEmpty(t, endpoint)

	healthResp, err := http.Get(fmt.Sprintf("http://%s/minio/health/live", endpoint))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, healthResp.StatusCode)
	require.NoError(t, healthResp.Body.Close())

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	require.NoError(t, err)

	exists, err := client.BucketExists(context.Background(), "it-bucket")
	require.NoError(t, err)
	require.True(t, exists)
}
