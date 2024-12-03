package logging

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/sethvargo/go-envconfig"
	log "github.com/sirupsen/logrus"
)

type LogConfig struct {
	Level  string `yaml:"level,omitempty" env:"LOG_LEVEL, default=info"`
	Format string `yaml:"format,omitempty" env:"LOG_FORMAT, default=text"`
	Output string `yaml:"output,omitempty" env:"LOG_OUTPUT, default=stderr"`
}

var logInit sync.Once

// Инициализирует логгер на основе предоставленной конфигурации.
func InitLogger(config *LogConfig) {
	logInit.Do(func() {
		SetLevel(config.Level)
		SetFormat(config.Format)
		SetOutput(config.Output)
	})
}

// Инициализирует логгер, загружая конфигурацию из переменных окружения.
func InitLoggerEnv(ctx context.Context) {
	config := &LogConfig{}
	if err := envconfig.Process(ctx, config); err != nil {
		log.Fatalf("Failed to load environment config: %v", err)
	}
	InitLogger(config)
}

// задает уровень логирования
func SetLevel(lvl string) {
	logLevel, err := log.ParseLevel(lvl)
	if err != nil {
		log.Warnf("Invalid log level '%s', defaulting to 'info'", lvl)
		logLevel = log.InfoLevel
	}

	if logLevel == log.TraceLevel {
		log.SetReportCaller(true)
	}

	log.SetLevel(logLevel)
}

// задает формат логирования
func SetFormat(format string) {
	var logFormat log.Formatter

	switch strings.ToLower(format) {
	case "json", "struct", "structured":
		logFormat = &log.JSONFormatter{}
	default:
		logFormat = &log.TextFormatter{}
	}

	log.SetFormatter(logFormat)
}

// открывает и задает файл для записи логов
func SetOutput(destination string) error {
	var output *os.File

	switch strings.ToLower(destination) {
	case "stdout", "out", "console":
		output = os.Stdout

	case "stderr", "err", "":
		output = os.Stderr

	default:
		var err error
		output, err = os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
		if err != nil {
			return fmt.Errorf("cannot open log file '%s': %v", destination, err)
		}
	}

	log.SetOutput(output)
	return nil
}
