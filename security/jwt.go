package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims defines what you put inside the token.
// Add more fields later (e.g., roles) as needed.
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}


// JWTManager signs and verifies JWTs.
type JWTManager struct {
	secret []byte
	ttl    time.Duration
	iss    string
}



// NewJWTManager creates a signer/verifier.
// `hours` is the token lifetime (defaults to 72 if <= 0).
func NewJWTManager(secret string, hours int) *JWTManager {
	if hours <= 0 {
		hours = 72
	}
	return &JWTManager{
		secret: []byte(secret),
		ttl:    time.Duration(hours) * time.Hour,
		iss:    "go-fiber-api", 		
	}
}



// Sign creates a JWT for the given user id + email.
func (j *JWTManager) Sign(userID int, email string) (string, error) {
	now := time.Now()
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprint(userID),            // user id as string
			Issuer:    j.iss,                         // issuer
			IssuedAt:  jwt.NewNumericDate(now),       // iat
			ExpiresAt: jwt.NewNumericDate(now.Add(j.ttl)), //exp
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}



// Parse verifies the token and returns its claims.
func (j *JWTManager) Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		// ensure HMAC (HS256/384/512) is used
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if c, ok := token.Claims.(*Claims); ok && token.Valid {
		return c, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
