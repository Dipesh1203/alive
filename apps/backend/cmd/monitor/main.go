package main

import (
	"backend/db"
	_ "backend/docs"
	"backend/internal/services"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

// @title My API Name
// @version 1.0
// @host localhost:3001
// @BasePath /
func main() {
	//load env
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	log.Println("Worker: Starting Monitoring Service...")

	services.StartMonitoring(ctx, client)

	log.Println("Worker: Shutting down gracefully")
}
