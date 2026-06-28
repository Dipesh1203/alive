package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/Dipesh1203/alive/apps/backend/docs"
	"github.com/Dipesh1203/alive/apps/backend/internal/services"

	"github.com/Dipesh1203/alive/apps/backend/db"

	"github.com/joho/godotenv"
)

// @title My API Name
// @version 1.0
// @host localhost:3001
// @BasePath /
func main() {
	//load env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			log.Fatal("Error disconnecting: %v", err)
		}
	}()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go services.StartMonitoring(ctx, client)

}
