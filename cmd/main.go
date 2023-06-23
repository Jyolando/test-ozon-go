package main

import (
	"github.com/jyolando/test-ozon-go/internal/handlers"
	"github.com/jyolando/test-ozon-go/pkg/helpers"
)

const (
	RunError = "run error"
)

func main() {
	//var errInfo string

	logger := helpers.NewLogger()

	logger.Info("starting")
	if server, err := handlers.NewServer(logger); err != nil {
		logger.Error(err)
	} else if err := server.Run(); err != nil {
		logger.Fatal(RunError)
	}

}
