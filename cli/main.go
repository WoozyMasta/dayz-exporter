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
	parseArgs()

	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Load config failed: %s", err)
	}

	connection, err := setupConnection(config)
	if err != nil {
		log.Fatalf("Failed to establish connections to RCON port %s:%d", config.Rcon.IP, config.Rcon.Port)
	}

	// init router
	mux := http.NewServeMux()

	// handle metrics
	mux.Handle(config.Listen.Endpoint, connection.metricsHandler())

	// handle probes
	mux.HandleFunc("/", connection.rootHandler)
	mux.HandleFunc("/health", connection.livenessHandler)
	mux.HandleFunc("/health/liveness", connection.livenessHandler)
	mux.HandleFunc("/health/readiness", connection.readinessHandler)

	// add auth middleware if password set
	var handler http.Handler = mux
	if config.Listen.Password != "" {
		handler = basicAuthMiddleware(mux, config.Listen)
	}

	// serve
	addr := fmt.Sprintf("%s:%d", config.Listen.IP, config.Listen.Port)
	log.Infof("Starting metrics server at %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
