package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/woozymasta/dayz-exporter/internal/service"
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
		log.Fatal().Msgf("Load config failed: %s", err)
	}

	connection, err := setupConnection(config)
	if err != nil {
		log.Fatal().Msgf("Failed to establish connections to RCON port %s:%d", config.Rcon.IP, config.Rcon.Port)
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
	log.Info().Msgf("Starting metrics server at %s", addr)

	// create HTTP-server with timeouts
	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 3 * time.Second,
		IdleTimeout:       60 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadTimeout:       5 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err)
	}
}
