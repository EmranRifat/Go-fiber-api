// models/booking.go
package models

import "time"

type Booking struct {
	ID uint `json:"id" gorm:"primaryKey"`

	ListingID     string `json:"listing_id"`
	PaymentMethod string `json:"payment_method"`

	CheckIn  time.Time `json:"check_in"`
	CheckOut time.Time `json:"check_out"`

	Adults      int     `json:"adults"`
	Children    int     `json:"children"`
	TotalAmount float64 `json:"total_amount"`
	Currency    string  `json:"currency"`

	BillingStreet  string `json:"billing_street"`
	BillingCity    string `json:"billing_city"`
	BillingZip     string `json:"billing_zip"`
	BillingCountry string `json:"billing_country"`

	CardLast4      string `json:"card_last4"`
	CardExpiration string `json:"card_expiration"`
	TermsAccepted  bool   `json:"terms_accepted"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}