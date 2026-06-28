package config

// Redis represents Redis configuration
type Redis struct {
	Host     string `envconfig:"REDIS_HOST" default:""`
	Port     string `envconfig:"REDIS_PORT" default:"6379"`
	Password string `envconfig:"REDIS_PASSWORD" default:""`
	DB       int    `envconfig:"REDIS_DB" default:"0"`
}
