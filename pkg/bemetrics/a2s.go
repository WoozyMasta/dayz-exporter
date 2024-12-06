package bemetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rumblefrog/go-a2s"
	"github.com/woozymasta/dayz-exporter/pkg/gametype"
)

// initialize a2s server metrics
func (mc *MetricsCollector) InitServerMetrics() {
	labels := mc.customLabels.Keys()

	if mc.serverPlayersOnline == nil {
		mc.serverPlayersOnline = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "a2s_info_players_online",
				Help: "Online players.",
			},
			labels,
		)
	}

	if mc.serverPlayersSlots == nil {
		mc.serverPlayersSlots = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "a2s_info_players_slots",
				Help: "Players slots count.",
			},
			labels,
		)
	}

	if mc.serverPlayersQueue == nil {
		mc.serverPlayersQueue = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "a2s_info_players_queue",
				Help: "Players wait in queue.",
			},
			labels,
		)
	}

	if mc.serverTime == nil {
		mc.serverTime = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "a2s_info_time",
				Help: "Duration of day time on server.",
			},
			labels,
		)
	}
}

// update a2s server metrics
func (mc *MetricsCollector) UpdateServerMetrics(serverInfo *a2s.ServerInfo) {
	var extendedInfo gametype.DayZ
	extendedInfo.Parse(serverInfo.ExtendedServerInfo.Keywords)

	values := mc.customLabels.Values()

	if mc.serverPlayersOnline != nil {
		mc.serverPlayersOnline.WithLabelValues(values...).Set(float64(serverInfo.Players))
	}

	if mc.serverPlayersSlots != nil {
		mc.serverPlayersSlots.WithLabelValues(values...).Set(float64(serverInfo.MaxPlayers))
	}

	if mc.serverPlayersQueue != nil {
		mc.serverPlayersQueue.WithLabelValues(values...).Set(float64(extendedInfo.PlayersQueue))
	}

	if mc.serverTime != nil {
		mc.serverTime.WithLabelValues(values...).Set(float64(extendedInfo.Time))
	}
}
