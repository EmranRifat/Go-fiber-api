package controllers

import (
	"go-fiber-api/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ProductCategories(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var categories []models.ProductCategory

		if err := db.Find(&categories).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":      "error",
				"status_code": fiber.StatusInternalServerError,
				"message":     "Failed to fetch product categories",
				"error":       err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":      "success",
			"status_code": fiber.StatusOK,
			"message":     "Product categories fetched successfully",
			"data":        categories,
		})
	}
}







// ==================================================================
                // Single category API handler
// ==================================================================

func SingleProductCategory(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		var category models.ProductCategory
		if err := db.First(&category, "id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status":      "error",
					"status_code": fiber.StatusNotFound,
					"message":     "Product category not found",
					"error":       err.Error(),
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":      "error",
				"status_code": fiber.StatusInternalServerError,
				"message":     "Failed to fetch product category",
				"error":       err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":      "success",
			"status_code": fiber.StatusOK,
			"message":     "Product category fetched successfully",
			"data":        category,
		})
	}
}