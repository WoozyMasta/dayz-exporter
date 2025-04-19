package main

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

// Middleware for Basic Auth
func basicAuthMiddleware(next http.Handler, config Listen) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// exclude probes if config.HealthAuth == false
		if !config.HealthAuth && strings.HasPrefix(r.URL.Path, "/health") {
			log.Trace().Msg("Health check endpoint, skipping auth")
			next.ServeHTTP(w, r)
			return
		}

		if !config.InfoAuth && strings.HasPrefix(r.URL.Path, "/info") {
			log.Trace().Msg("Info endpoint, skipping auth")
			next.ServeHTTP(w, r)
			return
		}

		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Basic ") {
			log.Warn().Msg("No or invalid 'Authorization' header, returning 401")
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		payload, err := base64.StdEncoding.DecodeString(auth[len("Basic "):])
		if err != nil {
			log.Error().Err(err).Msg("Failed to decode base64 authorization header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			log.Debug().Msg("Malformed auth header, returning 401")
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		username, password := pair[0], pair[1]
		if username != config.Username || password != config.Password {
			log.Warn().
				Str("username", username).
				Msg("Authorization failed, wrong username or password")
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Debug().Msg("User authorized successfully")
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler, domains string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", domains)
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
