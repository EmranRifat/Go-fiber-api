package models

import "time"

type Review struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	ListingID uint      `json:"listing_id"`

	Rating  int    `json:"rating"`
	Comment string `json:"comment"`

	Status string `json:"status"` // Pending, Approved, Rejected

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}