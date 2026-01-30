package metric

import "github.com/prometheus/client_golang/prometheus"

func newImageProcessesTotal() *prometheus.CounterVec {
	return prometheus.V2.NewCounterVec(prometheus.CounterVecOpts{
		CounterOpts: prometheus.CounterOpts{
			Namespace: "imageer",
			Subsystem: "processor",
			Name:      "image_processes_total",
			Help:      "Total count of image processing operations",
		},
		VariableLabels: prometheus.UnconstrainedLabels{
			"success",
			"input_format",
			"output_format",
		},
	})
}

func newImageProcessDurationSeconds() *prometheus.HistogramVec {
	return prometheus.V2.NewHistogramVec(prometheus.HistogramVecOpts{
		HistogramOpts: prometheus.HistogramOpts{
			Namespace: "imageer",
			Subsystem: "processor",
			Name:      "image_process_duration_seconds",
			Help:      "Duration of image processing operations in seconds",
			Buckets:   []float64{.01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 25},
		},
		VariableLabels: prometheus.UnconstrainedLabels{
			"success",
			"input_format",
			"output_format",
		},
	})
}
