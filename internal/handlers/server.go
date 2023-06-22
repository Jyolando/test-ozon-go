package handlers

import (
	"context"
	"github.com/jyolando/test-ozon-go/internal/database"
	"github.com/jyolando/test-ozon-go/internal/entities"
	api "github.com/jyolando/test-ozon-go/pkg/api/proto"
	"github.com/jyolando/test-ozon-go/pkg/helpers"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"os"
)

var log = logrus.New()

type Server struct {
	api.URLShortenerServer
	storage entities.Database
}

func NewServer() (*Server, error) {
	var err error

	mode := os.Getenv("STORAGE_TYPE")
	if len(mode) == 0 {
		return nil, entities.MissingStorageTypeError{}
	}

	server := &Server{}

	if mode == "memory" {
		server.storage = database.NewMemoryStorage()
	} else if mode == "postgresql" {
		if server.storage, err = database.NewPsqlStorage(); err != nil {
			return nil, entities.IncorrectPsqlStorage{}
		}
	} else {
		return nil, entities.MissingStorageTypeError{}
	}

	return server, nil
}

func (s *Server) AddURL(ctx context.Context, in *api.AddURLRequest) (*api.AddURLResponse, error) {
	log.Info(in)

	if originalLink := in.Url; !helpers.IsURL(originalLink) || len(originalLink) == 0 {
		return nil, status.Error(codes.InvalidArgument, entities.HTTP400)
	} else {
		if links, err := s.storage.AddURL(ctx, in); err != nil {
			helpers.HandleErrors(err, log)
		} else {
			return links, nil
		}
	}

	return nil, status.Error(codes.NotFound, entities.HTTP500)
}

func (s *Server) GetURL(ctx context.Context, request *api.GetURLRequest) (*api.GetURLResponse, error) {
	log.Info(request)

	if links, err := s.storage.GetURL(ctx, request); err != nil {
		helpers.HandleErrors(err, log)
	} else {
		return links, nil
	}

	return nil, status.Error(codes.NotFound, entities.HTTP500)
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return entities.RunError{}
	}

	grpcServer := grpc.NewServer()
	api.RegisterURLShortenerServer(grpcServer, s)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return entities.RunError{}
	}
	return nil
}
