package dashboard

import (
	"strconv"

	"go-fiber-api/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GET /api/reviews?page=1&limit=10
func GetAllReviews(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))

		if page < 1 {
			page = 1
		}

		if limit < 1 {
			limit = 10
		}

		offset := (page - 1) * limit

		var reviews []models.Review
		var total int64


		// Count reviews
		if err := db.Model(&models.Review{}).Count(&total).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":      "error",
				"status_code": fiber.StatusInternalServerError,
				"message":     "Failed to count reviews",
				"error":       err.Error(),
			})
		}


		// Fetch paginated reviews
		if err := db.
			Limit(limit).
			Offset(offset).
			Order("created_at DESC").
			Find(&reviews).Error; err != nil {

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":      "error",
				"status_code": fiber.StatusInternalServerError,
				"message":     "Failed to fetch reviews",
				"error":       err.Error(),
			})
		}


		// Format response
		result := make([]fiber.Map, 0, len(reviews))

		for _, review := range reviews {
			result = append(result, fiber.Map{
				"id":          review.ID,
				"user_id":     review.UserID,
				"listing_id":  review.ListingID,
				"rating":      review.Rating,
				"comment":     review.Comment,
				"status":      review.Status,
				"created_at":  review.CreatedAt,
				"updated_at":  review.UpdatedAt,
			})
		}


		totalPages := int((total + int64(limit) - 1) / int64(limit))


		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Reviews fetched successfully",

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