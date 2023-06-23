package handlers

import (
	"context"
	"github.com/jyolando/test-ozon-go/internal/database"
	"github.com/jyolando/test-ozon-go/internal/entities"
	api "github.com/jyolando/test-ozon-go/pkg/api/proto"
	"github.com/jyolando/test-ozon-go/pkg/helpers"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
)

var (
	port = ":9000"
)

type Server struct {
	api.URLShortenerServer
	storage entities.Database
	logger  *log.Logger
}

func NewServer(logger *log.Logger) (*Server, error) {
	var err error
	mode := os.Getenv("STORAGE_TYPE")
	if len(mode) == 0 {
		return nil, entities.MissingStorageTypeError
	}

	server := &Server{
		logger: logger,
	}

	if mode == "memory" {
		server.storage = database.NewMemoryStorage(logger)
		logger.Info("storage type: in-memory")
	} else if mode == "postgresql" {
		if server.storage, err = database.NewPsqlStorage(logger); err != nil {
			return nil, entities.IncorrectPsqlStorageError
		}
		logger.Info("Storage: Postgresql")
	} else {
		return nil, entities.MissingStorageTypeError
	}

	return server, nil
}

func (s *Server) AddURL(ctx context.Context, in *api.AddURLRequest) (*api.AddURLResponse, error) {
	if originalLink := in.Url; !helpers.IsURL(originalLink) || len(originalLink) == 0 {
		return nil, entities.NotFound
	} else {
		if links, err := s.storage.AddURL(ctx, in); err != nil {
			helpers.HandleErrors(err, s.logger, s.storage.GetStorageType(), in)
		} else {
			return links, nil
		}
	}

	return nil, entities.NotFound
}

func (s *Server) GetURL(ctx context.Context, request *api.GetURLRequest) (*api.GetURLResponse, error) {
	if links, err := s.storage.GetURL(ctx, request); err != nil {
		helpers.HandleErrors(err, s.logger, s.storage.GetStorageType(), request)
	} else {
		return links, nil
	}

	return nil, entities.ServerError
}

func (s *Server) Run() error {
	s.logger.Info("starting on 127.0.0.1" + port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		s.logger.Error("failed to net listen: %v", err)
		return entities.RunError
	}
	s.logger.Info("started on 127.0.0.1" + port)

	grpcServer := grpc.NewServer()
	api.RegisterURLShortenerServer(grpcServer, s)
	s.logger.Info("starting grpc server")
	if err := grpcServer.Serve(lis); err != nil {
		s.logger.Error("failed to grpc serve: %v", err)
		return entities.RunError
	}
	return nil
}
