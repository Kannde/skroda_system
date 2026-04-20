-- +migrate Up
CREATE TYPE payment_status AS ENUM (
    'pending',
    'processing',
    'completed',
    'failed',
    'refunded',
    'released'
);

CREATE TYPE payment_type AS ENUM (
    'escrow_deposit',
    'seller_release',
    'refund',
    'fee_collection'
);

CREATE TYPE payment_method AS ENUM ('momo_mtn', 'momo_telecel', 'momo_at', 'bank_transfer', 'card');

CREATE TABLE payments (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id      UUID NOT NULL REFERENCES transactions(id),
    payment_type        payment_type NOT NULL,
    payment_method      payment_method NOT NULL,
    status              payment_status NOT NULL DEFAULT 'pending',

    amount              DECIMAL(12,2) NOT NULL,
    currency            VARCHAR(3) NOT NULL DEFAULT 'GHS',

    payer_id            UUID NOT NULL REFERENCES users(id),
    payee_id            UUID REFERENCES users(id),

    provider_reference  VARCHAR(255),
    provider_metadata   JSONB DEFAULT '{}'::jsonb,

    idempotency_key     VARCHAR(64) UNIQUE NOT NULL,

    initiated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at        TIMESTAMPTZ,
    failed_at           TIMESTAMPTZ,
    failure_reason      TEXT,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_txn ON payments(transaction_id);
CREATE INDEX idx_payments_payer ON payments(payer_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_idempotency ON payments(idempotency_key);
