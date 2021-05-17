package main

import (
	"net/http"

	"github.com/hotstar/kubetimer/pkg/client"
	"github.com/hotstar/kubetimer/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	go startPromServer()
	client.Run()
}

func startPromServer() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(1)
	}
	defer logger.Sync()

	undo := zap.ReplaceGlobals(logger)
	defer undo()

	registry := prometheus.NewRegistry()
	registry.MustRegister(metrics.PodContainersReadyTime)
	http.Handle("/metrics", promhttp.HandlerFor(metrics.NewRegistry(), promhttp.HandlerOpts{}))
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {w.WriteHeader(http.StatusOK)})
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(1)
	}
}
