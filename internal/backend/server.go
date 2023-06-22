package backend

import (
	"context"
	"github.com/jyolando/test-ozon-go/internal/entities"
	api "github.com/jyolando/test-ozon-go/pkg/api/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
)

var log = logrus.New()

type Server struct {
	api.URLShortenerServer
}

func NewServer() (*Server, error) {
	mode := os.Getenv("STORAGE_TYPE")
	if len(mode) == 0 {
		return nil, entities.MissingStorageTypeError{}
	}

	server := &Server{}

	if mode == "memory" {
		log.Info("memory storage")
	} else if mode == "postgresql" {
		log.Info("psql storage")
	} else {
		return nil, entities.MissingStorageTypeError{}
	}

	return server, nil
}

func (s *Server) AddURL(ctx context.Context, in *api.AddURLRequest) (*api.AddURLResponse, error) {
	return &api.AddURLResponse{Url: &api.ShortenedURL{OriginalURL: in.GetUrl(), ShortenedURL: in.GetUrl() + "lol"}}, nil
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return entities.RunError{}
	}

	grpcServer := grpc.NewServer()
	api.RegisterURLShortenerServer(grpcServer, &Server{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return entities.RunError{}
	}
	return nil
}
