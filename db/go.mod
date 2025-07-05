module github.com/sonirico/vago/db

go 1.23.0

toolchain go1.23.9

require (
	github.com/ClickHouse/clickhouse-go/v2 v2.37.2
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang-migrate/migrate/v4 v4.18.3
	github.com/jackc/pgx/v5 v5.7.5
	github.com/sonirico/vago v0.7.2
	github.com/sonirico/vago/lol v0.0.0-20250705171120-b1db77f5b640
	github.com/stretchr/testify v1.10.0
	go.elastic.co/apm/module/apmgoredisv8/v2 v2.7.1
	go.elastic.co/apm/module/apmmongo/v2 v2.7.1
	go.elastic.co/apm/module/apmpgxv5/v2 v2.7.1
	go.elastic.co/apm/module/apmsql/v2 v2.7.1
	go.elastic.co/apm/v2 v2.7.1
	go.mongodb.org/mongo-driver v1.17.4
)

require (
	github.com/ClickHouse/ch-go v0.66.1 // indirect
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/elastic/go-sysinfo v1.15.3 // indirect
	github.com/elastic/go-windows v1.0.2 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/joeshaw/multierror v0.0.0-20140124173710-69b34d4ec901 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/paulmach/orb v0.11.1 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/procfs v0.17.0 // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.elastic.co/apm/module/apmzerolog/v2 v2.7.1 // indirect
	go.elastic.co/fastjson v1.5.1 // indirect
	go.opentelemetry.io/otel v1.37.0 // indirect
	go.opentelemetry.io/otel/trace v1.37.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/sync v0.15.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	howett.net/plist v1.0.1 // indirect
)

replace github.com/sonirico/vago => ../

replace github.com/sonirico/vago/lol => ../lol
