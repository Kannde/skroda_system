-- +migrate Up
CREATE TYPE transaction_status AS ENUM (
    'draft',
    'negotiation',
    'agreed',
    'funded',
    'shipped',
    'agent_received',
    'delivered',
    'inspection',
    'completed',
    'disputed',
    'cancelled',
    'refunded'
);

CREATE TYPE transaction_type AS ENUM ('goods', 'services');

CREATE TABLE transactions (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reference_code      VARCHAR(20) NOT NULL UNIQUE,
    transaction_type    transaction_type NOT NULL DEFAULT 'goods',
    title               VARCHAR(255) NOT NULL,
    description         TEXT,
    status              transaction_status NOT NULL DEFAULT 'draft',

    -- Parties
    seller_id           UUID NOT NULL REFERENCES users(id),
    buyer_id            UUID REFERENCES users(id),
    agent_id            UUID REFERENCES users(id),

    -- Financial
    amount              DECIMAL(12,2) NOT NULL,
    currency            VARCHAR(3) NOT NULL DEFAULT 'GHS',
    fee_amount          DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    fee_paid_by         VARCHAR(10) NOT NULL DEFAULT 'seller',

    -- Locations
    seller_city         VARCHAR(100) NOT NULL,
    buyer_city          VARCHAR(100),

    -- Inspection
    inspection_hours    INTEGER NOT NULL DEFAULT 48,
    inspection_ends_at  TIMESTAMPTZ,

    -- State timestamps
    agreed_at           TIMESTAMPTZ,
    funded_at           TIMESTAMPTZ,
    shipped_at          TIMESTAMPTZ,
    agent_received_at   TIMESTAMPTZ,
    delivered_at        TIMESTAMPTZ,
    completed_at        TIMESTAMPTZ,
    cancelled_at        TIMESTAMPTZ,

    -- Invite
    invite_token        VARCHAR(64) UNIQUE,
    invite_expires_at   TIMESTAMPTZ,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_txn_seller ON transactions(seller_id);
CREATE INDEX idx_txn_buyer ON transactions(buyer_id);
CREATE INDEX idx_txn_agent ON transactions(agent_id);
CREATE INDEX idx_txn_status ON transactions(status);
CREATE INDEX idx_txn_reference ON transactions(reference_code);
CREATE INDEX idx_txn_invite ON transactions(invite_token);
