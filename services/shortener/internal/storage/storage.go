package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"log"
	"url-shortener/services/shortener/internal/core"
)

type (
	Storage struct {
		db *sqlx.DB
	}

	Interface interface {
		Probe(ctx context.Context) error

		Get(ctx context.Context, shortUrl *core.ShortUrl) (*core.FullUrl, error)
		Set(ctx context.Context, fullUrl *core.FullUrl) (*core.ShortUrl, error)
	}
)

func New(db *sqlx.DB) Storage {
	return Storage{
		db: db,
	}
}

func (s *Storage) Probe(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		log.Fatalf("cannot connect to database: %v", err)
		return err
	}

	return nil
}
