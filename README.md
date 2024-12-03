# DayZ Prometheus Metrics Exporter

![logo]

Collects and publishes Prometheus metrics from Battleye RCON and
Steam A2S Query for DayZ server.

![dashboard]

## Configuration

```txt
Usage:
  ./dayz-exporter [option] [config.yaml]

Available options:
  -y, --get-yaml    Prints an example YAML configuration file.
  -e, --get-env     Prints an example .env file.
  -v, --version     Show version, commit, and build time.
  -h, --help        Prints this help message.
```

Run the program with a custom configuration file:

```bash
./dayz-exporter config.yaml
```

Configuration File Lookup:

* The program first checks for a configuration file specified by the
  variable `BERCON_EXPORTER_CONFIG_PATH`.
* If the config variable is not set, it will attempt to use a configuration
  file passed as a command-line argument.
* If no argument is provided, it will attempt to read from the default
  file `config.yaml` in the current directory.

You can get a sample YAML configuration here [example.config.yaml], or
save a sample YAML configuration to a file from within the application
itself by running:

```bash
./dayz-exporter --get-yaml > config.yaml
```

You can get a sample of environment variables here [example.env], or
save to a file from within the application itself by running:

```bash
./dayz-exporter --get-env > dayz-exporter.env
```

YAML config options have higher priority, for mixed use with variables
comment out the overridden options.

```bash
BERCON_EXPORTER_RCON_PASSWORD=strong %[1]s config.yaml
```

For more information on configuration parameters, refer to the example
configuration files (YAML and .env).

## Metrics

### Steam Query A2S INFO metrics

* **`a2s_info_players_online`** — Online players;
* **`a2s_info_players_slots`** — Players slots count;
* **`a2s_info_players_queue`** — Players wait in queue;
* **`a2s_info_time`** — Duration of day time on server;

### Battleye RCON players metrics

* **`bercon_player_ping_seconds`** — Ping of players in seconds.
  Extra labels: `name`, `ip`, `guid`, `lobby`, `country` [ℹ️](#labels)
* **`bercon_players_total`** — Total count of players;
* **`bercon_players_online`** — Count of players online;
* **`bercon_players_lobby`** — Count of players in lobby;
* **`bercon_players_invalid`** — Count of invalid players.

### Battleye RCON bans metrics (optional)

> [!TIP]  
> By default these metrics are disabled, they should be enabled separately
> in the settings. Can create a large number of metrics if you have
> a large ban list

* **`bercon_ban_guid_time_seconds`** — Time left for GUID bans in seconds.
  Extra labels: `reason`, `guid`;
* **`bercon_ban_guid_total`** — Total count of GUID bans;
* **`bercon_ban_ip_time_seconds`** — Time left for IP bans in seconds.
  Extra labels: `reason`, `ip`, `country` [ℹ️](#labels)
* **`bercon_ban_ip_total`** — Total count of IP bans;

### Labels

* **`server`** — Server name;
* **`map`** — Server map name;
* **`game`** — Game name;
* **`os`** — Server platform OS name;
* **`version`** — Game server version;
* Any static additional labels can also be installed via the application configuration.

> [!TIP]  
> `country` label show country code name only if GeoIP database configured
> in server settings

## Grafana Dashboard

Use dashboard from file [dayz-rcon.json] or [22457] from grafana

## Setup

### SystemD simple

```ini
[Unit]
Description=DayZ Prometheus Metrics Exporter
Documentation=https://woozymasta.github.io/dayz-exporter/
Wants=network-online.target
After=network-online.target dayz-server.target

[Service]
EnvironmentFile=-/etc/dayz-exporter.env
ExecStart=/usr/bin/dayz-exporter
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

### SystemD `--user` multiple servers

```ini
[Unit]
Description=DayZ Prometheus Metrics Exporter 'Server %I'
Documentation=https://woozymasta.github.io/dayz-exporter/
Wants=network-online.target
After=network-online.target dayz-server@%i.target

[Service]
WorkingDirectory=%h/dayz/
Environment="BERCON_EXPORTER_LISTEN_PORT=809%i"
Environment="BERCON_EXPORTER_QUERY_PORT=2702%i"
Environment="BERCON_EXPORTER_RCON_PORT=230%i5"
Environment="BERCON_EXPORTER_RCON_EXPOSE_BANS=true"
Environment="BERCON_EXPORTER_RCON_PASSWORD=strong"
Environment="BERCON_EXPORTER_GEOIP_DB=%h/dayz/GeoLite2-Country.mmdb"
ExecStart=%h/.local/bin/dayz-exporter
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=default.target
```

## Extra metrics

Uptime panel in main dashboard use metrics from [process-exporter]

Use extra dashboard from file [system-process.json] or [22458] from grafana

[process-exporter] configuration example for collect
DayZ servers process metrics

```yaml
process_names:
  - name: "dayz-{{.Matches.ConfigName}}"
    exe:
      - "DayZServer"
    cmdline:
      - '-config=(.*\/)?(?P<ConfigName>.*)\.cfg'
```

## Support me 💖

If you enjoy my projects and want to support further development,
feel free to donate! Every contribution helps to keep the work going.
Thank you!

### Crypto Donations

* **BTC**: `1Jb6vZAMVLQ9wwkyZfx2XgL5cjPfJ8UU3c`
* **USDT (TRC20)**: `TN99xawQTZKraRyvPAwMT4UfoS57hdH8Kz`
* **TON**: `UQBB5D7cL5EW3rHM_44rur9RDMz_fvg222R4dFiCAzBO_ptH`

Your support is greatly appreciated!

<!-- Links -->
[logo]: assets/dayz-exporter.png
[dashboard]: assets/dashboard.png
[example.config.yaml]: cli/example.config.yaml
[example.env]: cli/example.env
[dayz-rcon.json]: grafana/dayz-rcon.json
[system-process.json]: grafana/system-process.json

[process-exporter]: https://github.com/ncabatoff/process-exporter
[22457]: https://grafana.com/grafana/dashboards/22457 "DayZ Prometheus Metrics Exporter Dashboard"
[22458]: https://grafana.com/grafana/dashboards/22458 "System Processes Metrics Dashboard"
