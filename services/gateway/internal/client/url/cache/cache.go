package cache

import (
	"context"
	"math/rand"
	"sync"
	"time"
	pb "url-shortener/pb/shortener"
)

// Custom LRU implementation with TTL 1 hour + JITTER

type (
	Interface interface {
		Get(ctx context.Context, shortUrl *pb.ShortUrl) (*pb.FullUrl, error)
		Set(ctx context.Context, fullUrl *pb.FullUrl) (*pb.ShortUrl, error)
	}

	cacheValue struct {
		fullUrl string
	}

	cacheItem struct {
		cacheValue
		usedTimes  int
		expiration time.Time
	}

	cache map[string]cacheItem

	URLCache struct {
		rpc       pb.ShortenerClient
		cache     cache
		cacheSize int
		ttl       time.Duration
		jitter    time.Duration
		mutex     sync.Mutex
	}
)

func New(rpc pb.ShortenerClient, cacheSize int, ttl, jitter time.Duration) *URLCache {
	c := make(cache, cacheSize)
	return &URLCache{
		rpc:       rpc,
		cache:     c,
		cacheSize: cacheSize,
		ttl:       ttl,
		jitter:    jitter,
	}
}

func (c *URLCache) getFromCache(shortUrl string) (string, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, ok := c.cache[shortUrl]
	if !ok {
		return "", false
	}

	if time.Now().After(item.expiration) {
		c.invalidate(shortUrl)
		return "", false
	}

	item.usedTimes += 1
	c.cache[shortUrl] = item

	return item.fullUrl, true
}

func addJitter(maxJitter time.Duration) time.Duration {
	return time.Duration(rand.Int63n(int64(maxJitter)*2)) - maxJitter
}

func (c *URLCache) setToCache(shortUrl string, cv *cacheValue) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Calculate random jitter within the range [-jitter, +jitter]
	jitter := addJitter(c.jitter)

	// Calculate expiration time with jitter
	expiration := time.Now().Add(c.ttl + jitter)

	if len(c.cache) == c.cacheSize {
		c.replace()
	}

	c.cache[shortUrl] = cacheItem{
		cacheValue: *cv,
		expiration: expiration,
		usedTimes:  0,
	}

}

func (c *URLCache) invalidate(shortUrl string) {
	delete(c.cache, shortUrl)
}

func (c *URLCache) replace() {
	// LRU replacement algorithm
	var lowest cacheItem
	var lowestShortUrl string
	i := 0
	for sn, v := range c.cache {
		if time.Now().After(v.expiration) {
			c.invalidate(sn)
			return
		}
		if i == 0 {
			lowest = v
			lowestShortUrl = sn
		} else {
			if v.usedTimes < lowest.usedTimes {
				lowest = v
				lowestShortUrl = sn
			}
		}
		i++
	}
	delete(c.cache, lowestShortUrl)

	// optimization to avoid big integers
	if lowest.usedTimes > 100_000 {
		for _, v := range c.cache {
			v.usedTimes -= lowest.usedTimes
		}
	}
}

// gRPC methods

func (c *URLCache) Get(ctx context.Context, shortUrl *pb.ShortUrl) (*pb.FullUrl, error) {
	// read through

	cacheV, exists := c.getFromCache(shortUrl.ShortUrl)
	if exists {
		return &pb.FullUrl{FullUrl: cacheV}, nil
	}

	fullUrl, err := c.rpc.Get(ctx, shortUrl)
	if err != nil {
		return nil, err
	}

	c.setToCache(shortUrl.ShortUrl, &cacheValue{
		fullUrl: fullUrl.FullUrl,
	})

	return fullUrl, nil
}

func (c *URLCache) Set(ctx context.Context, fullUrl *pb.FullUrl) (*pb.ShortUrl, error) {
	// write through

	shortUrl, err := c.rpc.Set(ctx, fullUrl)
	if err != nil {
		return nil, err
	}

	c.setToCache(shortUrl.ShortUrl, &cacheValue{
		fullUrl: fullUrl.FullUrl,
	})

	return shortUrl, nil
}
