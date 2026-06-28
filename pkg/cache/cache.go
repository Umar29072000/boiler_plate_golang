package cache

import (
	"boiler_plate_be_golang/pkg/redis"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// CacheService provides high-level caching operations
type CacheService struct {
	DefaultTTL time.Duration
	HitCount   int64
	MissCount  int64
}

// NewCacheService creates a new cache service
func NewCacheService(defaultTTL time.Duration) *CacheService {
	if defaultTTL == 0 {
		defaultTTL = 5 * time.Minute
	}
	return &CacheService{
		DefaultTTL: defaultTTL,
	}
}

// Get retrieves a value from cache
func (c *CacheService) Get(key string, dest interface{}) error {
	if !redis.IsConnected() {
		c.MissCount++
		return fmt.Errorf("redis not connected")
	}

	data, err := redis.Get(key)
	if err != nil {
		c.MissCount++
		return err
	}

	c.HitCount++
	return json.Unmarshal([]byte(data), dest)
}

// Set stores a value in cache with TTL
func (c *CacheService) Set(key string, value interface{}, ttl time.Duration) error {
	if !redis.IsConnected() {
		return fmt.Errorf("redis not connected")
	}

	if ttl == 0 {
		ttl = c.DefaultTTL
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return redis.Set(key, string(data), ttl)
}

// Del deletes a cache key
func (c *CacheService) Del(key string) error {
	if !redis.IsConnected() {
		return fmt.Errorf("redis not connected")
	}
	return redis.Del(key)
}

// DelPattern deletes all keys matching a pattern
func (c *CacheService) DelPattern(pattern string) error {
	if !redis.IsConnected() {
		return fmt.Errorf("redis not connected")
	}
	_, err := redis.DelPattern(pattern)
	return err
}

// Exists checks if a key exists in cache
func (c *CacheService) Exists(key string) bool {
	if !redis.IsConnected() {
		return false
	}
	exists, _ := redis.Exists(key)
	return exists
}

// Flush clears all cache entries
func (c *CacheService) Flush() error {
	if !redis.IsConnected() {
		return fmt.Errorf("redis not connected")
	}
	return redis.FlushDB()
}

// GetStats returns cache hit/miss statistics
func (c *CacheService) GetStats() map[string]interface{} {
	total := c.HitCount + c.MissCount
	hitRate := 0.0
	if total > 0 {
		hitRate = float64(c.HitCount) / float64(total) * 100
	}

	return map[string]interface{}{
		"hits":     c.HitCount,
		"misses":   c.MissCount,
		"total":    total,
		"hit_rate": fmt.Sprintf("%.2f%%", hitRate),
	}
}

// InvalidateByPrefix invalidates all cache entries with a given prefix
func (c *CacheService) InvalidateByPrefix(prefix string) error {
	pattern := fmt.Sprintf("%s*", prefix)
	deleted, err := redis.DelPattern(pattern)
	if err != nil {
		return err
	}
	log.Printf("Cache invalidation: deleted %d keys with prefix '%s'", deleted, prefix)
	return nil
}

// WarmCache pre-loads cache with data (useful for frequently accessed data)
func (c *CacheService) WarmCache(key string, value interface{}, ttl time.Duration) error {
	return c.Set(key, value, ttl)
}
