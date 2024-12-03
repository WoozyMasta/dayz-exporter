package bemetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/woozymasta/bercon-cli/pkg/beparser"
)

// InitBansMetrics инициализирует метрики для банов.
func (mc *MetricsCollector) InitBansMetrics() {
	labels := mc.customLabels.Keys()

	if mc.banGuidTimeMetric == nil {
		mc.banGuidTimeMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "bercon_ban_guid_time_seconds",
				Help: "Time left for GUID bans in seconds.",
			},
			append(labels, "reason", "guid"),
		)
	}

	if mc.banGuidTotal == nil {
		mc.banGuidTotal = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "bercon_ban_guid_total",
				Help: "Total count of GUID bans.",
			},
			labels,
		)
	}

	if mc.banIpTimeMetric == nil {
		mc.banIpTimeMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "bercon_ban_ip_time_seconds",
				Help: "Time left for IP bans in seconds.",
			},
			append(labels, "reason", "ip", "country"),
		)
	}

	if mc.banIpTotal == nil {
		mc.banIpTotal = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "bercon_ban_ip_total",
				Help: "Total count of IP bans.",
			},
			labels,
		)
	}
}

// UpdateBansMetrics обновляет метрики для банов (GUID и IP).
func (mc *MetricsCollector) UpdateBansMetrics(bans *beparser.Bans) {
	values := mc.customLabels.Values()

	// Обновляем метрики для GUID банов
	if mc.banGuidTimeMetric != nil {
		mc.banGuidTimeMetric.Reset() // количество метрик меняется, каждый раз сбрасываем

		for _, ban := range bans.GUIDBans {
			mc.banGuidTimeMetric.WithLabelValues(append(values, ban.Reason, ban.GUID)...).Set(banSeconds(ban.MinutesLeft))
		}
	}

	if mc.banGuidTotal != nil {
		mc.banGuidTotal.WithLabelValues(values...).Set(float64(len(bans.GUIDBans)))
	}

	// Обновляем метрики для IP банов
	if mc.banIpTimeMetric != nil {
		mc.banIpTimeMetric.Reset() // количество метрик меняется, каждый раз сбрасываем

		for _, ban := range bans.IPBans {
			mc.banIpTimeMetric.WithLabelValues(append(values, ban.Reason, ban.IP, ban.Country)...).Set(banSeconds(ban.MinutesLeft))
		}
	}

	if mc.banIpTotal != nil {
		mc.banIpTotal.WithLabelValues(values...).Set(float64(len(bans.IPBans)))
	}
}

func banSeconds(minutes int) float64 {
	if minutes > 0 {
		return float64(minutes * 60)
	}
	return float64(minutes)
}
