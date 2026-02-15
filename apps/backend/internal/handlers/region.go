package handlers

import (
	"backend/db"
	"backend/internal/services"
	"backend/internal/utils"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type CreateRegionRequest struct {
	RegionName string `json:"regionName"`
}

type UpdateRegionReq struct {
	RegionName *string `json:"regionName,omitempty"`
}

// CreateRegion godoc
// @Summary      Create a new region
// @Description  Takes a region name and saves it to the DB via Prisma
// @Tags         regions
// @Accept       json
// @Produce      json
// @Param        body  body      CreateRegionRequest  true  "Create region request"
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/regions [post]
func CreateRegion(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var req CreateRegionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.RegionName == "" {
			http.Error(w, "regionName is required", http.StatusBadRequest)
			return
		}

		region, err := services.CreateRegion(ctx, database, req.RegionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, http.StatusCreated, region)
	}
}

// ListRegions godoc
// @Summary      List regions
// @Description  Fetches all regions
// @Tags         regions
// @Produce      json
// @Failure      500   {object}  map[string]string
// @Router       /api/regions [get]
func ListRegions(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		regions, err := services.ListRegions(ctx, database)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, regions)
	}
}

// GetRegion godoc
// @Summary      Get a region
// @Description  Fetches a region by id
// @Tags         regions
// @Produce      json
// @Param        id   path      string  true  "Region ID"
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/regions/{id} [get]
func GetRegion(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		id := mux.Vars(r)["id"]
		if id == "" {
			http.Error(w, "Region id required", http.StatusBadRequest)
			return
		}

		region, err := services.GetRegion(ctx, database, id)
		if err != nil {
			if db.IsErrNotFound(err) {
				http.Error(w, "Region not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, region)
	}
}

// UpdateRegion godoc
// @Summary      Update a region
// @Description  Updates region name by id
// @Tags         regions
// @Accept       json
// @Produce      json
// @Param        id    path      string            true  "Region ID"
// @Param        body  body      UpdateRegionReq   true  "Update region request"
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/regions/{id} [put]
func UpdateRegion(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		id := mux.Vars(r)["id"]
		if id == "" {
			http.Error(w, "Region id required", http.StatusBadRequest)
			return
		}

		var req UpdateRegionReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid req body", http.StatusBadRequest)
			return
		}

		updates := []db.RegionSetParam{}
		if req.RegionName != nil {
			updates = append(updates, db.Region.RegionName.Set(*req.RegionName))
		}

		if len(updates) == 0 {
			http.Error(w, "No fields to update", http.StatusBadRequest)
			return
		}

		region, err := services.UpdateRegion(ctx, database, id, updates)
		if err != nil {
			if db.IsErrNotFound(err) {
				http.Error(w, "Region not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, region)
	}
}

// DeleteRegion godoc
// @Summary      Delete a region
// @Description  Deletes a region by id
// @Tags         regions
// @Produce      json
// @Param        id   path      string  true  "Region ID"
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/regions/{id} [delete]
func DeleteRegion(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		id := mux.Vars(r)["id"]
		if id == "" {
			http.Error(w, "Region id required", http.StatusBadRequest)
			return
		}

		region, err := services.DeleteRegion(ctx, database, id)
		if err != nil {
			if db.IsErrNotFound(err) {
				http.Error(w, "Region not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, region)
	}

}