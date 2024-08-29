package url

import (
	"context"
	pb "url-shortener/pb/shortener"
)

func (c Client) Get(ctx context.Context, shortUrl *pb.ShortUrl) (*pb.FullUrl, error) {
	return c.cache.Get(ctx, shortUrl)
}

func (c Client) Set(ctx context.Context, fullUrl *pb.FullUrl) (*pb.ShortUrl, error) {
	return c.cache.Set(ctx, fullUrl)
}
