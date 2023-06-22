package main

import (
	"errors"
	"github.com/jyolando/test-ozon-go/internal/backend"
	"github.com/jyolando/test-ozon-go/internal/entities"
	"log"
)

const (
	RunError                = "run error"
	MissingStorageTypeError = "missing storage type"
)

func main() {
	var errInfo string
	//lis, err := net.Listen("tcp", "localhost:9000")
	//if err != nil {
	//	log.Fatalf("failed to listen: %v", err)
	//}
	//
	//s := grpc.NewServer()
	//api.RegisterURLShortenerServer(s, &backend.Server{})
	//if err := s.Serve(lis); err != nil {
	//	log.Fatalf("failed to serve: %v", err)
	//}

	if server, err := backend.NewServer(); err != nil {
		if errors.As(err, &entities.MissingStorageTypeError{}) {
			errInfo = MissingStorageTypeError
		}
		log.Fatal(errInfo)
	} else if err := server.Run(); err != nil {
		log.Fatal(RunError)
	}

}
