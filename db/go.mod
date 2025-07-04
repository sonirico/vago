module github.com/sonirico/vago/db

go 1.23.0

toolchain go1.23.9

require (
	github.com/ClickHouse/clickhouse-go/v2 v2.25.0
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang-migrate/migrate/v4 v4.17.1
	github.com/jackc/pgx/v5 v5.5.5
	github.com/sonirico/vago v0.0.0-20240701000000-000000000000
	github.com/sonirico/vago/lol v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.10.0
	go.elastic.co/apm/module/apmgoredisv8/v2 v2.6.0
	go.elastic.co/apm/module/apmmongo/v2 v2.6.0
	go.elastic.co/apm/module/apmpgxv5/v2 v2.6.0
	go.elastic.co/apm/module/apmsql/v2 v2.6.0
	go.elastic.co/apm/v2 v2.7.1
	go.mongodb.org/mongo-driver v1.16.0
)

require (
	github.com/ClickHouse/ch-go v0.61.5 // indirect
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/elastic/go-sysinfo v1.7.1 // indirect
	github.com/elastic/go-windows v1.0.0 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/joeshaw/multierror v0.0.0-20140124173710-69b34d4ec901 // indirect
	github.com/klauspost/compress v1.17.7 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/paulmach/orb v0.11.1 // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/procfs v0.0.0-20190425082905-87a4384529e0 // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	go.elastic.co/fastjson v1.5.1 // indirect
	go.opentelemetry.io/otel v1.26.0 // indirect
	go.opentelemetry.io/otel/trace v1.26.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	howett.net/plist v0.0.0-20181124034731-591f970eefbb // indirect
)

replace github.com/sonirico/vago => ../

replace github.com/sonirico/vago/lol => ../lol
