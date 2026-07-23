package dashboard

import (
	"go-fiber-api/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DeleteUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		id := c.Params("id")

		var user models.User

		// Find user
		if err := db.First(&user, id).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "User not found",
			})
		}

		// Delete user
		if err := db.Delete(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Failed to delete user",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "User deleted successfully",
			"data": fiber.Map{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
			},
		})
	}
}