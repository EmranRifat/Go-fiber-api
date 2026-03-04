package database

import (
	"encoding/json"
	"fmt"
	"go-fiber-api/logger"
	"go-fiber-api/models"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
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
	// Seed weather from JSON
	if err := SeedWeatherFromJSON(db); err != nil {
		logger.Error("Failed to seed weather", err)
		return err
	}

	logger.Success("✅ Database seeding completed successfully")
	return nil
}



// SeedProductsFromJSON reads products from assets/product.json and seeds them
func SeedProductsFromJSON(db *gorm.DB) error {
	logger.Success("📦 Seeding products from JSON...")

	var count int64
	if err := db.Model(&models.Product{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count products: %w", err)
	}

	if count > 0 {
		logger.Info(fmt.Sprintf("Products already seeded (%d records), skipping...", count))
		return nil
	}

	projectRoot, err := os.Getwd()
	if err != nil {
		return err
	}

	filePath := filepath.Join(projectRoot, "assets", "product.json")
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var productsData []struct {
		ID                string  `json:"id"`
		ProductCategoryID string  `json:"product_category_id"`
		Name              string  `json:"name"`
		Price             float64 `json:"price"`
		Image             string  `json:"image"`
		Description       string  `json:"description"`
		Manufacturer      string  `json:"manufacturer"`

		ProductCategory struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"product_category"`
	}

	if err := json.NewDecoder(file).Decode(&productsData); err != nil {
		return err
	}

	var products []models.Product

	for _, p := range productsData {

		// 🔹 Convert string → uuid.UUID
		productID, err := uuid.Parse(p.ID)
		if err != nil {
			return fmt.Errorf("invalid product uuid: %w", err)
		}

		categoryID, err := uuid.Parse(p.ProductCategoryID)
		if err != nil {
			return fmt.Errorf("invalid category uuid: %w", err)
		}

		// 🔹 Insert category safely
		db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true,
		}).Create(&models.ProductCategory{
			ID:   categoryID,
			Name: p.ProductCategory.Name,
		})

		// 🔹 Append product
		products = append(products, models.Product{
			ID:                productID,
			ProductCategoryID: categoryID,
			Name:              p.Name,
			Price:             p.Price,
			Image:             p.Image,
			Description:       p.Description,
			Manufacturer:      p.Manufacturer,
		})
	}

	// 🔹 Bulk Insert
	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&products).Error; err != nil {
		return err
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



 
// SeedWeatherFromJSON reads json_data/weather.json and seeds rows.
func SeedWeatherFromJSON(db *gorm.DB) error {
	logger.Success("📋 Seeding Weather from JSON...")

	// Read JSON file (relative to current working dir)
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}
	filePath := filepath.Join(projectRoot, "json_data", "weather.json")
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open weather.json: %w", err)
	}
	defer file.Close()

	// Decode JSON. Use float64 for numeric values to match the model types.
	var raw []struct {
		Division     string  `json:"division"`
		Lat          float64 `json:"lat"`
		Lon          float64 `json:"lon"`
		TemperatureC float64 `json:"temperature_c"`
		Humidity     int     `json:"humidity"`
		Condition    string  `json:"condition"`
		WindKph      float64 `json:"wind_kph"`
		VisibilityKm float64 `json:"visibility_km"`
		UpdatedAt    string  `json:"updated_at"` // "2025-11-23 16:00"
	}
	if err := json.NewDecoder(file).Decode(&raw); err != nil {
		return fmt.Errorf("failed to decode weather.json: %w", err)
	}

	// Map to model (preallocate exact length)
	weathers := make([]models.Weather, len(raw))
	for i, r := range raw {
		ts, err := parseTimestamp(r.UpdatedAt)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to parse timestamp for division %s", r.Division), err)
			ts = time.Now()
		}
		weathers[i] = models.Weather{
			Division:     r.Division,
			Lat:          r.Lat,
			Lon:          r.Lon,
			TemperatureC: r.TemperatureC,
			Humidity:     r.Humidity,
			Condition:    r.Condition,
			WindKph:      r.WindKph,
			VisibilityKm: r.VisibilityKm,
			UpdatedAt:    ts,
		}
	}

	// Insert with upsert-on-division (unique index on division ensures no duplicates).
	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "division"}},
		DoNothing: true,
	}).Create(&weathers).Error; err != nil {
		return fmt.Errorf("failed to seed weathers: %w", err)
	}

	logger.Success(fmt.Sprintf("✅ Successfully seeded %d weathers", len(weathers)))
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
		// support timestamps without seconds, e.g. "2025-11-23 16:00"
		"2006-01-02 15:04",
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