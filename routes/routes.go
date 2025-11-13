// routes/routes.go
package routes

import (
	"go-fiber-api/controllers"
	"go-fiber-api/handlers"
	"go-fiber-api/middleware"
	"go-fiber-api/security"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


func ManageRoutes(app *fiber.App, jwtm *security.JWTManager, db *gorm.DB) {
	
	api := app.Group("/api")
	// Auth (public)
	api.Post("/auth/register", controllers.RegisterDB(db))
	api.Post("/auth/login", controllers.LoginDB(jwtm, db))
	
	// Products (PUBLIC)
	api.Get("/product",  controllers.ListProductsDB(db))

	// Detail
	api.Get("/product/:id", controllers.GetProductByIDDB(db))

	// Create
	api.Post("/product",  controllers.CreateProductDB(db))
	 
	// (Optional) protected writes later:
	// api.Put("/product/:id",   middleware.Protect(jwtm), controllers.UpdateProductDB(db))
	// api.Patch("/product/:id", middleware.Protect(jwtm), controllers.PatchProductDB(db))
	// api.Delete("/product/:id",middleware.Protect(jwtm), controllers.DeleteProductDB(db))

	// Orders (PUBLIC)
	api.Get("/orders", handlers.GetAllOrders(db))
	api.Get("/orders/:order_id", handlers.GetOrderByID(db))

	// Who am I (protected)
	api.Get("/me", middleware.Protect(jwtm), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"subject": c.Locals("sub"),
			"email":   c.Locals("email"),
		})
	})
}
