package controllers

import (
	"math"
	"strconv"
	"strings"

	"go-fiber-api/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetBookings returns a paginated list of bookings with all DB columns.
//
// Optional query params:
//   - page           (default 1)
//   - limit          (default 10, max 100)
//   - listing_id     (filter by listing uuid)
//   - q              (search in guest name/email/phone/listing_id)
//   - payment_method (filter: manual / sslcommerz / ...)
func GetBookings(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// ---- Pagination defaults ----
		page := 1
		limit := 10
		const maxLimit = 100

		if p := c.Query("page"); p != "" {
			if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
				page = parsedPage
			}
		}

		if l := c.Query("limit"); l != "" {
			if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
				limit = parsedLimit
			}
		}

		if limit > maxLimit {
			limit = maxLimit
		}

		offset := (page - 1) * limit

		// ---- Build base query ----
		query := db.Model(&models.Booking{})

		if listingID := strings.TrimSpace(c.Query("listing_id")); listingID != "" {
			query = query.Where("listing_id = ?", listingID)
		}

		if pm := strings.TrimSpace(c.Query("payment_method")); pm != "" {
			query = query.Where("payment_method = ?", pm)
		}

		if q := strings.TrimSpace(c.Query("q")); q != "" {
			like := "%" + q + "%"
			query = query.Where(
				"user_name LIKE ? OR user_email LIKE ? OR user_phone LIKE ? OR listing_id LIKE ? OR product_title LIKE ? OR booking_id LIKE ?",
				like, like, like, like, like, like,
			)
		}

		// ---- Count total after filters ----
		var total int64
		if err := query.Count(&total).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to count bookings",
				"error":   err.Error(),
			})
		}

		// ---- Fetch page (all columns via Find) ----
		var bookings []models.Booking
		if err := query.
			Order("created_at DESC").
			Limit(limit).
			Offset(offset).
			Find(&bookings).Error; err != nil {

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to fetch bookings",
				"error":   err.Error(),
			})
		}

		totalPages := int(math.Ceil(float64(total) / float64(limit)))

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Bookings retrieved successfully",
			"data":    bookings,
			"pagination": fiber.Map{
				"total":       total,
				"page":        page,
				"limit":       limit,
				"totalPages":  totalPages,
				"hasNext":     page < totalPages,
				"hasPrevious": page > 1,
			},
		})
	}
}