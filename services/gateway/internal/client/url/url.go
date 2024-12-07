package url

import (
	"context"
	pb "url-shortener/pb/shortener"
)

func (c Client) Get(ctx context.Context, shortUrl *pb.ShortUrlMsg) (*pb.FullUrlMsg, error) {
	return c.cache.Get(ctx, shortUrl)
}

func (c Client) Set(ctx context.Context, fullUrl *pb.FullUrlMsg) (*pb.ShortUrlMsg, error) {
	return c.cache.Set(ctx, fullUrl)
}
