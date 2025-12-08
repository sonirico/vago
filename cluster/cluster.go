package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sonirico/vago/clock"
	"github.com/sonirico/vago/lol"
)

//go:generate easyjson

type (
	repo interface {
		Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
		Scan(ctx context.Context, key string) (map[string][]byte, error)
	}

	//easyjson:json
	Session struct {
		LastUpdated time.Time `json:"last_updated"`
		NodeID      string    `json:"node_id"`
		ServiceID   string    `json:"service_id"`
		// Interests   any `json:"interests"`
	}

	logger interface {
		Info(...any)
		Infof(string, ...any)
		Errorf(string, ...any)
	}

	node interface {
		ServiceID() string
		NodeID() string
	}

	Cluster struct {
		clock  clock.Clock
		logger logger
		repo   repo
		ns     string
		ttl    time.Duration
	}
)

func NewCluster(
	logger lol.Logger,
	cli redis.UniversalClient,
	namespace string,
	ttl time.Duration,
) *Cluster {
	return &Cluster{
		clock:  clock.New(),
		logger: logger,
		ns:     namespace,
		repo:   &redisRepo{cli: cli},
		ttl:    ttl,
	}
}

// Register starts a new session in the cluster and blocks until context is canceled.
func (c *Cluster) Register(ctx context.Context, node node) error {
	sleepage := c.ttl / 2
	ticker := time.NewTicker(sleepage)
	defer ticker.Stop() // Ensure that the timer's resources are freed

	c.logger.Infof("starting redis cluster, node=%s, interval=%v",
		node.NodeID(), sleepage)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case t, ok := <-ticker.C:
			if !ok {
				c.logger.Info("ticker was closed")
			}

			cx, cancel := context.WithTimeout(ctx, sleepage)
			c.logger.Infof("register at %s", t)

			if err := c.register(cx, node); err != nil {
				c.logger.Errorf("cannot register: %v", err)
			}

			cancel()
		}
	}
}

func (c *Cluster) DiscoverNodes(
	ctx context.Context,
	serviceID string,
) (map[string]Session, error) {
	return c.discover(ctx, c.serviceSessionPath(serviceID))
}

func (c *Cluster) Discover(ctx context.Context) (map[string]Session, error) {
	return c.discover(ctx, c.ns)
}

func (c *Cluster) discover(ctx context.Context, path string) (map[string]Session, error) {
	res, err := c.repo.Scan(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("%w: target key %s: %v", ErrCannotScan, c.ns, err)
	}

	if len(res) < 1 {
		return nil, ErrNoSessions
	}

	sessions := make(map[string]Session, len(res))

	for key, sessionRaw := range res {
		info := Session{}
		if err = json.Unmarshal(sessionRaw, &info); err != nil {
			return nil, err
		}
		sessions[key] = info
	}

	return sessions, nil
}

func (c *Cluster) register(ctx context.Context, node node) error {
	info := Session{
		LastUpdated: c.clock.Now(),
		NodeID:      node.NodeID(),
		ServiceID:   node.ServiceID(),
	}

	bts, err := json.Marshal(info)
	if err != nil {
		return err
	}

	if err = c.repo.Set(ctx, c.nodeSessionPath(node), bts, c.ttl); err != nil {
		c.logger.Errorf("cannot register: %q", err)
		return err
	}
	return nil
}

func (c *Cluster) nodeSessionNamespace(node node) string {
	return c.serviceSessionPath(node.ServiceID())
}

func (c *Cluster) nodeSessionPath(node node) string {
	return c.ns + ":" + node.ServiceID() + ":" + node.NodeID()
}

func (c *Cluster) serviceSessionPath(svc string) string {
	return c.ns + ":" + svc + ":"
}
