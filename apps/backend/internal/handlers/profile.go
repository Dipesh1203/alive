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
)

type UserProfileRequest struct {
	FirstName *string `json:"firstName,omitempty"`
	LastName  *string `json:"lastName,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Bio       *string `json:"bio,omitempty"`
	Avatar    *string `json:"avatar,omitempty"`
}

type UserProfileResponse struct {
	ID          string                 `json:"id"`
	UserID      string                 `json:"userId"`
	FirstName   string                 `json:"firstName"`
	LastName    string                 `json:"lastName"`
	Phone       string                 `json:"phone"`
	Bio         string                 `json:"bio"`
	Avatar      string                 `json:"avatar"`
	Preferences map[string]interface{} `json:"preferences,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

type PreferencesRequest struct {
	Preferences map[string]interface{} `json:"preferences"`
}

func buildUserProfileResponse(profile *db.UserProfileModel) UserProfileResponse {
	preferences := make(map[string]interface{})
	if preferencesJSON, ok := profile.Preferences(); ok {
		_ = json.Unmarshal(preferencesJSON, &preferences)
	}

	firstName, _ := profile.FirstName()
	lastName, _ := profile.LastName()
	phone, _ := profile.Phone()
	bio, _ := profile.Bio()
	avatar, _ := profile.Avatar()

	return UserProfileResponse{
		ID:          profile.ID,
		UserID:      profile.UserID,
		FirstName:   firstName,
		LastName:    lastName,
		Phone:       phone,
		Bio:         bio,
		Avatar:      avatar,
		Preferences: preferences,
		CreatedAt:   profile.CreatedAt,
		UpdatedAt:   profile.UpdatedAt,
	}
}

// GetMyProfile godoc
// @Summary      Get current user's profile
// @Description  Retrieves the authenticated user's profile information
// @Tags         profile
// @Produce      json
// @Param        Authorization  header   string  true  "Bearer token"
// @Success      200            {object} UserProfileResponse
// @Failure      401            {object} ErrorResponse
// @Failure      404            {object} ErrorResponse
// @Failure      500            {object} ErrorResponse
// @Router       /api/profile [get]
func GetMyProfile(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		userID := r.Context().Value("userID").(string)

		profile, err := services.GetUserProfile(ctx, database, userID)
		if err != nil {
			utils.WriteJSON(w, http.StatusNotFound, ErrorResponse{
				Error: "Profile not found",
			})
			return
		}

		utils.WriteJSON(w, http.StatusOK, buildUserProfileResponse(profile))
	}
}

// UpdateMyProfile godoc
// @Summary      Update current user's profile
// @Description  Updates the authenticated user's profile information
// @Tags         profile
// @Accept       json
// @Produce      json
// @Param        Authorization  header   string              true  "Bearer token"
// @Param        body           body     UserProfileRequest  true  "Profile update request"
// @Success      200            {object} UserProfileResponse
// @Failure      400            {object} ErrorResponse
// @Failure      401            {object} ErrorResponse
// @Failure      500            {object} ErrorResponse
// @Router       /api/profile [put]
func UpdateMyProfile(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		userID := r.Context().Value("userID").(string)

		var req UserProfileRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
				Error: "Invalid request body",
			})
			return
		}
		log.Printf("Received profile update request: %+v\n", req)

		profile, err := services.UpdateUserProfile(ctx, database, userID, req.FirstName, req.LastName, req.Phone, req.Bio, req.Avatar)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		utils.WriteJSON(w, http.StatusOK, buildUserProfileResponse(profile))
	}
}

// UpdateMyPreferences godoc
// @Summary      Update user preferences
// @Description  Updates the authenticated user's notification and app preferences
// @Tags         profile
// @Accept       json
// @Produce      json
// @Param        Authorization  header   string              true  "Bearer token"
// @Param        body           body     PreferencesRequest  true  "Preferences request"
// @Success      200            {object} UserProfileResponse
// @Failure      400            {object} ErrorResponse
// @Failure      401            {object} ErrorResponse
// @Failure      500            {object} ErrorResponse
// @Router       /api/profile/preferences [put]
func UpdateMyPreferences(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		userID := r.Context().Value("userID").(string)

		var req PreferencesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
				Error: "Invalid request body",
			})
			return
		}

		profile, err := services.UpdateUserPreferences(ctx, database, userID, req.Preferences)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		utils.WriteJSON(w, http.StatusOK, buildUserProfileResponse(profile))
	}
}

// GetMyPreferences godoc
// @Summary      Get user preferences
// @Description  Retrieves the authenticated user's preferences
// @Tags         profile
// @Produce      json
// @Param        Authorization  header   string  true  "Bearer token"
// @Success      200            {object} map[string]interface{}
// @Failure      401            {object} ErrorResponse
// @Failure      500            {object} ErrorResponse
// @Router       /api/profile/preferences [get]
func GetMyPreferences(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		userID := r.Context().Value("userID").(string)

		preferences, err := services.GetUserPreferences(ctx, database, userID)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		utils.WriteJSON(w, http.StatusOK, preferences)
	}
}
