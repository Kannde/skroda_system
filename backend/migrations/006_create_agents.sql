CREATE TYPE agent_status AS ENUM ('active', 'inactive', 'suspended');

CREATE TABLE agents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id),
    city VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL,
    bio TEXT,
    rating NUMERIC(3, 2) NOT NULL DEFAULT 0.00,
    total_handled INT NOT NULL DEFAULT 0,
    status agent_status NOT NULL DEFAULT 'active',
    commission_rate NUMERIC(5, 2) NOT NULL DEFAULT 2.50,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_agents_city ON agents(city);
CREATE INDEX idx_agents_status ON agents(status);
