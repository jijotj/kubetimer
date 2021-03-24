package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	Buckets = []float64{
		1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610,
	}

	PodScheduledTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "pod_scheduled_time",
		Help:    "Time taken to schedule pods",
		Buckets: Buckets,
	}, []string{"namespace", "service"})
	PodInitTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "pod_init_time",
		Help:    "Time taken for pods to initialize",
		Buckets: Buckets,
	}, []string{"namespace", "service"})
	PodContainersReadyTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "pod_containers_ready_time",
		Help:    "Time taken for containers in pods to be ready",
		Buckets: Buckets,
	}, []string{"namespace", "service"})
	PodReadyTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "pod_ready_time",
		Help:    "Time taken for pods to be ready",
		Buckets: Buckets,
	}, []string{"namespace", "service"})
)
