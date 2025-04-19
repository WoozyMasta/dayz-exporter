package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/hlog"
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

	// create mux
	mux := http.NewServeMux()

	// handle metrics
	mux.Handle(config.Listen.Endpoint, connection.metricsHandler())

	// handle probes
	mux.HandleFunc("/", connection.rootHandler)
	mux.HandleFunc("/health", connection.livenessHandler)
	mux.HandleFunc("/health/liveness", connection.livenessHandler)
	mux.HandleFunc("/health/readiness", connection.readinessHandler)

	if config.Listen.ExposeInfo {
		mux.HandleFunc("/info", connection.infoHandler)
	}

	var handler http.Handler = mux

	// enable CORS
	if config.Listen.CORSDomains != "" {
		handler = corsMiddleware(handler, config.Listen.CORSDomains)
	}

	// add basic auth if password is set
	if config.Listen.Password != "" {
		handler = basicAuthMiddleware(handler, config.Listen)
	}

	// wrap all with zerolog/hlog
	// hlog.NewHandler -> hlog.AccessHandler -> hlog.RemoteAddrHandler -> ...
	handler = hlog.NewHandler(log.Logger)(
		hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			log.Info().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Str("remote", r.RemoteAddr).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Msg("HTTP request completed")
		})(
			hlog.RemoteAddrHandler("ip")(
				hlog.UserAgentHandler("user_agent")(
					hlog.RefererHandler("referer")(handler),
				),
			),
		),
	)

	addr := fmt.Sprintf("%s:%d", config.Listen.IP, config.Listen.Port)
	log.Info().Msgf("Starting metrics server at %s", addr)

	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 3 * time.Second,
		IdleTimeout:       60 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadTimeout:       5 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Metrics server failed")
	}
}
