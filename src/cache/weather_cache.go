// The package cache provides helper functions based upon a TTL and a non-TTL cache.
package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// The cache.Weather interface provides cached access to weather information based on (or disregarding) TTL.
type Weather interface {
	Get(key string) (any, bool)
	GetIgnoreTTL(key string) (any, bool)
	Set(key string, value any)
}

type DefaultWeatherCache struct {
	ttlCache    *cache.Cache
	nonTTLCache *cache.Cache
}

var _ Weather = (*DefaultWeatherCache)(nil)

// NewWeatherCache internally creates a TTL and a non-TTL cache (used in the failure edge case).
func NewWeatherCache(ttl time.Duration) *DefaultWeatherCache {
	return &DefaultWeatherCache{
		ttlCache:    cache.New(ttl, 10*ttl),
		nonTTLCache: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

// Get wraps the cache.Get method, returning a value if the TTL has not expired.
func (w *DefaultWeatherCache) Get(key string) (any, bool) {
	return w.ttlCache.Get(key)
}

// GetIgnoreTTL wraps the cache.Get method, with values read from the non-TTL cache.
func (w *DefaultWeatherCache) GetIgnoreTTL(key string) (any, bool) {
	return w.nonTTLCache.Get(key)
}

// Set wraps the cache.SetDefault method, storing the values into the TTL and the non-TTL cache.
func (w *DefaultWeatherCache) Set(key string, value any) {
	w.ttlCache.SetDefault(key, value)
	w.nonTTLCache.SetDefault(key, value)
}
