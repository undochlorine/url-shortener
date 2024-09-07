package url

import (
	"context"
	pb "url-shortener/pb/shortener"
	"url-shortener/services/gateway/internal/client/url"
)

type (
	Interface interface {
		Get(ctx context.Context, shortUrl *pb.ShortUrlMsg) (*pb.FullUrlMsg, error)
		Set(ctx context.Context, fullUrl *pb.FullUrlMsg) (*pb.ShortUrlMsg, error)
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
