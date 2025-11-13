package models

import "time"

type Order struct {
	ID           uint      `json:"id"            gorm:"primaryKey"`
	OrderID      string    `json:"order_id"      gorm:"size:100;uniqueIndex;not null"`
	CustomerName string    `json:"customer_name" gorm:"size:150;not null"`
	Status       string    `json:"status"        gorm:"size:50;not null"`
	TotalAmount  float64   `json:"total_amount"  gorm:"not null;check:total_amount >= 0"`
	Currency     string    `json:"currency"      gorm:"size:10;not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
}
