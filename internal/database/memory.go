package database

import (
	"context"
	"github.com/jyolando/test-ozon-go/internal/entities"
	api "github.com/jyolando/test-ozon-go/pkg/api/proto"
	"github.com/jyolando/test-ozon-go/pkg/helpers"
	log "github.com/sirupsen/logrus"
	"sync"
)

type MemoryStorage struct {
	originalAsKey map[string]string
	shortAsKey    map[string]string

	logger *log.Entry

	sync.RWMutex
}

func NewMemoryStorage(l *log.Logger) *MemoryStorage {
	logger := l.WithField("storage", "memory")

	return &MemoryStorage{
		originalAsKey: make(map[string]string),
		shortAsKey:    make(map[string]string),
		logger:        logger,
	}
}

func (m *MemoryStorage) GetStorageType() string {
	return "memory"
}

func (m *MemoryStorage) AddURL(ctx context.Context, request *api.AddURLRequest) (*api.AddURLResponse, error) {
	var (
		response = &api.AddURLResponse{}
		err      error
	)

	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()

	hashOriginalLink := helpers.GetMD5Hash(request.GetUrl())
	if savedLink, ok := m.originalAsKey[hashOriginalLink]; ok {
		response, err = &api.AddURLResponse{Url: &api.ShortenedURL{OriginalURL: request.GetUrl(), ShortenedURL: savedLink}}, nil
	} else if shortLink, err := helpers.GenToken(10); err == nil {
		hashShortLink := helpers.GetMD5Hash(shortLink)
		m.originalAsKey[hashOriginalLink] = shortLink
		m.shortAsKey[hashShortLink] = request.GetUrl()
		response, err = &api.AddURLResponse{Url: &api.ShortenedURL{OriginalURL: request.GetUrl(), ShortenedURL: shortLink}}, nil
	} else {
		response, err = nil, entities.ServerError
	}

	if response != nil {
		m.logger.WithFields(log.Fields{
			"request":  request,
			"response": response,
			"code":     0,
		}).Info("addUrl request status: OK")
	}
	return response, err
}

func (m *MemoryStorage) GetURL(ctx context.Context, request *api.GetURLRequest) (*api.GetURLResponse, error) {
	var (
		response = &api.GetURLResponse{}
		err      error
	)

	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()

	hashShortLink := helpers.GetMD5Hash(request.GetUrl())
	if originalLink, ok := m.shortAsKey[hashShortLink]; ok {
		response, err = &api.GetURLResponse{Url: &api.ShortenedURL{OriginalURL: originalLink, ShortenedURL: request.Url}}, nil
	} else {
		response, err = nil, entities.NotFound
	}

	if response != nil {
		m.logger.WithFields(log.Fields{
			"request":  request,
			"response": response,
			"code":     0,
		}).Info("getUrl request status: OK")
	}
	return response, err
}
