package router

import (
	"backend/db"
	"backend/internal/handlers"

	"github.com/gorilla/mux"
)

func Router(database *db.PrismaClient) *mux.Router {
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/api/health", handlers.HealthCheck).Methods("GET")

	// Website endpoints
	router.HandleFunc("/api/websites", handlers.CreateWebsite(database)).Methods("POST")

	return router
}