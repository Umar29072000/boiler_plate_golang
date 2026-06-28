package migrations

import (
	"boiler_plate_be_golang/internal/models"
	"log"

	"gorm.io/gorm"
)

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&models.User{},
		// Add other models here
	)

	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}
