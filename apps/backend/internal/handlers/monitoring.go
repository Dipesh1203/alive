package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Dipesh1203/alive/backend/db"

	"github.com/Dipesh1203/alive/backend/internal/services"
	"github.com/Dipesh1203/alive/backend/internal/utils"

	"github.com/gorilla/mux"
)

// ToggleMonitoring godoc
// @Summary      Toggle website monitoring
// @Description  Enables or disables monitoring for a given website
// @Tags         monitoring
// @Accept       json
// @Produce      json
// @Param        Authorization  header   string  true  "Bearer token"
// @Param        id             path     string  true  "Website ID"
// @Success      201            {object} db.WebsiteModel
// @Failure      400            {object} map[string]string
// @Failure      500            {object} map[string]string
// @Router       /api/monitoring/{id} [post]
func ToggleMonitoring(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[MONITORING] New ToggleMonitoring request received")
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		id := mux.Vars(r)["id"]
		userID := r.Context().Value("userID").(string)
		log.Printf("[MONITORING] Toggling monitoring for website ID: %s by user: %s", id, userID)

		if id == "" {
			log.Printf("[MONITORING] ERROR: Invalid website ID")
			http.Error(w, "Invalid website ID", http.StatusBadRequest)
			return
		}

		website, err := services.GetWebsite(ctx, database, id, userID)
		if err != nil {
			log.Printf("[MONITORING] ERROR: Failed to get website %s - %v", id, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		updates := []db.WebsiteSetParam{}
		newStatus := false
		if website.IsMonitoringEnabled != false {
			updates = append(updates, db.Website.IsMonitoringEnabled.Set(false))
			log.Printf("[MONITORING] Disabling monitoring for website %s", id)
		} else {
			updates = append(updates, db.Website.IsMonitoringEnabled.Set(true))
			newStatus = true
			log.Printf("[MONITORING] Enabling monitoring for website %s", id)
		}
		start, err := services.UpdateWebsite(ctx, database, id, userID, updates)
		if err != nil {
			log.Printf("[MONITORING] ERROR: Failed to update monitoring status - %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("[MONITORING] Monitoring toggled successfully for website %s (new status: %v)", id, newStatus)
		utils.WriteJSON(w, http.StatusCreated, start)
	}
}
