package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// Root represents root configuration combining all configs
type Root struct {
	App       App
	Postgres  Postgres
	Redis     Redis
	JWT       JWT
	Email     Email
	CORS      CORS
	RateLimit RateLimit
}

// CORS represents CORS configuration
type CORS struct {
	AllowedOrigins string `envconfig:"CORS_ALLOWED_ORIGINS" default:"*"`
}

// RateLimit represents rate limit configuration
type RateLimit struct {
	Max      int           `envconfig:"RATE_LIMIT_MAX" default:"100"`
	Duration time.Duration `envconfig:"RATE_LIMIT_DURATION" default:"15m"`
}

// Load loads configuration from environment variables
func Load(envFilePath string) Root {
	// Load .env file if exists
	if err := godotenv.Load(envFilePath); err != nil {
		logrus.Warn("No .env file found, using environment variables")
	}

	var cfg Root

	// Process each config section
	if err := envconfig.Process("", &cfg.App); err != nil {
		logrus.Fatalf("Failed to process App config: %v", err)
	}

	if err := envconfig.Process("", &cfg.Postgres); err != nil {
		logrus.Fatalf("Failed to process Postgres config: %v", err)
	}

	if err := envconfig.Process("", &cfg.Redis); err != nil {
		logrus.Fatalf("Failed to process Redis config: %v", err)
	}

	if err := envconfig.Process("", &cfg.JWT); err != nil {
		logrus.Fatalf("Failed to process JWT config: %v", err)
	}

	if err := envconfig.Process("", &cfg.Email); err != nil {
		logrus.Fatalf("Failed to process Email config: %v", err)
	}

	if err := envconfig.Process("", &cfg.CORS); err != nil {
		logrus.Fatalf("Failed to process CORS config: %v", err)
	}

	if err := envconfig.Process("", &cfg.RateLimit); err != nil {
		logrus.Fatalf("Failed to process RateLimit config: %v", err)
	}

	return cfg
}

// Getenv gets environment variable with fallback
func Getenv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
