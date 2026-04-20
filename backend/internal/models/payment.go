package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PayPending    PaymentStatus = "pending"
	PayProcessing PaymentStatus = "processing"
	PayCompleted  PaymentStatus = "completed"
	PayFailed     PaymentStatus = "failed"
	PayRefunded   PaymentStatus = "refunded"
	PayReleased   PaymentStatus = "released"
)

type PaymentType string

const (
	PayTypeEscrowDeposit PaymentType = "escrow_deposit"
	PayTypeSellerRelease PaymentType = "seller_release"
	PayTypeRefund        PaymentType = "refund"
	PayTypeFeeCollection PaymentType = "fee_collection"
)

type PaymentMethod string

const (
	MethodMoMoMTN     PaymentMethod = "momo_mtn"
	MethodMoMoTelecel PaymentMethod = "momo_telecel"
	MethodMoMoAT      PaymentMethod = "momo_at"
	MethodBank        PaymentMethod = "bank_transfer"
	MethodCard        PaymentMethod = "card"
)

type Payment struct {
	ID                uuid.UUID     `json:"id" db:"id"`
	TransactionID     uuid.UUID     `json:"transaction_id" db:"transaction_id"`
	PaymentType       PaymentType   `json:"payment_type" db:"payment_type"`
	PaymentMethod     PaymentMethod `json:"payment_method" db:"payment_method"`
	Status            PaymentStatus `json:"status" db:"status"`
	Amount            float64       `json:"amount" db:"amount"`
	Currency          string        `json:"currency" db:"currency"`
	PayerID           uuid.UUID     `json:"payer_id" db:"payer_id"`
	PayeeID           *uuid.UUID    `json:"payee_id,omitempty" db:"payee_id"`
	ProviderReference *string       `json:"provider_reference,omitempty" db:"provider_reference"`
	ProviderMetadata  interface{}   `json:"provider_metadata,omitempty" db:"provider_metadata"`
	IdempotencyKey    string        `json:"idempotency_key" db:"idempotency_key"`
	InitiatedAt       time.Time     `json:"initiated_at" db:"initiated_at"`
	CompletedAt       *time.Time    `json:"completed_at,omitempty" db:"completed_at"`
	FailedAt          *time.Time    `json:"failed_at,omitempty" db:"failed_at"`
	FailureReason     *string       `json:"failure_reason,omitempty" db:"failure_reason"`
	CreatedAt         time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at" db:"updated_at"`
}

type InitiatePaymentRequest struct {
	TransactionID uuid.UUID `json:"transaction_id" binding:"required"`
	PaymentMethod string    `json:"payment_method" binding:"required,oneof=momo_mtn momo_telecel momo_at bank_transfer card"`
	Phone         string    `json:"phone" binding:"required_if=PaymentMethod momo_mtn,required_if=PaymentMethod momo_telecel,required_if=PaymentMethod momo_at"`
}
