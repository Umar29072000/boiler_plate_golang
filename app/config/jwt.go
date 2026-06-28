package config

import "time"

// JWT represents JWT configuration
type JWT struct {
	Secret     string        `envconfig:"JWT_SECRET" default:"your-super-secret-jwt-key"`
	Expiration time.Duration `envconfig:"JWT_EXPIRATION" default:"24h"`
}
