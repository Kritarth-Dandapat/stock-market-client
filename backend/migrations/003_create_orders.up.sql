-- Migration: 003_create_orders.sql

CREATE TABLE IF NOT EXISTS orders (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID NOT NULL REFERENCES users(id),
    scheme_code    TEXT NOT NULL,
    scheme_name    TEXT NOT NULL DEFAULT '',
    order_type     TEXT NOT NULL,
    amount         NUMERIC(18, 4) NOT NULL,
    units          NUMERIC(18, 4),
    nav            NUMERIC(18, 4),
    status         TEXT NOT NULL DEFAULT 'pending',
    bse_order_id   TEXT,
    payment_id     TEXT,
    failure_reason TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT orders_type_check CHECK (order_type IN ('lumpsum', 'sip', 'redeem')),
    CONSTRAINT orders_status_check CHECK (
        status IN ('pending', 'submitted', 'processing', 'completed', 'failed', 'cancelled')
    )
);

CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders (user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders (status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders (created_at DESC);

CREATE TABLE IF NOT EXISTS sip_mandates (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id),
    scheme_code TEXT NOT NULL,
    scheme_name TEXT NOT NULL DEFAULT '',
    amount      NUMERIC(18, 4) NOT NULL,
    frequency   TEXT NOT NULL DEFAULT 'MONTHLY',
    start_date  DATE NOT NULL,
    end_date    DATE,
    mandate_id  TEXT NOT NULL,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sip_mandates_user_id ON sip_mandates (user_id);
