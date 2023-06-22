package database

import (
	"context"
	"github.com/jyolando/test-ozon-go/internal/entities"
	api "github.com/jyolando/test-ozon-go/pkg/api/proto"
	"github.com/jyolando/test-ozon-go/pkg/helpers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MemoryStorage struct {
	originalAsKey map[string]string
	shortAsKey    map[string]string
}

func (m *MemoryStorage) AddURL(ctx context.Context, request *api.AddURLRequest) (*api.AddURLResponse, error) {
	var (
		response = &api.AddURLResponse{}
		err      error
	)

	hashOriginalLink := helpers.GetMD5Hash(request.GetUrl())
	if savedLink, ok := m.originalAsKey[hashOriginalLink]; ok {
		response, err = &api.AddURLResponse{Url: &api.ShortenedURL{OriginalURL: request.GetUrl(), ShortenedURL: savedLink}}, nil
	} else if shortLink, err := helpers.GenToken(); err == nil {
		hashShortLink := helpers.GetMD5Hash(shortLink)
		m.originalAsKey[hashOriginalLink] = shortLink
		m.shortAsKey[hashShortLink] = request.GetUrl()
		response, err = &api.AddURLResponse{Url: &api.ShortenedURL{OriginalURL: request.GetUrl(), ShortenedURL: shortLink}}, nil
	} else {
		response, err = nil, status.Error(codes.Internal, entities.HTTP500)
	}

	return response, err
}

func (m *MemoryStorage) GetURL(ctx context.Context, request *api.GetURLRequest) (*api.GetURLResponse, error) {
	var (
		response = &api.GetURLResponse{}
		err      error
	)

	hashShortLink := helpers.GetMD5Hash(request.GetUrl())
	if originalLink, ok := m.shortAsKey[hashShortLink]; ok {
		response, err = &api.GetURLResponse{Url: &api.ShortenedURL{OriginalURL: originalLink, ShortenedURL: request.Url}}, nil
	} else {
		response, err = nil, status.Error(codes.NotFound, entities.HTTP404)
	}

	return response, err
}

func NewMemoryStorage() *MemoryStorage {
	ms := &MemoryStorage{}
	ms.originalAsKey = make(map[string]string)
	ms.shortAsKey = make(map[string]string)

	return ms
}
