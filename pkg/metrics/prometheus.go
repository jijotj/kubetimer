package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var prometheusMetricsSync sync.Once

func NewRegistry() prometheus.Gatherer {
	registry := prometheus.NewRegistry()
	prometheusMetricsSync.Do(func() {
		registry.MustRegister(PodScheduledTime)
		registry.MustRegister(PodInitTime)
		registry.MustRegister(PodContainersReadyTime)
		registry.MustRegister(PodReadyTime)
	})
	return registry
}
