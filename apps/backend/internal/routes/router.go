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

	// Monitoring endpoint
	router.HandleFunc("/api/monitoring/{id}", handlers.ToggleMonitoring(database)).Methods("POST")

	// Region endpoints
	router.HandleFunc("/api/regions", handlers.CreateRegion(database)).Methods("POST")
	router.HandleFunc("/api/regions", handlers.ListRegions(database)).Methods("GET")
	router.HandleFunc("/api/regions/{id}", handlers.GetRegion(database)).Methods("GET")
	router.HandleFunc("/api/regions/{id}", handlers.UpdateRegion(database)).Methods("PUT")
	router.HandleFunc("/api/regions/{id}", handlers.DeleteRegion(database)).Methods("DELETE")

	//Detailed Information about a website endpoints
	router.HandleFunc("/api/websites/{id}/details", handlers.GetDetailsWebsite(database)).Methods("GET")

	router.HandleFunc("/api/websites/{id}/regions", handlers.AssignWebsiteRegions(database)).Methods("PUT")

	router.HandleFunc("/api/incidents", handlers.ListIncidents(database)).Methods("GET")

	// TODO(frontend): Needed for Notifications screen
	// router.HandleFunc("/api/notifications/channels", handlers.CreateNotificationChannel(database)).Methods("POST")
	// router.HandleFunc("/api/notifications/channels", handlers.ListNotificationChannels(database)).Methods("GET")
	// router.HandleFunc("/api/notifications/channels/{id}", handlers.UpdateNotificationChannel(database)).Methods("PUT")
	// router.HandleFunc("/api/notifications/channels/{id}", handlers.DeleteNotificationChannel(database)).Methods("DELETE")
	// router.HandleFunc("/api/notifications/preferences", handlers.UpdateNotificationPreferences(database)).Methods("PUT")

	// TODO(frontend): Needed for Team screen
	// router.HandleFunc("/api/team/members", handlers.ListTeamMembers(database)).Methods("GET")
	// router.HandleFunc("/api/team/invitations", handlers.CreateTeamInvitation(database)).Methods("POST")
	// router.HandleFunc("/api/team/members/{id}", handlers.UpdateTeamMemberRole(database)).Methods("PUT")
	// router.HandleFunc("/api/team/members/{id}", handlers.RemoveTeamMember(database)).Methods("DELETE")

	// TODO(frontend): Needed for Settings screen
	// router.HandleFunc("/api/settings/profile", handlers.GetProfileSettings(database)).Methods("GET")
	// router.HandleFunc("/api/settings/profile", handlers.UpdateProfileSettings(database)).Methods("PUT")
	// router.HandleFunc("/api/settings/workspace", handlers.GetWorkspaceSettings(database)).Methods("GET")
	// router.HandleFunc("/api/settings/workspace", handlers.UpdateWorkspaceSettings(database)).Methods("PUT")
	// router.HandleFunc("/api/settings/monitoring-defaults", handlers.GetMonitoringDefaults(database)).Methods("GET")
	// router.HandleFunc("/api/settings/monitoring-defaults", handlers.UpdateMonitoringDefaults(database)).Methods("PUT")
	// router.HandleFunc("/api/settings/api-keys", handlers.ListAPIKeys(database)).Methods("GET")
	// router.HandleFunc("/api/settings/api-keys", handlers.CreateAPIKey(database)).Methods("POST")
	// router.HandleFunc("/api/settings/api-keys/{id}", handlers.RevokeAPIKey(database)).Methods("DELETE")

	router.HandleFunc("/api/test", handlers.GetTest(database)).Methods("POST")

	return router
}
