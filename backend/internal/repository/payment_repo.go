package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skroda/backend/internal/models"
)

type PaymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(ctx context.Context, p *models.Payment) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO payments (id, transaction_id, payer_id, amount, currency, provider, provider_ref, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		p.ID, p.TransactionID, p.PayerID, p.Amount, p.Currency, p.Provider, p.ProviderRef, p.Status, p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func (r *PaymentRepository) UpdateStatus(ctx context.Context, providerRef string, status models.PaymentStatus) error {
	_, err := r.db.Exec(ctx, `UPDATE payments SET status = $1, updated_at = NOW() WHERE provider_ref = $2`, status, providerRef)
	return err
}
