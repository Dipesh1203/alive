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

type CreateRegionRequest struct {
	RegionName string `json:"regionName"`
}

type UpdateRegionReq struct {
	RegionName *string `json:"regionName,omitempty"`
}

// CreateRegion godoc
// @Summary      Create a new region
// @Description  Creates a new monitoring region
// @Tags         regions
// @Accept       json
// @Produce      json
// @Param        Authorization  header   string                 true  "Bearer token"
// @Param        body           body     CreateRegionRequest    true  "Create region request"
// @Success      201            {object} db.RegionModel
// @Failure      400            {object} map[string]string
// @Failure      500            {object} map[string]string
// @Router       /api/regions [post]
func CreateRegion(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[REGION] New CreateRegion request received")
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var req CreateRegionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("[REGION] ERROR: Failed to decode request body - %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.RegionName == "" {
			log.Printf("[REGION] ERROR: Region name is empty")
			http.Error(w, "regionName is required", http.StatusBadRequest)
			return
		}

		log.Printf("[REGION] Creating region: %s", req.RegionName)
		region, err := services.CreateRegion(ctx, database, req.RegionName)
		if err != nil {
			log.Printf("[REGION] ERROR: Failed to create region - %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("[REGION] Region created successfully with ID: %s", region.RegionID)
		utils.WriteJSON(w, http.StatusCreated, region)
	}
}

// ListRegions godoc
// @Summary      List regions
// @Description  Retrieves all available monitoring regions
// @Tags         regions
// @Produce      json
// @Param        Authorization  header   string  true  "Bearer token"
// @Success      200            {array}  db.RegionModel
// @Failure      500            {object} map[string]string
// @Router       /api/regions [get]
func ListRegions(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[REGION] New ListRegions request received")
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		log.Printf("[REGION] Fetching all regions")
		regions, err := services.ListRegions(ctx, database)
		if err != nil {
			log.Printf("[REGION] ERROR: Failed to list regions - %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("[REGION] Retrieved %d regions", len(regions))
		utils.WriteJSON(w, http.StatusOK, regions)
	}
}

// GetRegion godoc
// @Summary      Get a region
// @Description  Retrieves a specific region by ID
// @Tags         regions
// @Produce      json
// @Param        Authorization  header   string  true  "Bearer token"
// @Param        id             path     string  true  "Region ID"
// @Success      200            {object} db.RegionModel
// @Failure      400            {object} map[string]string
// @Failure      404            {object} map[string]string
// @Failure      500            {object} map[string]string
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
// @Description  Updates a region's name
// @Tags         regions
// @Accept       json
// @Produce      json
// @Param        Authorization  header   string           true  "Bearer token"
// @Param        id             path     string           true  "Region ID"
// @Param        body           body     UpdateRegionReq  true  "Update region request"
// @Success      200            {object} db.RegionModel
// @Failure      400            {object} map[string]string
// @Failure      404            {object} map[string]string
// @Failure      500            {object} map[string]string
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
// @Description  Removes a region from the system
// @Tags         regions
// @Produce      json
// @Param        Authorization  header   string  true  "Bearer token"
// @Param        id             path     string  true  "Region ID"
// @Success      200            {object} db.RegionModel
// @Failure      400            {object} map[string]string
// @Failure      404            {object} map[string]string
// @Failure      500            {object} map[string]string
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
