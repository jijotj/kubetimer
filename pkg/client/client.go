package client

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/hotstar/kubetimer/pkg/controller"
	"github.com/hotstar/kubetimer/pkg/handler"
)

// Run runs the event loop processing with given handler
func Run() {
	zap.S().Info("Starting")
	fmt.Print("Starting\n")
	var eventHandler = handler.NewTimer()
	controller.Start(eventHandler)
}
