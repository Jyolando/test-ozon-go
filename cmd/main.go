package main

import (
	"github.com/jyolando/test-ozon-go/internal/entities"
	"github.com/jyolando/test-ozon-go/internal/handlers"
	"github.com/jyolando/test-ozon-go/pkg/helpers"
)

func main() {
	logger := helpers.NewLogger()

	logger.Info("starting")
	if server, err := handlers.NewServer(logger); err != nil {
		logger.Error(err)
	} else if err := server.Run(); err != nil {
		logger.Fatal(entities.RunError)
	}
}
