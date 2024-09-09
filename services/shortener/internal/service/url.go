package service

import (
	"context"
	"url-shortener/services/shortener/internal/core"
)

func (s *Service) Get(ctx context.Context, shortUrl *core.ShortUrl) (*core.FullUrl, error) {
	return s.db.Get(ctx, shortUrl)
}

func (s *Service) Set(ctx context.Context, fullUrl *core.FullUrl) (*core.ShortUrl, error) {
	pair := &core.Pair{
		ShortUrl: core.Shorten(fullUrl.FullUrl),
		FullUrl:  fullUrl.FullUrl,
	}

	err := s.db.Set(ctx, pair)
	if err != nil {
		return nil, err
	}

	return &core.ShortUrl{ShortUrl: pair.ShortUrl}, nil
}
