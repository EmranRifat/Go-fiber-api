package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"go-fiber-api/controllers"
	"go-fiber-api/middleware"
	"go-fiber-api/security"
)

func Register(app *fiber.App, jwtm *security.JWTManager, db *gorm.DB) {
	api := app.Group("/api")

	// Auth (public)
	auth := api.Group("/auth")
	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login(jwtm))

	// Products (reads are public)
	product := api.Group("/product")
	product.Get("/", controllers.GetProducts)
	product.Get("/:id", controllers.GetProductByID)

	// CREATE (PUBLIC) — no middleware here
	product.Post("/", controllers.CreateProductDB(db))

	// Protected writes — attach middleware per-route
	product.Put("/:id",   middleware.Protect(jwtm), controllers.UpdateProduct) // swap to DB when ready
	product.Patch("/:id", middleware.Protect(jwtm), controllers.PatchProduct)  // swap to DB when ready
	product.Delete("/:id",middleware.Protect(jwtm), controllers.DeleteProduct) // swap to DB when ready

	
	// Who am I (protected)
	api.Get("/me", middleware.Protect(jwtm), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"subject": c.Locals("sub"),
			"email":   c.Locals("email"),
		})
	})
}
