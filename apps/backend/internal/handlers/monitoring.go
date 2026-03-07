package handlers

import (
	"backend/db"
	"backend/internal/services"
	"backend/internal/utils"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// ToggleMonitoring godoc
// @Summary      Toggle monitoring for a website
// @Description  Enables or disables monitoring for a given website
// @Tags         monitoring
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "Website ID"
// @Param        body  body      CreateWebsiteRequest  true  "Create website request"
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/monitoring [post]
func ToggleMonitoring(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		id := mux.Vars(r)["id"]
		fmt.Printf("Toggling monitoring for website ID: %s\n", id)

		if id == "" {
			http.Error(w, "Invalid website ID", http.StatusBadRequest)
			return
		}

		website, err := services.GetWebsite(ctx, database, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		updates := []db.WebsiteSetParam{}
		if website.IsMonitoringEnabled != false {
			updates = append(updates, db.Website.IsMonitoringEnabled.Set(false))
		} else {
			updates = append(updates, db.Website.IsMonitoringEnabled.Set(true))
		}
		start, err := services.UpdateWebsite(ctx, database, id, updates)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, http.StatusCreated, start)
	}
}
