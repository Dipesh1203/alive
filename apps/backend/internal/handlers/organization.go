package handlers

import (
	"backend/db"
	"backend/internal/services"
	"backend/internal/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type CreateOrganizationRequest struct {
	Name string `json:"name"`
}

type UpdateOrganizationRequest struct {
	Name *string `json:"name,omitempty"`
}

type AddMemberRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UpdateMemberRoleRequest struct {
	Role *string `json:"role,omitempty"`
}

type OrganizationResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	AdminID   string    `json:"adminId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type OrganizationMemberResponse struct {
	ID             string    `json:"id"`
	UserID         string    `json:"userId"`
	Name           string    `json:"name,omitempty"`
	Email          string    `json:"email,omitempty"`
	FirstName      string    `json:"firstName,omitempty"`
	LastName       string    `json:"lastName,omitempty"`
	Avatar         *string   `json:"avatar,omitempty"`
	OrganizationID string    `json:"organizationId"`
	Role           string    `json:"role"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type UserDetails struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	FirstName string  `json:"firstName,omitempty"`
	LastName  string  `json:"lastName,omitempty"`
	Avatar    *string `json:"avatar,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func buildOrganizationMemberResponse(
	ctx context.Context,
	database *db.PrismaClient,
	member db.OrganizationMemberModel,
	knownUser *db.UserModel,
) OrganizationMemberResponse {
	response := OrganizationMemberResponse{
		ID:             member.ID,
		UserID:         member.UserID,
		OrganizationID: member.OrganizationID,
		Role:           member.Role,
		CreatedAt:      member.CreatedAt,
		UpdatedAt:      member.UpdatedAt,
	}

	user := knownUser
	if user == nil {
		fetchedUser, err := services.GetUserByID(ctx, database, member.UserID)
		if err == nil {
			user = fetchedUser
		}
	}

	if user != nil {
		response.Email = user.Email
	}

	profile, err := services.GetUserProfile(ctx, database, member.UserID)
	if err == nil && profile != nil {
		firstName, _ := profile.FirstName()
		lastName, _ := profile.LastName()
		avatar, hasAvatar := profile.Avatar()

		response.FirstName = firstName
		response.LastName = lastName
		response.Name = strings.TrimSpace(firstName + " " + lastName)
		if hasAvatar && strings.TrimSpace(avatar) != "" {
			response.Avatar = &avatar
		}
	}

	if response.Name == "" && response.Email != "" {
		parts := strings.SplitN(response.Email, "@", 2)
		if len(parts) > 0 {
			response.Name = parts[0]
		}
	}

	return response
}

// CreateOrganization godoc
// @Summary      Create a new organization
// @Description  Creates a new organization with the authenticated user as admin
// @Tags         organizations
// @Accept       json
// @Produce      json
// @Param        Authorization  header   string                      true  "Bearer token"
// @Param        body           body     CreateOrganizationRequest   true  "Organization request"
// @Success      201            {object} OrganizationResponse
// @Failure      400            {object} ErrorResponse
// @Failure      401            {object} ErrorResponse
// @Failure      500            {object} ErrorResponse
// @Router       /api/organizations [post]
func CreateOrganization(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[ORG] New CreateOrganization request received")
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		userID := r.Context().Value("userID").(string)
		log.Printf("[ORG] Processing request by user: %s", userID)

		var req CreateOrganizationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("[ORG] ERROR: Failed to decode request body - %v", err)
			utils.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
				Error: "Invalid request body",
			})
			return
		}

		if req.Name == "" {
			log.Printf("[ORG] ERROR: Organization name is empty")
			utils.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
				Error: "Organization name is required",
			})
			return
		}

		log.Printf("[ORG] Creating organization: %s by user %s", req.Name, userID)
		org, err := services.CreateOrganization(ctx, database, req.Name, userID)
		if err != nil {
			log.Printf("[ORG] ERROR: Failed to create organization - %v", err)
			utils.WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		log.Printf("[ORG] Organization created successfully with ID: %s", org.ID)
		utils.WriteJSON(w, http.StatusCreated, OrganizationResponse{
			ID:        org.ID,
			Name:      org.Name,
			AdminID:   org.AdminID,
			CreatedAt: org.CreatedAt,
			UpdatedAt: org.UpdatedAt,
		})
	}
}

// ListUserOrganizations godoc
// @Summary      List user's organizations
// @Description  Lists all organizations where the user is a member
// @Tags         organizations
// @Produce      json
// @Param        Authorization  header   string  true  "Bearer token"
// @Success      200            {array}  OrganizationResponse
// @Failure      401            {object} ErrorResponse
// @Failure      500            {object} ErrorResponse
// @Router       /api/organizations [get]
func ListUserOrganizations(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[ORG] New ListUserOrganizations request received")
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		userID := r.Context().Value("userID").(string)
		log.Printf("[ORG] Fetching organizations for user: %s", userID)

		orgs, err := services.ListUserOrganizations(ctx, database, userID)
		if err != nil {
			log.Printf("[ORG] ERROR: Failed to list organizations - %v", err)
			utils.WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		response := make([]OrganizationResponse, 0)
		for _, org := range orgs {
			response = append(response, OrganizationResponse{
				ID:        org.ID,
				Name:      org.Name,
				AdminID:   org.AdminID,
				CreatedAt: org.CreatedAt,
				UpdatedAt: org.UpdatedAt,
			})
		}

		log.Printf("[ORG] Retrieved %d organizations for user %s", len(response), userID)
		utils.WriteJSON(w, http.StatusOK, response)
	}
}

// GetOrganization godoc
// @Summary      Get organization details
// @Description  Retrieves details of a specific organization
// @Tags         organizations
// @Produce      json
// @Param        id             path     string  true  "Organization ID"
// @Param        Authorization  header   string  true  "Bearer token"
// @Success      200            {object} OrganizationResponse
// @Failure      404            {object} ErrorResponse
// @Failure      500            {object} ErrorResponse
// @Router       /api/organizations/{id} [get]
func GetOrganization(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[ORG] New GetOrganization request received")
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		orgID := vars["id"]

		log.Printf("[ORG] Fetching organization: %s", orgID)
		org, err := services.GetOrganizationByID(ctx, database, orgID)
		if err != nil {
			log.Printf("[ORG] ERROR: Organization %s not found - %v", orgID, err)
			utils.WriteJSON(w, http.StatusNotFound, ErrorResponse{
				Error: "Organization not found",
			})
			return
		}

		log.Printf("[ORG] Retrieved organization: %s (%s)", org.ID, org.Name)
		utils.WriteJSON(w, http.StatusOK, OrganizationResponse{
			ID:        org.ID,
			Name:      org.Name,
			AdminID:   org.AdminID,
			CreatedAt: org.CreatedAt,
			UpdatedAt: org.UpdatedAt,
		})
	}
}

// UpdateOrganization godoc
// @Summary      Update organization
// @Description  Updates organization details (admin only)
// @Tags         organizations
// @Accept       json
// @Produce      json
// @Param        id             path     string                      true  "Organization ID"
// @Param        Authorization  header   string                      true  "Bearer token"
// @Param        body           body     UpdateOrganizationRequest   true  "Update request"
// @Success      200            {object} OrganizationResponse
// @Failure      400            {object} ErrorResponse
// @Failure      403            {object} ErrorResponse
// @Failure      500            {object} ErrorResponse
// @Router       /api/organizations/{id} [put]
func UpdateOrganization(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[ORG] New UpdateOrganization request received")
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		orgID := vars["id"]
		userID := r.Context().Value("userID").(string)

		log.Printf("[ORG] Checking admin permission for user %s on organization %s", userID, orgID)
		// Check if user is admin
		role, err := services.GetMemberRole(ctx, database, orgID, userID)
		if err != nil || role != "admin" {
			log.Printf("[ORG] ERROR: Access denied (role: %s) for user %s on organization %s", role, userID, orgID)
			utils.WriteJSON(w, http.StatusForbidden, ErrorResponse{
				Error: "You don't have permission to update this organization",
			})
			return
		}

		var req UpdateOrganizationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("[ORG] ERROR: Failed to decode request body - %v", err)
			utils.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
				Error: "Invalid request body",
			})
			return
		}

		if req.Name == nil {
			log.Printf("[ORG] ERROR: No fields to update")
			utils.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
				Error: "No fields to update",
			})
			return
		}

		log.Printf("[ORG] Updating organization %s with name: %s", orgID, *req.Name)
		org, err := services.UpdateOrganization(ctx, database, orgID, req.Name)
		if err != nil {
			log.Printf("[ORG] ERROR: Failed to update organization - %v", err)
			utils.WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		log.Printf("[ORG] Organization %s updated successfully", orgID)
		utils.WriteJSON(w, http.StatusOK, OrganizationResponse{
			ID:        org.ID,
			Name:      org.Name,
			AdminID:   org.AdminID,
			CreatedAt: org.CreatedAt,
			UpdatedAt: org.UpdatedAt,
		})
	}
}

// ListOrganizationMembers godoc
// @Summary      List organization members
// @Description  Lists all members of an organization
// @Tags         organizations
// @Produce      json
// @Param        id             path     string  true  "Organization ID"
// @Param        Authorization  header   string  true  "Bearer token"
// @Success      200            {array}  OrganizationMemberResponse
// @Failure      401            {object} ErrorResponse
// @Failure      500            {object} ErrorResponse
// @Router       /api/organizations/{id}/members [get]
func ListOrganizationMembers(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		orgID := vars["id"]

		members, err := services.ListOrganizationMembers(ctx, database, orgID)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		response := make([]OrganizationMemberResponse, 0)
		for _, member := range members {
			response = append(response, buildOrganizationMemberResponse(ctx, database, member, nil))
		}

		utils.WriteJSON(w, http.StatusOK, response)
	}
}

// AddOrganizationMember godoc
// @Summary      Add member to organization
// @Description  Adds a user to an organization by email (admin only)
// @Tags         organizations
// @Accept       json
// @Produce      json
// @Param        id             path     string            true  "Organization ID"
// @Param        Authorization  header   string            true  "Bearer token"
// @Param        body           body     AddMemberRequest  true  "Add member request"
// @Success      201            {object} OrganizationMemberResponse
// @Failure      400            {object} ErrorResponse
// @Failure      403            {object} ErrorResponse
// @Failure      500            {object} ErrorResponse
// @Router       /api/organizations/{id}/members [post]
func AddOrganizationMember(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		orgID := vars["id"]
		userID := r.Context().Value("userID").(string)

		// Check if user is admin
		role, err := services.GetMemberRole(ctx, database, orgID, userID)
		if err != nil || role != "admin" {
			utils.WriteJSON(w, http.StatusForbidden, ErrorResponse{
				Error: "You don't have permission to add members",
			})
			return
		}

		var req AddMemberRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
				Error: "Invalid request body",
			})
			return
		}

		if req.Email == "" {
			utils.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
				Error: "Email is required",
			})
			return
		}

		// Get user by email
		user, err := services.GetUserByEmail(ctx, database, req.Email)
		if err != nil {
			utils.WriteJSON(w, http.StatusNotFound, ErrorResponse{
				Error: "User not found",
			})
			return
		}

		memberRole := req.Role
		if memberRole == "" {
			memberRole = "viewer"
		}

		member, err := services.AddMemberToOrganization(ctx, database, orgID, user.ID, memberRole)
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		utils.WriteJSON(w, http.StatusCreated, buildOrganizationMemberResponse(ctx, database, *member, user))
	}
}

// UpdateMemberRole godoc
// @Summary      Update member role
// @Description  Updates a member's role in an organization (admin only)
// @Tags         organizations
// @Accept       json
// @Produce      json
// @Param        id             path     string                      true  "Organization ID"
// @Param        memberId       path     string                      true  "Member ID / User ID"
// @Param        Authorization  header   string                      true  "Bearer token"
// @Param        body           body     UpdateMemberRoleRequest     true  "Update role request"
// @Success      200            {object} OrganizationMemberResponse
// @Failure      400            {object} ErrorResponse
// @Failure      403            {object} ErrorResponse
// @Failure      500            {object} ErrorResponse
// @Router       /api/organizations/{id}/members/{memberId}/role [put]
func UpdateMemberRole(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		orgID := vars["id"]
		memberUserID := vars["memberId"]
		currentUserID := r.Context().Value("userID").(string)

		// Check if current user is admin
		role, err := services.GetMemberRole(ctx, database, orgID, currentUserID)
		if err != nil || role != "admin" {
			utils.WriteJSON(w, http.StatusForbidden, ErrorResponse{
				Error: "You don't have permission to update member roles",
			})
			return
		}

		var req UpdateMemberRoleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
				Error: "Invalid request body",
			})
			return
		}
		if req.Role == nil {
			utils.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
				Error: "No fields to update",
			})
			return
		}

		member, err := services.UpdateMemberRole(ctx, database, orgID, memberUserID, req.Role)
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		utils.WriteJSON(w, http.StatusOK, buildOrganizationMemberResponse(ctx, database, *member, nil))
	}
}

// RemoveOrganizationMember godoc
// @Summary      Remove member from organization
// @Description  Removes a user from an organization (admin only)
// @Tags         organizations
// @Produce      json
// @Param        id             path     string  true  "Organization ID"
// @Param        memberId       path     string  true  "Member ID / User ID"
// @Param        Authorization  header   string  true  "Bearer token"
// @Success      200            {object} map[string]string
// @Failure      403            {object} ErrorResponse
// @Failure      500            {object} ErrorResponse
// @Router       /api/organizations/{id}/members/{memberId} [delete]
func RemoveOrganizationMember(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		orgID := vars["id"]
		memberUserID := vars["memberId"]
		currentUserID := r.Context().Value("userID").(string)

		// Check if current user is admin
		role, err := services.GetMemberRole(ctx, database, orgID, currentUserID)
		if err != nil || role != "admin" {
			utils.WriteJSON(w, http.StatusForbidden, ErrorResponse{
				Error: "You don't have permission to remove members",
			})
			return
		}

		err = services.RemoveMemberFromOrganization(ctx, database, orgID, memberUserID)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{
			"message": "Member removed successfully",
		})
	}
}
