package main

import (
	"errors"
	"github.com/jyolando/test-ozon-go/internal/entities"
	"github.com/jyolando/test-ozon-go/internal/handlers"
	"log"
)

const (
	RunError                = "run error"
	MissingStorageTypeError = "missing storage type"
)

func main() {
	var errInfo string
	if server, err := handlers.NewServer(); err != nil {
		if errors.As(err, &entities.MissingStorageTypeError{}) {
			errInfo = MissingStorageTypeError
		}
		log.Fatal(errInfo)
	} else if err := server.Run(); err != nil {
		log.Fatal(RunError)
	}

}
