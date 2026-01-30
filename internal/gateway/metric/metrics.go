package metric

import "github.com/prometheus/client_golang/prometheus"

func newHTTPRequestsTotal() *prometheus.CounterVec {
	return prometheus.V2.NewCounterVec(prometheus.CounterVecOpts{
		CounterOpts: prometheus.CounterOpts{
			Namespace: "imageer",
			Subsystem: "gateway",
			Name:      "http_requests_total",
			Help:      "Total count of HTTP requests",
		},
		VariableLabels: prometheus.UnconstrainedLabels{
			"method",
			"path",
			"status",
		},
	})
}

func newHTTPRequestDurationSeconds() *prometheus.HistogramVec {
	return prometheus.V2.NewHistogramVec(prometheus.HistogramVecOpts{
		HistogramOpts: prometheus.HistogramOpts{
			Namespace: "imageer",
			Subsystem: "gateway",
			Name:      "http_request_duration_seconds",
			Help:      "Duration of HTTP requests in seconds",
			Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 25},
		},
		VariableLabels: prometheus.UnconstrainedLabels{
			"method",
			"path",
			"status",
		},
	})
}
