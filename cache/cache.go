// Package cache provides a simple in-memory cache for storing the current temperature and weekly forecast for a given postalCode
package cache

import (
	"sync"
	"time"

	"github.com/mfryhover/weather/api"
)

var (
	cacheInstance *Cache
	once          sync.Once
)

// Value holds the timestamp when the data was cached, the weekly forecast, and the current temperature.
type Value struct {
	// timestamp is when the data was added to the cache.
	timestamp time.Time
	// weeklyForecast contains the weekly weather forecast.
	weeklyForecast api.WeeklyForecast
	// currentTemp is the current temperature.
	currentTemp float64
}

// Cache provides an in-memory store with thread-safe access and entry expiration.
type Cache struct {
	// data stores the cached values mapped by a string key.
	data map[string]Value
	// mu protects concurrent access to the cache.
	mu sync.RWMutex
	// entryTTL defines the time-to-live for each cache entry.
	entryTTL time.Duration
}

// GetCacheInstance returns the singleton instance of the Cache.
// If the cache has already been initialized, it returns the existing instance.
// The cache is initialized with a default entryTTL of 30 minutes.
func GetCacheInstance() *Cache {
	once.Do(
		func() {
			cacheInstance = &Cache{
				data:     make(map[string]Value),
				entryTTL: 30 * time.Minute,
			}
		})

	return cacheInstance
}

// SetEntryTTL sets the time-to-live duration for cache entries.
// It is safe for concurrent use. Note that changing the TTL affects all existing entries and may lead to
// unexpected expiration times.
func (c *Cache) SetEntryTTL(entryTTL time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entryTTL = entryTTL
}

// PurgeCache removes expired entries from the cache based on the entry time-to-live.
// It is safe for concurrent use.
func (c *Cache) PurgeCache() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.data {
		if time.Now().After(v.timestamp.Add(c.entryTTL)) {
			delete(c.data, k)
		}
	}
}

// StartAutoPurge starts a background goroutine that purges expired entries at the specified interval.
// It is safe for concurrent use.
func (c *Cache) StartAutoPurge(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			c.PurgeCache()
		}
	}()
}

// Add inserts a new entry into the cache with the specified key, current temperature, and weekly forecast.
// It is safe for concurrent use.
func (c *Cache) Add(key string, currentTemp float64, weeklyForecast api.WeeklyForecast) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = Value{
		timestamp:      time.Now(),
		weeklyForecast: weeklyForecast,
		currentTemp:    currentTemp,
	}
}

// Get retrieves the current temperature and weekly forecast for the given key.
// It returns false if the key is not found or the entry has expired.
// It is safe for concurrent use.
func (c *Cache) Get(key string) (float64, api.WeeklyForecast, bool) {
	c.mu.RLock()
	value, ok := c.data[key]
	c.mu.RUnlock()

	if !ok {
		return 0, api.WeeklyForecast{}, false
	}

	if time.Since(value.timestamp) > c.entryTTL {
		c.Delete(key) // Safe to call; it acquires the write lock internally
		return 0, api.WeeklyForecast{}, false
	}

	return value.currentTemp, value.weeklyForecast, ok
}

// Delete removes the entry associated with the key from the cache.
// It is safe for concurrent use.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}
