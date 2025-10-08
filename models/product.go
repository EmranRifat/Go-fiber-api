package models

import "time"

type Product struct {
	ID        uint      `json:"id"        gorm:"primaryKey"`
	Name      string    `json:"name"      gorm:"size:120;not null"`
	Price     float64   `json:"price"     gorm:"not null;check:price >= 0"`
	InStock   bool      `json:"in_stock"  gorm:"not null;default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
