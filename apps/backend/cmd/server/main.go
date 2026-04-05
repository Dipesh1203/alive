package main

import (
	"backend/db"
	_ "backend/docs"
	router "backend/internal/routes"
	"backend/internal/utils"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

// @title My API Name
// @version 1.0
// @host localhost:8000
// @BasePath /
func main() {
	log.Printf("[MAIN] Starting backend server initialization...")

	//load env
	if err := godotenv.Load(); err != nil {
		log.Printf("[MAIN] Warning: No .env file found - using system environment variables")
	} else {
		log.Printf("[MAIN] Successfully loaded .env file")
	}

	// Database Setup
	log.Printf("[MAIN] Initializing database connection...")
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		log.Fatalf("[MAIN] FATAL: Failed to connect to database - %v", err)
	}
	log.Printf("[MAIN] Database connection established successfully")

	defer func() {
		log.Printf("[MAIN] Closing database connection...")
		if err := client.Prisma.Disconnect(); err != nil {
			log.Printf("[MAIN] Error disconnecting from database: %v", err)
		}
		log.Printf("[MAIN] Database connection closed")
	}()

	// RabbitMQ
	// internal.SetupRabbitMq()
	// ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	// defer stop()
	// go services.StartMonitoring(ctx, client)
	log.Printf("[MAIN] Setting up router...")
	r := router.Router(client)
	handler := utils.CorsMiddleware(r)
	log.Printf("[MAIN] Router setup complete")

	log.Printf("[MAIN] ========================================")
	log.Printf("[MAIN] Backend server running on :8000")
	log.Printf("[MAIN] ========================================")

	if err := http.ListenAndServe(":8000", handler); err != nil {
		log.Fatalf("[MAIN] FATAL: Server error - %v", err)
	}
}
