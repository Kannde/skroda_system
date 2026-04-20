CREATE TYPE dispute_status AS ENUM ('open', 'under_review', 'resolved_buyer', 'resolved_seller', 'escalated');
CREATE TYPE dispute_reason AS ENUM (
    'item_not_as_described',
    'item_damaged',
    'item_not_delivered',
    'wrong_item',
    'partial_delivery',
    'seller_no_ship',
    'other'
);

CREATE TABLE disputes (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id      UUID NOT NULL REFERENCES transactions(id),
    raised_by           UUID NOT NULL REFERENCES users(id),
    reason              dispute_reason NOT NULL,
    description         TEXT NOT NULL,
    status              dispute_status NOT NULL DEFAULT 'open',

    evidence            JSONB NOT NULL DEFAULT '[]'::jsonb,

    resolved_by         UUID REFERENCES users(id),
    resolution_notes    TEXT,
    resolved_at         TIMESTAMPTZ,

    response_deadline   TIMESTAMPTZ NOT NULL,
    resolution_deadline TIMESTAMPTZ NOT NULL,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE dispute_messages (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    dispute_id          UUID NOT NULL REFERENCES disputes(id) ON DELETE CASCADE,
    sender_id           UUID NOT NULL REFERENCES users(id),
    message             TEXT NOT NULL,
    attachments         JSONB DEFAULT '[]'::jsonb,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_disputes_txn ON disputes(transaction_id);
CREATE INDEX idx_disputes_status ON disputes(status);
CREATE INDEX idx_dispute_msgs_dispute ON dispute_messages(dispute_id);
