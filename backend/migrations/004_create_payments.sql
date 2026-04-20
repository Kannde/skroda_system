CREATE TYPE payment_status AS ENUM ('pending', 'success', 'failed', 'refunded');
CREATE TYPE payment_provider AS ENUM ('momo', 'stripe', 'manual');

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL REFERENCES transactions(id),
    payer_id UUID NOT NULL REFERENCES users(id),
    amount NUMERIC(15, 2) NOT NULL,
    currency CHAR(3) NOT NULL,
    provider payment_provider NOT NULL,
    provider_ref VARCHAR(255),
    status payment_status NOT NULL DEFAULT 'pending',
    failure_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_transaction ON payments(transaction_id);
CREATE INDEX idx_payments_provider_ref ON payments(provider_ref);
