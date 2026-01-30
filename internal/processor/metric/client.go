package metric

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type client struct {
	gatherer prometheus.Gatherer

	imageProcessesTotal         *prometheus.CounterVec
	imageProcessDurationSeconds *prometheus.HistogramVec
}

func Init() {
	c := &client{
		gatherer:                    prometheus.DefaultGatherer,
		imageProcessesTotal:         newImageProcessesTotal(),
		imageProcessDurationSeconds: newImageProcessDurationSeconds(),
	}

	prometheus.MustRegister(c.imageProcessesTotal)
	prometheus.MustRegister(c.imageProcessDurationSeconds)

	globalObserver = c
}

func (c *client) observeImageProcess(inputFormat, outputFormat string, success bool,
	duration time.Duration,
) {
	successStr := strconv.FormatBool(success)
	c.imageProcessesTotal.WithLabelValues(successStr, inputFormat, outputFormat).Inc()
	c.imageProcessDurationSeconds.WithLabelValues(successStr, inputFormat, outputFormat).
		Observe(duration.Seconds())
}
