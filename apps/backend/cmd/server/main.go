package main

import (
	"log"
	"net/http"

	_ "github.com/Dipesh1203/alive/apps/backend/docs"
	"github.com/Dipesh1203/alive/apps/backend/internal/config"
	router "github.com/Dipesh1203/alive/apps/backend/internal/routes"
	"github.com/Dipesh1203/alive/apps/backend/internal/utils"

	"github.com/Dipesh1203/alive/apps/backend/db"
)

// @title My API Name
// @version 1.0
// @host localhost:8000
// @BasePath /
func main() {
	config.InitConfig()
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
