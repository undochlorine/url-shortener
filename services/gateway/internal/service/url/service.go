package url

import (
	"context"
	pb "url-shortener/pb/shortener"
	"url-shortener/services/gateway/internal/client/url"
)

type (
	Interface interface {
		Get(ctx context.Context, shortUrl *pb.ShortUrl) (*pb.FullUrl, error)
		Set(ctx context.Context, fullUrl *pb.FullUrl) (*pb.ShortUrl, error)
	}

	Service struct {
		client url.IClient
	}
)

func New(client url.IClient) *Service {
	return &Service{
		client: client,
	}
}
