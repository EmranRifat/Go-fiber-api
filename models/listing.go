package models

import (
	"time"

	"github.com/google/uuid"
)

type Listing struct {
	ID                        uuid.UUID         `json:"id" gorm:"type:uuid;primaryKey"`
	HostID                    uint              `json:"host_id"`
	Title                     string            `json:"title"`
	Description               string            `json:"description"`
	PricePerNight             float64           `json:"price_per_night"`
	Currency                  string            `json:"currency"`
	City                      string            `json:"city"`
	Country                   string            `json:"country"`
	Location                  Location          `json:"location" gorm:"column:location_json;serializer:json;type:jsonb"`
	Address                   string            `json:"address"`
	Image                     string            `json:"image"`
	Images                    []string          `json:"images" gorm:"serializer:json;type:jsonb"`
	Category                  string            `json:"category"`
	Status                    string            `json:"status" gorm:"size:30;default:'Pending'"`
	Rating                    float64           `json:"rating"`
	ReviewsCount              int               `json:"reviews_count"`
	HostName                  string            `json:"host_name"`
	IsSuperhost               bool              `json:"is_superhost"`
	Host                      Host              `json:"host" gorm:"serializer:json;type:jsonb"`
	Details                   Details           `json:"details" gorm:"serializer:json;type:jsonb"`
	Kitchens                  int               `json:"kitchens"`
	CheckIn                   string            `json:"checkIn"`
	CheckOut                  string            `json:"checkOut"`
	Facilities                map[string]string `json:"facilities" gorm:"serializer:json;type:jsonb"`
	Amenities                 []string          `json:"amenities" gorm:"serializer:json;type:jsonb"`
	HouseRules                []string          `json:"house_rules" gorm:"serializer:json;type:jsonb"`
	Availability              bool              `json:"availability"`
	AvailabilitySelectionMode string            `json:"availabilitySelectionMode"`
	AvailableFrom             *time.Time        `json:"availableFrom"`
	AvailableTo               *time.Time        `json:"availableTo"`
	CreatedAt                 time.Time         `json:"created_at"`
	UpdatedAt                 time.Time         `json:"updated_at"`
}

type Host struct {
	Name        string `json:"name"`
	IsSuperhost bool   `json:"is_superhost"`
}

type Details struct {
	Guests    int `json:"guests"`
	Bedrooms  int `json:"bedrooms"`
	Beds      int `json:"beds"`
	Bathrooms int `json:"bathrooms"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
