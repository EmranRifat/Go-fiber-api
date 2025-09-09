package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go-fiber-api/types"
)

var products = []types.Product{
	{ID: 1, Name: "Mechanical Keyboard", Price: 49.99, InStock: true},
	{ID: 2, Name: "Wireless Mouse",      Price: 19.99, InStock: true},
	{ID: 3, Name: "27\" Monitor",        Price: 199.00, InStock: false},
}

// GET /api/v1/products
func GetProducts(c *fiber.Ctx) error {
	return c.JSON(products)
}

// GET /api/v1/products/:id
func GetProductByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	for _, p := range products {
		if p.ID == id {
			return c.JSON(p)
		}
	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "product not found"})
}
