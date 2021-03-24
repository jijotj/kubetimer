package main

import (
	"github.com/hotstar/kubetimer/pkg/client"
	"github.com/hotstar/kubetimer/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main() {
	go startPromServer()
	client.Run()
}

func startPromServer() {
	registry := prometheus.NewRegistry()
	registry.MustRegister(metrics.PodContainersReadyTime)
	http.Handle("/metrics", promhttp.HandlerFor(metrics.NewRegistry(), promhttp.HandlerOpts{}))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(1)
	}
}
