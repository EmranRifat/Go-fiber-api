package dashboard

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"go-fiber-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateHostListingStatusRequest struct {
	Status string `json:"status"`
}

func GetAdminHostListingsHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := requireAdmin(c, db); err != nil {
			return c.Status(adminErrorStatus(err)).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		query := db.Model(&models.HostListing{})

		if idParam := strings.TrimSpace(c.Query("id")); idParam != "" {
			hostListingID, err := uuid.Parse(idParam)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"success": false,
					"message": "Invalid host listing id",
					"error":   err.Error(),
				})
			}
			query = query.Where("id = ?", hostListingID)
		}

		if statusParam := strings.TrimSpace(c.Query("status")); statusParam != "" {
			status, err := normalizeHostListingStatus(statusParam)
			if err != nil {
				return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
					"success": false,
					"message": "Validation failed",
					"error":   err.Error(),
				})
			}
			query = query.Where("status = ?", status)
		}

		// ordering: support `order=created|updated` (default created)
		orderBy := strings.ToLower(strings.TrimSpace(c.Query("order")))
		if orderBy == "updated" || orderBy == "updated_at" {
			query = query.Order("updated_at DESC")
		} else {
			query = query.Order("created_at DESC")
		}

		// if latest=true (or latest=1) then limit to 1
		latest := strings.TrimSpace(c.Query("latest"))

		var hostListings []models.HostListing
		if latest == "true" || latest == "1" {
			query = query.Limit(1)
		}

		if err := query.Find(&hostListings).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Failed to load host listings",
				"error":   err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"count":   len(hostListings),
			"data":    hostListings,
		})
	}
}






func UpdateHostListingStatusHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := requireAdmin(c, db); err != nil {
			return c.Status(adminErrorStatus(err)).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		hostListingID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Invalid host listing id",
				"error":   err.Error(),
			})
		}

		var req UpdateHostListingStatusRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Invalid request body",
				"error":   err.Error(),
			})
		}

		statusValue := strings.TrimSpace(req.Status)
		if statusValue == "" {
			statusValue = strings.TrimSpace(c.Query("status"))
		}
		if statusValue == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Status is required",
			})
		}

		status, err := normalizeHostListingStatus(statusValue)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"message": "Validation failed",
				"error":   err.Error(),
			})
		}

		var hostListing models.HostListing
		if err := db.First(&hostListing, "id = ?", hostListingID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"success": false,
					"message": "Host listing not found",
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Failed to load host listing",
				"error":   err.Error(),
			})
		}

		var listing *models.Listing
		err = db.Transaction(func(tx *gorm.DB) error {
			hostListing.Status = status

			if status == models.HostListingStatusApproved && hostListing.ListingID == nil {
				var host models.User
				if err := tx.First(&host, hostListing.HostID).Error; err != nil {
					return err
				}

				newListing, err := listingFromHostListing(hostListing, host)
				if err != nil {
					return err
				}

				newListing.Status = status

				if err := tx.Create(&newListing).Error; err != nil {
					return err
				}

				hostListing.ListingID = &newListing.ID
				listing = &newListing
			} else if hostListing.ListingID != nil {
				var existingListing models.Listing
				err := tx.First(&existingListing, "id = ?", *hostListing.ListingID).Error
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						hostListing.ListingID = nil
						if status == models.HostListingStatusApproved {
							var host models.User
							if err := tx.First(&host, hostListing.HostID).Error; err != nil {
								return err
							}

							newListing, err := listingFromHostListing(hostListing, host)
							if err != nil {
								return err
							}

							newListing.Status = status
							if err := tx.Create(&newListing).Error; err != nil {
								return err
							}

							hostListing.ListingID = &newListing.ID
							listing = &newListing
						}
						// continue and save hostListing below
					} else {
						return err
					}
				}

				existingListing.Status = status
				if err := tx.Save(&existingListing).Error; err != nil {
					return err
				}

				listing = &existingListing
			}

			if err := tx.Save(&hostListing).Error; err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Failed to update host listing status",
				"error":   err.Error(),
			})
		}

		updatedHostListing := hostListing
		if err := db.First(&updatedHostListing, "id = ?", hostListingID).Error; err == nil {
			response := fiber.Map{
				"success":     true,
				"message":     "Host listing status updated",
				"hostListing": updatedHostListing,
			}
			if listing != nil {
				response["listing"] = listing
			}
			return c.JSON(response)
		}

		return c.JSON(fiber.Map{
			"success":     true,
			"message":     "Host listing status updated",
			"hostListing": hostListing,
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

func requireAdmin(c *fiber.Ctx, db *gorm.DB) error {
	userID, err := userIDFromContext(c)
	if err != nil {
		return fmt.Errorf("unauthorized: %w", err)
	}

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return errors.New("unauthorized: admin user not found")
	}

	role := strings.TrimSpace(user.Role)

	switch {
	case strings.EqualFold(role, "admin"),
		strings.EqualFold(role, "superadmin"):
		return nil
	default:
		return errors.New("only admin or superadmin can access host listings")
	}
}

func adminErrorStatus(err error) int {
	message := strings.ToLower(err.Error())
	if strings.Contains(message, "unauthorized") || strings.Contains(message, "missing") || strings.Contains(message, "invalid") {
		return fiber.StatusUnauthorized
	}

	return fiber.StatusForbidden
}

func normalizeHostListingStatus(status string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "pending":
		return models.HostListingStatusPending, nil
	case "approved", "approve":
		return models.HostListingStatusApproved, nil
	case "rejected", "reject":
		return models.HostListingStatusRejected, nil
	default:
		return "", errors.New("status must be Pending, Approved, or Rejected")
	}
}

