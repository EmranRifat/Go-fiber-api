// routes/routes.go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"go-fiber-api/controllers"
	"go-fiber-api/middleware"
	"go-fiber-api/security"
)

func ManageRoutes(app *fiber.App, jwtm *security.JWTManager, db *gorm.DB) {
	api := app.Group("/api")

	// Auth (public)
	api.Post("/auth/register", controllers.RegisterDB(db))
	api.Post("/auth/login", controllers.Login(jwtm))

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




	
	// Who am I (protected)
	api.Get("/me", middleware.Protect(jwtm), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"subject": c.Locals("sub"),
			"email":   c.Locals("email"),
		})
	})
}
