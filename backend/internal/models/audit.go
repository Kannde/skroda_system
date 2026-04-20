package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditEntry struct {
	ID          uuid.UUID   `json:"id" db:"id"`
	EntityType  string      `json:"entity_type" db:"entity_type"`
	EntityID    uuid.UUID   `json:"entity_id" db:"entity_id"`
	Action      string      `json:"action" db:"action"`
	OldValue    interface{} `json:"old_value,omitempty" db:"old_value"`
	NewValue    interface{} `json:"new_value,omitempty" db:"new_value"`
	PerformedBy *uuid.UUID  `json:"performed_by,omitempty" db:"performed_by"`
	IPAddress   *string     `json:"ip_address,omitempty" db:"ip_address"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
}
