package server

import (
	"context"
	"log"
	proto "url-shortener/pb/shortener"
	"url-shortener/services/shortener/internal/core"
)

func (s Server) Get(ctx context.Context, shortUrl *proto.ShortUrlMsg) (*proto.FullUrlMsg, error) {
	fullUrl, err := s.service.Get(ctx, core.ShortUrlFromPb(shortUrl))
	if err != nil {
		log.Fatalf("cannot get url: %s", err)
		return nil, err
	}

	return core.FullUrlToPb(fullUrl), nil
}

func (s Server) Set(ctx context.Context, fullUrl *proto.FullUrlMsg) (*proto.ShortUrlMsg, error) {
	shortUrl, err := s.service.Set(ctx, core.FullUrlFromPb(fullUrl))
	if err != nil {
		log.Fatalf("cannot set url: %s", err)
		return nil, err
	}

	return core.ShortUrlToPb(shortUrl), nil
}
