package url

import (
	"context"
	pb "url-shortener/pb/shortener"
	"url-shortener/services/gateway/internal/client/url/cache"
)

type (
	IClient interface {
		Get(ctx context.Context, shortUrl *pb.ShortUrl) (*pb.FullUrl, error)
		Set(ctx context.Context, fullUrl *pb.FullUrl) (*pb.ShortUrl, error)
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
