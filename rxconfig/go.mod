module github.com/sonirico/vago/rxconfig

go 1.25.5

replace github.com/sonirico/vago => ../

replace github.com/sonirico/vago/lol => ../lol

replace github.com/sonirico/vago/codec => ../codec

require (
	github.com/go-redis/redis/v8 v8.11.5
	github.com/mailru/easyjson v0.9.1
	github.com/pkg/errors v0.9.1
	github.com/sonirico/vago/codec v0.0.0-00010101000000-000000000000
	github.com/sonirico/vago/lol v0.0.0-00010101000000-000000000000
)

require (
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/elastic/go-sysinfo v1.15.4 // indirect
	github.com/elastic/go-windows v1.0.2 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/prometheus/procfs v0.19.2 // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	go.elastic.co/apm/module/apmzerolog/v2 v2.7.2 // indirect
	go.elastic.co/apm/v2 v2.7.2 // indirect
	go.elastic.co/fastjson v1.5.1 // indirect
	golang.org/x/sys v0.39.0 // indirect
	howett.net/plist v1.0.1 // indirect
)
