package database

import (
	"fmt"
	"go-fiber-api/models"
)

// autoMigrate runs auto migration for all models
func autoMigrate() error {
	// Migrate all models
	models := []interface{}{
		&models.User{},
		&models.Product{},
		&models.Order{},
	}

	for _, model := range models {
		if err := DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %w", model, err)
		}
	}

	return nil
}
