package core

import pb "url-shortener/pb/shortener"

type ShortUrl struct {
	ShortUrl string `db:"short_url"`
}

func ShortUrlToPb(s *ShortUrl) *pb.ShortUrlMsg {
	return &pb.ShortUrlMsg{
		ShortUrl: s.ShortUrl,
	}
}

func ShortUrlFromPb(s *pb.ShortUrlMsg) *ShortUrl {
	return &ShortUrl{
		ShortUrl: s.ShortUrl,
	}
}

type FullUrl struct {
	FullUrl string `db:"full_url"`
}

func FullUrlToPb(f *FullUrl) *pb.FullUrlMsg {
	return &pb.FullUrlMsg{
		FullUrl: f.FullUrl,
	}
}

func FullUrlFromPb(f *pb.FullUrlMsg) *FullUrl {
	return &FullUrl{
		FullUrl: f.FullUrl,
	}
}

type Pair struct {
	ShortUrl string `db:"short_url"`
	FullUrl  string `db:"full_url"`
}
