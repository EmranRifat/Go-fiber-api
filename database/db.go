package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: Error loading .env file")
	}
	
	// Get database configuration from environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	sslmode := os.Getenv("DB_SSLMODE")

	// Set default sslmode if not provided
	if sslmode == "" {
		sslmode = "disable"
	}

	// Build PostgreSQL DSN string (Data Source Name)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, database, sslmode)

	fmt.Println("🚀 Connecting to database...")

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect DB:", err)
		return nil, err
	}
	fmt.Println("✅ Successfully connected DB ")

	// Auto migrate models
	if err := autoMigrate(); err != nil {
		fmt.Println("Failed to migrate models:", err)
		return nil, err
	}
	fmt.Println("✅ Database migration completed successfully ")


	// Seed database with initial data
	if err := SeedData(DB); err != nil {
		fmt.Println("Warning: Failed to seed database:", err)
		// Don't return error, seeding is optional
	}

	return DB, nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// ConnectDB is a legacy function for backward compatibility
func ConnectDB() (*gorm.DB, error) {
	return InitDB()
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
