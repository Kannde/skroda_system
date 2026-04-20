package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotifEmail NotificationType = "email"
	NotifSMS   NotificationType = "sms"
	NotifPush  NotificationType = "push"
)

type NotificationStatus string

const (
	NotifPending NotificationStatus = "pending"
	NotifSent    NotificationStatus = "sent"
	NotifFailed  NotificationStatus = "failed"
)

type Notification struct {
	ID         uuid.UUID          `json:"id" db:"id"`
	UserID     uuid.UUID          `json:"user_id" db:"user_id"`
	Type       NotificationType   `json:"type" db:"type"`
	Subject    string             `json:"subject" db:"subject"`
	Body       string             `json:"body" db:"body"`
	Status     NotificationStatus `json:"status" db:"status"`
	SentAt     *time.Time         `json:"sent_at,omitempty" db:"sent_at"`
	CreatedAt  time.Time          `json:"created_at" db:"created_at"`
}
