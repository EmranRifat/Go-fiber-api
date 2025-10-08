package controllers

import (
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"go-fiber-api/security"
	"go-fiber-api/types"
)

var (
	userMu     sync.RWMutex
	usersByID  = map[int]*types.User{}
	usersByEM  = map[string]*types.User{} // key = lowercase email
	nextUserID = 1
)

func normalizeEmail(e string) string { return strings.ToLower(strings.TrimSpace(e)) }



// POST /api/auth/register
func Register(c *fiber.Ctx) error {
	var in types.RegisterInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	in.Email = normalizeEmail(in.Email)
	if in.Name == "" || !strings.Contains(in.Email, "@") || len(in.Password) < 6 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			fiber.Map{"error": "name required, valid email, password>=6"},
		)
	}

	userMu.Lock()
	defer userMu.Unlock()

	if _, exists := usersByEM[in.Email]; exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email already registered"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "hashing error"})
	}

	u := &types.User{
		ID:           nextUserID,
		Name:         in.Name,
		Email:        in.Email,
		PasswordHash: string(hash),
	}
	nextUserID++
	usersByID[u.ID] = u
	usersByEM[u.Email] = u

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user": fiber.Map{"id": u.ID, "name": u.Name, "email": u.Email},
	})
}



// POST /api/auth/login
func Login(jwtm *security.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in types.LoginInput
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}
		email := normalizeEmail(in.Email)

		userMu.RLock()
		u := usersByEM[email]
		userMu.RUnlock()

		if u == nil || bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)) != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}

		tok, err := jwtm.Sign(u.ID, u.Email)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "token issue"})
		}

		return c.JSON(fiber.Map{
			"token": tok,
			"user":  fiber.Map{"id": u.ID, "name": u.Name, "email": u.Email},
		})
	}
}
