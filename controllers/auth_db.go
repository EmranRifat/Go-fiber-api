package controllers

import (
	"strings"
	// "sync"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"go-fiber-api/models"
	"go-fiber-api/security"
	"go-fiber-api/types"
)

// var (
// 	userMu     sync.RWMutex
// 	usersByID  = map[int]*types.User{}
// 	usersByEM  = map[string]*types.User{} // key = lowercase email
// 	nextUserID = 1
// )


func normalizeEmail(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}


// POST /api/auth/register (DB)
func RegisterDB(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var in types.RegisterInput

		// ❌ Invalid Body
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":      "error",
				"status_code": fiber.StatusBadRequest,
				"message":     "Invalid request body",
				"error":       err.Error(),
			})
		}

		in.Email = normalizeEmail(in.Email)

		// ❌ Validation Error
		if in.Name == "" || !strings.Contains(in.Email, "@") || len(in.Password) < 6 {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":      "error",
				"status_code": fiber.StatusUnprocessableEntity,
				"message":     "Validation failed",
				"error":       "name required, valid email, password >= 6",
			})
		}

		// 🔐 Hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":      "error",
				"status_code": fiber.StatusInternalServerError,
				"message":     "Failed to process password",
				"error":       err.Error(),
			})
		}

		u := models.User{
			Name:         in.Name,
			Email:        in.Email,
			PasswordHash: string(hash),
			Role:         "user",
		}

		// 🔥 Insert user
		if err := db.Create(&u).Error; err != nil {

			// ❌ Duplicate Email
			if strings.Contains(err.Error(), "duplicate") ||
				strings.Contains(err.Error(), "unique") {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"status":      "error",
					"status_code": fiber.StatusConflict,
					"message":     "Email already registered",
					"error":       err.Error(),
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":      "error",
				"status_code": fiber.StatusInternalServerError,
				"message":     "Failed to register user",
				"error":       err.Error(),
			})
		}

		// ✅ Success Response (NO status_code)
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status":  "success",
			"message": "User registered successfully",
			"data": fiber.Map{
				"user": fiber.Map{
					"id":    u.ID,
					"name":  u.Name,
					"email": u.Email,
					"role":  u.Role,
				},
			},
		})
	}
}



 
// POST /api/auth/login  (DB-backed)
func LoginDB(jwtm *security.JWTManager, db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error { 
		var in types.LoginInput
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}
		email := normalizeEmail(in.Email)
		if !strings.Contains(email, "@") || len(in.Password) == 0 {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "invalid credentials"})
		}

		// 1) find the user by email (case-insensitive)
		var u models.User
		if err := db.Where("LOWER(email) = ?", email).First(&u).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
			}
			return c.Status(500).JSON(fiber.Map{"error": "db error"})
		}
		
		// 2) verify password
		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}

		// 3) issue JWT
		tok, err := jwtm.Sign(int(u.ID), u.Email)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "token issue"})
		}

			return c.JSON(fiber.Map{
			"token": tok,
			"user":  fiber.Map{
				"id": u.ID, 
				"name": u.Name, 
				"email": u.Email,
				"role":  u.Role,
			},
		})
	}
}