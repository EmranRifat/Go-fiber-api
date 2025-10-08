package main

import (
	"fmt"

	_ "github.com/joho/godotenv/autoload"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"go-fiber-api/config"
	"go-fiber-api/database"
	"go-fiber-api/routes"
	"go-fiber-api/security"
)


func main() {
	cfg := config.Load()
	jwtm := security.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiresHours)

	// Connect DB + ping
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	if err := database.Ping(db); err != nil {
		panic("DB ping failed: " + err.Error())
	}
	fmt.Println("DB Connection is OK...✅")


	// Fiber app
	app := fiber.New(fiber.Config{AppName: "Go Fiber API"})
	app.Use(logger.New())
	app.Use(cors.New())

	// Basic routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("🚀 Go Fiber API running...")
	})
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})


	// DB ping route
	app.Get("/api/db/ping", func(c *fiber.Ctx) error {
		if err := database.Ping(db); err != nil {
			return c.Status(500).JSON(fiber.Map{"db": "down", "detail": err.Error()})
		}
		return c.JSON(fiber.Map{"db": "ok"})
	})


	// App routes (auth/products)
	// routes.Register(app, jwtm)
	routes.Register(app, jwtm, db)  // <-- add db here

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	if err := app.Listen(addr); err != nil {
		panic(err)
	}
}
