package database

import (
	"context"
	api "github.com/jyolando/test-ozon-go/pkg/api/proto"
	log "github.com/sirupsen/logrus"
	"reflect"
	"sync"
	"testing"
)

func TestMemoryStorage_AddURL(t *testing.T) {
	type fields struct {
		originalAsKey map[string]string
		shortAsKey    map[string]string
		logger        *log.Entry
		RWMutex       sync.RWMutex
	}
	type args struct {
		ctx     context.Context
		request *api.AddURLRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *api.AddURLResponse
		wantErr bool
	}{
		{
			name: "Valid (First Original Link)",
			fields: fields{
				originalAsKey: make(map[string]string),
				shortAsKey:    make(map[string]string),
				logger:        log.New().WithField("random", "random"),
			},
			args:    args{ctx: context.TODO(), request: &api.AddURLRequest{Url: "http://ozon.ru"}},
			want:    &api.AddURLResponse{Url: &api.ShortenedURL{OriginalURL: "http://ozon.ru", ShortenedURL: "random"}},
			wantErr: false,
		},
		{
			name: "Valid (Second Original Link)",
			fields: fields{
				originalAsKey: make(map[string]string),
				shortAsKey:    make(map[string]string),
				logger:        log.New().WithField("random", "random"),
			},
			args:    args{ctx: context.TODO(), request: &api.AddURLRequest{Url: "http://google.com"}},
			want:    &api.AddURLResponse{Url: &api.ShortenedURL{OriginalURL: "http://google.com", ShortenedURL: "random"}},
			wantErr: false,
		},
	}

	for _, tt := range tests { //nolint:govet
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryStorage{
				originalAsKey: tt.fields.originalAsKey,
				shortAsKey:    tt.fields.shortAsKey,
				logger:        tt.fields.logger,
				RWMutex:       sync.RWMutex{},
			}
			got, err := m.AddURL(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Url.OriginalURL, tt.want.Url.OriginalURL) {
				t.Errorf("AddURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryStorage_GetURL(t *testing.T) {
	type fields struct {
		originalAsKey map[string]string
		shortAsKey    map[string]string
		logger        *log.Entry
		RWMutex       sync.RWMutex
	}
	type args struct {
		ctx     context.Context
		request *api.GetURLRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *api.GetURLResponse
		wantErr bool
	}{
		{
			name: "Valid",
			fields: fields{
				originalAsKey: map[string]string{"http://ozon.ru": "random"},
				shortAsKey:    map[string]string{"random": "http://ozon.ru"},
				logger:        log.New().WithField("random", "random"),
			},
			args:    args{ctx: context.TODO(), request: &api.GetURLRequest{Url: "random"}},
			want:    &api.GetURLResponse{Url: &api.ShortenedURL{OriginalURL: "http://ozon.ru", ShortenedURL: "random"}},
			wantErr: false,
		},
		{
			name: "Non valid (empty map)",
			fields: fields{
				originalAsKey: map[string]string{},
				shortAsKey:    map[string]string{},
				logger:        log.New().WithField("random", "random"),
			},
			args:    args{ctx: context.TODO(), request: &api.GetURLRequest{Url: "random"}},
			want:    &api.GetURLResponse{Url: &api.ShortenedURL{OriginalURL: "http://google.com", ShortenedURL: "random"}},
			wantErr: true,
		},
	}

	for _, tt := range tests { //nolint:govet
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryStorage{
				originalAsKey: tt.fields.originalAsKey,
				shortAsKey:    tt.fields.shortAsKey,
				logger:        tt.fields.logger,
				RWMutex:       sync.RWMutex{},
			}
			got, err := m.GetURL(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.name == "Valid" {
				if !reflect.DeepEqual(got.Url.OriginalURL, tt.want.Url.OriginalURL) {
					t.Errorf("GetURL() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
