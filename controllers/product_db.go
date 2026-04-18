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

var (
	mu sync.RWMutex
	products = map[int]types.Product{
		1: {ID: 1, Name: "Mechanical Keyboard", Price: 49.99, InStock: true},
		2: {ID: 2, Name: "Wireless Mouse", Price: 19.99, InStock: true},
		3: {ID: 3, Name: "27\" Monitor", Price: 199.00, InStock: false},
	}
	nextID = 4
)


func GetListingDataDB(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page := c.QueryInt("page", 1)
		limit := c.QueryInt("limit", 10)
		q := strings.TrimSpace(c.Query("q"))
		category := strings.TrimSpace(c.Query("category"))

		if page < 1 {
			page = 1
		}
		if limit < 1 {
			limit = 10
		}

		offset := (page - 1) * limit

		var listings []models.Listing
		var total int64

		query := db.Model(&models.Listing{})

		if q != "" {
			search := "%" + strings.ToLower(q) + "%"
			query = query.Where(
				`LOWER(title) LIKE ? OR LOWER(description) LIKE ? OR LOWER(city) LIKE ? OR LOWER(country) LIKE ? OR LOWER(host_name) LIKE ?`,
				search, search, search, search, search,
			)
		}

		if category != "" {
			query = query.Where("LOWER(category) = ?", strings.ToLower(category))
		}

		if err := query.Count(&total).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":     "error",
				"statusCode": fiber.StatusInternalServerError,
				"message":    "Failed to count listings",
				"error":      err.Error(),
			})
		}

		if err := query.
			Offset(offset).
			Limit(limit).
			Find(&listings).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":     "error",
				"statusCode": fiber.StatusInternalServerError,
				"message":    "Failed to fetch listings",
				"error":      err.Error(),
			})
		}

		totalPages := int((total + int64(limit) - 1) / int64(limit))

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":     "success",
			"statusCode": fiber.StatusOK,
			"message":    "Listings fetched successfully",
			"data":       listings,
			"meta": fiber.Map{
				"page":       page,
				"limit":      limit,
				"total":      total,
				"totalPages": totalPages,
			},
		})
	}
}





func GetListingByIDDB(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		var p models.Listing

		if err := db.Preload("ProductCategory").
			First(&p, "id = ?", id).Error; err != nil {

			if err == gorm.ErrRecordNotFound {
				return c.Status(404).JSON(fiber.Map{"error": "product not found"})
			}
			return c.Status(500).JSON(fiber.Map{"error": "db error"})
		}

		return c.JSON(p)
	}
}



// CreateProductDB returns a Fiber handler that creates a product in DB.
func CreateListingDB(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in types.ProductInput
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}
		if in.Name == "" || in.Price < 0 {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "name required, price >= 0"})
		}

		p := models.Listing{
			Title:    in.Name,
			PricePerNight:   in.Price,
			// InStock: in.InStock,
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
