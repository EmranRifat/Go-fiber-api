package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-fiber-api/models"
)

func Connect() (*gorm.DB, error) {

   dsn := fmt.Sprintf(
  "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
  "127.0.0.1", "postgres", "postgres", "go_fiber_api", "5432", "disable", "Asia/Dhaka",
)


	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// =================================================
	// migrate models to create tables if not exist 
	// =================================================
	if err := db.AutoMigrate(
		&models.Product{},
		&models.User{},

	);
	 err != nil {
		return nil, err
	}
		return db, nil
	}



// Ping checks the database connection
func Ping(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return sqlDB.PingContext(ctx)
}
