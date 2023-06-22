package database

import (
	"context"
	"github.com/jackc/pgx"
	"github.com/jyolando/test-ozon-go/internal/entities"
	api "github.com/jyolando/test-ozon-go/pkg/api/proto"
	"github.com/jyolando/test-ozon-go/pkg/helpers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PsqlStorage struct {
	pool *pgx.ConnPool
}

func NewPsqlStorage() (*PsqlStorage, error) {
	config, err := helpers.ParsePsqlConfig()
	if err != nil {
		return nil, err
	}

	if pool, err := pgx.NewConnPool(*config); err != nil {
		return nil, err
	} else {
		return &PsqlStorage{pool: pool}, nil
	}
}

func (p *PsqlStorage) checkURLExists(originalLink string) (*api.AddURLResponse, error) {
	qry := `
		SELECT short_link
		FROM links
		WHERE original_link = $1
	`

	conn, err := p.pool.Acquire()
	if err != nil {
		return nil, status.Error(codes.Internal, entities.HTTP500)
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
		return link, nil
	}

	conn, err := p.pool.Acquire()
	if err != nil {
		return nil, status.Error(codes.Internal, entities.HTTP500)
	}
	defer p.pool.Release(conn)

	if shortLink, err := helpers.GenToken(); err != nil {
		return nil, status.Error(codes.Internal, entities.HTTP500)
	} else if _, err := conn.Exec(qry, request.GetUrl(), shortLink); err != nil {
		return nil, status.Error(codes.Internal, entities.HTTP500)
	} else {
		return &api.AddURLResponse{Url: &api.ShortenedURL{ShortenedURL: shortLink, OriginalURL: request.GetUrl()}}, nil
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
		return nil, status.Error(codes.Internal, entities.HTTP500)
	}
	defer p.pool.Release(conn)

	var originalLink string
	err = conn.QueryRow(qry, request.GetUrl()).Scan(&originalLink)
	if err != nil {
		return nil, status.Error(codes.NotFound, entities.HTTP404)
	} else {
		return &api.GetURLResponse{Url: &api.ShortenedURL{ShortenedURL: request.GetUrl(), OriginalURL: originalLink}}, nil
	}
}
