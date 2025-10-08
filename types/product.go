package types

import "time"

type Products struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	InStock     bool    `json:"in_stock"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Brand       string  `json:"brand"`
}

type Product struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	InStock bool    `json:"in_stock"`
}

// For POST (create) and PUT (full update)
type ProductInput struct {
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	InStock bool    `json:"in_stock"`
}

// For PATCH (partial update)
type ProductPatch struct {
	Name    *string  `json:"name,omitempty"`
	Price   *float64 `json:"price,omitempty"`
	InStock *bool    `json:"in_stock,omitempty"`
}

// type Config struct {
// 	AppPort         string
// 	JWTSecret       string
// 	JWTExpiresHours int
// }

type JWTManager struct {
	secret []byte
	ttl    time.Duration
	iss    string
}