package main

import (
	"context"
	"os"

	"github.com/rumblefrog/go-a2s"
	"github.com/sethvargo/go-envconfig"
	log "github.com/sirupsen/logrus"
	"github.com/woozymasta/dayz-exporter/pkg/bemetrics"
	"github.com/woozymasta/dayz-exporter/pkg/logging"
	"gopkg.in/yaml.v3"
)

const defaultConfigPath = "config.yaml"

// Главная структура конфига.
type Config struct {
	Listen  Listen            `yaml:"listen,omitempty" env:", prefix=DAYZ_EXPORTER_LISTEN_"`
	Query   Query             `yaml:"query,omitempty" env:", prefix=DAYZ_EXPORTER_QUERY_"`
	Rcon    Rcon              `yaml:"rcon,omitempty" env:", prefix=DAYZ_EXPORTER_RCON_"`
	Labels  map[string]string `yaml:"labels,omitempty" env:"DAYZ_EXPORTER_LABELS"`
	GeoDB   string            `yaml:"geo_db,omitempty" env:"DAYZ_EXPORTER_GEOIP_DB"`
	Logging logging.LogConfig `yaml:"logging,omitempty" env:", prefix=DAYZ_EXPORTER_"`
}

// Структура для глобальных настроек (settings).
type Listen struct {
	IP       string `yaml:"ip,omitempty" env:"IP, default=0.0.0.0"`
	Port     uint16 `yaml:"port,omitempty" env:"PORT, default=8098"`
	Endpoint string `yaml:"endpoint,omitempty" env:"ENDPOINT, default=/metrics"`
}

type Query struct {
	IP   string `yaml:"ip,omitempty" env:"IP, default=127.0.0.1"`
	Port uint16 `yaml:"port,omitempty" env:"PORT, default=27016"`
}

type Rcon struct {
	IP               string `yaml:"ip,omitempty" env:"IP, default=127.0.0.1"`
	Port             uint16 `yaml:"port,omitempty" env:"PORT, default=2305"`
	Password         string `yaml:"password" env:"PASSWORD"`
	Bans             bool   `yaml:"expose_bans,omitempty" env:"EXPOSE_BANS, default=false"`
	BufferSize       int    `yaml:"buffer_size,omitempty" env:"BUFFER_SIZE, default=1024"`
	KeepaliveTimeout int    `yaml:"keepalive_timeout,omitempty" env:"KEEPALIVE_TIMEOUT, default=30"`
	DeadlineTimeout  int    `yaml:"deadline_timeout,omitempty" env:"DEADLINE_TIMEOUT, default=5"`
}

// Функция для загрузки конфига.
func loadConfig() (*Config, error) {
	var config Config

	// пытаемся загрузить конфигурацию из файла, если файл существует
	if path, ok := getConfigPath(); ok {
		configFile, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer configFile.Close()

		decoder := yaml.NewDecoder(configFile)
		err = decoder.Decode(&config)
		if err != nil {
			return nil, err
		}
	}

	// загружаем переменные окружения
	ctx := context.Background()
	if err := envconfig.Process(ctx, &config); err != nil {
		return nil, err
	}

	if config.Rcon.Password == "" {
		log.Fatalf("Missing required RCON password")
	}

	logging.InitLogger(&config.Logging)

	if config.GeoDB != "" {
		log.Tracef("Try find Geo DB file: %s", config.GeoDB)
		if _, err := os.Stat(config.GeoDB); err != nil {
			log.Warnf("Cant find GeoDB file '%s'", config.GeoDB)
			config.GeoDB = ""
		}
	}

	log.Tracef("Loaded config: %+v", config)
	return &config, nil
}

// Функция для получения пути к конфигу (переменная среды или аргументы).
func getConfigPath() (string, bool) {
	if path := os.Getenv("DAYZ_EXPORTER_CONFIG_PATH"); path != "" {
		return path, true
	}
	if len(os.Args) > 1 {
		return os.Args[1], true
	}
	if _, err := os.Stat(defaultConfigPath); err == nil {
		return defaultConfigPath, true
	}

	return "", false
}

func makeLabels(info *a2s.ServerInfo, extraLabels map[string]string) bemetrics.Labels {
	game := info.Game
	if game == "" {
		game = info.Folder
	}

	labels := []bemetrics.Label{
		{Key: "server", Value: info.Name},
		{Key: "map", Value: info.Map},
		{Key: "game", Value: game},
		{Key: "os", Value: info.ServerOS.String()},
		{Key: "version", Value: info.Version},
	}

	for k, v := range extraLabels {
		labels = append(labels, bemetrics.Label{Key: k, Value: v})
	}

	return labels
}
