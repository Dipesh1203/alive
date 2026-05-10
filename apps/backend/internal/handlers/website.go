package handlers

import (
	"backend/db"
	"backend/internal/services"
	"backend/internal/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// HealthCheck godoc
// @Summary      Health check
// @Description  Returns server health status
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /api/health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Printf("[HEALTH] Health check endpoint called")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().String(),
	})
	log.Printf("[HEALTH] Health check response sent")
}

type CreateWebsiteRequest struct {
	WebsiteName    string `json:"websiteName"`
	URL            string `json:"url"`
	OrganizationID string `json:"organizationId"`
}

type UpdateWebsiteReq struct {
	WebsiteName *string `json:"websiteName,omitempty"`
	URL         *string `json:"url,omitempty"`
}

// CreateWebsite godoc
// @Summary      Create a new website
// @Description  Creates a new website for monitoring in the specified organization
// @Tags         websites
// @Accept       json
// @Produce      json
// @Param        Authorization  header   string                 true  "Bearer token"
// @Param        body           body     CreateWebsiteRequest   true  "Create website request"
// @Success      201            {object} db.WebsiteModel
// @Failure      400            {object} map[string]string
// @Failure      403            {object} map[string]string
// @Failure      500            {object} map[string]string
// @Router       /api/websites [post]
func CreateWebsite(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[WEBSITE] New CreateWebsite request received")
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		userID := r.Context().Value("userID").(string)
		log.Printf("[WEBSITE] Processing request for user: %s", userID)

		var req CreateWebsiteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("[WEBSITE] ERROR: Failed to decode request body - %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.WebsiteName == "" || req.URL == "" {
			log.Printf("[WEBSITE] ERROR: websiteName or URL is empty")
			http.Error(w, "websiteName and url are required", http.StatusBadRequest)
			return
		}

		if req.OrganizationID == "" {
			log.Printf("[WEBSITE] ERROR: organizationId is empty")
			http.Error(w, "organizationId is required", http.StatusBadRequest)
			return
		}

		log.Printf("[WEBSITE] Validating access to organization: %s", req.OrganizationID)
		// Verify user has access to this organization
		role, err := services.GetMemberRole(ctx, database, req.OrganizationID, userID)
		if err != nil || role != "admin" {
			log.Printf("[WEBSITE] ERROR: Access denied to organization %s for user %s", req.OrganizationID, userID)
			http.Error(w, "Access denied to organization", http.StatusForbidden)
			return
		}

		log.Printf("[WEBSITE] Creating website: %s (URL: %s) in organization %s", req.WebsiteName, req.URL, req.OrganizationID)
		website, err := services.CreateWebsite(ctx, database, req.WebsiteName, req.URL, req.OrganizationID)
		if err != nil {
			log.Printf("[WEBSITE] ERROR: Failed to create website - %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("[WEBSITE] Website created successfully with ID: %s", website.ID)
		utils.WriteJSON(w, http.StatusCreated, website)
	}
}

// ListWebsites godoc
// @Summary      List websites
// @Description  Fetches all websites in organizations where user is a member
// @Tags         websites
// @Produce      json
// @Param        Authorization  header   string  true  "Bearer token"
// @Success      200            {array}  db.WebsiteModel
// @Failure      401            {object} map[string]string
// @Failure      500            {object} map[string]string
// @Router       /api/websites [get]
func ListWebsites(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[WEBSITE] New ListWebsites request received")
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		userID := r.Context().Value("userID").(string)
		log.Printf("[WEBSITE] Fetching websites for user: %s", userID)

		websites, err := services.ListWebsites(ctx, database, userID)
		if err != nil {
			log.Printf("[WEBSITE] ERROR: Failed to list websites for user %s - %v", userID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("[WEBSITE] Retrieved %d websites for user %s", len(websites), userID)
		utils.WriteJSON(w, http.StatusOK, websites)
	}
}

// GetWebsite godoc
// @Summary      Get a website
// @Description  Retrieves a specific website by ID
// @Tags         websites
// @Produce      json
// @Param        Authorization  header   string  true  "Bearer token"
// @Param        id             path     string  true  "Website ID"
// @Success      200            {object} db.WebsiteModel
// @Failure      400            {object} map[string]string
// @Failure      404            {object} map[string]string
// @Failure      500            {object} map[string]string
// @Router       /api/websites/{id} [get]
func GetWebsite(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[WEBSITE] New GetWebsite request received")
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		userID := r.Context().Value("userID").(string)
		id := mux.Vars(r)["id"]

		if id == "" {
			log.Printf("[WEBSITE] ERROR: Missing website ID")
			http.Error(w, "Website id required", http.StatusBadRequest)
			return
		}

		log.Printf("[WEBSITE] Fetching website ID: %s for user: %s", id, userID)
		website, err := services.GetWebsite(ctx, database, id, userID)
		if err != nil {
			log.Printf("[WEBSITE] ERROR: Failed to get website %s - %v", id, err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		log.Printf("[WEBSITE] Retrieved website: %s", website.ID)
		utils.WriteJSON(w, http.StatusOK, website)
	}
}

// UpdateWebsite godoc
// @Summary      Update a website
// @Description  Updates website name and/or URL
// @Tags         websites
// @Accept       json
// @Produce      json
// @Param        Authorization  header   string           true  "Bearer token"
// @Param        id             path     string           true  "Website ID"
// @Param        body           body     UpdateWebsiteReq true  "Update website request"
// @Success      200            {object} db.WebsiteModel
// @Failure      400            {object} map[string]string
// @Failure      500            {object} map[string]string
// @Router       /api/websites/{id} [put]
func UpdateWebsite(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[WEBSITE] New UpdateWebsite request received")
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		userID := r.Context().Value("userID").(string)
		id := mux.Vars(r)["id"]

		if id == "" {
			log.Printf("[WEBSITE] ERROR: Missing website ID")
			http.Error(w, "Website id required", http.StatusBadRequest)
			return
		}

		log.Printf("[WEBSITE] Updating website ID: %s for user: %s", id, userID)

		var req UpdateWebsiteReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("[WEBSITE] ERROR: Failed to decode request body - %v", err)
			http.Error(w, "Invalid req body", http.StatusBadRequest)
			return
		}

		updates := []db.WebsiteSetParam{}
		if req.WebsiteName != nil {
			updates = append(updates, db.Website.WebsiteName.Set(*req.WebsiteName))
			log.Printf("[WEBSITE] Update: WebsiteName = %s", *req.WebsiteName)
		}
		if req.URL != nil {
			updates = append(updates, db.Website.URL.Set(*req.URL))
			log.Printf("[WEBSITE] Update: URL = %s", *req.URL)
		}

		if len(updates) == 0 {
			log.Printf("[WEBSITE] ERROR: No fields to update")
			http.Error(w, "No fields to update", http.StatusBadRequest)
			return
		}

		website, err := services.UpdateWebsite(ctx, database, id, userID, updates)
		if err != nil {
			log.Printf("[WEBSITE] ERROR: Failed to update website %s - %v", id, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("[WEBSITE] Website %s updated successfully", id)
		utils.WriteJSON(w, http.StatusOK, website)
	}
}

// DeleteWebsite godoc
// @Summary      Delete a website
// @Description  Removes a website from the system
// @Tags         websites
// @Produce      json
// @Param        Authorization  header   string  true  "Bearer token"
// @Param        id             path     string  true  "Website ID"
// @Success      200            {object} db.WebsiteModel
// @Failure      400            {object} map[string]string
// @Failure      404            {object} map[string]string
// @Failure      500            {object} map[string]string
// @Router       /api/websites/{id} [delete]
func DeleteWebsite(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[WEBSITE] New DeleteWebsite request received")
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		userID := r.Context().Value("userID").(string)
		id := mux.Vars(r)["id"]

		if id == "" {
			log.Printf("[WEBSITE] ERROR: Missing website ID")
			http.Error(w, "Website id required", http.StatusBadRequest)
			return
		}

		log.Printf("[WEBSITE] Deleting website ID: %s for user: %s", id, userID)
		website, err := services.DeleteWebsite(ctx, database, id, userID)
		if err != nil {
			log.Printf("[WEBSITE] ERROR: Failed to delete website %s - %v", id, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("[WEBSITE] Website %s deleted successfully", id)
		utils.WriteJSON(w, http.StatusOK, website)
	}
}
