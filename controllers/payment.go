// controllers/payment.go
package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"go-fiber-api/models"
	"go-fiber-api/services"
)

type InitPaymentRequest struct {
	// Accept both numeric ("34") and string ("34") from clients.
	BookingID     string `json:"booking_id"`
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
	CustomerPhone string `json:"customer_phone"`
}

func InitSSLPayment(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req InitPaymentRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid request body",
			})
		}

		if strings.TrimSpace(req.BookingID) == "" {
			return c.Status(422).JSON(fiber.Map{
				"status":  "error",
				"message": "booking_id is required",
			})
		}

		// Resolve either DB id (uint) or client "booking_id" string.
		var booking models.Booking
		bookingDBID, err := strconv.ParseUint(strings.TrimSpace(req.BookingID), 10, 64)
		if err == nil && bookingDBID > 0 {
			if dbErr := db.First(&booking, uint(bookingDBID)).Error; dbErr != nil {
				return c.Status(404).JSON(fiber.Map{
					"status":  "error",
					"message": "Booking not found",
				})
			}
		} else {
			if dbErr := db.Where("booking_id = ?", req.BookingID).First(&booking).Error; dbErr != nil {
				return c.Status(404).JSON(fiber.Map{
					"status":  "error",
					"message": "Booking not found",
				})
			}
		}

		if strings.TrimSpace(strings.ToLower(booking.PaymentMethod)) != "sslcommerz" {
			// If the booking was created as manual but user is initiating an SSLCommerz
			// payment flow, update the booking's payment method instead of creating
			// a new booking row. This keeps a single booking record per reservation.
			booking.PaymentMethod = "sslcommerz"
			if err := db.Model(&booking).Update("payment_method", booking.PaymentMethod).Error; err != nil {
				return c.Status(500).JSON(fiber.Map{
					"status":  "error",
					"message": "Failed to update booking payment method",
				})
			}
		}

		var paidCount int64
		if err := db.Model(&models.Payment{}).
			Where("booking_id = ? AND status = ?", booking.ID, "paid").
			Count(&paidCount).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to check existing payments",
			})
		}
		if paidCount > 0 {
			return c.Status(409).JSON(fiber.Map{
				"status":  "error",
				"message": "Booking is already paid",
			})
		}

		amount := booking.TotalAmount
		if amount <= 0 {
			return c.Status(422).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid payment amount",
			})
		}

		currency := booking.Currency
		if currency == "" {
			currency = "BDT"
		}

		if err := db.Model(&models.Payment{}).
			Where("booking_id = ? AND status = ?", booking.ID, "pending").
			Update("status", "cancelled").Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to cancel pending payments",
			})
		}

		tranID := fmt.Sprintf("BOOKING_%d_%d", booking.ID, time.Now().Unix())

		payment := models.Payment{
			BookingID:     booking.ID,
			TransactionID: tranID,
			Amount:        amount,
			Currency:      currency,
			Status:        "pending",
			Gateway:       "sslcommerz",
		}

		if err := db.Create(&payment).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Payment record create failed",
			})
		}

		customerName := req.CustomerName
		if customerName == "" {
			customerName = booking.UserName
		}
		if customerName == "" {
			customerName = "Customer"
		}

		customerEmail := req.CustomerEmail
		if customerEmail == "" {
			customerEmail = "customer@example.com"
		}

		customerPhone := req.CustomerPhone
		if customerPhone == "" {
			customerPhone = "01700000000"
		}

		sslResp, err := services.CreateSSLSession(
			tranID,
			amount,
			currency,
			customerName,
			customerEmail,
			customerPhone,
		)
		if err != nil {
			_ = db.Model(&payment).Update("status", "failed")
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"status":       "success",
			"gateway_url":  sslResp.GatewayPageURL,
			"tran_id":      tranID,
			"payment_id":   payment.ID,
			"booking_id":   booking.ID,
			"amount":       amount,
			"currency":     currency,
		})
	}
}
 