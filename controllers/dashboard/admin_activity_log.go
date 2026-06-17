package dashboard

import (
	"strconv"
	"strings"
	"time"

	"go-fiber-api/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAdminActivityLogsHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := requireAdmin(c, db); err != nil {
			return c.Status(adminErrorStatus(err)).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		page := parsePositiveInt(c.Query("page"), 1)
		limit := parsePositiveInt(c.Query("limit"), 50)
		if limit > 200 {
			limit = 200
		}

		query := db.Model(&models.APILogs{})

		if method := strings.ToUpper(strings.TrimSpace(c.Query("method"))); method != "" {
			query = query.Where("method = ?", method)
		}
		if operation := strings.ToUpper(strings.TrimSpace(c.Query("operation"))); operation != "" {
			query = query.Where("operation = ?", operation)
		}
		if isError := strings.ToLower(strings.TrimSpace(c.Query("isError"))); isError != "" {
			switch isError {
			case "true", "1", "yes":
				query = query.Where("is_error = ?", true)
			case "false", "0", "no":
				query = query.Where("is_error = ?", false)
			default:
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"success": false,
					"message": "isError must be true or false",
				})
			}
		}
		if errorType := strings.ToUpper(strings.TrimSpace(c.Query("errorType"))); errorType != "" {
			query = query.Where("error_type = ?", errorType)
		}
		if path := strings.TrimSpace(c.Query("path")); path != "" {
			query = query.Where("path ILIKE ?", "%"+path+"%")
		}
		if status := strings.TrimSpace(c.Query("status")); status != "" {
			statusCode, err := strconv.Atoi(status)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"success": false,
					"message": "status must be a number",
				})
			}
			query = query.Where("status_code = ?", statusCode)
		}
		if userID := strings.TrimSpace(c.Query("userId")); userID != "" {
			query = query.Where("user_id = ?", userID)
		}
		if userEmail := strings.TrimSpace(c.Query("userEmail")); userEmail != "" {
			query = query.Where("user_email ILIKE ?", "%"+userEmail+"%")
		}
		if from := strings.TrimSpace(c.Query("from")); from != "" {
			fromTime, err := time.Parse(time.RFC3339, from)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"success": false,
					"message": "from must be RFC3339 datetime",
				})
			}
			query = query.Where("created_at >= ?", fromTime)
		}
		if to := strings.TrimSpace(c.Query("to")); to != "" {
			toTime, err := time.Parse(time.RFC3339, to)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"success": false,
					"message": "to must be RFC3339 datetime",
				})
			}
			query = query.Where("created_at <= ?", toTime)
		}

		var total int64
		if err := query.Count(&total).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Failed to count activity logs",
				"error":   err.Error(),
			})
		}

		var logs []models.APILogs
		offset := (page - 1) * limit
		if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Failed to load activity logs",
				"error":   err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"page":    page,
			"limit":   limit,
			"total":   total,
			"count":   len(logs),
			"data":    logs,
		})
	}
}

func parsePositiveInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || parsed <= 0 {
		return fallback
	}

	return parsed
}
