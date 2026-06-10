
package database

import (
	"fmt"
	"go-fiber-api/models"
)

func autoMigrate() error {

	// AutoMigrate all models here 
	modelList := []interface{}{
		&models.User{},
		&models.ProductCategory{},
		&models.Listing{},
		&models.Order{},
		&models.Weather{},
		&models.Booking{}, 
		&models.Payment{},
		&models.HostListing{},


	}

	for _, model := range modelList {
		if err := DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %w", model, err)
		}
	}

	return nil
}