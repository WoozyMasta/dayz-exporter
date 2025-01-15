package bemetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// responsible for storing various metrics
type MetricsCollector struct {
	playerPingMetric    *prometheus.GaugeVec
	playersTotal        *prometheus.GaugeVec
	playersOnline       *prometheus.GaugeVec
	playersLobby        *prometheus.GaugeVec
	playersInvalid      *prometheus.GaugeVec
	banGUIDTimeMetric   *prometheus.GaugeVec
	banGUIDTotal        *prometheus.GaugeVec
	banIPTimeMetric     *prometheus.GaugeVec
	banIPTotal          *prometheus.GaugeVec
	serverPing          *prometheus.GaugeVec
	serverPlayersOnline *prometheus.GaugeVec
	serverPlayersSlots  *prometheus.GaugeVec
	serverPlayersQueue  *prometheus.GaugeVec
	serverTime          *prometheus.GaugeVec
	customLabels        Labels
}

// creates an empty MetricsCollector instance
func NewMetricsCollector(customLabels Labels) *MetricsCollector {
	return &MetricsCollector{
		customLabels: customLabels,
	}
}

// returns all metrics from the MetricsCollector structure
func (mc *MetricsCollector) getAllMetrics() []prometheus.Collector {
	return []prometheus.Collector{
		// players
		mc.playerPingMetric,
		mc.playersTotal,
		mc.playersOnline,
		mc.playersInvalid,
		mc.playersLobby,
		// bans
		mc.banGUIDTimeMetric,
		mc.banGUIDTotal,
		mc.banIPTimeMetric,
		mc.banIPTotal,
		// server
		mc.serverPing,
		mc.serverPlayersOnline,
		mc.serverPlayersSlots,
		mc.serverPlayersQueue,
		mc.serverTime,
	}
}

// register only initialized metrics
func (mc *MetricsCollector) RegisterMetrics() {
	for _, metric := range mc.getAllMetrics() {
		if gaugeVec, ok := metric.(*prometheus.GaugeVec); ok {
			if gaugeVec != nil {
				prometheus.MustRegister(gaugeVec)
			}
		}
	}
}

// resets all initialized metrics
func (mc *MetricsCollector) ResetMetrics() {
	for _, metric := range mc.getAllMetrics() {
		if gaugeVec, ok := metric.(*prometheus.GaugeVec); ok {
			if gaugeVec != nil {
				gaugeVec.Reset()
			}
		}
	}
}
