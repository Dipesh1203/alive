package awsservice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	snssdk "github.com/aws/aws-sdk-go-v2/service/sns"

	snsservice "github.com/Dipesh1203/alive/backend/internal/services/aws/sns"
)

type NotificationService struct {
	sesClient      *sesv2.Client
	snsActions     snsservice.SnsActions
	platformAppArn string
	fromEmail      string
}

type PushPayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func NewNotificationService(ctx context.Context) (*NotificationService, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	return &NotificationService{
		sesClient:      sesv2.NewFromConfig(cfg),
		snsActions:     snsservice.SnsActions{SnsClient: snssdk.NewFromConfig(cfg)},
		platformAppArn: os.Getenv("AWS_SNS_PLATFORM_APPLICATION_ARN"),
		fromEmail:      os.Getenv("AWS_SES_FROM_EMAIL"),
	}, nil
}

func (s *NotificationService) SendEmail(ctx context.Context, toEmail string, subject string, body string) error {
	if s.fromEmail == "" {
		return errors.New("AWS_SES_FROM_EMAIL is not configured")
	}

	_, err := s.sesClient.SendEmail(ctx, &sesv2.SendEmailInput{
		FromEmailAddress: &s.fromEmail,
		Destination: &types.Destination{
			ToAddresses: []string{toEmail},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{Data: &subject},
				Body: &types.Body{
					Text: &types.Content{Data: &body},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send email with ses: %w", err)
	}

	return nil
}

func (s *NotificationService) SendHTMLEmail(ctx context.Context, toEmail string, subject string, htmlBody string) error {
	if s.fromEmail == "" {
		return errors.New("AWS_SES_FROM_EMAIL is not configured")
	}

	_, err := s.sesClient.SendEmail(ctx, &sesv2.SendEmailInput{
		FromEmailAddress: &s.fromEmail,
		Destination: &types.Destination{
			ToAddresses: []string{toEmail},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{Data: &subject},
				Body: &types.Body{
					Html: &types.Content{Data: &htmlBody},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send html email with ses: %w", err)
	}

	return nil
}

func (s *NotificationService) SendTemplateEmail(ctx context.Context, toEmail string, subject string, templateName string, data map[string]any) error {
	htmlBody, err := renderEmailTemplate(templateName, data)
	if err != nil {
		return err
	}

	return s.SendHTMLEmail(ctx, toEmail, subject, htmlBody)
}

func renderEmailTemplate(templateName string, data map[string]any) (string, error) {
	if templateName == "" {
		return "", errors.New("template name is required")
	}

	if templateName != filepath.Base(templateName) {
		return "", errors.New("template name is invalid")
	}

	templateDir := os.Getenv("EMAIL_TEMPLATE_DIR")
	if templateDir == "" {
		templateDir = filepath.Join("internal", "templates", "email")
	}

	templatePath := filepath.Join(templateDir, templateName)
	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse email template: %w", err)
	}

	var out bytes.Buffer
	if err := tpl.Execute(&out, data); err != nil {
		return "", fmt.Errorf("failed to render email template: %w", err)
	}

	return out.String(), nil
}

func (s *NotificationService) RegisterPushEndpoint(ctx context.Context, fcmToken string) (string, error) {
	if s.platformAppArn == "" {
		return "", errors.New("AWS_SNS_PLATFORM_APPLICATION_ARN is not configured")
	}

	endpointArn, err := s.snsActions.CreateUserEndpoint(ctx, s.platformAppArn, fcmToken)
	if err != nil {
		return "", fmt.Errorf("failed to create sns push endpoint: %w", err)
	}

	return endpointArn, nil
}

func (s *NotificationService) PublishPushToEndpoint(ctx context.Context, endpointArn string, payload PushPayload) error {
	message := payload.Body
	if payload.Title != "" {
		message = payload.Title + ": " + payload.Body
	}

	if payload.Title != "" || payload.Body != "" {
		gcmPayload := map[string]map[string]string{
			"notification": {
				"title": payload.Title,
				"body":  payload.Body,
			},
		}
		if payload.Title != "" && payload.Body != "" {
			if encoded, err := json.Marshal(gcmPayload); err == nil {
				message = string(encoded)
			}
		}
	}

	if err := s.snsActions.PublishToUser(ctx, endpointArn, message); err != nil {
		return fmt.Errorf("failed to publish sns push notification: %w", err)
	}

	return nil
}
