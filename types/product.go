package types
type Product struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	InStock bool    `json:"in_stock"`
}