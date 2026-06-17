package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go-fiber-api/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ActivityLogger(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		logEntry := models.APILogs{
			Method:     c.Method(),
			Operation:  crudOperation(c.Method()),
			Path:       c.Path(),
			Query:      string(c.Context().QueryArgs().QueryString()),
			StatusCode: c.Response().StatusCode(),
			IPAddress:  c.IP(),
			UserAgent:  c.Get("User-Agent"),
			LatencyMS:  time.Since(start).Milliseconds(),
		}
		applyResponseSummary(&logEntry, c.Response().Body())

		if sub, ok := c.Locals("sub").(string); ok {
			logEntry.UserID = strings.TrimSpace(sub)
		}
		if email, ok := c.Locals("email").(string); ok {
			logEntry.UserEmail = strings.TrimSpace(email)
		}
		if err != nil {
			logEntry.ErrorMessage = err.Error()
		}
		if logEntry.StatusCode >= fiber.StatusBadRequest {
			logEntry.IsError = true
			logEntry.ErrorType = errorType(logEntry.StatusCode)
			if strings.TrimSpace(logEntry.ErrorMessage) == "" {
				logEntry.ErrorMessage = fmt.Sprintf("%d %s", logEntry.StatusCode, http.StatusText(logEntry.StatusCode))
			}
		}

		if createErr := db.Create(&logEntry).Error; createErr != nil {
			return err
		}

		return err
	}
}

func applyResponseSummary(logEntry *models.APILogs, responseBody []byte) {
	if len(responseBody) == 0 {
		return
	}

	var body map[string]any
	if err := json.Unmarshal(responseBody, &body); err != nil {
		return
	}

	if status, ok := stringField(body, "status"); ok {
		logEntry.Status = status
	}
	if message, ok := stringField(body, "message"); ok {
		logEntry.Message = message
	}
	if statusCode, ok := intField(body, "statusCode"); ok {
		logEntry.StatusCode = statusCode
		return
	}
	if statusCode, ok := intField(body, "status_code"); ok {
		logEntry.StatusCode = statusCode
	}
}

func stringField(body map[string]any, key string) (string, bool) {
	value, ok := body[key]
	if !ok {
		return "", false
	}

	text, ok := value.(string)
	if !ok {
		return "", false
	}

	return strings.TrimSpace(text), strings.TrimSpace(text) != ""
}

func intField(body map[string]any, key string) (int, bool) {
	value, ok := body[key]
	if !ok {
		return 0, false
	}

	switch typed := value.(type) {
	case float64:
		return int(typed), true
	case int:
		return typed, true
	default:
		return 0, false
	}
}

func crudOperation(method string) string {
	switch strings.ToUpper(strings.TrimSpace(method)) {
	case fiber.MethodGet:
		return "READ"
	case fiber.MethodPost:
		return "CREATE"
	case fiber.MethodPut, fiber.MethodPatch:
		return "UPDATE"
	case fiber.MethodDelete:
		return "DELETE"
	default:
		return "OTHER"
	}
}

func errorType(statusCode int) string {
	if statusCode >= fiber.StatusInternalServerError {
		return "SERVER_ERROR"
	}

	return "CLIENT_ERROR"
}
