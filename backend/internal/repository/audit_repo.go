package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditRepository struct {
	pool *pgxpool.Pool
}

func NewAuditRepository(pool *pgxpool.Pool) *AuditRepository {
	return &AuditRepository{pool: pool}
}

func (r *AuditRepository) Log(ctx context.Context, entityType string, entityID uuid.UUID, action string, oldValue, newValue interface{}, performedBy *uuid.UUID, ipAddress *string) error {
	oldJSON, _ := json.Marshal(oldValue)
	newJSON, _ := json.Marshal(newValue)

	_, err := r.pool.Exec(ctx,
		`INSERT INTO audit_log (id, entity_type, entity_id, action, old_value, new_value, performed_by, ip_address)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		uuid.New(), entityType, entityID, action,
		oldJSON, newJSON, performedBy, ipAddress,
	)
	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}
	return nil
}
