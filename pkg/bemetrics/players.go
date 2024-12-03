package bemetrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/woozymasta/bercon-cli/pkg/beparser"
)

// InitPlayerMetrics инициализирует метрики для игроков.
func (mc *MetricsCollector) InitPlayerMetrics() {
	labels := mc.customLabels.Keys()

	if mc.playerPingMetric == nil {
		mc.playerPingMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "bercon_player_ping_seconds",
				Help: "Ping of players in seconds.",
			},
			append(labels, "name", "ip", "guid", "lobby", "country"),
		)
	}

	if mc.playersTotal == nil {
		mc.playersTotal = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "bercon_players_total",
				Help: "Total count of players.",
			},
			labels,
		)
	}

	if mc.playersOnline == nil {
		mc.playersOnline = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "bercon_players_online",
				Help: "Count of players online.",
			},
			labels,
		)
	}

	if mc.playersLobby == nil {
		mc.playersLobby = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "bercon_players_lobby",
				Help: "Count of players in lobby.",
			},
			labels,
		)
	}

	if mc.playersInvalid == nil {
		mc.playersInvalid = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "bercon_players_invalid",
				Help: "Count of invalid players.",
			},
			labels,
		)
	}
}

// UpdatePlayerMetrics обновляет метрики для списка игроков.
func (mc *MetricsCollector) UpdatePlayerMetrics(players *beparser.Players) {
	values := mc.customLabels.Values()

	if mc.playerPingMetric != nil {
		mc.playerPingMetric.Reset() // количество метрик меняется, каждый раз сбрасываем

		for _, player := range *players {
			lobby := fmt.Sprintf("%t", player.Lobby)
			playerLabels := append(values, player.Name, player.IP, player.GUID, lobby, player.Country)
			mc.playerPingMetric.WithLabelValues(playerLabels...).Set(float64(player.Ping))
		}
	}

	online, lobby, invalid := countPlayers(*players)

	if mc.playersTotal != nil {
		mc.playersTotal.WithLabelValues(values...).Set(float64(len(*players)))
	}

	if mc.playersOnline != nil {
		mc.playersOnline.WithLabelValues(values...).Set(online)
	}

	if mc.playersLobby != nil {
		mc.playersLobby.WithLabelValues(values...).Set(lobby)
	}

	if mc.playersInvalid != nil {
		mc.playersInvalid.WithLabelValues(values...).Set(invalid)
	}
}

func countPlayers(players []beparser.Player) (float64, float64, float64) {
	var online, lobby, invalid float64

	for _, player := range players {
		if !player.Valid {
			invalid++
		} else if player.Lobby {
			lobby++
		} else {
			online++
		}
	}

	return online, lobby, invalid
}
