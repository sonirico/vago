module github.com/sonirico/vago/db

go 1.25.3

replace github.com/sonirico/vago => ../

replace github.com/sonirico/vago/lol => ../lol

require (
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/sonirico/vago v0.0.0-00010101000000-000000000000
	github.com/sonirico/vago/lol v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/elastic/go-sysinfo v1.15.4 // indirect
	github.com/elastic/go-windows v1.0.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/procfs v0.19.2 // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	go.elastic.co/apm/module/apmzerolog/v2 v2.7.2 // indirect
	go.elastic.co/apm/v2 v2.7.2 // indirect
	go.elastic.co/fastjson v1.5.1 // indirect
	golang.org/x/sys v0.39.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	howett.net/plist v1.0.1 // indirect
)
