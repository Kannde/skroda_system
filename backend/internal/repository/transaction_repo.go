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

var (
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrInvalidTransition   = errors.New("invalid status transition")
)

type TransactionRepository struct {
	pool *pgxpool.Pool
}

func NewTransactionRepository(pool *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{pool: pool}
}

const txnSelectCols = `
	id, reference_code, transaction_type, title, description, status,
	seller_id, buyer_id, agent_id,
	amount, currency, fee_amount, fee_paid_by,
	seller_city, buyer_city,
	inspection_hours, inspection_ends_at,
	agreed_at, funded_at, shipped_at, agent_received_at,
	delivered_at, completed_at, cancelled_at,
	invite_token, invite_expires_at,
	created_at, updated_at`

func scanTxn(row interface{ Scan(...any) error }) (*models.Transaction, error) {
	txn := &models.Transaction{}
	err := row.Scan(
		&txn.ID, &txn.ReferenceCode, &txn.TransactionType, &txn.Title, &txn.Description, &txn.Status,
		&txn.SellerID, &txn.BuyerID, &txn.AgentID,
		&txn.Amount, &txn.Currency, &txn.FeeAmount, &txn.FeePaidBy,
		&txn.SellerCity, &txn.BuyerCity,
		&txn.InspectionHours, &txn.InspectionEndsAt,
		&txn.AgreedAt, &txn.FundedAt, &txn.ShippedAt, &txn.AgentReceivedAt,
		&txn.DeliveredAt, &txn.CompletedAt, &txn.CancelledAt,
		&txn.InviteToken, &txn.InviteExpiresAt,
		&txn.CreatedAt, &txn.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return txn, nil
}

func (r *TransactionRepository) Create(ctx context.Context, txn *models.Transaction) error {
	query := `
		INSERT INTO transactions (
			id, reference_code, transaction_type, title, description, status,
			seller_id, amount, currency, fee_amount, fee_paid_by,
			seller_city, buyer_city, inspection_hours,
			invite_token, invite_expires_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,
			$7,$8,$9,$10,$11,
			$12,$13,$14,
			$15,$16
		) RETURNING created_at, updated_at`

	txn.ID = uuid.New()

	return r.pool.QueryRow(ctx, query,
		txn.ID, txn.ReferenceCode, txn.TransactionType, txn.Title, txn.Description, txn.Status,
		txn.SellerID, txn.Amount, txn.Currency, txn.FeeAmount, txn.FeePaidBy,
		txn.SellerCity, txn.BuyerCity, txn.InspectionHours,
		txn.InviteToken, txn.InviteExpiresAt,
	).Scan(&txn.CreatedAt, &txn.UpdatedAt)
}

func (r *TransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT `+txnSelectCols+` FROM transactions WHERE id = $1`, id)
	txn, err := scanTxn(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTransactionNotFound
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	return txn, nil
}

func (r *TransactionRepository) GetByInviteToken(ctx context.Context, token string) (*models.Transaction, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT `+txnSelectCols+` FROM transactions WHERE invite_token = $1`, token)
	txn, err := scanTxn(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTransactionNotFound
		}
		return nil, fmt.Errorf("failed to get transaction by invite: %w", err)
	}
	return txn, nil
}

func (r *TransactionRepository) ListByUser(ctx context.Context, userID uuid.UUID, status string, limit, offset int) ([]models.Transaction, int, error) {
	countQuery := `SELECT COUNT(*) FROM transactions WHERE (seller_id = $1 OR buyer_id = $1)`
	countArgs := []interface{}{userID}

	if status != "" {
		countQuery += " AND status = $2"
		countArgs = append(countArgs, status)
	}

	var total int
	if err := r.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	query := `SELECT ` + txnSelectCols + ` FROM transactions WHERE (seller_id = $1 OR buyer_id = $1)`
	fetchArgs := []interface{}{userID}
	paramIdx := 2

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", paramIdx)
		fetchArgs = append(fetchArgs, status)
		paramIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", paramIdx, paramIdx+1)
	fetchArgs = append(fetchArgs, limit, offset)

	rows, err := r.pool.Query(ctx, query, fetchArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list transactions: %w", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		txn, err := scanTxn(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, *txn)
	}

	return transactions, total, nil
}

func (r *TransactionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.TransactionStatus, timestampField string) error {
	query := fmt.Sprintf(
		`UPDATE transactions SET status = $1, %s = NOW(), updated_at = NOW() WHERE id = $2`,
		timestampField,
	)
	result, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrTransactionNotFound
	}
	return nil
}

func (r *TransactionRepository) AssignBuyer(ctx context.Context, txnID, buyerID uuid.UUID, buyerCity string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE transactions SET buyer_id = $1, buyer_city = $2, status = 'negotiation', updated_at = NOW() WHERE id = $3`,
		buyerID, buyerCity, txnID)
	return err
}

func (r *TransactionRepository) AssignAgent(ctx context.Context, txnID, agentID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE transactions SET agent_id = $1, updated_at = NOW() WHERE id = $2`, agentID, txnID)
	return err
}
