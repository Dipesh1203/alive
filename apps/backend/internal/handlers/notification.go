package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Dipesh1203/alive/backend/db"

	awsservice "github.com/Dipesh1203/alive/backend/internal/services/aws"
	"github.com/Dipesh1203/alive/backend/internal/utils"
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

type TestEmailRequest struct {
	Email   string `json:"email"`
	Subject string `json:"subject,omitempty"`
	Message string `json:"message,omitempty"`
}

type TestTemplateEmailRequest struct {
	Email        string         `json:"email"`
	Subject      string         `json:"subject,omitempty"`
	TemplateName string         `json:"templateName"`
	Data         map[string]any `json:"data,omitempty"`
}

// CreateNotificationChannel godoc
// @Summary      Create a notification channel
// @Description  Sets up email or push notification channel for the user
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Param        Authorization  header   string                              true  "Bearer token"
// @Param        body           body     CreateNotificationChannelRequest   true  "Notification channel request"
// @Success      201            {object} CreateNotificationChannelResponse
// @Failure      400            {object} map[string]string
// @Failure      500            {object} map[string]string
// @Router       /api/notifications/channels [post]
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

// SendNotification godoc
// @Summary      Send a test notification
// @Description  Sends a test notification via email or push
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Param        Authorization  header   string                   true  "Bearer token"
// @Param        body           body     SendNotificationRequest  true  "Send notification request"
// @Success      200            {object} map[string]string
// @Failure      400            {object} map[string]string
// @Failure      500            {object} map[string]string
// @Router       /api/notifications/test [post]
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

// TestEmail godoc
// @Summary      Test SES email sending
// @Description  Sends a test email through AWS SES
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Param        Authorization  header   string            true  "Bearer token"
// @Param        body           body     TestEmailRequest   true  "Test email request"
// @Success      200            {object} map[string]string
// @Failure      400            {object} map[string]string
// @Failure      500            {object} map[string]string
// @Router       /api/notifications/test-email [post]
func TestEmail(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = database
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		var req TestEmailRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		email := strings.TrimSpace(req.Email)
		if email == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "email is required"})
			return
		}

		subject := strings.TrimSpace(req.Subject)
		if subject == "" {
			subject = "Alive SES Test Email"
		}

		message := strings.TrimSpace(req.Message)
		if message == "" {
			message = "This is a test email from Alive to verify AWS SES setup."
		}

		service, err := awsservice.NewNotificationService(ctx)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		if err := service.SendEmail(ctx, email, subject, message); err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "test email sent successfully"})
	}
}

// TestTemplateEmail godoc
// @Summary      Test template email sending
// @Description  Renders an HTML template and sends it through AWS SES
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Param        Authorization  header   string                     true  "Bearer token"
// @Param        body           body     TestTemplateEmailRequest   true  "Template email request"
// @Success      200            {object} map[string]string
// @Failure      400            {object} map[string]string
// @Failure      500            {object} map[string]string
// @Router       /api/notifications/test-template-email [post]
func TestTemplateEmail(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = database
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		var req TestTemplateEmailRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		email := strings.TrimSpace(req.Email)
		if email == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "email is required"})
			return
		}

		templateName := strings.TrimSpace(req.TemplateName)
		if templateName == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "templateName is required"})
			return
		}

		subject := strings.TrimSpace(req.Subject)
		if subject == "" {
			subject = "Alive Template Email"
		}

		service, err := awsservice.NewNotificationService(ctx)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		if err := service.SendTemplateEmail(ctx, email, subject, templateName, req.Data); err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "template email sent successfully"})
	}
}
