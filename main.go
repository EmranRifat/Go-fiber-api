// package main

// import (
// 	"fmt"

// 	_ "github.com/joho/godotenv/autoload" // loads .env automatically

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/fiber/v2/middleware/cors"
// 	"github.com/gofiber/fiber/v2/middleware/logger"

// 	"go-fiber-api/config"
// 	"go-fiber-api/database"
// )

// func main() {
// 	// Load config & connect DB
// 	cfg := config.Load()

// 	db, err := database.Connect(cfg)
// 	if err != nil {
// 		panic(err)
// 	}
// 	_ = db // db is ready (Connect ran AutoMigrate), we’ll use it once handlers exist

// 	// Fiber app
// 	app := fiber.New(fiber.Config{AppName: "Go Fiber API"})

// 	// Basic middleware (Fiber built-ins only; your custom mw isn’t created yet)
// 	app.Use(logger.New())
// 	app.Use(cors.New())

// 	// Health route
// 	app.Get("/api/v1/health", func(c *fiber.Ctx) error {
// 		return c.JSON(fiber.Map{"status": "ok"})
// 	})

// 	// Start server
// 	addr := fmt.Sprintf(":%s", cfg.AppPort)
// 	if err := app.Listen(addr); err != nil {
// 		panic(err)
// 	}
// }

// ////////////////////////////////////////////////////////////////////////////////////////////
// // The code below is the simplified version without DB connection and handlers/routes
// ////////////////////////////////////////////////////////////////////////////////////////////

package main

import (
	"fmt"

	_ "github.com/joho/godotenv/autoload"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"go-fiber-api/routes"
)

func main() {

	app := fiber.New(fiber.Config{AppName: "Go Fiber API"})

	app.Use(logger.New())
	app.Use(cors.New())

	// Root + health
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("🚀 Welcome to Go Fiber API, Go-Fiber server is Running..!")
	})



	app.Get("/api/v1/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Register product routes
	routes.Register(app)

	addr := fmt.Sprintf(":%s", "3001")
	if err := app.Listen(addr); err != nil {
		panic(err)
	}
}
