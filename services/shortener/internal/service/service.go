package service

import (
	"context"
	"url-shortener/services/shortener/internal/core"
	"url-shortener/services/shortener/internal/storage"
)

type (
	Service struct {
		db storage.Interface
	}

	Interface interface {
		Url
		ServiceOperations
	}

	Url interface {
		Get(ctx context.Context, shortUrl *core.ShortUrl) (*core.FullUrl, error)
		Set(ctx context.Context, fullUrl *core.FullUrl) (*core.ShortUrl, error)
	}

	ServiceOperations interface {
		Probe(ctx context.Context) error
	}
)

func New(db storage.Interface) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) Probe(ctx context.Context) error {
	return s.db.Probe(ctx)
}
