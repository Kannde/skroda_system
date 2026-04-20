package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationChannel string

const (
	ChannelSMS   NotificationChannel = "sms"
	ChannelEmail NotificationChannel = "email"
	ChannelPush  NotificationChannel = "push"
)

type Notification struct {
	ID            uuid.UUID           `json:"id" db:"id"`
	UserID        uuid.UUID           `json:"user_id" db:"user_id"`
	TransactionID *uuid.UUID          `json:"transaction_id,omitempty" db:"transaction_id"`
	Channel       NotificationChannel `json:"channel" db:"channel"`
	Subject       string              `json:"subject" db:"subject"`
	Body          string              `json:"body" db:"body"`
	Sent          bool                `json:"sent" db:"sent"`
	SentAt        *time.Time          `json:"sent_at,omitempty" db:"sent_at"`
	FailureReason *string             `json:"failure_reason,omitempty" db:"failure_reason"`
	CreatedAt     time.Time           `json:"created_at" db:"created_at"`
}
