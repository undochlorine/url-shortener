package service

import (
	"context"
	"errors"
	"net/url"
	"regexp"
	"strings"
	"url-shortener/services/shortener/internal/core"
)

func (s *Service) Get(ctx context.Context, shortUrl *core.ShortUrl) (*core.FullUrl, error) {
	return s.db.Get(ctx, shortUrl)
}

func (s *Service) Set(ctx context.Context, fullUrl *core.FullUrl) (*core.ShortUrl, error) {
	_, err := url.ParseRequestURI(fullUrl.FullUrl)
	if err != nil {
		return nil, errors.New("invalid url")
	}

	cleanFullUrl, err := cleanURL(fullUrl.FullUrl)
	if err != nil {
		return nil, errors.New("invalid url")
	}

	pair := &core.Pair{
		ShortUrl: core.Shorten(cleanFullUrl),
		FullUrl:  fullUrl.FullUrl,
	}

	err = s.db.Set(ctx, pair)
	if err != nil {
		return nil, err
	}

	return &core.ShortUrl{ShortUrl: pair.ShortUrl}, nil
}

// cleanURL removes the scheme (http/https), subdomains, and TLDs from a URL
func cleanURL(input string) (string, error) {
	// Parse the input URL
	parsedURL, err := url.Parse(input)
	if err != nil {
		return "", err
	}

	// Extract the host from the parsed URL (e.g., www.example.com)
	host := parsedURL.Host

	// Remove the port number if it exists
	host = strings.Split(host, ":")[0]

	// Regular expression to remove the TLD (like .com, .org, etc.)
	tldRegex := regexp.MustCompile(`(\.[a-z]{2,})+$`)
	cleanHost := tldRegex.ReplaceAllString(host, "")

	// Remove subdomains (keeping only the main domain)
	domainParts := strings.Split(cleanHost, ".")
	if len(domainParts) > 1 {
		cleanHost = domainParts[len(domainParts)-2]
	}

	return cleanHost, nil
}
