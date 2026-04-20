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

const txSelectCols = `id, reference_code, transaction_type, title, description, status,
	seller_id, buyer_id, agent_id,
	amount, currency, fee_amount, fee_paid_by,
	seller_city, buyer_city,
	inspection_hours, inspection_ends_at,
	agreed_at, funded_at, shipped_at, agent_received_at, delivered_at, completed_at, cancelled_at,
	invite_token, invite_expires_at,
	created_at, updated_at`

func scanTx(row interface{ Scan(...interface{}) error }) (*models.Transaction, error) {
	tx := &models.Transaction{}
	err := row.Scan(
		&tx.ID, &tx.ReferenceCode, &tx.TransactionType, &tx.Title, &tx.Description, &tx.Status,
		&tx.SellerID, &tx.BuyerID, &tx.AgentID,
		&tx.Amount, &tx.Currency, &tx.FeeAmount, &tx.FeePaidBy,
		&tx.SellerCity, &tx.BuyerCity,
		&tx.InspectionHours, &tx.InspectionEndsAt,
		&tx.AgreedAt, &tx.FundedAt, &tx.ShippedAt, &tx.AgentReceivedAt, &tx.DeliveredAt, &tx.CompletedAt, &tx.CancelledAt,
		&tx.InviteToken, &tx.InviteExpiresAt,
		&tx.CreatedAt, &tx.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *TransactionRepository) Create(ctx context.Context, tx *models.Transaction) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO transactions
			(id, reference_code, transaction_type, title, description, status,
			 seller_id, buyer_id, agent_id,
			 amount, currency, fee_amount, fee_paid_by,
			 seller_city, buyer_city,
			 inspection_hours, invite_token, invite_expires_at,
			 created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)`,
		tx.ID, tx.ReferenceCode, tx.TransactionType, tx.Title, tx.Description, tx.Status,
		tx.SellerID, tx.BuyerID, tx.AgentID,
		tx.Amount, tx.Currency, tx.FeeAmount, tx.FeePaidBy,
		tx.SellerCity, tx.BuyerCity,
		tx.InspectionHours, tx.InviteToken, tx.InviteExpiresAt,
		tx.CreatedAt, tx.UpdatedAt,
	)
	return err
}

func (r *TransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	row := r.db.QueryRow(ctx, `SELECT `+txSelectCols+` FROM transactions WHERE id = $1`, id)
	return scanTx(row)
}

func (r *TransactionRepository) GetByReference(ctx context.Context, ref string) (*models.Transaction, error) {
	row := r.db.QueryRow(ctx, `SELECT `+txSelectCols+` FROM transactions WHERE reference_code = $1`, ref)
	return scanTx(row)
}

func (r *TransactionRepository) GetByInviteToken(ctx context.Context, token string) (*models.Transaction, error) {
	row := r.db.QueryRow(ctx, `SELECT `+txSelectCols+` FROM transactions WHERE invite_token = $1`, token)
	return scanTx(row)
}

func (r *TransactionRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+txSelectCols+` FROM transactions
		WHERE seller_id = $1 OR buyer_id = $1
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []models.Transaction
	for rows.Next() {
		tx, err := scanTx(rows)
		if err != nil {
			return nil, err
		}
		txs = append(txs, *tx)
	}
	return txs, nil
}

func (r *TransactionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.TransactionStatus) error {
	_, err := r.db.Exec(ctx,
		`UPDATE transactions SET status = $1, updated_at = NOW() WHERE id = $2`, status, id)
	return err
}

func (r *TransactionRepository) AssignBuyer(ctx context.Context, id, buyerID uuid.UUID, buyerCity string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE transactions SET buyer_id = $1, buyer_city = $2, updated_at = NOW() WHERE id = $3`,
		buyerID, buyerCity, id)
	return err
}

func (r *TransactionRepository) AssignAgent(ctx context.Context, id, agentID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE transactions SET agent_id = $1, updated_at = NOW() WHERE id = $2`, agentID, id)
	return err
}
