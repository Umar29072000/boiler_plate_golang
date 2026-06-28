package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
)

// Compress creates a gzip compression middleware
func Compress() fiber.Handler {
	return compress.New(compress.Config{
		Level: compress.LevelDefault, // Default compression level (balance of speed and size)
	})
}

// CompressBest creates compression middleware with best compression
func CompressBest() fiber.Handler {
	return compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	})
}

// CompressFast creates compression middleware with fast compression
func CompressFast() fiber.Handler {
	return compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	})
}
