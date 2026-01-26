package database

import (
	"fmt"
	"go-fiber-api/models"
)


func autoMigrate() error {
	// Migrate all models
	models := []interface{}{
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.Weather{},
	}



	for _, model := range models {
		if err := DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %w", model, err)
		}
	}

	return nil
	
}
