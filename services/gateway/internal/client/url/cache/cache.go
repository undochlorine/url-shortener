package cache

import (
	"context"
	"math/rand"
	"sync"
	"time"
	pb "url-shortener/pb/shortener"
	"url-shortener/services/gateway/internal/client/url"
)

const CACHE_SIZE = 10

// Custom LRU implementation with TTL 1 hour + JITTER

type (
	CacheValue struct {
		FullUrl   string
		UsedTimes int
	}

	CacheItem struct {
		CacheValue
		expiration time.Time
	}

	Cache map[string]CacheItem

	URLCache struct {
		*url.Client
		cache  Cache
		ttl    time.Duration
		jitter time.Duration
		mutex  sync.Mutex
	}
)

func New(client *url.Client, ttl, jitter time.Duration) *URLCache {
	cache := make(Cache, CACHE_SIZE)
	return &URLCache{
		Client: client,
		cache:  cache,
		ttl:    ttl,
		jitter: jitter,
	}
}

func (c *URLCache) Get(shortUrl string) (string, bool) {
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

	item.UsedTimes += 1
	c.cache[shortUrl] = item

	return item.FullUrl, true
}

func addJitter(maxJitter time.Duration) time.Duration {
	return time.Duration(rand.Int63n(int64(maxJitter)*2)) - maxJitter
}

func (c *URLCache) set(shortUrl string, cv *CacheValue) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Calculate random jitter within the range [-jitter, +jitter]
	jitter := addJitter(c.jitter)

	// Calculate expiration time with jitter
	expiration := time.Now().Add(c.ttl + jitter)

	if len(c.cache) == CACHE_SIZE {
		c.replace()
	}

	c.cache[shortUrl] = CacheItem{
		CacheValue: *cv,
		expiration: expiration,
	}

}

func (c *URLCache) invalidate(shortUrl string) {
	delete(c.cache, shortUrl)
}

func (c *URLCache) replace() {
	// LRU replacement algorithm
	var lowest CacheItem
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
			if v.CacheValue.UsedTimes < lowest.CacheValue.UsedTimes {
				lowest = v
				lowestShortUrl = sn
			}
		}
		i++
	}
	delete(c.cache, lowestShortUrl)
}

// gRPC client methods

func (c *URLCache) GetUrl(ctx context.Context, shortUrl *pb.ShortUrl) (*pb.FullUrl, error) {
	fullName, exists := c.Get(shortUrl.ShortUrl)
	if exists {
		return &pb.FullUrl{FullUrl: fullName}, nil
	}

	return c.Client.GetUrl(ctx, shortUrl)
}

func (c *URLCache) Add(ctx context.Context, fullUrl *pb.FullUrl) (*pb.ShortUrl, error) {

	return c.Client.Add(ctx, fullUrl)
}
