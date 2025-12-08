module github.com/sonirico/vago/cqrs

go 1.25.5

require (
	github.com/google/uuid v1.6.0
	github.com/mailru/easyjson v0.9.1
	github.com/sonirico/stadio v0.8.0
	github.com/sonirico/vago v0.9.0
	github.com/sonirico/vago/lol v0.0.0-20251207192038-45d83c821566
	github.com/sonirico/vago/rp v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.11.1
	github.com/twmb/franz-go v1.20.5
	go.elastic.co/apm/v2 v2.7.2
)

require (
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/elastic/go-sysinfo v1.15.4 // indirect
	github.com/elastic/go-windows v1.0.2 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.18.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/procfs v0.19.2 // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	github.com/twmb/franz-go/pkg/kmsg v1.12.0 // indirect
	go.elastic.co/apm/module/apmhttp/v2 v2.7.2 // indirect
	go.elastic.co/apm/module/apmzerolog/v2 v2.7.2 // indirect
	go.elastic.co/fastjson v1.5.1 // indirect
	golang.org/x/sys v0.39.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	howett.net/plist v1.0.1 // indirect
)

replace github.com/sonirico/vago => ../

replace github.com/sonirico/vago/lol => ../lol

replace github.com/sonirico/vago/rp => ../rp

replace github.com/sonirico/vago/codec => ../codec
