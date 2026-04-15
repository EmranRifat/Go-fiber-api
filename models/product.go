package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	ProductCategoryID uuid.UUID `json:"product_category_id" gorm:"type:uuid;not null"`
	Name              string    `json:"name"`
	Price             float64   `json:"price"`
	Image             string    `json:"image"`
	Description       string    `json:"description"`
	Manufacturer      string    `json:"manufacturer"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	ProductCategory ProductCategory `json:"product_category" gorm:"foreignKey:ProductCategoryID;references:ID"`
}