package metrics

import "github.com/prometheus/client_golang/prometheus"
import "github.com/prometheus/client_golang/prometheus/promauto"

var m *MetricSet

type MetricSet struct {
	StartUpTime              *prometheus.HistogramVec
	ConfiguredUpstreams      prometheus.Gauge
	HealthyUpstreams         prometheus.Gauge
	EthSyncingLatency        *prometheus.HistogramVec
	EthGetBlockNumberLatency *prometheus.HistogramVec
	EthGetBalanceLatency     *prometheus.HistogramVec
}

func Metrics() *MetricSet {
	if m == nil {
		m = &MetricSet{
			StartUpTime: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:        "startup_time_ms",
					Help:        "time from program start until the server is ready to handle requests",
					ConstLabels: nil,
					Buckets:     prometheus.DefBuckets,
				}, nil),
			ConfiguredUpstreams: promauto.NewGauge(
				prometheus.GaugeOpts{
					Name: "upstreams_configured",
					Help: "configured upstream endpoints",
				},
			),
			HealthyUpstreams: promauto.NewGauge(
				prometheus.GaugeOpts{
					Name: "upstreams_healthy",
					Help: "currently healthy upstream endpoints",
				},
			),
			EthSyncingLatency: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "latency_eth_syncing",
					Help:    "time from program start until the server is ready to handle requests",
					Buckets: prometheus.DefBuckets,
				}, nil,
			),
			EthGetBlockNumberLatency: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "latency_eth_get_block_number",
					Help:    "time from program start until the server is ready to handle requests",
					Buckets: prometheus.DefBuckets,
				}, nil,
			),
			EthGetBalanceLatency: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "latency_eth_get_balance",
					Help:    "time from program start until the server is ready to handle requests",
					Buckets: prometheus.DefBuckets,
				}, nil,
			),
		}
	}
	return m
}
