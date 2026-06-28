package redis

import (
	"boiler_plate_be_golang/internal/config"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client
var ctx = context.Background()

// Connect initializes Redis connection
func Connect() error {
	// Skip Redis if not configured
	if config.App.Redis.Host == "" {
		log.Println("Redis not configured, skipping Redis connection")
		return nil
	}

	Client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", config.App.Redis.Host, config.App.Redis.Port),
		Password:     config.App.Redis.Password,
		DB:           config.App.Redis.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	// Test connection
	_, err := Client.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		log.Println("Rate limiting will use in-memory fallback")
		Client = nil
		return nil // Don't fail startup, just warn
	}

	log.Println("Redis connected successfully")
	return nil
}

// Close closes Redis connection
func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

// IsConnected checks if Redis is connected
func IsConnected() bool {
	if Client == nil {
		return false
	}
	_, err := Client.Ping(ctx).Result()
	return err == nil
}

// Get retrieves value from Redis
func Get(key string) (string, error) {
	if Client == nil {
		return "", fmt.Errorf("redis not connected")
	}
	return Client.Get(ctx, key).Result()
}

// Set stores value in Redis with expiration
func Set(key string, value interface{}, expiration time.Duration) error {
	if Client == nil {
		return fmt.Errorf("redis not connected")
	}
	return Client.Set(ctx, key, value, expiration).Err()
}

// Del deletes key from Redis
func Del(key string) error {
	if Client == nil {
		return fmt.Errorf("redis not connected")
	}
	return Client.Del(ctx, key).Err()
}

// Incr increments value in Redis
func Incr(key string) (int64, error) {
	if Client == nil {
		return 0, fmt.Errorf("redis not connected")
	}
	return Client.Incr(ctx, key).Result()
}

// Expire sets expiration for key
func Expire(key string, expiration time.Duration) error {
	if Client == nil {
		return fmt.Errorf("redis not connected")
	}
	return Client.Expire(ctx, key, expiration).Err()
}

// Exists checks if key exists
func Exists(key string) (bool, error) {
	if Client == nil {
		return false, fmt.Errorf("redis not connected")
	}
	result, err := Client.Exists(ctx, key).Result()
	return result > 0, err
}
