package url

import (
	"context"
	pb "url-shortener/pb/shortener"
)

func (s Service) GetUrl(ctx context.Context, shortUrl *pb.ShortUrl) (*pb.FullUrl, error) {
	return s.client.GetUrl(ctx, shortUrl)
}

func (s Service) Add(ctx context.Context, fullUrl *pb.FullUrl) (*pb.ShortUrl, error) {
	return s.client.Add(ctx, fullUrl)
}
