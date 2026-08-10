// Package testit provides the backend-agnostic docker resource pool for
// integration tests: running containers with dockertest, retrying until
// they answer, applying per-resource setup, and tearing everything down.
// It has no backend driver dependency: starting an actual backend is the
// job of the sibling backend modules below.
//
// Main types:
//   - Resource: a container recipe (run options, retry, migrate, env hooks).
//   - Pool, DockerResourcesPool: run resources on a private docker network.
//   - RunSafe: TestMain wrapper that tears the pool down even on panic.
//   - Logger: the minimal logging contract this package needs from a
//     caller-supplied logger.
//
// Backend modules:
//   - github.com/sonirico/vago/testit/postgres: Postgres resource, suite and SQL fixtures
//   - github.com/sonirico/vago/testit/clickhouse: ClickHouse resource, suite and fixtures
//   - github.com/sonirico/vago/testit/redis: Redis resource and suite
//   - github.com/sonirico/vago/testit/redpanda: Redpanda resource with cluster config support
package testit
