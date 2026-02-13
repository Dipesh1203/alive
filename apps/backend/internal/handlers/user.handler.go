package handlers

import (
	"backend/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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
		fmt.Println("PRint  ",req);

		
		website, err := database.Website.CreateOne(db.Website.WebsiteName.Set(req.WebsiteName),db.Website.URL.Set(req.URL)).Exec(ctx);
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(website)
	}
}
