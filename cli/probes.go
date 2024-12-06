package main

import (
	"net/http"

	"github.com/woozymasta/dayz-exporter/pkg/config"
)

const docsURL = "https://woozymasta.github.io/dayz-exporter/"

// check is alive
func (c *connection) livenessHandler(w http.ResponseWriter, r *http.Request) {
	if !c.rcon.IsAlive() {
		http.Error(w, "BattleEye RCON not connected", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// simple OK if up and ready to handle requests
func (c *connection) readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (c *connection) rootHandler(w http.ResponseWriter, r *http.Request) {
	game := c.info.Game
	if game == "" {
		game = c.info.Folder
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>DayZ Exporter</title>
</head>
<body>
	<h1>DayZ Exporter</h1>
	<h3>Collects and publishes Prometheus metrics from Battleye RCON and Steam A2S Query for DayZ server.</h3>
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
	server: ` + c.info.Name + `
	map: ` + c.info.Map + `
	game: ` + game + `
	os: ` + c.info.ServerOS.String() + `
	version: ` + c.info.Version + `
	</pre>
	<p>Exporter information:</p>
	<pre>
	version: ` + config.Version + `
	commit: ` + config.Commit + `
	built: ` + config.BuildTime + `
	</pre>
</body>
</html>
`))
}
