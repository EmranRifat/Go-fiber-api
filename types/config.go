package types

// Config holds app + auth + db settings
type Config struct {
	AppPort         string
	JWTSecret       string
	JWTExpiresHours int

	// DB
	DBHost    string
	DBPort    string
	DBUser    string
	DBPassword string
	DBName    string
	DBSSLMode string
	DBTZ      string
}
