package middleware

import (
	"boiler_plate_be_golang/internal/config"
	"boiler_plate_be_golang/pkg/redis"
	"boiler_plate_be_golang/pkg/utils"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

// RateLimiter configuration
type RateLimiterConfig struct {
	Max      int           // Maximum requests
	Duration time.Duration // Time window
}

// In-memory rate limiter (fallback when Redis unavailable)
type memoryStore struct {
	mu      sync.RWMutex
	clients map[string]*clientData
}

type clientData struct {
	count      int
	resetTime  time.Time
}

var memStore = &memoryStore{
	clients: make(map[string]*clientData),
}

// RateLimiter creates a rate limiting middleware with Redis backend
func RateLimiter(cfg RateLimiterConfig) fiber.Handler {
	// Default configuration
	if cfg.Max == 0 {
		cfg.Max = 100
	}
	if cfg.Duration == 0 {
		cfg.Duration = 15 * time.Minute
	}

	return func(c *fiber.Ctx) error {
		// Get client identifier (IP address)
		clientIP := c.IP()
		key := fmt.Sprintf("ratelimit:%s", clientIP)

		// Try Redis first if available
		if redis.IsConnected() {
			return handleRedisRateLimit(c, key, cfg)
		}

		// Fallback to in-memory rate limiting
		return handleMemoryRateLimit(c, clientIP, cfg)
	}
}

// handleRedisRateLimit uses Redis for distributed rate limiting
func handleRedisRateLimit(c *fiber.Ctx, key string, cfg RateLimiterConfig) error {
	// Increment request count
	count, err := redis.Incr(key)
	if err != nil {
		// If Redis fails, allow the request (fail open)
		return c.Next()
	}

	// Set expiration on first request
	if count == 1 {
		redis.Expire(key, cfg.Duration)
	}

	// Get remaining time
	ttl := cfg.Duration

	// Set rate limit headers
	c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", cfg.Max))
	c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", max(0, cfg.Max-int(count))))
	c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(ttl).Unix()))

	// Check if limit exceeded
	if count > int64(cfg.Max) {
		c.Set("Retry-After", fmt.Sprintf("%d", int(ttl.Seconds())))
		return utils.ErrorResponse(c, fiber.StatusTooManyRequests, "Rate limit exceeded. Please try again later.", nil)
	}

	return c.Next()
}

// handleMemoryRateLimit uses in-memory store as fallback
func handleMemoryRateLimit(c *fiber.Ctx, clientIP string, cfg RateLimiterConfig) error {
	memStore.mu.Lock()
	defer memStore.mu.Unlock()

	now := time.Now()
	client, exists := memStore.clients[clientIP]

	// Clean up expired entries periodically
	if len(memStore.clients) > 10000 {
		for ip, data := range memStore.clients {
			if now.After(data.resetTime) {
				delete(memStore.clients, ip)
			}
		}
	}

	if !exists || now.After(client.resetTime) {
		// New client or expired window
		memStore.clients[clientIP] = &clientData{
			count:     1,
			resetTime: now.Add(cfg.Duration),
		}
		client = memStore.clients[clientIP]
	} else {
		// Increment count
		client.count++
	}

	// Set rate limit headers
	remaining := max(0, cfg.Max-client.count)
	c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", cfg.Max))
	c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
	c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", client.resetTime.Unix()))

	// Check if limit exceeded
	if client.count > cfg.Max {
		resetIn := time.Until(client.resetTime).Seconds()
		c.Set("Retry-After", fmt.Sprintf("%d", int(resetIn)))
		return utils.ErrorResponse(c, fiber.StatusTooManyRequests, "Rate limit exceeded. Please try again later.", nil)
	}

	return c.Next()
}

// DefaultRateLimiter creates rate limiter with default config from env
func DefaultRateLimiter() fiber.Handler {
	return RateLimiter(RateLimiterConfig{
		Max:      config.App.RateLimit.Max,
		Duration: config.App.RateLimit.Duration,
	})
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
