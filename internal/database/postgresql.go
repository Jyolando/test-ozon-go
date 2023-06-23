package database

import (
	"context"
	"github.com/jackc/pgx"
	"github.com/jyolando/test-ozon-go/internal/entities"
	api "github.com/jyolando/test-ozon-go/pkg/api/proto"
	"github.com/jyolando/test-ozon-go/pkg/helpers"
	log "github.com/sirupsen/logrus"
)

type PsqlStorage struct {
	pool *pgx.ConnPool

	logger *log.Entry
}

func NewPsqlStorage(l *log.Logger) (*PsqlStorage, error) {
	config, err := helpers.ParsePsqlConfig()
	if err != nil {
		return nil, err
	}

	logger := l.WithField("storage", "postgresql")

	if pool, err := pgx.NewConnPool(*config); err != nil {
		return nil, err
	} else {
		return &PsqlStorage{
			pool:   pool,
			logger: logger,
		}, nil
	}
}

func (p *PsqlStorage) GetStorageType() string {
	return "postgresql"
}

func (p *PsqlStorage) checkURLExists(originalLink string) (*api.AddURLResponse, error) {
	qry := `
		SELECT short_link
		FROM links
		WHERE original_link = $1
	`

	conn, err := p.pool.Acquire()
	if err != nil {
		return nil, entities.ServerError
	}
	defer p.pool.Release(conn)

	var shortLink string
	err = conn.QueryRow(qry, originalLink).Scan(&shortLink)
	if err != nil {
		return nil, err
	} else {
		return &api.AddURLResponse{Url: &api.ShortenedURL{ShortenedURL: shortLink, OriginalURL: originalLink}}, nil
	}
}

func (p *PsqlStorage) AddURL(ctx context.Context, request *api.AddURLRequest) (*api.AddURLResponse, error) {
	qry := `
		INSERT INTO links
		(original_link, short_link)
		VALUES ($1, $2)
	`

	if link, err := p.checkURLExists(request.GetUrl()); err == nil {
		p.logger.WithFields(log.Fields{
			"request":  request,
			"response": link,
			"code":     0,
		}).Info("addUrl success")

		return link, nil
	}

	conn, err := p.pool.Acquire()
	if err != nil {
		return nil, entities.ServerError
	}
	defer p.pool.Release(conn)

	if shortLink, err := p.getSecureToken(10); err != nil {
		return nil, entities.ServerError
	} else if _, err := conn.Exec(qry, request.GetUrl(), shortLink); err != nil {
		return nil, entities.ServerError
	} else {
		response := &api.AddURLResponse{Url: &api.ShortenedURL{ShortenedURL: shortLink, OriginalURL: request.GetUrl()}}

		p.logger.WithFields(log.Fields{
			"request":  request,
			"response": response,
			"code":     0,
		}).Info("addUrl success")

		return response, nil
	}
}

func (p *PsqlStorage) GetURL(ctx context.Context, request *api.GetURLRequest) (*api.GetURLResponse, error) {
	qry := `
		SELECT original_link
		FROM links
		WHERE short_link = $1
	`

	conn, err := p.pool.Acquire()
	if err != nil {
		return nil, entities.ServerError
	}
	defer p.pool.Release(conn)

	var originalLink string
	err = conn.QueryRow(qry, request.GetUrl()).Scan(&originalLink)
	if err != nil {
		return nil, entities.NotFound
	} else {
		response := &api.GetURLResponse{Url: &api.ShortenedURL{ShortenedURL: request.GetUrl(), OriginalURL: originalLink}}

		p.logger.WithFields(log.Fields{
			"request":  request,
			"response": response,
			"code":     0,
		}).Info("getUrl success")

		return response, nil
	}
}

func (p *PsqlStorage) getSecureToken(length int) (string, error) {
	qry := `
		SELECT original_link
		FROM links
		WHERE short_link = $1
	`

	for {
		conn, err := p.pool.Acquire()
		if err != nil {
			return "", entities.ServerError
		}

		token, err := helpers.GenToken(length)
		if err != nil {
			return "", err
		}

		var link string
		err = conn.QueryRow(qry, token).Scan(&link)
		if err != nil {
			return token, nil
		}

		p.pool.Release(conn)
	}
}
