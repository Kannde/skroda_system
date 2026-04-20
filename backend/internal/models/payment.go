package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentSuccess   PaymentStatus = "success"
	PaymentFailed    PaymentStatus = "failed"
	PaymentRefunded  PaymentStatus = "refunded"
)

type PaymentProvider string

const (
	ProviderMoMo   PaymentProvider = "momo"
	ProviderStripe PaymentProvider = "stripe"
	ProviderManual PaymentProvider = "manual"
)

type Payment struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	TransactionID   uuid.UUID       `json:"transaction_id" db:"transaction_id"`
	PayerID         uuid.UUID       `json:"payer_id" db:"payer_id"`
	Amount          float64         `json:"amount" db:"amount"`
	Currency        string          `json:"currency" db:"currency"`
	Provider        PaymentProvider `json:"provider" db:"provider"`
	ProviderRef     string          `json:"provider_ref,omitempty" db:"provider_ref"`
	Status          PaymentStatus   `json:"status" db:"status"`
	FailureReason   string          `json:"failure_reason,omitempty" db:"failure_reason"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

type InitiatePaymentRequest struct {
	TransactionID uuid.UUID       `json:"transaction_id" validate:"required"`
	Provider      PaymentProvider `json:"provider" validate:"required,oneof=momo stripe manual"`
	PhoneNumber   string          `json:"phone_number"`
}
