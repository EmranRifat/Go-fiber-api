package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"go-fiber-api/config"
	"go-fiber-api/database"
	"go-fiber-api/logger"
	"go-fiber-api/routes"
	"go-fiber-api/security"
)

func main() {

	// Load environment variables properly
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️.env file not loaded")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Config load failed:", err)
	}
	jwtm := security.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiresHours)

	// DB connect
	db, err := database.ConnectDB()
	if err != nil {
		logger.Error("Failed to connect to database", err)
		return
	}

	if err := database.Ping(db); err != nil {
		logger.Error("DB ping failed", err)
		return
	}

	logger.Success("DB Connection OK 👍")
	// Auto migrate
	// Setup HTML Engine
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		AppName: "Go Fiber API",
		Views:   engine,
	})

	app.Use(fiberlogger.New())
	app.Use(cors.New())


	// Routes
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// ✅ Home route to render HTML
	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.Render("index", fiber.Map{})
	// })

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("🚀 Go Fiber API running...")
	})


	app.Get("/api/db/ping", func(c *fiber.Ctx) error {
		if err := database.Ping(db); err != nil {
			return c.Status(500).JSON(fiber.Map{"db": "down", "detail": err.Error()})
		}
		return c.JSON(fiber.Map{"db": "ok"})
	})


	routes.ManageRoutes(app, jwtm, db)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	logger.Success(fmt.Sprintf("🚀Server is running at http://localhost%s", addr))
	
	if err := app.Listen(addr); err != nil {
		logger.Error("Failed to start server", err)
		return
	}
}
