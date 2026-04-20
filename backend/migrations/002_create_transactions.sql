CREATE TYPE transaction_status AS ENUM (
    'pending', 'funded', 'in_progress', 'inspection',
    'completed', 'disputed', 'cancelled', 'refunded'
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    buyer_id UUID NOT NULL REFERENCES users(id),
    seller_id UUID NOT NULL REFERENCES users(id),
    agent_id UUID REFERENCES users(id),
    title VARCHAR(200) NOT NULL,
    description TEXT,
    amount NUMERIC(15, 2) NOT NULL CHECK (amount > 0),
    currency CHAR(3) NOT NULL DEFAULT 'NGN',
    status transaction_status NOT NULL DEFAULT 'pending',
    buyer_city VARCHAR(100),
    seller_city VARCHAR(100) NOT NULL,
    inspection_ends_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transactions_buyer ON transactions(buyer_id);
CREATE INDEX idx_transactions_seller ON transactions(seller_id);
CREATE INDEX idx_transactions_status ON transactions(status);
