package url

import (
	"context"
	pb "url-shortener/pb/shortener"
)

func (s Service) Get(ctx context.Context, shortUrl *pb.ShortUrl) (*pb.FullUrl, error) {
	return s.client.Get(ctx, shortUrl)
}

func (s Service) Set(ctx context.Context, fullUrl *pb.FullUrl) (*pb.ShortUrl, error) {
	return s.client.Set(ctx, fullUrl)
}
