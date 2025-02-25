module github.com/woozymasta/dayz-exporter

go 1.24

require (
	github.com/oschwald/geoip2-golang v1.11.0
	github.com/prometheus/client_golang v1.21.0
	github.com/rs/zerolog v1.33.0
	github.com/sethvargo/go-envconfig v1.1.1
	github.com/woozymasta/a2s v0.2.2
	github.com/woozymasta/bercon-cli v0.3.1
	golang.org/x/sys v0.30.0
	golang.org/x/term v0.29.0
	gopkg.in/yaml.v3 v3.0.1
	internal/vars v0.0.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/oschwald/maxminddb-golang v1.13.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/rs/xid v1.6.0 // indirect
	github.com/woozymasta/steam v0.1.3 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)

replace internal/vars => ./internal/vars
