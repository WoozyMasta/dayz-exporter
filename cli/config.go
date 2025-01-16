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

// main config structure
type Config struct {
	Labels  map[string]string `yaml:"labels,omitempty" env:"DAYZ_EXPORTER_LABELS"`
	Logging Logging           `yaml:"logging,omitempty" env:", prefix=DAYZ_EXPORTER_LOG_"`
	GeoDB   string            `yaml:"geo_db,omitempty" env:"DAYZ_EXPORTER_GEOIP_DB"`
	Listen  Listen            `yaml:"listen,omitempty" env:", prefix=DAYZ_EXPORTER_LISTEN_"`
	Query   Query             `yaml:"query,omitempty" env:", prefix=DAYZ_EXPORTER_QUERY_"`
	Rcon    Rcon              `yaml:"rcon,omitempty" env:", prefix=DAYZ_EXPORTER_RCON_"`
}

// listen settings for exporter
type Listen struct {
	IP         string `yaml:"ip,omitempty" env:"IP, default=0.0.0.0"`
	Endpoint   string `yaml:"endpoint,omitempty" env:"ENDPOINT, default=/metrics"`
	Username   string `yaml:"username,omitempty" env:"USERNAME, default=metrics"`
	Password   string `yaml:"password,omitempty" env:"PASSWORD"`
	Port       uint16 `yaml:"port,omitempty" env:"PORT, default=8098"`
	HealthAuth bool   `yaml:"health_auth,omitempty" env:"HEALTH_AUTH, default=false"`
}

// Steam A2S query connection settings
type Query struct {
	IP   string `yaml:"ip,omitempty" env:"IP, default=127.0.0.1"`
	Port int    `yaml:"port,omitempty" env:"PORT, default=27016"`
}

// BattleEye RCON connection settings
type Rcon struct {
	IP               string `yaml:"ip,omitempty" env:"IP, default=127.0.0.1"`
	Password         string `yaml:"password" env:"PASSWORD"`
	KeepaliveTimeout int    `yaml:"keepalive_timeout,omitempty" env:"KEEPALIVE_TIMEOUT, default=30"`
	DeadlineTimeout  int    `yaml:"deadline_timeout,omitempty" env:"DEADLINE_TIMEOUT, default=5"`
	Port             int    `yaml:"port,omitempty" env:"PORT, default=2305"`
	BufferSize       uint16 `yaml:"buffer_size,omitempty" env:"BUFFER_SIZE, default=1024"`
	Bans             bool   `yaml:"expose_bans,omitempty" env:"EXPOSE_BANS, default=false"`
}

// Steam A2S query connection settings
type Logging struct {
	Level  string `yaml:"level,omitempty" env:"LEVEL, default=info"`
	Format string `yaml:"format,omitempty" env:"FORMAT, default=text"`
	Output string `yaml:"output,omitempty" env:"OUTPUT, default=stdout"`
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
		log.Trace().Msgf("Try find Geo DB file: %s", config.GeoDB)
		if _, err := os.Stat(config.GeoDB); err != nil {
			log.Warn().Msgf("Cant open GeoDB file '%s'", config.GeoDB)
			config.GeoDB = ""
		}
	}

	log.Trace().Msgf("Loaded config: %+v", config)
	return &config, nil
}

// get path to configuration file from variables, argument or use default
func getConfigPath() (string, bool) {
	if path := os.Getenv("DAYZ_EXPORTER_CONFIG_PATH"); path != "" {
		log.Trace().Msgf("Use config form variable DAYZ_EXPORTER_CONFIG_PATH: %s", path)
		return path, true
	}
	if len(os.Args) > 1 {
		log.Trace().Msgf("Use config form argument: %s", os.Args[1])
		return os.Args[1], true
	}
	if _, err := os.Stat(defaultConfigPath); err == nil {
		log.Trace().Msgf("Use default config path: %s", defaultConfigPath)
		return defaultConfigPath, true
	}

	log.Trace().Msg("Config file not provided and not found default, work with defaults and environment variables")
	return "", false
}

// configure logging
func (c *Config) setupLogging() {
	// setup log level
	if logLevel, err := zerolog.ParseLevel(c.Logging.Level); err == nil {
		log.Trace().Msgf("Setup log level to: %s", c.Logging.Level)
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
		file, err := os.OpenFile(c.Logging.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to open log file: %s", c.Logging.Output)
		}
		writer = file
	}

	// check need enable colors
	var useColors bool
	if f, ok := writer.(*os.File); ok {
		useColors = term.IsTerminal(int(f.Fd()))
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
