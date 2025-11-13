package database

import (
	"encoding/json"
	"fmt"
	"go-fiber-api/logger"
	"go-fiber-api/models"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SeedData seeds the database with initial data from JSON files
func SeedData(db *gorm.DB) error {
	logger.Success("🌱 Starting database seeding...")

	// Seed products from JSON
	if err := SeedProductsFromJSON(db); err != nil {
		logger.Error("Failed to seed products", err)
		return err
	}

	// Seed orders from JSON
	if err := SeedOrdersFromJSON(db); err != nil {
		logger.Error("Failed to seed orders", err)
		return err
	}

	logger.Success("✅ Database seeding completed successfully")
	return nil
}

// SeedProductsFromJSON reads products from assets/product.json and seeds them
func SeedProductsFromJSON(db *gorm.DB) error {
	logger.Success("📦 Seeding products from JSON...")

	// Check if products already exist
	var count int64
	if err := db.Model(&models.Product{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count products: %w", err)
	}

	if count > 0 {
		logger.Info(fmt.Sprintf("Products already seeded (%d records), skipping...", count))
		return nil
	}

	// Read JSON file
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	filePath := filepath.Join(projectRoot, "assets", "product.json")
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open product.json: %w", err)
	}
	defer file.Close()

	// Decode JSON
	var productsData []struct {
		ID          uint    `json:"id"`
		Name        string  `json:"name"`
		Price       float64 `json:"price"`
		Category    string  `json:"category"`
		Description string  `json:"description"`
		Brand       string  `json:"brand"`
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&productsData); err != nil {
		return fmt.Errorf("failed to decode product.json: %w", err)
	}

	// Convert to Product models
	products := make([]models.Product, len(productsData))
	for i, p := range productsData {
		products[i] = models.Product{
			ID:      p.ID,
			Name:    p.Name,
			Price:   p.Price,
			InStock: true,
		}
	}

	// Insert products with upsert
	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&products).Error; err != nil {
		return fmt.Errorf("failed to seed products: %w", err)
	}

	logger.Success(fmt.Sprintf("✅ Successfully seeded %d products", len(products)))
	return nil
}

// SeedOrdersFromJSON reads orders from json_data/order.json and seeds them
func SeedOrdersFromJSON(db *gorm.DB) error {
	logger.Success("📋 Seeding orders from JSON...")

	// Check if orders already exist
	var count int64
	if err := db.Model(&models.Order{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count orders: %w", err)
	}

	if count > 0 {
		logger.Info(fmt.Sprintf("Orders already seeded (%d records), skipping...", count))
		return nil
	}

	// Read JSON file
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	filePath := filepath.Join(projectRoot, "json_data", "order.json")
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open order.json: %w", err)
	}
	defer file.Close()

	// Decode JSON
	var ordersData []struct {
		OrderID      string  `json:"order_id"`
		CustomerName string  `json:"customer_name"`
		Status       string  `json:"status"`
		TotalAmount  float64 `json:"total_amount"`
		Currency     string  `json:"currency"`
		CreatedAt    string  `json:"created_at"`
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&ordersData); err != nil {
		return fmt.Errorf("failed to decode order.json: %w", err)
	}

	// Convert to Order models
	orders := make([]models.Order, len(ordersData))
	for i, o := range ordersData {
		// Parse created_at timestamp
		createdAt, err := parseTimestamp(o.CreatedAt)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to parse timestamp for order %s", o.OrderID), err)
			createdAt = time.Now()
		}
 
		orders[i] = models.Order{
			OrderID:      o.OrderID,
			CustomerName: o.CustomerName,
			Status:       o.Status,
			TotalAmount:  o.TotalAmount,
			Currency:     o.Currency,
			CreatedAt:    createdAt,
		}
	}

	// Insert orders with upsert
	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "order_id"}},
		DoNothing: true,
	}).Create(&orders).Error; err != nil {
		return fmt.Errorf("failed to seed orders: %w", err)
	}

	logger.Success(fmt.Sprintf("✅ Successfully seeded %d orders", len(orders)))
	return nil
}

// parseTimestamp parses various timestamp formats
func parseTimestamp(timestamp string) (time.Time, error) {
	// Try RFC3339 format first (e.g., "2025-02-10T14:10:00Z")
	t, err := time.Parse(time.RFC3339, timestamp)
	if err == nil {
		return t, nil
	}

	// Try other common formats if needed
	formats := []string{
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		t, err := time.Parse(format, timestamp)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse timestamp: %s", timestamp)
}