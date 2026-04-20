CREATE TYPE term_status AS ENUM ('proposed', 'accepted', 'rejected', 'counter_proposed');

CREATE TABLE transaction_terms (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id      UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    version             INTEGER NOT NULL DEFAULT 1,
    proposed_by         UUID NOT NULL REFERENCES users(id),
    status              term_status NOT NULL DEFAULT 'proposed',

    item_description    TEXT NOT NULL,
    item_condition      VARCHAR(50) NOT NULL DEFAULT 'new',
    quantity            INTEGER NOT NULL DEFAULT 1,
    amount              DECIMAL(12,2) NOT NULL,
    inspection_hours    INTEGER NOT NULL DEFAULT 48,
    delivery_deadline   INTEGER NOT NULL DEFAULT 72,

    acceptance_criteria TEXT,
    image_urls          JSONB NOT NULL DEFAULT '[]'::jsonb,
    rejection_reason    TEXT,
    responded_at        TIMESTAMPTZ,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_terms_txn ON transaction_terms(transaction_id);
CREATE INDEX idx_terms_txn_version ON transaction_terms(transaction_id, version);
