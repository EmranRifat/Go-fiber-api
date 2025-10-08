package middleware

import (
	"strings"

	"go-fiber-api/security"

	"github.com/gofiber/fiber/v2"
)


func Protect(jwtm *security.JWTManager) fiber.Handler {
	
	return func(c *fiber.Ctx) error {
		h := c.Get("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing or malformed Authorization header"})
		}
		token := strings.TrimPrefix(h, "Bearer ")

		claims, err := jwtm.Parse(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
		}

		c.Locals("sub", claims.Subject)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}
