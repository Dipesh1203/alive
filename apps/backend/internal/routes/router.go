package router

import (
	"log"

	"github.com/Dipesh1203/alive/apps/backend/internal/handlers"
	"github.com/Dipesh1203/alive/apps/backend/internal/middleware"

	"github.com/Dipesh1203/alive/apps/backend/db"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Router(database *db.PrismaClient) *mux.Router {
	log.Printf("[ROUTER] Initializing router...")
	router := mux.NewRouter()

	log.Printf("[ROUTER] Registering swagger endpoint...")
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Health check endpoint
	log.Printf("[ROUTER] Registering health check endpoint: GET /api/health")
	router.HandleFunc("/api/health", handlers.HealthCheck).Methods("GET")

	// Auth endpoints (no auth required)
	log.Printf("[ROUTER] Registering auth endpoints...")
	router.HandleFunc("/api/auth/signup", handlers.Signup(database)).Methods("POST")
	log.Printf("[ROUTER] Registered: POST /api/auth/signup")
	router.HandleFunc("/api/auth/login", handlers.Login(database)).Methods("POST")
	log.Printf("[ROUTER] Registered: POST /api/auth/login")

	// Public landing endpoints (no auth required)
	log.Printf("[ROUTER] Registering public landing endpoints...")
	router.HandleFunc("/api/public/landing", handlers.GetLandingOverview(database)).Methods("GET")
	router.HandleFunc("/api/public/landing/testimonials", handlers.ListLandingTestimonials(database)).Methods("GET")
	router.HandleFunc("/api/public/landing/pricing", handlers.ListLandingPricing(database)).Methods("GET")
	router.HandleFunc("/api/public/landing/faqs", handlers.ListLandingFAQs(database)).Methods("GET")

	// Protected routes - apply auth middleware
	log.Printf("[ROUTER] Initializing protected routes with auth middleware...")
	protectedRouter := router.PathPrefix("/api").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware)

	// Profile endpoints (protected)
	log.Printf("[ROUTER] Registering profile endpoints...")
	protectedRouter.HandleFunc("/profile", handlers.GetMyProfile(database)).Methods("GET")
	protectedRouter.HandleFunc("/profile", handlers.UpdateMyProfile(database)).Methods("PUT")
	protectedRouter.HandleFunc("/profile/preferences", handlers.GetMyPreferences(database)).Methods("GET")
	protectedRouter.HandleFunc("/profile/preferences", handlers.UpdateMyPreferences(database)).Methods("PUT")

	// Organization endpoints (protected)
	log.Printf("[ROUTER] Registering organization endpoints...")
	protectedRouter.HandleFunc("/organizations", handlers.CreateOrganization(database)).Methods("POST")
	protectedRouter.HandleFunc("/organizations", handlers.ListUserOrganizations(database)).Methods("GET")
	protectedRouter.HandleFunc("/organizations/{id}", handlers.GetOrganization(database)).Methods("GET")
	protectedRouter.HandleFunc("/organizations/{id}", handlers.UpdateOrganization(database)).Methods("PUT")
	protectedRouter.HandleFunc("/organizations/{id}/members", handlers.ListOrganizationMembers(database)).Methods("GET")
	protectedRouter.HandleFunc("/organizations/{id}/members", handlers.AddOrganizationMember(database)).Methods("POST")
	protectedRouter.HandleFunc("/organizations/{id}/members/{memberId}/role", handlers.UpdateMemberRole(database)).Methods("PUT")
	protectedRouter.HandleFunc("/organizations/{id}/members/{memberId}", handlers.RemoveOrganizationMember(database)).Methods("DELETE")

	// Website endpoints (protected - with org filtering)
	log.Printf("[ROUTER] Registering website endpoints...")
	protectedRouter.HandleFunc("/websites", handlers.CreateWebsite(database)).Methods("POST")
	protectedRouter.HandleFunc("/websites", handlers.ListWebsites(database)).Methods("GET")
	protectedRouter.HandleFunc("/websites/{id}", handlers.GetWebsite(database)).Methods("GET")
	protectedRouter.HandleFunc("/websites/{id}", handlers.UpdateWebsite(database)).Methods("PUT")
	protectedRouter.HandleFunc("/websites/{id}", handlers.DeleteWebsite(database)).Methods("DELETE")

	// Monitoring endpoint (protected)
	log.Printf("[ROUTER] Registering monitoring endpoint...")
	protectedRouter.HandleFunc("/monitoring/{id}", handlers.ToggleMonitoring(database)).Methods("POST")

	// Region endpoints (protected)
	log.Printf("[ROUTER] Registering region endpoints...")
	protectedRouter.HandleFunc("/regions", handlers.CreateRegion(database)).Methods("POST")
	protectedRouter.HandleFunc("/regions", handlers.ListRegions(database)).Methods("GET")
	protectedRouter.HandleFunc("/regions/{id}", handlers.GetRegion(database)).Methods("GET")
	protectedRouter.HandleFunc("/regions/{id}", handlers.UpdateRegion(database)).Methods("PUT")
	protectedRouter.HandleFunc("/regions/{id}", handlers.DeleteRegion(database)).Methods("DELETE")

	// Detailed Information about a website endpoints (protected)
	log.Printf("[ROUTER] Registering website details endpoints...")
	protectedRouter.HandleFunc("/websites/{id}/details", handlers.GetDetailsWebsite(database)).Methods("GET")
	protectedRouter.HandleFunc("/websites/{id}/regions", handlers.AssignWebsiteRegions(database)).Methods("PUT")

	// Incidents endpoints (protected)
	log.Printf("[ROUTER] Registering incidents endpoint...")
	protectedRouter.HandleFunc("/incidents", handlers.ListIncidents(database)).Methods("GET")

	// Notifications endpoints (protected)
	log.Printf("[ROUTER] Registering notification endpoints...")
	protectedRouter.HandleFunc("/notifications/channels", handlers.CreateNotificationChannel(database)).Methods("POST")
	protectedRouter.HandleFunc("/notifications/test", handlers.SendNotification(database)).Methods("POST")
	protectedRouter.HandleFunc("/notifications/test-email", handlers.TestEmail(database)).Methods("POST")
	protectedRouter.HandleFunc("/notifications/test-template-email", handlers.TestTemplateEmail(database)).Methods("POST")

	// Test endpoint (protected)
	log.Printf("[ROUTER] Registering test endpoint...")
	protectedRouter.HandleFunc("/test", handlers.GetTest(database)).Methods("POST")

	log.Printf("[ROUTER] Router initialization complete - total endpoints registered")
	return router
}
