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

	"github.com/golang-jwt/jwt/v5"
)

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Message string        `json:"message"`
	User    *db.UserModel `json:"user,omitempty"`
	Token   string        `json:"token,omitempty"`
	Error   string        `json:"error,omitempty"`
}

// Signup godoc
// @Summary      Sign up a new user
// @Description  Creates a new user account with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      SignupRequest  true  "Signup request"
// @Success      201   {object}  AuthResponse
// @Failure      400   {object}  AuthResponse
// @Failure      500   {object}  AuthResponse
// @Router       /api/auth/signup [post]
func Signup(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[SIGNUP] New signup request received")
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var req SignupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("[SIGNUP] ERROR: Failed to decode request body - %v", err)
			utils.WriteJSON(w, http.StatusBadRequest, AuthResponse{
				Error: "Invalid request body",
			})
			return
		}

		log.Printf("[SIGNUP] Processing signup for email: %s", req.Email)

		// Validate input
		if req.Email == "" || req.Password == "" {
			log.Printf("[SIGNUP] ERROR: Email or password is empty")
			utils.WriteJSON(w, http.StatusBadRequest, AuthResponse{
				Error: "Email and password are required",
			})
			return
		}

		if len(req.Password) < 6 {
			log.Printf("[SIGNUP] ERROR: Password too short for email %s", req.Email)
			utils.WriteJSON(w, http.StatusBadRequest, AuthResponse{
				Error: "Password must be at least 6 characters long",
			})
			return
		}

		log.Printf("[SIGNUP] Validation passed for %s, creating user...", req.Email)

		// Create user
		user, err := services.UserSignup(ctx, database, req.Email, req.Password)
		if err != nil {
			log.Printf("[SIGNUP] ERROR: Failed to create user %s - %v", req.Email, err)
			utils.WriteJSON(w, http.StatusBadRequest, AuthResponse{
				Error: err.Error(),
			})
			return
		}

		log.Printf("[SIGNUP] User created successfully with ID: %s", user.ID)

		// Generate JWT token
		token, err := generateJWTToken(user.ID, user.Email)
		if err != nil {
			log.Printf("[SIGNUP] ERROR: Failed to generate JWT token for user %s - %v", user.ID, err)
			utils.WriteJSON(w, http.StatusInternalServerError, AuthResponse{
				Error: "Failed to generate token",
			})
			return
		}

		log.Printf("[SIGNUP] JWT token generated successfully for user %s", user.ID)

		// Clear password from response
		user.Password = ""

		utils.WriteJSON(w, http.StatusCreated, AuthResponse{
			Message: "User created successfully",
			User:    user,
			Token:   token,
		})
		log.Printf("[SIGNUP] Signup completed successfully for user %s", user.ID)
	}
}

// Login godoc
// @Summary      Log in a user
// @Description  Authenticates a user with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      LoginRequest   true  "Login request"
// @Success      200   {object}  AuthResponse
// @Failure      400   {object}  AuthResponse
// @Failure      401   {object}  AuthResponse
// @Failure      500   {object}  AuthResponse
// @Router       /api/auth/login [post]
func Login(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[LOGIN] New login request received")
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("[LOGIN] ERROR: Failed to decode request body - %v", err)
			utils.WriteJSON(w, http.StatusBadRequest, AuthResponse{
				Error: "Invalid request body",
			})
			return
		}

		log.Printf("[LOGIN] Processing login for email: %s", req.Email)

		// Validate input
		if req.Email == "" || req.Password == "" {
			log.Printf("[LOGIN] ERROR: Email or password is empty")
			utils.WriteJSON(w, http.StatusBadRequest, AuthResponse{
				Error: "Email and password are required",
			})
			return
		}

		log.Printf("[LOGIN] Attempting to authenticate user %s", req.Email)

		// Authenticate user
		user, err := services.UserLogin(ctx, database, req.Email, req.Password)
		if err != nil {
			log.Printf("[LOGIN] ERROR: Authentication failed for %s - %v", req.Email, err)
			utils.WriteJSON(w, http.StatusUnauthorized, AuthResponse{
				Error: err.Error(),
			})
			return
		}

		log.Printf("[LOGIN] Authentication successful for user %s (ID: %s)", req.Email, user.ID)

		// Generate JWT token
		token, err := generateJWTToken(user.ID, user.Email)
		if err != nil {
			log.Printf("[LOGIN] ERROR: Failed to generate JWT token for user %s - %v", user.ID, err)
			utils.WriteJSON(w, http.StatusInternalServerError, AuthResponse{
				Error: "Failed to generate token",
			})
			return
		}

		log.Printf("[LOGIN] JWT token generated successfully for user %s", user.ID)

		// Clear password from response
		user.Password = ""

		utils.WriteJSON(w, http.StatusOK, AuthResponse{
			Message: "Login successful",
			User:    user,
			Token:   token,
		})
		log.Printf("[LOGIN] Login completed successfully for user %s", user.ID)
	}
}

// generateJWTToken creates a JWT token for the given user
func generateJWTToken(userID string, email string) (string, error) {
	// Secret key - should be loaded from environment variable in production
	var SecretKey = utils.GoGetEnv("JWT_SECRET")
	log.Printf("[JWT] jwt SecretKey: %s", SecretKey)
	secretKey := []byte(SecretKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"email":  email,
		"exp":    time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days expiration
		"iat":    time.Now().Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
