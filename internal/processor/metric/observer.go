package metric

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var globalObserver observer = noopObserver{}

func Gatherer() prometheus.Gatherer {
	if o, ok := globalObserver.(*client); ok {
		return o.gatherer
	}
	return nil
}

func ObserveImageProcess(inputFormat, outputFormat string, success bool, duration time.Duration) {
	globalObserver.observeImageProcess(inputFormat, outputFormat, success, duration)
}

type observer interface {
	observeImageProcess(inputFormat, outputFormat string, success bool, duration time.Duration)
}

type noopObserver struct{}

func (noopObserver) observeImageProcess(inputFormat, outputFormat string, success bool,
	duration time.Duration) {
}
