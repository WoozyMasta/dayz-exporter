package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	parseArgs()

	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Load config failed: %s", err)
	}

	connection, err := setupConnection(config)
	if err != nil {
		log.Fatalf("Failed to establish connections to RCON port %s:%d", config.Rcon.IP, config.Rcon.Port)
	}

	serve(config, connection)
}

// Запуск HTTP сервера для экспорта метрик
func serve(config *Config, conn *connection) {
	http.Handle(config.Listen.Endpoint, conn.metricsHandler())
	addr := fmt.Sprintf("%s:%d", config.Listen.IP, config.Listen.Port)
	log.Infof("Starting metrics server at %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
