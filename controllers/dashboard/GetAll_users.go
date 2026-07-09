package dashboard

import (
	"strconv"

	"go-fiber-api/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GET /api/users?page=1&limit=10
func GetAllUsers(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Get page and limit from query params
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))

		if page < 1 {
			page = 1
		}

		if limit < 1 {
			limit = 10
		}

		offset := (page - 1) * limit

		var users []models.User
		var total int64

		// Count total users
		if err := db.Model(&models.User{}).Count(&total).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":      "error",
				"status_code": fiber.StatusInternalServerError,
				"message":     "Failed to count users",
				"error":       err.Error(),
			})
		}

		// Fetch paginated users
		if err := db.
			Limit(limit).
			Offset(offset).
			Find(&users).Error; err != nil {

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":      "error",
				"status_code": fiber.StatusInternalServerError,
				"message":     "Failed to fetch users",
				"error":       err.Error(),
			})
		}

		// Remove password from response
		result := make([]fiber.Map, 0, len(users))
		for _, u := range users {
			result = append(result, fiber.Map{
				"id":         u.ID,
				"name":       u.Name,
				"email":      u.Email,
				"role":       u.Role,
				"created_at": u.CreatedAt,
				"updated_at": u.UpdatedAt,
			})
		}

		totalPages := int((total + int64(limit) - 1) / int64(limit))

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Users fetched successfully",

				"pagination": fiber.Map{
				"page":       page,
				"limit":      limit,
				"total":      total,
				"totalPages": totalPages,
	},

			"data": result,
		})
	}
}