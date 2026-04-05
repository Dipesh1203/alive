package handlers

import (
	"backend/db"
	"backend/internal/services"
	"backend/internal/utils"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// GetDetailsWebsite godoc
// @Summary      Get detailed website metrics
// @Description  Retrieves uptime ticks and latency data for a website within a date range
// @Tags         websites
// @Produce      json
// @Param        Authorization  header   string  true  "Bearer token"
// @Param        id             path     string  true  "Website ID"
// @Param        skip           query    integer false "Number of records to skip"
// @Param        take           query    integer false "Number of records to take"
// @Param        startDate      query    string  false "Start date (YYYY-MM-DD)"
// @Param        endDate        query    string  false "End date (YYYY-MM-DD)"
// @Success      200            {array}  db.WebsiteTicksModel
// @Failure      400            {object} map[string]string
// @Failure      404            {object} map[string]string
// @Failure      500            {object} map[string]string
// @Router       /api/websites/{id}/details [get]
func GetDetailsWebsite(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		query := r.URL.Query()
		id := mux.Vars(r)["id"]
		if id == "" {
			http.Error(w, "Website id required", http.StatusBadRequest)
			return
		}

		var skip *int
		if skipStr := query.Get("skip"); skipStr != "" {
			if val, err := strconv.Atoi(skipStr); err == nil {
				skip = &val
			}
		}

		var take *int
		if takeStr := query.Get("take"); takeStr != "" {
			if val, err := strconv.Atoi(takeStr); err == nil {
				take = &val
			}
		}
		startDate := query.Get("startDate")
		endDate := query.Get("endDate")

		_, err := services.GetWebsite(ctx, database, id, r.Context().Value("userID").(string))
		if err != nil {
			if db.IsErrNotFound(err) {
				http.Error(w, "Website not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ticks, err2 := services.ListTicks(ctx, database, id, skip, take, &startDate, &endDate)

		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
			return
		}
		utils.WriteJSON(w, http.StatusOK, ticks)
	}
}
