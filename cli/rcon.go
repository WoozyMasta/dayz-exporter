package main

import (
	"fmt"
	"net/http"

	"github.com/oschwald/geoip2-golang"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/woozymasta/a2s/pkg/a2s"
	"github.com/woozymasta/bercon-cli/pkg/beparser"
	"github.com/woozymasta/bercon-cli/pkg/bercon"
	"github.com/woozymasta/dayz-exporter/pkg/bemetrics"
)

type connection struct {
	rcon       *bercon.Connection          // connection to BattleEye RCON server
	query      *a2s.Client                 // connection to A2S Steam Query
	collector  *bemetrics.MetricsCollector // metrics collector
	geo        *geoip2.Reader              // reader for geoip DB
	info       *a2s.Info                   // server information
	bans       bool                        // flag for enable/disable bans metrics
	exposeInfo bool                        // flag for enable/disable /info json endpoint
}

// create connection manager
func setupConnection(cfg *Config) (*connection, error) {
	// create connection to BattleEye RCON
	rcon, err := bercon.Open(fmt.Sprintf("%s:%d", cfg.Rcon.IP, cfg.Rcon.Port), cfg.Rcon.Password)
	if err != nil {
		return nil, fmt.Errorf("open RCON connection: %v", err)
	}

	rconVersion, err := rcon.Send("version")
	if err != nil {
		return nil, fmt.Errorf("get RCON version: %v", err)
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
	query, err := a2s.New(cfg.Query.IP, cfg.Query.Port)
	if err != nil {
		return nil, fmt.Errorf("open A2S connection: %v", err)
	}
	info, err := query.GetInfo()
	if err != nil {
		return nil, fmt.Errorf("get A2S_INFO: %v", err)
	}

	log.Info().
		Str("ip", cfg.Query.IP).
		Int("query port", cfg.Query.Port).
		Str("version", info.Version).
		Str("name", info.Name).
		Str("map", info.Map).
		Int("rcon port", cfg.Rcon.Port).
		Str("rcon version", string(rconVersion)).
		Msg("Connected to server")

	if info.ID != 221100 && info.ID != 1024020 {
		log.Error().Msg("Game ID on the server does not match DayZ, this looks like a configuration issue")
	}

	// create bemetrics metrics collector
	collector := bemetrics.NewMetricsCollector(makeLabels(info, cfg.Labels))

	var geoDB *geoip2.Reader
	if cfg.GeoDB != "" {
		geoDB, err = geoip2.Open(cfg.GeoDB)
		if err != nil {
			return nil, fmt.Errorf("open GeoIP DB: %v", err)
		}
		log.Trace().Msgf("GeoDB loaded success")
	}

	// init connection structure
	connection := connection{
		rcon:       rcon,
		query:      query,
		collector:  collector,
		bans:       cfg.Rcon.Bans,
		info:       info,
		geo:        geoDB,
		exposeInfo: cfg.Listen.ExposeInfo,
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
	info, err := c.query.GetInfo()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get A2S info query response")
		return err
	}

	c.collector.UpdateServerMetrics(info)
	c.info = info
	log.Trace().Msg("Server A2S metrics updated")

	return nil
}

// get and update players metrics from BattleEye RCON
func (c *connection) updatePlayersMetrics() error {
	data, err := c.rcon.Send("players")
	if err != nil {
		log.Error().Err(err).Msg("Failed to send 'players' command")
		return err
	}

	playersData := beparser.Parse(data, "players")
	if players, ok := playersData.(*beparser.Players); ok {
		if c.geo != nil {
			players.SetCountryCode(c.geo)
		}
		c.collector.UpdatePlayerMetrics(players)
		log.Trace().Msg("Player metrics updated")

		return nil
	}

	log.Warn().Msg("Unexpected data type for 'players' response")
	return fmt.Errorf("unexpected data type for 'players' response")
}

// get and update bans metrics from BattleEye RCON
func (c *connection) updateBansMetrics() error {
	if !c.bans {
		log.Debug().Msg("Ban metrics disabled, skipping update")
		return nil
	}

	_, err := c.rcon.Send("loadBans")
	if err != nil {
		log.Error().Err(err).Msg("Failed to send 'loadBans' command")
		return err
	}

	data, err := c.rcon.Send("bans")
	if err != nil {
		log.Error().Err(err).Msg("Failed to send 'bans' command")
		return err
	}

	bansData := beparser.Parse(data, "bans")
	if bans, ok := bansData.(*beparser.Bans); ok {
		if c.geo != nil {
			bans.SetCountryCode(c.geo)
		}
		c.collector.UpdateBansMetrics(bans)
		log.Trace().Msg("Bans metrics updated")
		return nil
	}

	log.Warn().Msg("Unexpected data type for 'bans' response")
	return fmt.Errorf("unexpected data type for 'bans' response")
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

		log.Debug().Msgf("Metrics updated")
		// pass control over the standard handler prometheus
		promhttp.Handler().ServeHTTP(w, r)
	}
}

// error handler in http connection
func (c *connection) handleError(w http.ResponseWriter, err error, context string) {
	if err != nil {
		log.Error().Err(err).Str("context", context).Msg("Error updating metrics")
	}

	defer func() {
		log.Debug().Msg("Resetting metrics and closing connections")
		c.collector.ResetMetrics()

		if err := c.query.Close(); err != nil {
			log.Error().Msg("Cant close query connection")
		}
		if err := c.rcon.Close(); err != nil {
			log.Error().Msg("Cant close rcon connection")
		}
		if c.geo != nil {
			if err := c.geo.Close(); err != nil {
				log.Error().Msg("Cant close geo ip database file")
			}
		}
	}()

	http.Error(w, fmt.Sprintf("Error updating metrics (%s)", context), http.StatusInternalServerError)
	log.Fatal().Msgf("Failed to update metrics (%s)", context)
}
