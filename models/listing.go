package models

import (
	"time"

	"github.com/google/uuid"
)

type Listing struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	PricePerNight float64   `json:"price_per_night"`
	City          string    `json:"city"`
	Country       string    `json:"country"`
	Image         string    `json:"image"`
	Category      string    `json:"category"`
	Rating        float64   `json:"rating"`
	ReviewsCount  int       `json:"reviews_count"`
	HostName      string    `json:"host_name"`
	IsSuperhost   bool      `json:"is_superhost"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}