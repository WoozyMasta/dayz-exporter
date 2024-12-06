package main

import (
	"fmt"
	"net/http"

	"github.com/oschwald/geoip2-golang"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rumblefrog/go-a2s"
	log "github.com/sirupsen/logrus"
	"github.com/woozymasta/bercon-cli/pkg/beparser"
	"github.com/woozymasta/bercon-cli/pkg/bercon"
	"github.com/woozymasta/dayz-exporter/pkg/bemetrics"
)

type connection struct {
	rcon      *bercon.Connection          // connection to BattleEye RCON server
	query     *a2s.Client                 // connection to A2S Steam Query
	collector *bemetrics.MetricsCollector // metrics collector
	geo       *geoip2.Reader              // reader for geoip DB
	info      *a2s.ServerInfo             // server information
	bans      bool                        // flag for enable/disable bans metrics
}

// create connection manager
func setupConnection(cfg *Config) (*connection, error) {
	// create connection to BattleEye RCON
	rcon, err := bercon.Open(fmt.Sprintf("%s:%d", cfg.Rcon.IP, cfg.Rcon.Port), cfg.Rcon.Password)
	if err != nil {
		return nil, err
	}

	// setup connection
	if cfg.Rcon.KeepaliveTimeout != 0 {
		rcon.SetKeepaliveTimeout(cfg.Rcon.KeepaliveTimeout)
	}
	rcon.SetDeadlineTimeout(cfg.Rcon.DeadlineTimeout)
	rcon.SetBufferSize(cfg.Rcon.BufferSize)

	// start keepalive for BattleEye RCON connections
	rcon.StartKeepAlive()

	// create connection to Steam A2S Query
	query, err := a2s.NewClient(fmt.Sprintf("%s:%d", cfg.Query.IP, cfg.Query.Port))
	if err != nil {
		return nil, err
	}
	info, err := query.QueryInfo()
	if err != nil {
		return nil, err
	}

	log.Infof("Connected to \"%s\" (%s) server on %s:%d", info.Name, info.Map, cfg.Query.IP, cfg.Query.Port)

	// create bemetrics metrics collector
	collector := bemetrics.NewMetricsCollector(makeLabels(info, cfg.Labels))

	var geoDB *geoip2.Reader
	if cfg.GeoDB != "" {
		geoDB, err = geoip2.Open(cfg.GeoDB)
		if err != nil {
			log.Errorf("Cant open GeoDB %e", err)
		}
		log.Traceln("GeoDB loaded success")
	}

	// init connection structure
	connection := connection{
		rcon:      rcon,
		query:     query,
		collector: collector,
		bans:      cfg.Rcon.Bans,
		info:      info,
		geo:       geoDB,
	}

	// initialize metrics
	connection.collector.InitServerMetrics()
	connection.collector.InitPlayerMetrics()
	if cfg.Rcon.Bans {
		connection.collector.InitBansMetrics()
	}

	// register metrics
	connection.collector.RegisterMetrics()

	return &connection, nil
}

// get and update server metrics from Steam A2S Query
func (c *connection) updateServerMetrics() error {
	info, err := c.query.QueryInfo()
	if err != nil {
		return fmt.Errorf("failed to get A2S info query response: %w", err)
	}

	c.collector.UpdateServerMetrics(info)

	log.Traceln("metrics updated: server A2S info")
	return nil
}

// get and update players metrics from BattleEye RCON
func (c *connection) updatePlayersMetrics() error {
	data, err := c.rcon.Send("players")
	if err != nil {
		return fmt.Errorf("failed to send 'players' command: %w", err)
	}

	playersData := beparser.Parse(data, "players")
	if players, ok := playersData.(*beparser.Players); ok {
		if c.geo != nil {
			players.SetCountryCode(c.geo)
		}
		c.collector.UpdatePlayerMetrics(players)
	} else {
		return fmt.Errorf("unexpected data type for 'players' response")
	}

	log.Traceln("metrics updated: players data")
	return nil
}

// get and update bans metrics from BattleEye RCON
func (c *connection) updateBansMetrics() error {
	if !c.bans {
		return nil
	}

	_, err := c.rcon.Send("loadBans")
	if err != nil {
		return fmt.Errorf("failed to send 'loadBans' command: %w", err)
	}

	data, err := c.rcon.Send("bans")
	if err != nil {
		return fmt.Errorf("failed to send 'bans' command: %w", err)
	}

	bansData := beparser.Parse(data, "bans")
	if bans, ok := bansData.(*beparser.Bans); ok {
		if c.geo != nil {
			bans.SetCountryCode(c.geo)
		}
		c.collector.UpdateBansMetrics(bans)
	} else {
		return fmt.Errorf("unexpected data type for 'bans' response")
	}

	log.Traceln("metrics updated: bans data")
	return nil
}

// http handler for update metrics for each request
func (c *connection) metricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := c.updateServerMetrics(); err != nil {
			c.handleError(w, err, "server")
			return
		}
		if err := c.updatePlayersMetrics(); err != nil {
			c.handleError(w, err, "player")
			return
		}
		if err := c.updateBansMetrics(); err != nil {
			c.handleError(w, err, "bans")
			return
		}

		log.Debugf("Metrics updated")
		// pass control over the standard handler prometheus
		promhttp.Handler().ServeHTTP(w, r)
	}
}

// error handler in http connection
func (c *connection) handleError(w http.ResponseWriter, err error, context string) {
	c.collector.ResetMetrics()
	c.query.Close()
	c.rcon.Close()
	c.geo.Close()

	http.Error(w, fmt.Sprintf("Error updating metrics (%s)", context), http.StatusInternalServerError)
	log.WithError(err).Fatalf("Failed to update metrics (%s)", context)
}
