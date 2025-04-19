package bemetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/woozymasta/bercon-cli/pkg/beparser"
)

// InitBansMetrics initialize bercon ban metrics
func (mc *MetricsCollector) InitBansMetrics() {
	labels := mc.customLabels.Keys()

	if mc.banGUIDTimeMetric == nil {
		mc.banGUIDTimeMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "bercon_ban_guid_time_seconds",
				Help: "Time left for GUID bans in seconds.",
			},
			append(labels, "reason", "guid"),
		)
	}

	if mc.banGUIDTotal == nil {
		mc.banGUIDTotal = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "bercon_ban_guid_total",
				Help: "Total count of GUID bans.",
			},
			labels,
		)
	}

	if mc.banIPTimeMetric == nil {
		mc.banIPTimeMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "bercon_ban_ip_time_seconds",
				Help: "Time left for IP bans in seconds.",
			},
			append(labels, "reason", "ip", "country"),
		)
	}

	if mc.banIPTotal == nil {
		mc.banIPTotal = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "bercon_ban_ip_total",
				Help: "Total count of IP bans.",
			},
			labels,
		)
	}
}

// UpdateBansMetrics use for update ban metrics (GUID and IP)
func (mc *MetricsCollector) UpdateBansMetrics(bans *beparser.Bans) {
	values := mc.customLabels.Values()

	// update GUID bans
	if mc.banGUIDTimeMetric != nil {
		mc.banGUIDTimeMetric.Reset() // count of metrics is dynamic, reset it always

		for _, ban := range bans.GUIDBans {
			mc.banGUIDTimeMetric.WithLabelValues(append(values, ban.Reason, ban.GUID)...).Set(banSeconds(ban.MinutesLeft))
		}
	}

	if mc.banGUIDTotal != nil {
		mc.banGUIDTotal.WithLabelValues(values...).Set(float64(len(bans.GUIDBans)))
	}

	// update IP bans
	if mc.banIPTimeMetric != nil {
		mc.banIPTimeMetric.Reset() // count of metrics is dynamic, reset it always

		for _, ban := range bans.IPBans {
			mc.banIPTimeMetric.WithLabelValues(append(values, ban.Reason, ban.IP, ban.Country)...).Set(banSeconds(ban.MinutesLeft))
		}
	}

	if mc.banIPTotal != nil {
		mc.banIPTotal.WithLabelValues(values...).Set(float64(len(bans.IPBans)))
	}
}

func banSeconds(minutes int) float64 {
	if minutes > 0 {
		return float64(minutes * 60)
	}
	return float64(minutes)
}
