package url

import (
	"context"
	pb "url-shortener/pb/shortener"
)

func (c Client) GetUrl(ctx context.Context, shortUrl *pb.ShortUrl) (*pb.FullUrl, error) {
	return nil, nil
}

func (c Client) Add(ctx context.Context, fullUrl *pb.FullUrl) (*pb.ShortUrl, error) {
	return nil, nil
}
