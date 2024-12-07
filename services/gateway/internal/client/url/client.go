package url

import (
	"context"
	pb "url-shortener/pb/shortener"
	"url-shortener/services/gateway/internal/client/url/cache"
)

type (
	IClient interface {
		Get(ctx context.Context, shortUrl *pb.ShortUrlMsg) (*pb.FullUrlMsg, error)
		Set(ctx context.Context, fullUrl *pb.FullUrlMsg) (*pb.ShortUrlMsg, error)
	}

	Client struct {
		cache cache.Interface
	}
)

func New(cache cache.Interface) *Client {
	return &Client{
		cache: cache,
	}
}
