package handlers

import (
	"go-fiber-api/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


// GetAllOrders retrieves all orders from the database
func GetAllOrders(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var orders []models.Order
		
		if err := db.Find(&orders).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to retrieve orders",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"count":   len(orders),
			"data":    orders, 
		})
	}
}



// GetOrderByID retrieves a single order by its id
func GetOrderByID(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orderID := c.Params("id")
		println("Fetching order with ID:", orderID)
		var order models.Order
		if err := db.Where("id = ?", orderID).First(&order).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(404).JSON(fiber.Map{
					"error": "Order not found",
				})
			}
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to retrieve order",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    order,
		})
	}
}
