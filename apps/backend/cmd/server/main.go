package main

import (
	"backend/db"
	_ "backend/docs"
	"backend/internal"
	router "backend/internal/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

// @title My API Name
// @version 1.0
// @host localhost:3000
// @BasePath /
func main() {
	//load env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// Database Setup
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			log.Printf("Error disconnecting: %v", err)
		}
	}()

	// RabbitMQ
	internal.SetupRabbitMq()

	fmt.Println("Backend Running on :3000")

	r := router.Router(client)

	log.Fatal(http.ListenAndServe(":3000", r))
}
