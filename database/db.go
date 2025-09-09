package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-fiber-api/models"
)

func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		"localhost", "postgres", "postgres", "go_fiber_api", "5432", "disable", "Asia/Dhaka",
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// demo migration
	if err := db.AutoMigrate(&models.Book{}); err != nil {
		return nil, err
	}
	return db, nil
}
