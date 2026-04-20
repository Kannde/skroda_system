CREATE TABLE terms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    buyer_agreed BOOLEAN NOT NULL DEFAULT FALSE,
    seller_agreed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_terms_transaction ON terms(transaction_id);
