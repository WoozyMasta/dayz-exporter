package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/woozymasta/dayz-exporter/pkg/service"
)

func main() {
	if service.IsServiceMode() {
		service.RunAsService(runApp)
		return
	}

	runApp()
}

// load config, init connections and serve metrics
func runApp() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Load config failed: %s", err)
	}

	connection, err := setupConnection(config)
	if err != nil {
		log.Fatalf("Failed to establish connections to RCON port %s:%d", config.Rcon.IP, config.Rcon.Port)
	}

	// handle metrics
	http.Handle(config.Listen.Endpoint, connection.metricsHandler())

	// handle probes
	http.HandleFunc("/", connection.rootHandler)
	http.HandleFunc("/health", connection.livenessHandler)
	http.HandleFunc("/health/liveness", connection.livenessHandler)
	http.HandleFunc("/health/readiness", connection.readinessHandler)

	// serve
	addr := fmt.Sprintf("%s:%d", config.Listen.IP, config.Listen.Port)
	log.Infof("Starting metrics server at %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
