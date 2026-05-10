package middleware

import (
	"backend/internal/utils"
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = utils.GoGetEnv("JWT_SECRET")

// AuthMiddleware validates JWT tokens and extracts user information
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[AUTH] Processing authorization for %s %s", r.Method, r.URL.Path)

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			log.Printf("[AUTH] ERROR: Missing authorization header for %s", r.URL.Path)
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		// Extract bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("[AUTH] ERROR: Invalid authorization header format for %s", r.URL.Path)
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		log.Printf("[AUTH] Validating token for %s", r.URL.Path)

		log.Printf("[JWT] jwt SecretKey: %s", SecretKey)
		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(SecretKey), nil
		})

		if err != nil || !token.Valid {
			log.Printf("[AUTH] ERROR: Invalid or expired token for %s - %v", r.URL.Path, err)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Printf("[AUTH] ERROR: Invalid token claims for %s", r.URL.Path)
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["userID"].(string)
		if !ok {
			log.Printf("[AUTH] ERROR: Missing userID in token claims for %s", r.URL.Path)
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			log.Printf("[AUTH] ERROR: Missing email in token claims for %s", r.URL.Path)
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		log.Printf("[AUTH] Token validated successfully for user %s (email: %s) - %s %s", userID, email, r.Method, r.URL.Path)

		// Add to context
		ctx := context.WithValue(r.Context(), "userID", userID)
		ctx = context.WithValue(ctx, "email", email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuthMiddleware is like AuthMiddleware but doesn't fail if token is missing
func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[AUTH] Processing optional authorization for %s %s", r.Method, r.URL.Path)

		authHeader := r.Header.Get("Authorization")

		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString := parts[1]
				log.Printf("[AUTH-OPT] Found token, attempting validation for %s", r.URL.Path)

				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, jwt.ErrSignatureInvalid
					}
					return []byte(SecretKey), nil
				})

				if err == nil && token.Valid {
					claims, ok := token.Claims.(jwt.MapClaims)
					if ok {
						userID, _ := claims["userID"].(string)
						email, _ := claims["email"].(string)

						log.Printf("[AUTH-OPT] Optional token validated for user %s (email: %s) - %s %s", userID, email, r.Method, r.URL.Path)

						ctx := context.WithValue(r.Context(), "userID", userID)
						ctx = context.WithValue(ctx, "email", email)
						r = r.WithContext(ctx)
					}
				} else {
					log.Printf("[AUTH-OPT] Optional token validation failed for %s - %v", r.URL.Path, err)
				}
			}
		} else {
			log.Printf("[AUTH-OPT] No authorization header, proceeding unauthenticated for %s", r.URL.Path)
		}

		next.ServeHTTP(w, r)
	})
}

// GetUserID extracts the user ID from request context (requires user to be authenticated)
func GetUserID(r *http.Request) string {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		return ""
	}
	return userID
}

// GetEmail extracts the email from request context (requires user to be authenticated)
func GetEmail(r *http.Request) string {
	email, ok := r.Context().Value("email").(string)
	if !ok {
		return ""
	}
	return email
}
