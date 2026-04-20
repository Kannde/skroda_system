package services

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/skroda/backend/internal/models"
)

type NotificationService struct {
	resendAPIKey string
	smsAPIKey    string
	smsAPIURL    string
}

func NewNotificationService(resendAPIKey, smsAPIKey, smsAPIURL string) *NotificationService {
	return &NotificationService{resendAPIKey: resendAPIKey, smsAPIKey: smsAPIKey, smsAPIURL: smsAPIURL}
}

func (s *NotificationService) SendEmail(ctx context.Context, userID uuid.UUID, subject, body string) error {
	log.Printf("[email] to=%s subject=%s", userID, subject)
	_ = &models.Notification{
		UserID:  userID,
		Type:    models.NotifEmail,
		Subject: subject,
		Body:    body,
	}
	return nil
}

func (s *NotificationService) SendSMS(ctx context.Context, userID uuid.UUID, body string) error {
	log.Printf("[sms] to=%s body=%s", userID, body)
	return nil
}
