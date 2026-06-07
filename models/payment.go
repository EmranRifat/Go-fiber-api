// models/payment.go
package models

import "time"

type Payment struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	BookingID     uint      `json:"booking_id"`
	TransactionID string    `gorm:"uniqueIndex" json:"transaction_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `gorm:"default:BDT" json:"currency"`
	Status        string    `gorm:"default:pending" json:"status"` // pending, paid, failed, cancelled
	Gateway       string    `gorm:"default:sslcommerz" json:"gateway"`
	ValidationID  string    `json:"validation_id"`
	UpdatedAt     time.Time `json:"updated_at"`
}