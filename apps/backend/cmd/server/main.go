package main

import (
	"backend/db"
	router "backend/internal/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: No .env file found")
    }

    client := db.NewClient()
    if err := client.Prisma.Connect(); err != nil {
        log.Fatal("Failed to connect to DB:", err)
    }

    defer func() {
        if err := client.Prisma.Disconnect(); err != nil {
            log.Printf("Error disconnecting: %v", err)
        }
    }()

    fmt.Println("Backend Running on :3000")

    r := router.Router(client) 
    
    log.Fatal(http.ListenAndServe(":3000", r))
}