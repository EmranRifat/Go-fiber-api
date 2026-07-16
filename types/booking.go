// types/booking.go
package types

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// StringOrNumber allows JSON string or numeric values for fields that are
// ultimately stored as strings in the backend.
type StringOrNumber string

func (s *StringOrNumber) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		*s = ""
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*s = StringOrNumber(str)
		return nil
	}

	var num json.Number
	if err := json.Unmarshal(data, &num); err == nil {
		*s = StringOrNumber(num.String())
		return nil
	}

	return fmt.Errorf("invalid string or numeric value")
}

// CreateBookingRequest matches the client payload sent to POST /api/bookings.
// Fields are intentionally nested to match the frontend contract.
type CreateBookingRequest struct {
	BookingID     StringOrNumber  `json:"booking_id"` // client-generated optional id
	ListingID     string          `json:"listing_id"`
	PaymentMethod string  `json:"payment_method"`
	CheckIn       string  `json:"check_in"`
	CheckOut      string  `json:"check_out"`
	Adults        int     `json:"adults"`
	Children      int     `json:"children"`
	TotalAmount   float64 `json:"total_amount"`
	Currency      string  `json:"currency"`

	// Listing snapshot
	ProductTitle   string `json:"product_title"`
	ProductImage   string `json:"product_image"`
	Category       string `json:"category"`
	ProductAddress string `json:"product_address"`

	// Guest information (nested object on the wire)
	UserInformation struct {
		Name  string `json:"name"`
		Role  string `json:"role"`
		Phone string `json:"phone"`
		Email string `json:"email"`
	} `json:"user_information"`

	// Guest address (nested object on the wire)
	UserAddres struct {
		Street  string `json:"street"`
		City    string `json:"city"`
		Zip     string `json:"zip"`
		Country string `json:"country"`
	} `json:"user_addres"`

	// Payment
	CardLast4      string `json:"card_last4"`
	CardExpiration string `json:"card_expiration"`
	TermsAccepted  bool   `json:"terms_accepted"`
}