package router

import (
	"backend/db"
	"backend/internal/handlers"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Router(database *db.PrismaClient) *mux.Router {
	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	// Health check endpoint
	router.HandleFunc("/api/health", handlers.HealthCheck).Methods("GET")

	// Website endpoints
	router.HandleFunc("/api/websites", handlers.CreateWebsite(database)).Methods("POST")
	router.HandleFunc("/api/websites", handlers.ListWebsites(database)).Methods("GET")
	router.HandleFunc("/api/websites/{id}", handlers.GetWebsite(database)).Methods("GET")
	router.HandleFunc("/api/websites/{id}", handlers.UpdateWebsite(database)).Methods("PUT")
	router.HandleFunc("/api/websites/{id}", handlers.DeleteWebsite(database)).Methods("DELETE")

	// Region endpoints
	router.HandleFunc("/api/regions", handlers.CreateRegion(database)).Methods("POST")
	router.HandleFunc("/api/regions", handlers.ListRegions(database)).Methods("GET")
	router.HandleFunc("/api/regions/{id}", handlers.GetRegion(database)).Methods("GET")
	router.HandleFunc("/api/regions/{id}", handlers.UpdateRegion(database)).Methods("PUT")
	router.HandleFunc("/api/regions/{id}", handlers.DeleteRegion(database)).Methods("DELETE")



	return router
}