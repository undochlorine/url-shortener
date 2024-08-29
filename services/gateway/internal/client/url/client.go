package url

import (
	"context"
	pb "url-shortener/pb/shortener"
)

type (
	IClient interface {
		GetUrl(ctx context.Context, shortUrl *pb.ShortUrl) (*pb.FullUrl, error)
		Add(ctx context.Context, fullUrl *pb.FullUrl) (*pb.ShortUrl, error)
	}

	Client struct {
		rpc pb.ShortenerClient
	}
)

func New(rpc pb.ShortenerClient) *Client {
	return &Client{
		rpc: rpc,
	}
}
