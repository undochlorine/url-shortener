package storage

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"url-shortener/services/shortener/internal/core"
)

func (s *Storage) Get(ctx context.Context, shortUrl *core.ShortUrl) (*core.FullUrl, error) {
	full_url := make([]*core.FullUrl, 0)

	sqlQuery := `
		SELECT full_url
		FROM urls
		WHERE short_url = ?
		LIMIT 1`

	q, sqlArgs, err := sqlx.In(sqlQuery, shortUrl.ShortUrl)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	q = s.db.Rebind(q)

	if err := s.db.SelectContext(ctx, &full_url, q, sqlArgs...); err != nil {
		fmt.Println(err.Error())
		log.Printf("cannot get full_url with alias: %v", shortUrl.ShortUrl)
		return nil, err
	}

	if len(full_url) < 1 {
		return nil, nil
	}

	return full_url[0], nil
}

func (s *Storage) Set(ctx context.Context, pair *core.Pair) error {
	sqlQuery := `
		INSERT INTO urls
		(full_url, short_url)
		VALUES (?, ?)
		`

	q, sqlArgs, err := sqlx.In(sqlQuery, pair.FullUrl, pair.ShortUrl)
	if err != nil {
		return err
	}

	q = s.db.Rebind(q)

	if _, err = s.db.ExecContext(ctx, q, sqlArgs...); err != nil {
		log.Printf("cannot set url: %v as: %v", pair.FullUrl, pair.ShortUrl)
		return err
	}

	return nil
}
