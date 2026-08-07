module github.com/sonirico/vago/db/postgres

go 1.25.3

replace github.com/sonirico/vago => ../../

replace github.com/sonirico/vago/lol => ../../lol

replace github.com/sonirico/vago/db => ../

require (
	github.com/golang-migrate/migrate/v4 v4.19.1
	github.com/jackc/pgx/v5 v5.7.6
	github.com/sonirico/vago/db v0.2.0
	github.com/sonirico/vago/lol v0.1.0
	github.com/stretchr/testify v1.11.1
	go.elastic.co/apm/module/apmpgxv5/v2 v2.7.2
	go.elastic.co/apm/module/apmsql/v2 v2.7.2
)

require (
	github.com/DATA-DOG/go-sqlmock v1.5.2 // indirect
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/docker/docker v28.5.2+incompatible // indirect
	github.com/docker/go-connections v0.6.0 // indirect
	github.com/elastic/go-sysinfo v1.15.4 // indirect
	github.com/elastic/go-windows v1.0.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/moby/term v0.5.2 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/procfs v0.19.2 // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	github.com/sonirico/vago v0.12.0 // indirect
	go.elastic.co/apm/module/apmzerolog/v2 v2.7.2 // indirect
	go.elastic.co/apm/v2 v2.7.2 // indirect
	go.elastic.co/fastjson v1.5.1 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.63.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	golang.org/x/crypto v0.45.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	howett.net/plist v1.0.1 // indirect
)
