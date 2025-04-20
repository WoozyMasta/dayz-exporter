package main

import (
	"context"
	"errors"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

const defaultConfigPath = "config.yaml"

// Config represents the main configuration structure for the exporter.
type Config struct {
	Labels  map[string]string `yaml:"labels,omitempty" env:"DAYZ_EXPORTER_LABELS"`
	Logging Logging           `yaml:"logging,omitempty" env:", prefix=DAYZ_EXPORTER_LOG_"`
	GeoDB   string            `yaml:"geo_db,omitempty" env:"DAYZ_EXPORTER_GEOIP_DB"`
	Listen  Listen            `yaml:"listen,omitempty" env:", prefix=DAYZ_EXPORTER_LISTEN_"`
	Query   Query             `yaml:"query,omitempty" env:", prefix=DAYZ_EXPORTER_QUERY_"`
	Rcon    Rcon              `yaml:"rcon,omitempty" env:", prefix=DAYZ_EXPORTER_RCON_"`
}

// Listen contains settings for the exporter's HTTP server.
type Listen struct {
	IP          string `yaml:"ip,omitempty" env:"IP, default=0.0.0.0"`
	Endpoint    string `yaml:"endpoint,omitempty" env:"ENDPOINT, default=/metrics"`
	Username    string `yaml:"username,omitempty" env:"USERNAME, default=metrics"`
	Password    string `yaml:"password,omitempty" env:"PASSWORD"`
	CORSDomains string `yaml:"cors_domains,omitempty" env:"CORS_DOMAINS"`
	Port        uint16 `yaml:"port,omitempty" env:"PORT, default=8098"`
	ExposeInfo  bool   `yaml:"expose_info,omitempty" env:"EXPOSE_INFO, default=false"`
	InfoAuth    bool   `yaml:"info_auth,omitempty" env:"INFO_AUTH, default=false"`
	HealthAuth  bool   `yaml:"health_auth,omitempty" env:"HEALTH_AUTH, default=false"`
}

// Query contains Steam A2S query connection settings.
type Query struct {
	IP   string `yaml:"ip,omitempty" env:"IP, default=127.0.0.1"`
	Port int    `yaml:"port,omitempty" env:"PORT, default=27016"`
}

// Rcon contains BattleEye RCON connection settings.
type Rcon struct {
	IP               string `yaml:"ip,omitempty" env:"IP, default=127.0.0.1"`
	Password         string `yaml:"password" env:"PASSWORD"`
	KeepaliveTimeout int    `yaml:"keepalive_timeout,omitempty" env:"KEEPALIVE_TIMEOUT, default=30"`
	DeadlineTimeout  int    `yaml:"deadline_timeout,omitempty" env:"DEADLINE_TIMEOUT, default=5"`
	Port             int    `yaml:"port,omitempty" env:"PORT, default=2305"`
	BufferSize       uint16 `yaml:"buffer_size,omitempty" env:"BUFFER_SIZE, default=1024"`
	Bans             bool   `yaml:"expose_bans,omitempty" env:"EXPOSE_BANS, default=false"`
}

// Logging contains configuration for log output.
type Logging struct {
	Level     string `yaml:"level,omitempty" env:"LEVEL, default=info"`
	Format    string `yaml:"format,omitempty" env:"FORMAT, default=text"`
	Output    string `yaml:"output,omitempty" env:"OUTPUT, default=stdout"`
	NoMetrics bool   `yaml:"metrics_disabled,omitempty" env:"METRICS_DISABLED"`
	NoHealth  bool   `yaml:"health_disabled,omitempty" env:"HEALTH_DISABLED"`
}

// config loader
func loadConfig() (*Config, error) {
	var config Config

	// initial prepare log level from env for debuging purpose
	if lvl := os.Getenv("DAYZ_EXPORTER_LOG_LEVEL"); lvl != "" {
		if logLevel, err := zerolog.ParseLevel(lvl); err == nil {
			log.Logger = log.Level(logLevel)
		}
	} else {
		log.Logger = log.Level(zerolog.InfoLevel)
	}

	// load config from file if is exists
	if path, ok := getConfigPath(); ok {
		configFile, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer func() {
			if err := configFile.Close(); err != nil {
				log.Error().Msgf("Cant close config file")
			}
		}()

		decoder := yaml.NewDecoder(configFile)
		err = decoder.Decode(&config)
		if err != nil {
			return nil, err
		}
	}

	// load environment variables
	ctx := context.Background()
	if err := envconfig.Process(ctx, &config); err != nil {
		return nil, err
	}

	config.setupLogging()

	if config.Rcon.Password == "" {
		return nil, errors.New("missing required RCON password")
	}

	if config.GeoDB != "" {
		log.Trace().Str("file", config.GeoDB).Msg("Try find Geo DB file")
		if _, err := os.Stat(config.GeoDB); err != nil {
			log.Warn().Str("file", config.GeoDB).Msg("Cant open GeoDB file")
			config.GeoDB = ""
		}
	}

	log.Trace().Any("config", config).Msg("Config loaded")
	return &config, nil
}

// get path to configuration file from variables, argument or use default
func getConfigPath() (string, bool) {
	if path := os.Getenv("DAYZ_EXPORTER_CONFIG_PATH"); path != "" {
		log.Trace().Str("file", path).Msg("Use config form variable DAYZ_EXPORTER_CONFIG_PATH")
		return path, true
	}
	if len(os.Args) > 1 {
		log.Trace().Str("file", os.Args[1]).Msg("Use config form argument")
		return os.Args[1], true
	}
	if _, err := os.Stat(defaultConfigPath); err == nil {
		log.Trace().Str("file", defaultConfigPath).Msg("Use default config path")
		return defaultConfigPath, true
	}

	log.Trace().Msg("Config file not provided and not found default, work with defaults and environment variables")
	return "", false
}

// configure logging
func (c *Config) setupLogging() {
	// setup log level
	if logLevel, err := zerolog.ParseLevel(c.Logging.Level); err == nil {
		log.Trace().Str("level", c.Logging.Level).Msg("Setup log level")
		log.Logger = log.Level(logLevel)
	} else {
		log.Warn().Msgf("Log level %s unknown, fallback to Info level", c.Logging.Level)
		log.Logger = log.Level(zerolog.InfoLevel)
	}

	// setup log output
	var writer io.Writer
	switch c.Logging.Output {
	case "stdout", "out", "1":
		writer = os.Stdout
	case "stderr", "err", "2":
		writer = os.Stderr
	default:
		file, err := os.OpenFile(c.Logging.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatal().Err(err).Str("file", c.Logging.Output).Msg("Failed to open log file")
		}
		writer = file
	}

	// check need enable colors
	var useColors bool
	if f, ok := writer.(*os.File); ok {
		// fd is smallest numbers
		useColors = term.IsTerminal(int(f.Fd())) // #nosec G115
	} else {
		useColors = false
	}

	// setup log format
	switch c.Logging.Format {
	case "text":
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        writer,
			TimeFormat: time.RFC3339,
			NoColor:    !useColors,
		})
	case "json":
		fallthrough
	default:
		log.Logger = log.Output(writer)
	}
}
