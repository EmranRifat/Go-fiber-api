package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"go-fiber-api/models"
	"go-fiber-api/types"
)

// var Products = []types.Products{
// 	{
// 		ID:          1,
// 		Name:        "Mechanical Keyboard",
// 		Price:       49.99,
// 		Category:    "Accessories",
// 		Description: "Durable mechanical keyboard with RGB backlight and tactile switches for gaming and productivity.",
// 		Brand:       "Keychron",
// 	},
// 	{
// 		ID:          2,
// 		Name:        "Wireless Mouse",
// 		Price:       19.99,
// 		Category:    "Accessories",
// 		Description: "Lightweight wireless mouse with ergonomic design and long-lasting battery life.",
// 		Brand:       "Logitech",
// 	},

// }

// const productsFile = "assets/product.json" // <- keep the filename consistent

// thread-safe in-memory store
var (
	mu sync.RWMutex
	products = map[int]types.Product{
		1: {ID: 1, Name: "Mechanical Keyboard", Price: 49.99, InStock: true},
		2: {ID: 2, Name: "Wireless Mouse", Price: 19.99, InStock: true},
		3: {ID: 3, Name: "27\" Monitor", Price: 199.00, InStock: false},
	}
	nextID = 4
)


func ListProductsDB(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var products []models.Product

		// Query params
		q := strings.TrimSpace(c.Query("q", ""))
		println("Search query:", q)
		page := c.QueryInt("page", 1)
		limit := c.QueryInt("limit", 20)

		// Clamp values
		if page < 1 {
			page = 1
		}
		if limit < 1 || limit > 100 {
			limit = 20
		}
		offset := (page - 1) * limit

		// Build query
		tx := db.Model(&models.Product{})
		if q != "" {
			tx = tx.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(q)+"%")
		}

		// Count total records
		var total int64
		if err := tx.Count(&total).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to count products",
			})
		}

		// Apply pagination and fetch results
		if err := tx.Order("id ASC").Limit(limit).Offset(offset).Find(&products).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to retrieve products",
			})
		}

		// Response
		return c.JSON(fiber.Map{
			"success": true,
			"count":   len(products),
			"total":   total,
			"data":    products,
		})
	}
}

// GET /api/product  -> list (with optional ?q= and pagination)
func ListProductsDB1(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var items []models.Product
		
		// optional filters
		q := c.Query("q")
		println("Search query:", q)
		page, limit := c.QueryInt("page", 1), c.QueryInt("limit", 20)
		if page < 1 { page = 1 }
		if limit < 1 || limit > 100 { limit = 20 }
		offset := (page - 1) * limit

		tx := db.Model(&models.Product{})
		if q != "" {
			tx = tx.Where("name ILIKE ?", "%"+q+"%")
		}

		if err := tx.Order("id ASC").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "db error"})
		}
		return c.JSON(items)
	}
}


func ListProductsDB2(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var products []models.Product

		if err := db.Find(&products).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to retrieve products",

			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"count":   len(products),
			"data":    products,
		})
	}
}




// GET /api/product/:id -> detail
func GetProductByIDDB(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil || id < 1 {
			return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
		}
		var p models.Product

		if err := db.First(&p, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(404).JSON(fiber.Map{"error": "product not found"})
			}
			return c.Status(500).JSON(fiber.Map{"error": "db error"})
		}
		return c.JSON(p)
	}
}




// CreateProductDB returns a Fiber handler that creates a product in DB.
func CreateProductDB(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in types.ProductInput
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}
		if in.Name == "" || in.Price < 0 {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "name required, price >= 0"})
		}

		p := models.Product{
			Name:    in.Name,
			Price:   in.Price,
			InStock: in.InStock,
		}
		
		if err := db.Create(&p).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "db error"})
		}

		// Optional: tell the client where the new resource lives
		c.Location(fmt.Sprintf("/api/product/%d", p.ID))

		
		// Success message + data
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Product created successfully",
			"data": p,
		})
	}
}




// PUT /api/product/:id (full replace)
func UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var in types.ProductInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if in.Name == "" || in.Price < 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "name required, price >= 0"})
	}

	mu.Lock()
	defer mu.Unlock()
	if _, ok := products[id]; !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "product not found"})
	}
	products[id] = types.Product{
		ID:      id,
		Name:    in.Name,
		Price:   in.Price,
		InStock: in.InStock,
	}
	return c.JSON(products[id])
}


// PATCH /api/product/:id (partial update)
func PatchProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var patch types.ProductPatch
	if err := c.BodyParser(&patch); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	mu.Lock()
	defer mu.Unlock()
	existing, ok := products[id]
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "product not found"})
	}

	if patch.Name != nil {
		existing.Name = *patch.Name
	}
	if patch.Price != nil {
		if *patch.Price < 0 {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "price >= 0"})
		}
		existing.Price = *patch.Price
	}
	if patch.InStock != nil {
		existing.InStock = *patch.InStock
	}
	products[id] = existing
	return c.JSON(existing)
}


// DELETE /api/product/:id
func DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := products[id]; !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "product not found"})
	}
	delete(products, id)
	return c.JSON(fiber.Map{"message": "deleted"})

}
