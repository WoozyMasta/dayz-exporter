package bemetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Отвечает за хранение различных метрик
type MetricsCollector struct {
	playerPingMetric    *prometheus.GaugeVec
	playersTotal        *prometheus.GaugeVec
	playersOnline       *prometheus.GaugeVec
	playersLobby        *prometheus.GaugeVec
	playersInvalid      *prometheus.GaugeVec
	banGuidTimeMetric   *prometheus.GaugeVec
	banGuidTotal        *prometheus.GaugeVec
	banIpTimeMetric     *prometheus.GaugeVec
	banIpTotal          *prometheus.GaugeVec
	serverPlayersOnline *prometheus.GaugeVec
	serverPlayersSlots  *prometheus.GaugeVec
	serverPlayersQueue  *prometheus.GaugeVec
	serverTime          *prometheus.GaugeVec
	customLabels        Labels
}

// Создает пустой экземпляр MetricsCollector
func NewMetricsCollector(customLabels Labels) *MetricsCollector {
	return &MetricsCollector{
		customLabels: customLabels,
	}
}

// Метод отдает все метрики из структуры MetricsCollector
func (mc *MetricsCollector) getAllMetrics() []prometheus.Collector {
	return []prometheus.Collector{
		// players
		mc.playerPingMetric,
		mc.playersTotal,
		mc.playersOnline,
		mc.playersInvalid,
		mc.playersLobby,
		// bans
		mc.banGuidTimeMetric,
		mc.banGuidTotal,
		mc.banIpTimeMetric,
		mc.banIpTotal,
		// server
		mc.serverPlayersOnline,
		mc.serverPlayersSlots,
		mc.serverPlayersQueue,
		mc.serverTime,
	}
}

// Регистрирует только инициализированные метрики
func (mc *MetricsCollector) RegisterMetrics() {
	for _, metric := range mc.getAllMetrics() {
		if gaugeVec, ok := metric.(*prometheus.GaugeVec); ok {
			if gaugeVec != nil {
				prometheus.MustRegister(gaugeVec)
			}
		}
	}
}

// Сбрасывает все инициализированные метрики
func (mc *MetricsCollector) ResetMetrics() {
	for _, metric := range mc.getAllMetrics() {
		if gaugeVec, ok := metric.(*prometheus.GaugeVec); ok {
			if gaugeVec != nil {
				gaugeVec.Reset()
			}
		}
	}
}
