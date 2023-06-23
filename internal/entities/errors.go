package entities

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InvalidLink = status.Error(codes.InvalidArgument, "invalid link")
	NotFound    = status.Error(codes.NotFound, "element not found")
	ServerError = status.Error(codes.Internal, "server error")

	MissingStorageTypeError   = errors.New("missing storage type")
	IncorrectPsqlStorageError = errors.New("psql authentication data is incorrect")
	RunError                  = errors.New("internal error, shutting down")
)
