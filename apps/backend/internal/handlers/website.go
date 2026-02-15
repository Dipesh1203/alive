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

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().String(),
	})
}

type CreateWebsiteRequest struct {
	WebsiteName string `json:"websiteName"`
	URL         string `json:"url"`
}

type UpdateWebsiteReq struct {
	WebsiteName *string `json:"websiteName,omitempty"`
	URL         *string `json:"url,omitempty"`
}

// CreateWebsite godoc
	// @Summary      Create a new website
	// @Description  Takes a website name/URL and saves it to the DB via Prisma
	// @Tags         websites
	// @Accept       json
	// @Produce      json
	// @Param        body  body      CreateWebsiteRequest  true  "Create website request"
	// @Failure      400   {object}  map[string]string
	// @Failure      500   {object}  map[string]string
	// @Router       /api/websites [post]
func CreateWebsite(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var req CreateWebsiteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.WebsiteName == "" || req.URL == "" {
			http.Error(w, "websiteName and url are required", http.StatusBadRequest)
			return
		}
		
		website, err := services.CreateWebsite(ctx, database, req.WebsiteName, req.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, http.StatusCreated, website)
	}
}

// ListWebsites godoc
// @Summary      List websites
// @Description  Fetches all websites
// @Tags         websites
// @Produce      json
// @Failure      500   {object}  map[string]string
// @Router       /api/websites [get]
func ListWebsites(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		websites, err := services.ListWebsites(ctx, database)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, websites)
	}
}

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
func GetWebsite(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		id := mux.Vars(r)["id"]
		if id == "" {
			http.Error(w, "Website id required", http.StatusBadRequest)
			return
		}

		website, err := services.GetWebsite(ctx, database, id)
		if err != nil {
			if db.IsErrNotFound(err) {
				http.Error(w, "Website not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, website)
	}
}


// UpdateWebsite godoc
// @Summary      Update a website
// @Description  Updates website name and/or URL by id
// @Tags         websites
// @Accept       json
// @Produce      json
// @Param        body  body      UpdateWebsiteReq  true  "Update website request"
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/websites [put]
func UpdateWebsite(database *db.PrismaClient)http.HandlerFunc{
	
	return func(w http.ResponseWriter, r *http.Request) {
		ctx,cancel := context.WithTimeout(r.Context(),10*time.Second);
		defer cancel();

		id := mux.Vars(r)["id"]
		if id == "" {
			http.Error(w, "Website id required", http.StatusBadRequest)
			return
		}

		var req UpdateWebsiteReq
		if err:=json.NewDecoder(r.Body).Decode(&req); err!=nil {
			http.Error(w,"Invalid req body",http.StatusBadRequest);
			return;
		}

		updates := []db.WebsiteSetParam{}
		if req.WebsiteName != nil {
			updates = append(updates, db.Website.WebsiteName.Set(*req.WebsiteName))
		}
		if req.URL != nil {
			updates = append(updates, db.Website.URL.Set(*req.URL))
		}

		if len(updates) == 0 {
			http.Error(w, "No fields to update", http.StatusBadRequest)
			return
		}

		website, err := services.UpdateWebsite(ctx, database, id, updates)
		if err != nil {
			if db.IsErrNotFound(err) {
				http.Error(w, "Website not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, website)

	}
}

// DeleteWebsite godoc
// @Summary      Delete a website
// @Description  Deletes a website by id
// @Tags         websites
// @Produce      json
// @Param        id   path      string  true  "Website ID"
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/websites/{id} [delete]
func DeleteWebsite(database *db.PrismaClient)http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		id := mux.Vars(r)["id"]
		if id == "" {
			http.Error(w, "Website id required", http.StatusBadRequest)
			return
		}

		website, err := services.DeleteWebsite(ctx, database, id)
		if err != nil {
			if db.IsErrNotFound(err) {
				http.Error(w, "Website not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, website)
	}
}
