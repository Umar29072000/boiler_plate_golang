package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"boiler_plate_be_golang/pkg/redis"

	"github.com/gofiber/fiber/v2"
)

// CacheConfig configuration for cache middleware
type CacheConfig struct {
	TTL              time.Duration // Cache TTL
	KeyGenerator     func(c *fiber.Ctx) string
	SkipCondition    func(c *fiber.Ctx) bool
	CacheableMethods []string
}

// Cache creates a response caching middleware
func Cache(config CacheConfig) fiber.Handler {
	// Default configuration
	if config.TTL == 0 {
		config.TTL = 5 * time.Minute
	}
	if config.KeyGenerator == nil {
		config.KeyGenerator = defaultCacheKeyGenerator
	}
	if config.SkipCondition == nil {
		config.SkipCondition = func(c *fiber.Ctx) bool { return false }
	}
	if len(config.CacheableMethods) == 0 {
		config.CacheableMethods = []string{"GET"}
	}

	return func(c *fiber.Ctx) error {
		// Skip if Redis not connected
		if !redis.IsConnected() {
			return c.Next()
		}

		// Skip if condition is met
		if config.SkipCondition(c) {
			return c.Next()
		}

		// Only cache specified methods
		if !contains(config.CacheableMethods, c.Method()) {
			return c.Next()
		}

		// Generate cache key
		cacheKey := config.KeyGenerator(c)

		// Try to get from cache
		cachedData, err := redis.Get(cacheKey)
		if err == nil {
			// Cache hit
			c.Set("X-Cache", "HIT")
			c.Set("Content-Type", "application/json")
			return c.SendString(cachedData)
		}

		// Cache miss
		c.Set("X-Cache", "MISS")

		// Execute the handler
		if err := c.Next(); err != nil {
			return err
		}

		// Only cache successful responses (2xx status codes)
		if c.Response().StatusCode() >= 200 && c.Response().StatusCode() < 300 {
			responseBody := string(c.Response().Body())

			// Store in cache (async)
			go func() {
				if err := redis.Set(cacheKey, responseBody, config.TTL); err != nil {
					// Log error but don't fail the request
					fmt.Printf("Cache set error: %v\n", err)
				}
			}()
		}

		return nil
	}
}

// DefaultCache creates cache middleware with default 5-minute TTL
func DefaultCache() fiber.Handler {
	return Cache(CacheConfig{
		TTL: 5 * time.Minute,
	})
}

// defaultCacheKeyGenerator generates a cache key from request method, path, and query
func defaultCacheKeyGenerator(c *fiber.Ctx) string {
	// Create a unique key from method, path, and query string
	key := fmt.Sprintf("%s:%s", c.Method(), c.Path())

	// Include query parameters if present
	if len(c.Request().URI().QueryArgs().String()) > 0 {
		queryHash := hashString(c.Request().URI().QueryArgs().String())
		key = fmt.Sprintf("%s:q:%s", key, queryHash)
	}

	return fmt.Sprintf("cache:%s", key)
}

// hashString creates a SHA256 hash of a string
func hashString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:8]) // Use first 8 bytes for shorter key
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// InvalidateCache middleware to invalidate cache for specific routes
func InvalidateCache(patterns ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Execute the handler first
		if err := c.Next(); err != nil {
			return err
		}

		// If request was successful, invalidate cache
		if c.Response().StatusCode() >= 200 && c.Response().StatusCode() < 300 {
			for _, pattern := range patterns {
				go func(p string) {
					if _, err := redis.DelPattern(fmt.Sprintf("cache:%s*", p)); err != nil {
						fmt.Printf("Cache invalidation error: %v\n", err)
					}
				}(pattern)
			}
		}

		return nil
	}
}
