package controllers

import (
	"errors"
	"math"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"go-fiber-api/models"
	"go-fiber-api/services"
)

var (
	errTranMismatch    = errors.New("transaction id mismatch")
	errPaymentNotFound = errors.New("payment not found")
	errAmountMismatch  = errors.New("amount mismatch")
	errCurrencyMismatch = errors.New("currency mismatch")
)

func callbackValID(c *fiber.Ctx) string {
	if v := c.FormValue("val_id"); v != "" {
		return v
	}
	return c.Query("val_id")
}

func callbackTranID(c *fiber.Ctx) string {
	if v := c.FormValue("tran_id"); v != "" {
		return v
	}
	return c.Query("tran_id")
}

func amountsMatch(expected float64, actual string) bool {
	parsed, err := strconv.ParseFloat(actual, 64)
	if err != nil {
		return false
	}
	return math.Abs(expected-parsed) < 0.01
}

func finalizeSSLPayment(db *gorm.DB, valID, tranID string) (*models.Payment, error) {
	validated, err := services.ValidateSSLTransaction(valID)
	if err != nil {
		return nil, err
	}

	if validated.TranID != tranID {
		return nil, errTranMismatch
	}

	var payment models.Payment
	if err := db.Where("transaction_id = ?", tranID).First(&payment).Error; err != nil {
		return nil, errPaymentNotFound
	}

	if payment.Status == "paid" {
		return &payment, nil
	}

	if !amountsMatch(payment.Amount, validated.Amount) {
		return nil, errAmountMismatch
	}

	if validated.Currency != "" && payment.Currency != "" && validated.Currency != payment.Currency {
		return nil, errCurrencyMismatch
	}

	status := "paid"
	if validated.RiskLevel == "1" {
		status = "on_hold"
	}

	payment.Status = status
	payment.ValidationID = validated.ValID
	if err := db.Save(&payment).Error; err != nil {
		return nil, err
	}

	return &payment, nil
}

func markPaymentStatus(db *gorm.DB, tranID, status string) error {
	return db.Model(&models.Payment{}).
		Where("transaction_id = ? AND status = ?", tranID, "pending").
		Update("status", status).Error
}

func frontendRedirect(path string, tranID string) string {
	base := os.Getenv("FRONTEND_URL")
	if base == "" {
		base = "http://localhost:3000"
	}
	url := base + path
	if tranID != "" {
		url += "?tran_id=" + tranID
	}
	return url
}

func SSLSuccess(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		valID := callbackValID(c)
		tranID := callbackTranID(c)

		if valID == "" || tranID == "" {
			return c.Redirect(frontendRedirect("/payment/failed", tranID), fiber.StatusFound)
		}

		if _, err := finalizeSSLPayment(db, valID, tranID); err != nil {
			return c.Redirect(frontendRedirect("/payment/failed", tranID), fiber.StatusFound)
		}

		return c.Redirect(frontendRedirect("/payment/success", tranID), fiber.StatusFound)
	}
}

func SSLFail(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tranID := callbackTranID(c)
		_ = markPaymentStatus(db, tranID, "failed")
		return c.Redirect(frontendRedirect("/payment/failed", tranID), fiber.StatusFound)
	}
}

func SSLCancel(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tranID := callbackTranID(c)
		_ = markPaymentStatus(db, tranID, "cancelled")
		return c.Redirect(frontendRedirect("/payment/cancelled", tranID), fiber.StatusFound)
	}
}

func SSLIPN(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		valID := callbackValID(c)
		tranID := callbackTranID(c)
		status := c.FormValue("status")

		if valID == "" || tranID == "" {
			return c.SendString("INVALID")
		}

		if status == "FAILED" {
			_ = markPaymentStatus(db, tranID, "failed")
			return c.SendString("FAILED")
		}

		if status == "CANCELLED" {
			_ = markPaymentStatus(db, tranID, "cancelled")
			return c.SendString("CANCELLED")
		}

		if _, err := finalizeSSLPayment(db, valID, tranID); err != nil {
			return c.SendString("INVALID")
		}

		return c.SendString("VALID")
	}
}
