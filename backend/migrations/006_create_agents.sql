-- +migrate Up
CREATE TYPE agent_status AS ENUM ('pending_approval', 'active', 'suspended', 'inactive');
CREATE TYPE agent_tier AS ENUM ('bronze', 'silver', 'gold');

CREATE TABLE agent_profiles (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id             UUID NOT NULL UNIQUE REFERENCES users(id),
    business_name       VARCHAR(255) NOT NULL,
    business_address    TEXT NOT NULL,
    business_phone      VARCHAR(20) NOT NULL,
    city                VARCHAR(100) NOT NULL,
    region              VARCHAR(100) NOT NULL,
    gps_lat             DECIMAL(10,8),
    gps_lng             DECIMAL(11,8),

    business_reg_number VARCHAR(100),
    id_document_url     TEXT,
    verified            BOOLEAN NOT NULL DEFAULT FALSE,
    verified_at         TIMESTAMPTZ,

    status              agent_status NOT NULL DEFAULT 'pending_approval',
    tier                agent_tier NOT NULL DEFAULT 'bronze',
    rating              DECIMAL(3,2) NOT NULL DEFAULT 0.00,
    total_deliveries    INTEGER NOT NULL DEFAULT 0,
    successful_deliveries INTEGER NOT NULL DEFAULT 0,
    bond_amount         DECIMAL(12,2) NOT NULL DEFAULT 0.00,

    max_concurrent      INTEGER NOT NULL DEFAULT 5,
    active_deliveries   INTEGER NOT NULL DEFAULT 0,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_agent_city ON agent_profiles(city);
CREATE INDEX idx_agent_status ON agent_profiles(status);
CREATE INDEX idx_agent_tier ON agent_profiles(tier);
CREATE INDEX idx_agent_user ON agent_profiles(user_id);

CREATE TABLE audit_log (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    entity_type         VARCHAR(50) NOT NULL,
    entity_id           UUID NOT NULL,
    action              VARCHAR(100) NOT NULL,
    old_value           JSONB,
    new_value           JSONB,
    performed_by        UUID REFERENCES users(id),
    ip_address          INET,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_entity ON audit_log(entity_type, entity_id);
CREATE INDEX idx_audit_created ON audit_log(created_at);
