package controllers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"go-fiber-api/models"
)

// GET /api/weather  -> list (filters + pagination)
func ListWeatherDB(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var rows []models.Weather
		
		// ---- query params ----
		q        := strings.TrimSpace(c.Query("q"))              // matches division (ILIKE)
		minTemp  := c.QueryFloat("min_temp", -1e9)               // optional
		maxTemp  := c.QueryFloat("max_temp",  1e9)               // optional
		page     := c.QueryInt("page", 1)
		limit    := c.QueryInt("limit", 20)
		sort     := c.Query("sort", "updated_at_desc")           // or "updated_at_asc","id_asc","id_desc"

		if page < 1 { page = 1 }
		if limit < 1 || limit > 100 { limit = 20 }
		offset := (page - 1) * limit

		tx := db.Model(&models.Weather{})

		if q != "" {
			tx = tx.Where("division ILIKE ?", "%"+q+"%")
		}
		// temp range (inclusive)
		tx = tx.Where("temperature_c >= ? AND temperature_c <= ?", minTemp, maxTemp)

		switch sort {
		case "updated_at_asc":
			tx = tx.Order("updated_at ASC")
		case "id_asc":
			tx = tx.Order("id ASC")
		case "id_desc":
			tx = tx.Order("id DESC")
		default: // "updated_at_desc"
			tx = tx.Order("updated_at DESC")
		}

		if err := tx.Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "db error"})
		}
		return c.JSON(rows)
	}
}

	// GET /api/weather/:id  -> by numeric id
	func GetWeatherByIDDB(db *gorm.DB) fiber.Handler {
		return func(c *fiber.Ctx) error {
			id, err := strconv.Atoi(c.Params("id"))
			if err != nil || id < 1 {
				return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
			}
			var w models.Weather
			if err := db.First(&w, id).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return c.Status(404).JSON(fiber.Map{"error": "weather not found"})
				}
				return c.Status(500).JSON(fiber.Map{"error": "db error"})
			}
			return c.JSON(w)
		}
	}

// GET /api/weather/division/:division  -> by unique division
func GetWeatherByDivisionDB(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		division := strings.TrimSpace(c.Params("division"))
		if division == "" {
			return c.Status(400).JSON(fiber.Map{"error": "invalid division"})
		}
		var w models.Weather
		if err := db.Where("division ILIKE ?", division).First(&w).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(404).JSON(fiber.Map{"error": "weather not found"})
			}
			return c.Status(500).JSON(fiber.Map{"error": "db error"})
		}
		return c.JSON(w)
	}
}
