package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims 
}

type JWTManager struct {
	secret []byte
	ttl    time.Duration
	iss    string
}

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

func (j *JWTManager) Sign(userID int, email string) (string, error) {
	now := time.Now()
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprint(userID), // ← subject = user id as string
			Issuer:    j.iss,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTManager) Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
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
