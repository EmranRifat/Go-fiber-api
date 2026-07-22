package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func RequireRoles(roles ...string) fiber.Handler {

	return func(c *fiber.Ctx) error {

		userRole := c.Locals("role")

		fmt.Println("USER ROLE FROM JWT:", userRole)
		fmt.Println("ALLOWED ROLES:", roles)


		roleString, ok := userRole.(string)

		if !ok || roleString == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error": "role not found",
			})
		}


		role := strings.ToLower(strings.TrimSpace(roleString))


		for _, allowed := range roles {

			if role == strings.ToLower(strings.TrimSpace(allowed)) {
				return c.Next()
			}

		}


		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": "access denied",
			"role":  role,
		})
	}
}