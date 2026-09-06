module github.com/sonirico/vago/lol/v2/zerolog/apm

go 1.25.0

require (
	github.com/rs/zerolog v1.34.0
	github.com/sonirico/vago/lol/v2 v2.0.0
	github.com/sonirico/vago/lol/v2/zerolog v0.1.0
	github.com/stretchr/testify v1.12.1
	go.elastic.co/apm/module/apmzerolog/v2 v2.7.2
	go.elastic.co/apm/v2 v2.7.2
)

require (
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/elastic/go-sysinfo v1.7.1 // indirect
	github.com/elastic/go-windows v1.0.0 // indirect
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/joeshaw/multierror v0.0.0-20140124173710-69b34d4ec901 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/procfs v0.0.0-20190425082905-87a4384529e0 // indirect
	go.elastic.co/fastjson v1.5.1 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sys v0.12.0 // indirect
	howett.net/plist v0.0.0-20181124034731-591f970eefbb // indirect
)

replace github.com/sonirico/vago/lol/v2 => ../../

replace github.com/sonirico/vago/lol/v2/zerolog => ../
