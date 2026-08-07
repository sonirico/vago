package db

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sonirico/vago/lol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// closeTrackingBody wraps an io.ReadCloser and records whether Close was called.
type closeTrackingBody struct {
	io.ReadCloser
	closed *bool
}

func (b closeTrackingBody) Close() error {
	*b.closed = true
	return b.ReadCloser.Close()
}

// closeTrackingTransport wraps http.RoundTripper so tests can observe whether
// the response body returned by request() was closed.
type closeTrackingTransport struct {
	base   http.RoundTripper
	closed *bool
}

func (t closeTrackingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	resp.Body = closeTrackingBody{ReadCloser: resp.Body, closed: t.closed}
	return resp, nil
}

func newTestClickHouseHttp(t *testing.T, handler http.HandlerFunc) (*ClickhouseHttp, *bool) {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	closed := new(bool)
	client := &http.Client{
		Transport: closeTrackingTransport{base: http.DefaultTransport, closed: closed},
	}

	return NewClickHouseHttp(lol.ZeroTestLogger, srv.URL, client), closed
}

func TestClickhouseHttp_Request_ClosesBodyOnUnexpectedStatus(t *testing.T) {
	t.Parallel()

	// Arrange
	cli, closed := newTestClickHouseHttp(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("boom"))
	})

	// Act
	err := cli.Ping(context.Background())

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code 500")
	assert.True(t, *closed, "response body must be closed on the error path")
}

func TestClickhouseHttp_Request_LeavesBodyOpenOnExpectedStatus(t *testing.T) {
	t.Parallel()

	// Arrange
	cli, closed := newTestClickHouseHttp(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("1\n"))
	})

	// Act
	body, err := cli.Query(context.Background(), "SELECT 1")

	// Assert
	require.NoError(t, err)
	assert.False(t, *closed, "success path must hand the caller an open body")
	require.NoError(t, body.Close())
}
