package config

import (
	"fmt"
	"os"
	"strconv"

	"go-fiber-api/types"
)

func Load() (*types.Config, error) {
	cfg := &types.Config{
		AppPort:         getenv("APP_PORT", "3001"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		JWTExpiresHours: getenvInt("JWT_EXPIRES_HOURS", 72),

		DBHost:     getenv("DB_HOST", "localhost"),
		DBPort:     getenv("DB_PORT", "5432"),
		DBUser:     getenv("DB_USER", "postgres"),
		DBPassword: getenv("DB_PASSWORD", "postgres"),
		DBName:     getenv("DB_NAME", "go_fiber_api"),
		DBSSLMode:  getenv("DB_SSLMODE", "disable"),
		DBTZ:       getenv("DB_TIMEZONE", "Asia/Dhaka"),
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

// 👇 MUST be outside Load()

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func getenvInt(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}







// package config

// import (
// 	"os"
// 	"strconv"
// 	"go-fiber-api/types"
// )


// func Load() *types.Config {
// 	return &types.Config{
// 		// app
// 		AppPort:         getenv("APP_PORT", "3001"),
// 		// jwt
// 		JWTSecret:       getenv("JWT_SECRET", ""),
// 		JWTExpiresHours: getenvInt("JWT_EXPIRES_HOURS", 72),
// 		// db
// 		DBHost:     getenv("DB_HOST", "localhost"),
// 		DBPort:     getenv("DB_PORT", "5432"),
// 		DBUser:     getenv("DB_USER", "postgres"),
// 		DBPassword: getenv("DB_PASSWORD", "postgres"),
// 		DBName:     getenv("DB_NAME", "go_fiber_api"),
// 		DBSSLMode:  getenv("DB_SSLMODE", "disable"),
// 		DBTZ:       getenv("DB_TIMEZONE", "Asia/Dhaka"),
// 	}
// }


// func getenv(k, def string) string {
// 	if v := os.Getenv(k); v != "" {
// 		return v
// 	}
// 	return def
// }


// func getenvInt(k string, def int) int {
// 	if v := os.Getenv(k); v != "" {
// 		if n, err := strconv.Atoi(v); err == nil {
// 			return n
// 		}
// 	}
// 	return def
// }




