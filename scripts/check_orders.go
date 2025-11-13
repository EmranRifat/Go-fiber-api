package main

import (
	"fmt"
	"go-fiber-api/database"
	"go-fiber-api/models"
	"log"
)

func main() {
	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Query all orders
	var orders []models.Order
	if err := db.Find(&orders).Error; err != nil {
		log.Fatal("Failed to query orders:", err)
	}

	// Display orders
	fmt.Printf("\n📋 Total Orders in Database: %d\n", len(orders))
	fmt.Println("==========================================")
	for _, order := range orders {
		fmt.Printf("\nOrder ID: %s\n", order.OrderID)
		fmt.Printf("Customer: %s\n", order.CustomerName)
		fmt.Printf("Status: %s\n", order.Status)
		fmt.Printf("Amount: %.2f %s\n", order.TotalAmount, order.Currency)
		fmt.Printf("Created: %s\n", order.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println("------------------------------------------")
	}
}
