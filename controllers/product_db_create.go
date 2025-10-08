package controllers

// import (
// 	"github.com/gofiber/fiber/v2"
// 	"gorm.io/gorm"

// 	"go-fiber-api/models"
// 	"go-fiber-api/types"
// )

// // CreateProductDB returns a Fiber handler that creates a product in DB.
// func CreateProductDB(db *gorm.DB) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		var in types.ProductInput
// 		if err := c.BodyParser(&in); err != nil {
// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
// 		}
// 		if in.Name == "" || in.Price < 0 {
// 			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "name required, price >= 0"})
// 		}

// 		p := models.Product{
// 			Name:    in.Name,
// 			Price:   in.Price,
// 			InStock: in.InStock,
// 		}
// 		if err := db.Create(&p).Error; err != nil {
// 			return c.Status(500).JSON(fiber.Map{"error": "db error"})
// 		}
// 		return c.Status(fiber.StatusCreated).JSON(p)
// 	}
// }
