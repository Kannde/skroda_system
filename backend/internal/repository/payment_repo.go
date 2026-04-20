package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skroda/backend/internal/models"
)

var ErrPaymentNotFound = errors.New("payment not found")

type PaymentRepository struct {
	pool *pgxpool.Pool
}

func NewPaymentRepository(pool *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{pool: pool}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	query := `
		INSERT INTO payments (
			id, transaction_id, payment_type, payment_method, status,
			amount, currency, payer_id, payee_id, idempotency_key
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING initiated_at, created_at, updated_at`

	payment.ID = uuid.New()

	return r.pool.QueryRow(ctx, query,
		payment.ID, payment.TransactionID, payment.PaymentType,
		payment.PaymentMethod, payment.Status,
		payment.Amount, payment.Currency,
		payment.PayerID, payment.PayeeID,
		payment.IdempotencyKey,
	).Scan(&payment.InitiatedAt, &payment.CreatedAt, &payment.UpdatedAt)
}

func (r *PaymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	query := `
		SELECT id, transaction_id, payment_type, payment_method, status,
		       amount, currency, payer_id, payee_id,
		       provider_reference, provider_metadata, idempotency_key,
		       initiated_at, completed_at, failed_at, failure_reason,
		       created_at, updated_at
		FROM payments WHERE id = $1`

	p := &models.Payment{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.TransactionID, &p.PaymentType, &p.PaymentMethod, &p.Status,
		&p.Amount, &p.Currency, &p.PayerID, &p.PayeeID,
		&p.ProviderReference, &p.ProviderMetadata, &p.IdempotencyKey,
		&p.InitiatedAt, &p.CompletedAt, &p.FailedAt, &p.FailureReason,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPaymentNotFound
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	return p, nil
}

func (r *PaymentRepository) GetByIdempotencyKey(ctx context.Context, key string) (*models.Payment, error) {
	query := `
		SELECT id, transaction_id, payment_type, payment_method, status,
		       amount, currency, payer_id, payee_id,
		       provider_reference, idempotency_key,
		       initiated_at, completed_at, failed_at, failure_reason,
		       created_at, updated_at
		FROM payments WHERE idempotency_key = $1`

	p := &models.Payment{}
	err := r.pool.QueryRow(ctx, query, key).Scan(
		&p.ID, &p.TransactionID, &p.PaymentType, &p.PaymentMethod, &p.Status,
		&p.Amount, &p.Currency, &p.PayerID, &p.PayeeID,
		&p.ProviderReference, &p.IdempotencyKey,
		&p.InitiatedAt, &p.CompletedAt, &p.FailedAt, &p.FailureReason,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPaymentNotFound
		}
		return nil, fmt.Errorf("failed to get payment by idempotency key: %w", err)
	}
	return p, nil
}

func (r *PaymentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.PaymentStatus, providerRef *string, failureReason *string) error {
	var err error
	switch status {
	case models.PayFailed:
		_, err = r.pool.Exec(ctx,
			`UPDATE payments SET status=$1, provider_reference=$2, failure_reason=$3, failed_at=NOW(), updated_at=NOW() WHERE id=$4`,
			status, providerRef, failureReason, id)
	case models.PayCompleted, models.PayReleased:
		_, err = r.pool.Exec(ctx,
			`UPDATE payments SET status=$1, provider_reference=$2, completed_at=NOW(), updated_at=NOW() WHERE id=$3`,
			status, providerRef, id)
	default:
		_, err = r.pool.Exec(ctx,
			`UPDATE payments SET status=$1, provider_reference=$2, updated_at=NOW() WHERE id=$3`,
			status, providerRef, id)
	}
	return err
}

func (r *PaymentRepository) ListByTransaction(ctx context.Context, txnID uuid.UUID) ([]models.Payment, error) {
	query := `
		SELECT id, transaction_id, payment_type, payment_method, status,
		       amount, currency, payer_id, payee_id,
		       provider_reference, idempotency_key,
		       initiated_at, completed_at, failed_at, failure_reason,
		       created_at, updated_at
		FROM payments WHERE transaction_id = $1 ORDER BY created_at ASC`

	rows, err := r.pool.Query(ctx, query, txnID)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var p models.Payment
		if err := rows.Scan(
			&p.ID, &p.TransactionID, &p.PaymentType, &p.PaymentMethod, &p.Status,
			&p.Amount, &p.Currency, &p.PayerID, &p.PayeeID,
			&p.ProviderReference, &p.IdempotencyKey,
			&p.InitiatedAt, &p.CompletedAt, &p.FailedAt, &p.FailureReason,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}
		payments = append(payments, p)
	}
	return payments, nil
}
