// types/booking.go
package types

type CreateBookingRequest struct {
	ListingID     string  `json:"listing_id"`
	PaymentMethod string  `json:"payment_method"`
	CheckIn       string  `json:"check_in"`
	CheckOut      string  `json:"check_out"`
	Adults        int     `json:"adults"`
	Children      int     `json:"children"`
	TotalAmount   float64 `json:"total_amount"`
	Currency      string  `json:"currency"`

	BillingAddress struct {
		Street  string `json:"street"`
		City    string `json:"city"`
		Zip     string `json:"zip"`
		Country string `json:"country"`
	} `json:"billing_address"`

	CardLast4      string `json:"card_last4"`
	CardExpiration string `json:"card_expiration"`
	TermsAccepted  bool   `json:"terms_accepted"`
}