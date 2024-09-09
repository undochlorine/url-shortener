package storage

import (
	"context"
	"url-shortener/services/shortener/internal/core"
)

func (s *Storage) Get(ctx context.Context, shortUrl *core.ShortUrl) (*core.FullUrl, error) {
	return nil, nil
}

func (s *Storage) Set(ctx context.Context, fullUrl *core.Pair) error {
	return nil
}
