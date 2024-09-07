package url

import (
	"context"
	pb "url-shortener/pb/shortener"
)

func (s Service) Get(ctx context.Context, shortUrl *pb.ShortUrlMsg) (*pb.FullUrlMsg, error) {
	return s.client.Get(ctx, shortUrl)
}

func (s Service) Set(ctx context.Context, fullUrl *pb.FullUrlMsg) (*pb.ShortUrlMsg, error) {
	return s.client.Set(ctx, fullUrl)
}