func listingFromHostListing(hostListing models.HostListing, host models.User) (models.Listing, error) {
	rentPerNight, err := strconv.ParseFloat(strings.TrimSpace(hostListing.RentPerNight), 64)
	if err != nil || rentPerNight < 0 {
		return models.Listing{}, errors.New("rentPerNight must be a valid positive number")
	}

	image := ""
	if len(hostListing.Photos) > 0 {
		image = hostListing.Photos[0]
	}

	isSuperhost := strings.EqualFold(strings.TrimSpace(host.Role), "host")

	return models.Listing{
		ID:                        uuid.New(),
		HostID:                    host.ID,
		Title:                     strings.TrimSpace(hostListing.Title),
		Description:               strings.TrimSpace(hostListing.Description),
		PricePerNight:             rentPerNight,
		Currency:                  "BDT",
		Address:                   strings.TrimSpace(hostListing.Location),
		Location:                  models.Location{Lat: hostListing.Latitude, Lng: hostListing.Longitude},
		Image:                     image,
		Images:                    hostListing.Photos,
		Category:                  strings.TrimSpace(hostListing.PropertyType),
		HostName:                  host.Name,
		IsSuperhost:               isSuperhost,
		Host:                      models.Host{Name: host.Name, IsSuperhost: isSuperhost},
		Details:                   models.Details{Bedrooms: hostListing.Bedrooms, Bathrooms: hostListing.Bathrooms},
		Kitchens:                  hostListing.Kitchens,
		CheckIn:                   strings.TrimSpace(hostListing.CheckIn),
		CheckOut:                  strings.TrimSpace(hostListing.CheckOut),
		Facilities:                mapAnyToString(hostListing.Facilities),
		Availability:              true,
		AvailabilitySelectionMode: hostListing.AvailabilitySelectionMode,
		AvailableFrom:             hostListing.AvailableFrom,
		AvailableTo:               hostListing.AvailableTo,
	}, nil
}

func mapAnyToString(input map[string]any) map[string]string {
	result := make(map[string]string)

	for key, value := range input {
		result[key] = fmt.Sprintf("%v", value)
	}

	return result
}
