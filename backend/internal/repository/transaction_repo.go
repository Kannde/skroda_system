package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skroda/backend/internal/models"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, tx *models.Transaction) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO transactions (id, buyer_id, seller_id, agent_id, title, description, amount, currency, status, seller_city, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		tx.ID, tx.BuyerID, tx.SellerID, tx.AgentID, tx.Title, tx.Description,
		tx.Amount, tx.Currency, tx.Status, tx.SellerCity, tx.CreatedAt, tx.UpdatedAt,
	)
	return err
}

func (r *TransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	tx := &models.Transaction{}
	err := r.db.QueryRow(ctx, `
		SELECT id, buyer_id, seller_id, agent_id, title, description, amount, currency, status, buyer_city, seller_city, inspection_ends_at, created_at, updated_at
		FROM transactions WHERE id = $1`, id).
		Scan(&tx.ID, &tx.BuyerID, &tx.SellerID, &tx.AgentID, &tx.Title, &tx.Description,
			&tx.Amount, &tx.Currency, &tx.Status, &tx.BuyerCity, &tx.SellerCity, &tx.InspectionEndsAt, &tx.CreatedAt, &tx.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *TransactionRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, buyer_id, seller_id, agent_id, title, description, amount, currency, status, buyer_city, seller_city, inspection_ends_at, created_at, updated_at
		FROM transactions WHERE buyer_id = $1 OR seller_id = $1
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		if err := rows.Scan(&tx.ID, &tx.BuyerID, &tx.SellerID, &tx.AgentID, &tx.Title, &tx.Description,
			&tx.Amount, &tx.Currency, &tx.Status, &tx.BuyerCity, &tx.SellerCity, &tx.InspectionEndsAt, &tx.CreatedAt, &tx.UpdatedAt); err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}
	return txs, nil
}

func (r *TransactionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.TransactionStatus) error {
	_, err := r.db.Exec(ctx, `UPDATE transactions SET status = $1, updated_at = NOW() WHERE id = $2`, status, id)
	return err
}
