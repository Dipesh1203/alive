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

// GetWebsite godoc
// @Summary      Get a website
// @Description  Fetches a website by id
// @Tags         websites
// @Produce      json
// @Param        id   path      string  true  "Website ID"
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/websites/{id} [get]
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

		_, err := services.GetWebsite(ctx, database, id)
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
