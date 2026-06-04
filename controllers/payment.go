// controllers/payment.go
package controllers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"go-fiber-api/models"
	"go-fiber-api/services"
)

type InitPaymentRequest struct {
	BookingID     uint   `json:"booking_id"`
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

		if req.BookingID == 0 {
			return c.Status(422).JSON(fiber.Map{
				"status":  "error",
				"message": "booking_id is required",
			})
		}

		var booking models.Booking
		if err := db.First(&booking, req.BookingID).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Booking not found",
			})
		}

		var paidCount int64
		db.Model(&models.Payment{}).
			Where("booking_id = ? AND status = ?", booking.ID, "paid").
			Count(&paidCount)
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

		db.Model(&models.Payment{}).
			Where("booking_id = ? AND status = ?", booking.ID, "pending").
			Update("status", "cancelled")

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
			customerName = booking.BillingCity
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
 