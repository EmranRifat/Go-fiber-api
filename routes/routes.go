package routes

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-api/controllers"
)

func Register(app *fiber.App) {
	api := app.Group("/api")

	// Products
	api.Get("/products", controllers.GetProducts)
	api.Get("/products/:id", controllers.GetProductByID)
}
