package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"go-fiber-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// func CreateHostListingHandler(db *gorm.DB) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		var req models.HostListing

// 		if err := c.BodyParser(&req); err != nil {
// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 				"success": false,
// 				"message": "Invalid request body",
// 				"error":   err.Error(),
// 			})
// 		}

// 		userID, err := userIDFromContext(c)
// 		if err != nil {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"success": false,
// 				"message": "Unauthorized",
// 				"error":   err.Error(),
// 			})
// 		}

// 		var user models.User
// 		if err := db.First(&user, userID).Error; err != nil {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"success": false,
// 				"message": "User not found",
// 			})
// 		}

// 		if err := validateHostListingRequest(&req); err != nil {
// 			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
// 				"success": false,
// 				"message": "Validation failed",
// 				"error":   err.Error(),
// 			})
// 		}

// 		rentPerNight, err := strconv.ParseFloat(strings.TrimSpace(req.RentPerNight), 64)
// 		if err != nil || rentPerNight < 0 {
// 			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
// 				"success": false,
// 				"message": "Validation failed",
// 				"error":   "rentPerNight must be a valid positive number",
// 			})
// 		}

// 		if req.AvailabilitySelectionMode == "single" {
// 			req.AvailableTo = req.AvailableFrom
// 		}

// 		hostListing := req
// 		hostListing.ID = uuid.New()
// 		hostListing.HostID = user.ID
// 		hostListing.ListingID = nil
// 		hostListing.Status = models.HostListingStatusPending

// 		if err := db.Create(&hostListing).Error; err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"success": false,
// 				"message": "Failed to create host listing",
// 				"error":   err.Error(),
// 			})
// 		}

// 		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 			"success": true,
// 			"message": "Host listing submitted for admin approval",
// 			"data":    hostListing,
// 		})
// 	}
// }



func CreateHostListingHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.HostListing

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Invalid request body",
				"error":   err.Error(),
			})
		}

		for _, photo := range req.Photos {
			if strings.HasPrefix(photo, "blob:") {
				return c.Status(400).JSON(fiber.Map{
					"success": false,
					"message": "Blob URLs are not supported. Upload actual image files as multipart/form-data.",
				})
			}
		}

		userID, err := userIDFromContext(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized",
			})
		}

		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "User not found",
			})
		}

		form, err := c.MultipartForm()
		if err == nil {
			if req.Facilities == nil {
				if facilitiesJSON, ok := form.Value["facilities"]; ok && len(facilitiesJSON) > 0 {
					var parsed map[string]any
					if err := json.Unmarshal([]byte(facilitiesJSON[0]), &parsed); err != nil {
						return c.Status(400).JSON(fiber.Map{
							"success": false,
							"message": "Invalid facilities JSON",
							"error":   err.Error(),
						})
					}
					req.Facilities = parsed
				}
			}

			if form.File["photos"] != nil {
				var photoURLs []string

				for _, file := range form.File["photos"] {
					ext := filepath.Ext(file.Filename)
					filename := uuid.New().String() + ext
					savePath := filepath.Join("uploads", filename)
					host := c.Get("Host")
					if strings.TrimSpace(host) == "" {
						host = c.Hostname()
					}
					publicURL := fmt.Sprintf("%s://%s/uploads/%s", c.Protocol(), host, filename)

					if err := c.SaveFile(file, savePath); err != nil {
						return c.Status(500).JSON(fiber.Map{
							"success": false,
							"message": "Failed to save image",
							"error":   err.Error(),
						})
					}

					photoURLs = append(photoURLs, publicURL)
				}

				req.Photos = photoURLs
			}
		}

		if req.AvailabilitySelectionMode == "single" {
			req.AvailableTo = req.AvailableFrom
		}

		req.HostID = user.ID
		if req.ID == uuid.Nil {
			req.ID = uuid.New()
		}
		req.ListingID = nil
		req.Status = models.HostListingStatusPending

		if err := validateHostListingRequest(&req); err != nil {
			return c.Status(422).JSON(fiber.Map{
				"success": false,
				"message": "Validation failed",
				"error":   err.Error(),
			})
		}

		if err := db.Create(&req).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Failed to create host listing",
				"error":   err.Error(),
			})
		}

		return c.Status(201).JSON(fiber.Map{
			"success": true,
			"message": "Host listing submitted for admin approval",
			"data":    req,
		})
	}
}








func userIDFromContext(c *fiber.Ctx) (uint, error) {
	sub, ok := c.Locals("sub").(string)
	if !ok || strings.TrimSpace(sub) == "" {
		return 0, errors.New("missing user subject")
	}

	id, err := strconv.ParseUint(sub, 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("invalid user subject")
	}

	return uint(id), nil
}

func validateHostListingRequest(req *models.HostListing) error {
	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	req.PropertyType = strings.TrimSpace(req.PropertyType)
	req.Location = strings.TrimSpace(req.Location)
	req.RentPerNight = strings.TrimSpace(req.RentPerNight)
	req.AvailabilitySelectionMode = strings.ToLower(strings.TrimSpace(req.AvailabilitySelectionMode))

	if req.Title == "" {
		return errors.New("title is required")
	}
	if req.PropertyType == "" {
		return errors.New("propertyType is required")
	}
	if req.Location == "" {
		return errors.New("location is required")
	}
	if req.RentPerNight == "" {
		return errors.New("rentPerNight is required")
	}
	if req.AvailableFrom == nil {
		return errors.New("availableFrom is required")
	}
	if req.AvailabilitySelectionMode == "" {
		req.AvailabilitySelectionMode = "range"
	}
	if req.AvailabilitySelectionMode != "range" && req.AvailabilitySelectionMode != "single" {
		return errors.New("availabilitySelectionMode must be single or range")
	}
	if req.AvailabilitySelectionMode == "range" {
		if req.AvailableTo == nil {
			return errors.New("availableTo is required for range availability")
		}
		if req.AvailableTo.Before(*req.AvailableFrom) {
			return errors.New("availableTo must be after availableFrom")
		}
	}

	if req.Facilities == nil {
		req.Facilities = map[string]any{}
	}

	if req.Photos == nil {
		req.Photos = []string{}
	}

	return nil
}
