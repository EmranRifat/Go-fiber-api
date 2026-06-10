package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	HostListingStatusPending  = "Pending"
	HostListingStatusApproved = "Approved"
	HostListingStatusRejected = "Rejected"
)

type HostListing struct {
	ID                        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	ListingID                 *uuid.UUID     `json:"listingId" gorm:"type:uuid"`
	HostID                    uint           `json:"host_id"`
	Title                     string         `json:"title" gorm:"type:varchar(255);not null"`
	Description               string         `json:"description" gorm:"type:text"`
	Status                    string         `json:"status" gorm:"size:30;default:'Pending'"`
	Photos                    []string       `json:"photos" gorm:"serializer:json;type:jsonb"`
	RentPerNight              string         `json:"rentPerNight" gorm:"column:rent_per_night"`
	Bathrooms                 int            `json:"bathrooms"`
	Bedrooms                  int            `json:"bedrooms"`
	CheckIn                   string         `json:"checkIn" gorm:"column:check_in"`
	CheckOut                  string         `json:"checkOut" gorm:"column:check_out"`
	Facilities                map[string]any `json:"facilities" gorm:"serializer:json;type:jsonb"`
	Kitchens                  int            `json:"kitchens"`
	Latitude                  float64        `json:"latitude"`
	Location                  string         `json:"location" gorm:"type:text"`
	Longitude                 float64        `json:"longitude"`
	PropertyType              string         `json:"propertyType" gorm:"column:property_type"`
	AvailabilitySelectionMode string         `json:"availabilitySelectionMode"`
	AvailableFrom             *time.Time     `json:"availableFrom"`
	AvailableTo               *time.Time     `json:"availableTo"`
	CreatedAt                 time.Time      `json:"created_at"`
	UpdatedAt                 time.Time      `json:"updated_at"`
}
