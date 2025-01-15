module github.com/woozymasta/dayz-exporter

go 1.23.1

require (
	github.com/oschwald/geoip2-golang v1.11.0
	github.com/prometheus/client_golang v1.20.5
	github.com/rs/zerolog v1.33.0
	github.com/sethvargo/go-envconfig v1.1.0
	github.com/woozymasta/a2s v0.2.1
	github.com/woozymasta/bercon-cli v0.3.0
	golang.org/x/sys v0.29.0
	golang.org/x/term v0.28.0
	gopkg.in/yaml.v3 v3.0.1
	internal/vars v0.0.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/oschwald/maxminddb-golang v1.13.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.61.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/woozymasta/steam v0.1.2 // indirect
	google.golang.org/protobuf v1.36.3 // indirect
)

replace internal/vars => ./internal/vars
