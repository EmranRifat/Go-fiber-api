// controllers/booking.go
package controllers

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"go-fiber-api/models"
	"go-fiber-api/types"
)


func CreateBookingDB(db *gorm.DB) fiber.Handler { 
	return func(c *fiber.Ctx) error {
		var in types.CreateBookingRequest

		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid body",
			})
		}

		if in.ListingID == "" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": "listing_id is required",
			})
		}
		if in.PaymentMethod == "" {
			in.PaymentMethod = "manual"
		}

		if in.Adults < 1 {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": "at least 1 adult is required",
			})
		}

		if in.TotalAmount < 0 {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": "total_amount must be >= 0",
			})
		}

		if in.Currency == "" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": "currency is required",
			})
		}

		if !in.TermsAccepted {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": "terms must be accepted",
			})
		}

		checkIn, err := time.Parse("2006-01-02", in.CheckIn)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid check_in format, use YYYY-MM-DD",
			})
		}

		checkOut, err := time.Parse("2006-01-02", in.CheckOut)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid check_out format, use YYYY-MM-DD",
			})
		}

		if !checkOut.After(checkIn) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": "check_out must be after check_in",
			})
		}

		// Check duplicate booking
		var existing models.Booking

		err = db.Where(
			"listing_id = ? AND check_in = ? AND check_out = ? AND adults = ? AND children = ? AND total_amount = ? AND currency = ?",
			in.ListingID,
			checkIn,
			checkOut,
			in.Adults,
			in.Children,
			in.TotalAmount,
			in.Currency,
		).First(&existing).Error

		if err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "booking already exists",
			})
		}

		if err != nil && err != gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to check existing booking",
			})
		}

		booking := models.Booking{
			ListingID:     in.ListingID,
			PaymentMethod: in.PaymentMethod,
			CheckIn:       checkIn,
			CheckOut:      checkOut,
			Adults:        in.Adults,
			Children:      in.Children,
			TotalAmount:   in.TotalAmount,
			Currency:      in.Currency,

			BillingStreet:  in.BillingAddress.Street,
			BillingCity:    in.BillingAddress.City,
			BillingZip:     in.BillingAddress.Zip,
			BillingCountry: in.BillingAddress.Country,

			CardLast4:      in.CardLast4,
			CardExpiration: in.CardExpiration,
			TermsAccepted:  in.TermsAccepted,
		}

		if err := db.Create(&booking).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to create booking",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "booking created successfully",
			"data":    booking,
		})
	}
}