package handlers

import (
	"backend/db"
	awsservice "backend/internal/services/aws"
	"backend/internal/utils"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type CreateNotificationChannelRequest struct {
	Type     string `json:"type"`
	Email    string `json:"email,omitempty"`
	FCMToken string `json:"fcmToken,omitempty"`
}

type CreateNotificationChannelResponse struct {
	Message     string `json:"message"`
	EndpointArn string `json:"endpointArn,omitempty"`
}

type SendNotificationRequest struct {
	Type        string `json:"type"`
	Email       string `json:"email,omitempty"`
	EndpointArn string `json:"endpointArn,omitempty"`
	Subject     string `json:"subject,omitempty"`
	Message     string `json:"message"`
	Title       string `json:"title,omitempty"`
}

func CreateNotificationChannel(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = database
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		var req CreateNotificationChannelRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		service, err := awsservice.NewNotificationService(ctx)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		switch strings.ToLower(strings.TrimSpace(req.Type)) {
		case "email":
			if strings.TrimSpace(req.Email) == "" {
				utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "email is required for type=email"})
				return
			}

			subject := "Alive Notifications Enabled"
			body := "You are now subscribed to Alive email notifications."
			if err := service.SendEmail(ctx, strings.TrimSpace(req.Email), subject, body); err != nil {
				utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}

			utils.WriteJSON(w, http.StatusCreated, CreateNotificationChannelResponse{
				Message: "email notification channel created",
			})
		case "push":
			if strings.TrimSpace(req.FCMToken) == "" {
				utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "fcmToken is required for type=push"})
				return
			}

			endpointArn, err := service.RegisterPushEndpoint(ctx, strings.TrimSpace(req.FCMToken))
			if err != nil {
				utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}

			utils.WriteJSON(w, http.StatusCreated, CreateNotificationChannelResponse{
				Message:     "push notification channel created",
				EndpointArn: endpointArn,
			})
		default:
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "type must be either 'email' or 'push'"})
		}
	}
}

func SendNotification(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = database
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		var req SendNotificationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		if strings.TrimSpace(req.Message) == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "message is required"})
			return
		}

		service, err := awsservice.NewNotificationService(ctx)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		switch strings.ToLower(strings.TrimSpace(req.Type)) {
		case "email":
			if strings.TrimSpace(req.Email) == "" {
				utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "email is required for type=email"})
				return
			}

			subject := strings.TrimSpace(req.Subject)
			if subject == "" {
				subject = "Alive Notification"
			}

			if err := service.SendEmail(ctx, strings.TrimSpace(req.Email), subject, req.Message); err != nil {
				utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}

			utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "email sent successfully"})
		case "push":
			if strings.TrimSpace(req.EndpointArn) == "" {
				utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "endpointArn is required for type=push"})
				return
			}

			if err := service.PublishPushToEndpoint(ctx, strings.TrimSpace(req.EndpointArn), awsservice.PushPayload{
				Title: strings.TrimSpace(req.Title),
				Body:  req.Message,
			}); err != nil {
				utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}

			utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "push notification sent successfully"})
		default:
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "type must be either 'email' or 'push'"})
		}
	}
}
