package client

import (
	"go.uber.org/zap"

	"github.com/hotstar/kubetimer/pkg/controller"
	"github.com/hotstar/kubetimer/pkg/handler"
)

// Run runs the event loop processing with given handler
func Run() {
	zap.S().Info("Starting")
	var eventHandler = handler.NewTimer()
	controller.Start(eventHandler)
}
