package dashboard

import (
	"go-fiber-api/dto"
	"go-fiber-api/models"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UpdateUserRole(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		id := c.Params("id")

		var req dto.UpdateUserRoleRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Invalid request body",
			})
		}

		role := strings.ToLower(strings.TrimSpace(req.Role))

		allowedRoles := map[string]bool{
			"user":  true,
			"host":  true,
			"admin": true,
		}

		if !allowedRoles[role] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Invalid role",
			})
		}

		var user models.User

		if err := db.First(&user, id).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "User not found",
			})
		}

		user.Role = role

		if err := db.Save(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Failed to update role",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "User role updated successfully",
			"data": fiber.Map{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
				"role":  user.Role,
			},
		})
	}
}