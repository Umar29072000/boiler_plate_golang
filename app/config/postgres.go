package config

import (
	"fmt"

	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Postgres represents PostgreSQL database configuration
type Postgres struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     string `envconfig:"DB_PORT" default:"5432"`
	User     string `envconfig:"DB_USER" default:"postgres"`
	Password string `envconfig:"DB_PASSWORD" default:"postgres"`
	Dbname   string `envconfig:"DB_NAME" default:"fiber_boilerplate"`
	SSLMode  string `envconfig:"DB_SSL_MODE" default:"disable"`
	Timezone string `envconfig:"DB_TIMEZONE" default:"Asia/Jakarta"`
}

// GetDSN returns database connection string
func (c *Postgres) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Dbname,
		c.SSLMode,
		c.Timezone,
	)
}

// OpenDatabaseConnection opens a PostgreSQL database connection
func OpenDatabaseConnection(postgres Postgres) (*gorm.DB, error) {
	dsn := postgres.GetDSN()
	
	db, err := gorm.Open(postgresDriver.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgresql database: %w", err)
	}
	
	return db, nil
}
