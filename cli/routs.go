package main

import (
	_ "embed"
	"internal/vars"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

//go:embed style.min.css
var styleCSS []byte

//go:generate minify style.css -o style.min.css

const docsURL = "https://woozymasta.github.io/dayz-exporter/"

// check is alive
func (c *connection) livenessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Debug().Str("method", r.Method).Msg("Method not allowed on liveness")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if !c.rcon.IsAlive() {
		log.Warn().Msg("BattleEye RCON not connected")
		http.Error(w, "BattleEye RCON not connected", http.StatusServiceUnavailable)
		return
	}

	log.Trace().Msg("Liveness check OK")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Error().Err(err).Msg("Error writing liveness response")
	}
}

// simple OK if up and ready to handle requests
func (c *connection) readinessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Debug().Str("method", r.Method).Msg("Method not allowed on readiness")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	log.Trace().Msg("Readiness check OK")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Error().Err(err).Msg("Error writing readiness response")
	}
}

func (c *connection) rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Debug().Str("method", r.Method).Msg("Method not allowed on index page")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var serverInfo strings.Builder
	serverInfo.WriteString("server: " + c.info.Name + "\n")
	if c.info.Game != "" {
		serverInfo.WriteString("description: " + c.info.Game + "\n")
	}
	serverInfo.WriteString(
		"map: " + c.info.Map + "\n" +
			"game: " + c.info.Folder + "\n" +
			"os: " + c.info.Environment.String() + "\n" +
			"version: " + c.info.Version,
	)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>DayZ Exporter</title>
	<style>` + string(styleCSS) + `</style>
</head>
<body>
	<div class="container">
		<h1>DayZ Exporter</h1>
		<h3>Collects and publishes Prometheus metrics from Battleye RCON and Steam A2S Query for DayZ server.</h3>
		<p>Source code: <a href="` + vars.URL + `" target="_blank">` + vars.URL + `</a></p>
		<p>Documentation: <a href="` + docsURL + `" target="_blank">` + docsURL + `</a></p>
		<hr/>
		<p>This application exposes the following endpoints:</p>
		<ul>
			<li><a href="/metrics">/metrics</a>: Exposes Prometheus metrics.</li>
			<li><a href="/health">/health</a>: General health check of the service;</li>
			<li><a href="/health/liveness">/health/liveness</a>: Checks if the service is alive (RCON connection);</li>
			<li><a href="/health/readiness">/health/readiness</a>: Checks if the service is ready (all required connections are established);</li>
		</ul>
		<hr/>
		<p>Game server information:</p>
		<pre>
` + serverInfo.String() + `
		</pre>
		<p>Exporter information:</p>
		<pre>
version: ` + vars.Version + `
commit: ` + vars.Commit + `
built: ` + vars.BuildTime + `
		</pre>
	</div>
</body>
</html>
`))

	if err != nil {
		log.Error().Msgf("index page: %v", err)
	}
}
