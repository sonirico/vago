package db

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"errors"

	"github.com/sonirico/vago/lol"
)

var (
	ErrReceivedFromServer = errors.New("error received from server")
)

type ClickhouseHttp struct {
	log    lol.Logger
	url    *url.URL
	client *http.Client
}

func NewClickHouseHttp(log lol.Logger, addr string, clientHttp *http.Client) *ClickhouseHttp {
	u, err := url.Parse(addr)
	if err != nil {
		log.Fatalf("failed to parse url for ch http %q: %v", addr, err)
	}

	cli := &ClickhouseHttp{
		log:    log,
		url:    u,
		client: clientHttp,
	}

	return cli
}

func (r *ClickhouseHttp) Query(ctx context.Context, query string) (res io.ReadCloser, err error) {
	r.log.Debugln(query)

	bodyReader := strings.NewReader(query)

	return r.request(ctx, http.MethodPost, "", bodyReader, http.StatusOK)
}

func (r *ClickhouseHttp) Ping(ctx context.Context) error {
	_, err := r.request(ctx, http.MethodGet, "/ping", nil, http.StatusOK)
	return err
}

func (r *ClickhouseHttp) request(
	ctx context.Context,
	method string,
	path string,
	bodyReader io.Reader,
	expectedStatusCode int,
) (res io.ReadCloser, err error) {
	uri := r.url.ResolveReference(&url.URL{Path: path}).String()

	req, err := http.NewRequestWithContext(ctx, method, uri, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error requesting %s: %w", uri, err)
	}

	if resp.StatusCode != expectedStatusCode {
		msg, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}

		return nil, fmt.Errorf(
			"unexpected status code %d, expected %d: %s",
			resp.StatusCode,
			expectedStatusCode,
			msg,
		)
	}

	return resp.Body, nil
}
