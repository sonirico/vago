module github.com/sonirico/vago/db/redis

go 1.25.3

replace github.com/sonirico/vago => ../../

replace github.com/sonirico/vago/lol => ../../lol

replace github.com/sonirico/vago/db => ../

require (
	github.com/go-redis/redis/v8 v8.11.5
	go.elastic.co/apm/module/apmgoredisv8/v2 v2.7.2
	go.elastic.co/apm/v2 v2.7.2
)

require (
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/elastic/go-sysinfo v1.15.4 // indirect
	github.com/elastic/go-windows v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/procfs v0.19.2 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	go.elastic.co/fastjson v1.5.1 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	howett.net/plist v1.0.1 // indirect
)
