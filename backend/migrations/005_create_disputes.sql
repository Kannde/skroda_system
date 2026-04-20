CREATE TYPE dispute_status AS ENUM ('open', 'under_review', 'resolved', 'closed');
CREATE TYPE dispute_resolution AS ENUM ('buyer', 'seller', 'split');

CREATE TABLE disputes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL REFERENCES transactions(id),
    raised_by_id UUID NOT NULL REFERENCES users(id),
    agent_id UUID REFERENCES users(id),
    reason TEXT NOT NULL,
    status dispute_status NOT NULL DEFAULT 'open',
    resolution dispute_resolution,
    resolution_note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ
);

CREATE INDEX idx_disputes_transaction ON disputes(transaction_id);
CREATE INDEX idx_disputes_status ON disputes(status);
