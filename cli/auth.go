package main

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// Middleware for Basic Auth
func basicAuthMiddleware(mux http.Handler, config Listen) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// exclude probes
		if !config.HealthAuth && strings.HasPrefix(r.URL.Path, "/health") {
			mux.ServeHTTP(w, r)
			return
		}

		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Basic ") {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		payload, err := base64.StdEncoding.DecodeString(auth[len("Basic "):])
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 || pair[0] != config.Username || pair[1] != config.Password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		mux.ServeHTTP(w, r)
	})
}
