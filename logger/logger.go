package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	errorLogger = log.New(os.Stderr, "❌ ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger  = log.New(os.Stdout, "ℹ️  INFO: ", log.Ldate|log.Ltime)
)

// Error logs error messages
func Error(message string, err error) {
	if err != nil {
		errorLogger.Printf("%s: %v\n", message, err)
	} else {
		errorLogger.Println(message)
	}
}

// Info logs info messages
func Info(message string) {
	infoLogger.Println(message)
}

// Success logs success messages
func Success(message string) {
	fmt.Printf("✅ %s\n", message)
}
