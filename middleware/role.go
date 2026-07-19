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

		if userRole == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "role not found",
			})
		}

		role := strings.ToLower(userRole.(string))

		for _, allowed := range roles {
			if role == strings.ToLower(allowed) {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "access denied",
			"role": role,
		})
	}
}